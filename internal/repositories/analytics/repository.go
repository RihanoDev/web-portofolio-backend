package analytics

import (
	"strings"
	"time"
	applog "web-porto-backend/common/logger"
	"web-porto-backend/internal/domain/models"

	"gorm.io/gorm"
)

type ViewStats struct {
	Total  int64
	Today  int64
	Week   int64
	Month  int64
	Unique int64
}

type TimeSeriesPoint struct {
	Bucket time.Time `json:"bucket"`
	Count  int64     `json:"count"`
}

type Repository interface {
	TrackView(v *models.PageView) error
	GetStats(page string) (*ViewStats, error)
	GetStatsWithFilter(page string, start, end *time.Time, country string) (*ViewStats, error)
	GetTimeSeries(page string, start, end time.Time, interval string) ([]TimeSeriesPoint, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) TrackView(v *models.PageView) error {
	log := applog.GetLogger().WithFields(applog.Fields{"repo": "analytics", "method": "TrackView"})
	// Basic bot filtering
	ua := strings.ToLower(v.UserAgent)
	if ua != "" {
		if strings.Contains(ua, "bot") || strings.Contains(ua, "crawl") || strings.Contains(ua, "spider") {
			log.Debug("skip bot user-agent", applog.Fields{"ua": v.UserAgent})
			return nil
		}
	}
	// Simple dedup window: 30s same visitor/page
	var recent int64
	window := time.Now().Add(-30 * time.Second)
	r.db.Model(&models.PageView{}).
		Where("page = ? AND visitor_id = ? AND timestamp >= ?", v.Page, v.VisitorID, window).
		Count(&recent)
	if recent > 0 {
		log.Debug("dedup skip within window", applog.Fields{"page": v.Page, "visitor": v.VisitorID})
		return nil
	}
	if err := r.db.Create(v).Error; err != nil {
		log.Error("db create failed", applog.Fields{"error": err.Error()})
		return err
	}

	// Invalidate cache entries to force refresh on next request
	// This avoids serving stale data after new views are added
	for key := range statsMemCache {
		if key == "all" || (v.Page != "" && key == v.Page) ||
			(strings.HasPrefix(key, "filter:") && strings.Contains(key, "p:"+v.Page+":")) {
			delete(statsMemCache, key)
		}
	}

	log.Info("view tracked, cache invalidated", applog.Fields{"page": v.Page, "visitor": v.VisitorID})
	return nil
}

// In-memory cache for analytics stats with 5 minute TTL
type statsCache struct {
	stats     *ViewStats
	timestamp time.Time
	key       string
}

// Global cache map (page -> stats)
var statsMemCache = map[string]statsCache{}

func (r *repository) GetStats(page string) (*ViewStats, error) {
	log := applog.GetLogger().WithFields(applog.Fields{"repo": "analytics", "method": "GetStats", "page": page})

	// Cache key
	cacheKey := "all"
	if page != "" {
		cacheKey = page
	}

	// Check cache first (5 minute TTL)
	if cache, ok := statsMemCache[cacheKey]; ok {
		if time.Since(cache.timestamp) < 5*time.Minute {
			log.Info("using cached stats", applog.Fields{"page": page, "age": time.Since(cache.timestamp).Seconds()})
			return cache.stats, nil
		}
	}

	var stats ViewStats
	now := time.Now()
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	weekAgo := now.AddDate(0, 0, -7)
	monthAgo := now.AddDate(0, -1, 0)

	// Use a single optimized query with CASE expressions for PostgreSQL
	// This query gets all counts in one go to avoid multiple round-trips
	query := `
		SELECT 
			COUNT(*) as total,
			COUNT(CASE WHEN timestamp >= ? THEN 1 END) as today,
			COUNT(CASE WHEN timestamp >= ? THEN 1 END) as week,
			COUNT(CASE WHEN timestamp >= ? THEN 1 END) as month,
			COUNT(DISTINCT visitor_id) as unique_visitors
		FROM page_views
	`

	args := []interface{}{todayStart, weekAgo, monthAgo}

	if page != "" {
		query += " WHERE page = ?"
		args = append(args, page)
	}

	// Define a temporary struct to receive the query results
	var result struct {
		Total          int64 `gorm:"column:total"`
		Today          int64 `gorm:"column:today"`
		Week           int64 `gorm:"column:week"`
		Month          int64 `gorm:"column:month"`
		UniqueVisitors int64 `gorm:"column:unique_visitors"`
	}

	// Run optimized query
	err := r.db.Raw(query, args...).Scan(&result).Error

	// Copy results to stats struct
	if err == nil {
		stats.Total = result.Total
		stats.Today = result.Today
		stats.Week = result.Week
		stats.Month = result.Month
		stats.Unique = result.UniqueVisitors
	}

	if err != nil {
		log.Error("stats query failed", applog.Fields{"error": err.Error()})
		return nil, err
	}

	// Cache the result
	statsMemCache[cacheKey] = statsCache{
		stats:     &stats,
		timestamp: now,
		key:       cacheKey,
	}

	log.Info("stats computed and cached", applog.Fields{"total": stats.Total, "unique": stats.Unique})
	return &stats, nil
}

// GetStatsWithFilter returns stats within a custom date range and optional filters
func (r *repository) GetStatsWithFilter(page string, start, end *time.Time, country string) (*ViewStats, error) {
	log := applog.GetLogger().WithFields(applog.Fields{"repo": "analytics", "method": "GetStatsWithFilter", "page": page})

	// Construct cache key for this filter combination
	cacheKey := "filter:"
	if page != "" {
		cacheKey += "p:" + page + ":"
	}
	if start != nil {
		cacheKey += "s:" + start.Format(time.RFC3339) + ":"
	}
	if end != nil {
		cacheKey += "e:" + end.Format(time.RFC3339) + ":"
	}
	if country != "" {
		cacheKey += "c:" + country
	}

	// Check cache first (2 minute TTL for filtered results)
	if cache, ok := statsMemCache[cacheKey]; ok {
		if time.Since(cache.timestamp) < 2*time.Minute {
			log.Info("using cached filtered stats", applog.Fields{
				"page": page,
				"age":  time.Since(cache.timestamp).Seconds(),
			})
			return cache.stats, nil
		}
	}

	var stats ViewStats
	now := time.Now()
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	weekAgo := now.AddDate(0, 0, -7)
	monthAgo := now.AddDate(0, -1, 0)

	// Build the dynamic SQL query with all conditions
	query := `
		SELECT 
			COUNT(*) as total,
			COUNT(CASE WHEN timestamp >= ? THEN 1 END) as today,
			COUNT(CASE WHEN timestamp >= ? THEN 1 END) as week,
			COUNT(CASE WHEN timestamp >= ? THEN 1 END) as month,
			COUNT(DISTINCT visitor_id) as unique_visitors
		FROM page_views
		WHERE 1=1
	`

	args := []interface{}{todayStart, weekAgo, monthAgo}

	// Add filters
	if page != "" {
		query += " AND page = ?"
		args = append(args, page)
	}
	if start != nil {
		query += " AND timestamp >= ?"
		args = append(args, *start)
	}
	if end != nil {
		query += " AND timestamp <= ?"
		args = append(args, *end)
	}
	if country != "" {
		query += " AND country = ?"
		args = append(args, country)
	}

	// Define a temporary struct to receive the query results
	var result struct {
		Total          int64 `gorm:"column:total"`
		Today          int64 `gorm:"column:today"`
		Week           int64 `gorm:"column:week"`
		Month          int64 `gorm:"column:month"`
		UniqueVisitors int64 `gorm:"column:unique_visitors"`
	}

	// Run query
	err := r.db.Raw(query, args...).Scan(&result).Error
	if err != nil {
		log.Error("filtered stats query failed", applog.Fields{"error": err.Error()})
		return nil, err
	}

	// Copy results
	stats.Total = result.Total
	stats.Today = result.Today
	stats.Week = result.Week
	stats.Month = result.Month
	stats.Unique = result.UniqueVisitors

	// Cache results
	statsMemCache[cacheKey] = statsCache{
		stats:     &stats,
		timestamp: now,
		key:       cacheKey,
	}

	log.Info("filtered stats computed and cached", applog.Fields{"total": stats.Total, "unique": stats.Unique})
	return &stats, nil
}

// GetTimeSeries returns aggregated counts by hour or day using Postgres date_trunc
func (r *repository) GetTimeSeries(page string, start, end time.Time, interval string) ([]TimeSeriesPoint, error) {
	log := applog.GetLogger().WithFields(applog.Fields{"repo": "analytics", "method": "GetTimeSeries", "page": page, "interval": interval})
	// Validate interval
	trunc := "day"
	switch strings.ToLower(interval) {
	case "hour":
		trunc = "hour"
	case "day":
		trunc = "day"
	default:
		trunc = "day"
	}

	var points []TimeSeriesPoint
	baseSQL := "SELECT date_trunc(?, timestamp) AS bucket, COUNT(*) AS count FROM page_views WHERE timestamp BETWEEN ? AND ?"
	args := []interface{}{trunc, start, end}
	if page != "" {
		baseSQL += " AND page = ?"
		args = append(args, page)
	}
	baseSQL += " GROUP BY bucket ORDER BY bucket"

	if err := r.db.Raw(baseSQL, args...).Scan(&points).Error; err != nil {
		log.Error("timeseries query failed", applog.Fields{"error": err.Error()})
		return nil, err
	}
	log.Info("timeseries computed", applog.Fields{"points": len(points)})
	return points, nil
}
