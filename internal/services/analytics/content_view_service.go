package analytics

import (
	"time"
	"web-porto-backend/internal/domain/models"

	"gorm.io/gorm"
)

// ContentType represents the type of content for view tracking
type ContentType string

const (
	ContentTypeArticle ContentType = "article"
	ContentTypeProject ContentType = "project"
	ContentTypePage    ContentType = "page"
)

// ContentView represents a view for a specific content item
type ContentView struct {
	ID        uint   `gorm:"primarykey"`
	ContentID string `gorm:"index;not null"`
	Type      string `gorm:"index;not null"` // article, project, page
	VisitorID string `gorm:"index"`
	UserAgent string
	Referrer  string
	IP        string
	Timestamp time.Time `gorm:"index;default:CURRENT_TIMESTAMP"`
}

// ViewCountResponse represents the response for view counts
type ViewCountResponse struct {
	Count int `json:"count"`
}

// AnalyticsDataPoint represents a data point for analytics
type AnalyticsDataPoint struct {
	Date  string `json:"date"`
	Count int    `json:"count"`
}

// ContentViewService handles operations for content view tracking
type ContentViewService struct {
	db *gorm.DB
}

// NewContentViewService creates a new content view service
func NewContentViewService(db *gorm.DB) *ContentViewService {
	// Auto-migrate the content view model
	db.AutoMigrate(&ContentView{})
	return &ContentViewService{db: db}
}

// TrackContentView records a new content view
func (s *ContentViewService) TrackContentView(
	contentID string,
	contentType ContentType,
	visitorID string,
	userAgent string,
	referrer string,
	ip string,
) error {
	view := ContentView{
		ContentID: contentID,
		Type:      string(contentType),
		VisitorID: visitorID,
		UserAgent: userAgent,
		Referrer:  referrer,
		IP:        ip,
		Timestamp: time.Now(),
	}

	return s.db.Create(&view).Error
}

// GetContentViewCount returns the view count for a specific content
func (s *ContentViewService) GetContentViewCount(contentID string, contentType ContentType) (int, error) {
	var count int64
	err := s.db.Model(&ContentView{}).
		Where("content_id = ? AND type = ?", contentID, string(contentType)).
		Count(&count).Error

	return int(count), err
}

// GetContentViewAnalytics returns analytics data for a specific content
func (s *ContentViewService) GetContentViewAnalytics(
	contentID string,
	contentType ContentType,
	period string,
	limit int,
) ([]AnalyticsDataPoint, error) {
	var subQuery string

	// Determine date format and subquery based on period
	switch period {
	case "day":
		subQuery = "DATE(timestamp) as date"
	case "week":
		subQuery = "CONCAT(DATE_PART('year', timestamp), '-W', LPAD(DATE_PART('week', timestamp)::text, 2, '0')) as date"
	case "month":
		subQuery = "TO_CHAR(timestamp, 'YYYY-MM') as date"
	default:
		subQuery = "DATE(timestamp) as date"
	}

	// Raw SQL query to get daily/weekly/monthly counts
	var analytics []struct {
		Date  string
		Count int
	}

	query := s.db.Table("content_views").
		Select(subQuery+", COUNT(*) as count").
		Where("content_id = ? AND type = ?", contentID, string(contentType)).
		Group("date").
		Order("date DESC").
		Limit(limit)

	if err := query.Scan(&analytics).Error; err != nil {
		return nil, err
	}

	// Map to response format
	dataPoints := make([]AnalyticsDataPoint, len(analytics))
	for i, point := range analytics {
		dataPoints[i] = AnalyticsDataPoint{
			Date:  point.Date,
			Count: point.Count,
		}
	}

	return dataPoints, nil
}

// TrackPageView is a compatibility function for the existing PageView model
func (s *ContentViewService) TrackPageView(view *models.PageView) (struct {
	Page           string `json:"page"`
	TotalViews     int    `json:"totalViews"`
	UniqueVisitors int    `json:"uniqueVisitors"`
}, error) {
	contentView := ContentView{
		ContentID: view.Page,
		Type:      string(ContentTypePage),
		VisitorID: view.VisitorID,
		UserAgent: view.UserAgent,
		Referrer:  view.Referrer,
		IP:        view.IP,
		Timestamp: view.Timestamp,
	}

	err := s.db.Create(&contentView).Error
	if err != nil {
		return struct {
			Page           string `json:"page"`
			TotalViews     int    `json:"totalViews"`
			UniqueVisitors int    `json:"uniqueVisitors"`
		}{}, err
	}

	// Get total views for compatibility
	var totalViews int64
	s.db.Model(&ContentView{}).Where("content_id = ? AND type = ?", view.Page, string(ContentTypePage)).Count(&totalViews)

	// Get unique visitors for compatibility
	var uniqueVisitors int64
	s.db.Model(&ContentView{}).
		Where("content_id = ? AND type = ?", view.Page, string(ContentTypePage)).
		Distinct("visitor_id").
		Count(&uniqueVisitors)

	return struct {
		Page           string `json:"page"`
		TotalViews     int    `json:"totalViews"`
		UniqueVisitors int    `json:"uniqueVisitors"`
	}{
		Page:           view.Page,
		TotalViews:     int(totalViews),
		UniqueVisitors: int(uniqueVisitors),
	}, nil
}
