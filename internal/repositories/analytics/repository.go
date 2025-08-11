package analytics

import (
	"strings"
	"time"
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
	// Basic bot filtering
	ua := strings.ToLower(v.UserAgent)
	if ua != "" {
		if strings.Contains(ua, "bot") || strings.Contains(ua, "crawl") || strings.Contains(ua, "spider") {
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
		return nil
	}
	return r.db.Create(v).Error
}

func (r *repository) GetStats(page string) (*ViewStats, error) {
	var stats ViewStats
	now := time.Now()
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	weekAgo := now.AddDate(0, 0, -7)
	monthAgo := now.AddDate(0, -1, 0)
	const wherePage = "page = ?"
	const whereSince = "timestamp >= ?"

	query := r.db.Model(&models.PageView{})
	if page != "" {
		query = query.Where(wherePage, page)
	}

	// Total
	if err := query.Count(&stats.Total).Error; err != nil {
		return nil, err
	}

	// Today
	qToday := r.db.Model(&models.PageView{}).Where(whereSince, todayStart)
	if page != "" {
		qToday = qToday.Where(wherePage, page)
	}
	if err := qToday.Count(&stats.Today).Error; err != nil {
		return nil, err
	}

	// Week
	qWeek := r.db.Model(&models.PageView{}).Where(whereSince, weekAgo)
	if page != "" {
		qWeek = qWeek.Where(wherePage, page)
	}
	if err := qWeek.Count(&stats.Week).Error; err != nil {
		return nil, err
	}

	// Month
	qMonth := r.db.Model(&models.PageView{}).Where(whereSince, monthAgo)
	if page != "" {
		qMonth = qMonth.Where(wherePage, page)
	}
	if err := qMonth.Count(&stats.Month).Error; err != nil {
		return nil, err
	}

	// Unique by visitor
	qUnique := r.db.Model(&models.PageView{})
	if page != "" {
		qUnique = qUnique.Where(wherePage, page)
	}
	if err := qUnique.Distinct("visitor_id").Count(&stats.Unique).Error; err != nil {
		return nil, err
	}

	return &stats, nil
}

// GetStatsWithFilter returns stats within a custom date range and optional filters
func (r *repository) GetStatsWithFilter(page string, start, end *time.Time, country string) (*ViewStats, error) {
	var stats ViewStats
	const wherePage = "page = ?"
	const whereSince = "timestamp >= ?"
	const whereBefore = "timestamp <= ?"
	const whereCountry = "country = ?"

	base := r.db.Model(&models.PageView{})
	if page != "" {
		base = base.Where(wherePage, page)
	}
	if start != nil {
		base = base.Where(whereSince, *start)
	}
	if end != nil {
		base = base.Where(whereBefore, *end)
	}
	if country != "" {
		base = base.Where(whereCountry, country)
	}

	// Total in range
	if err := base.Count(&stats.Total).Error; err != nil {
		return nil, err
	}

	// Today in range (if start/end don't already constrain today)
	// Compute today start and ensure it's within [start,end] window
	now := time.Now()
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	qToday := base
	qToday = qToday.Where(whereSince, todayStart)
	if err := qToday.Count(&stats.Today).Error; err != nil {
		return nil, err
	}

	// Week in range
	weekAgo := now.AddDate(0, 0, -7)
	qWeek := base.Where(whereSince, weekAgo)
	if err := qWeek.Count(&stats.Week).Error; err != nil {
		return nil, err
	}

	// Month in range
	monthAgo := now.AddDate(0, -1, 0)
	qMonth := base.Where(whereSince, monthAgo)
	if err := qMonth.Count(&stats.Month).Error; err != nil {
		return nil, err
	}

	// Unique by visitor within range
	qUnique := base.Distinct("visitor_id")
	if err := qUnique.Count(&stats.Unique).Error; err != nil {
		return nil, err
	}

	return &stats, nil
}

// GetTimeSeries returns aggregated counts by hour or day using Postgres date_trunc
func (r *repository) GetTimeSeries(page string, start, end time.Time, interval string) ([]TimeSeriesPoint, error) {
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
		return nil, err
	}
	return points, nil
}
