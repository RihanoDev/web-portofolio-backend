package migrations

import (
	"log"
	"math/rand"
	"time"

	"web-porto-backend/internal/domain/models"

	"gorm.io/gorm"
)

// Initialize random number source for Go 1.20+ compatibility
var rng = rand.New(rand.NewSource(time.Now().UnixNano()))

// Sample data generators
func getSampleVisitorIDs() []string {
	return []string{
		"visitor_1660000000_abcdef",
		"visitor_1670000000_ghijkl",
		"visitor_1680000000_mnopqr",
		"visitor_1690000000_stuvwx",
		"visitor_1700000000_yzabcd",
		"visitor_1710000000_efghij",
		"visitor_1720000000_klmnop",
		"visitor_1730000000_qrstuv",
	}
}

func getSampleUserAgents() []string {
	return []string{
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.0 Safari/605.1.15",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 14_6 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.0 Mobile/15E148 Safari/604.1",
		"Mozilla/5.0 (iPad; CPU OS 14_6 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.0 Mobile/15E148 Safari/604.1",
		"Mozilla/5.0 (Linux; Android 11; SM-G991B) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.120 Mobile Safari/537.36",
	}
}

func getSampleLocations() []struct {
	Country string
	City    string
} {
	return []struct {
		Country string
		City    string
	}{
		{"Indonesia", "Jakarta"},
		{"Indonesia", "Bandung"},
		{"Indonesia", "Surabaya"},
		{"Singapore", "Singapore"},
		{"Malaysia", "Kuala Lumpur"},
		{"United States", "San Francisco"},
		{"United Kingdom", "London"},
		{"Japan", "Tokyo"},
	}
}

func getSamplePages() []string {
	return []string{
		"/",
		"/about",
		"/projects",
		"/articles",
		"/contact",
		"/projects/web-portfolio",
		"/projects/e-commerce-platform",
		"/projects/machine-learning-app",
		"/articles/getting-started-with-go",
		"/articles/react-best-practices",
	}
}

// Generate a single page view with random data
func generatePageView(date time.Time, visitorIDs []string, pages []string, userAgents []string, locations []struct {
	Country string
	City    string
}) *models.PageView {
	// Random hour and minute
	hour := rng.Intn(24)
	minute := rng.Intn(60)
	second := rng.Intn(60)
	timestamp := time.Date(date.Year(), date.Month(), date.Day(), hour, minute, second, 0, date.Location())

	// Generate random visitor and page
	visitor := visitorIDs[rng.Intn(len(visitorIDs))]
	page := pages[rng.Intn(len(pages))]
	userAgent := userAgents[rng.Intn(len(userAgents))]
	location := locations[rng.Intn(len(locations))]

	// Home page should have more views
	if rng.Intn(10) < 3 {
		page = "/"
	}

	// Create page view
	return &models.PageView{
		Page:      page,
		VisitorID: visitor,
		UserAgent: userAgent,
		IP:        "192.168.1." + string(rune(rng.Intn(254)+1)),
		Country:   location.Country,
		City:      location.City,
		Timestamp: timestamp,
	}
}

// SeedAnalytics creates sample analytics data if the page_views table is empty
func SeedAnalytics(db *gorm.DB) error {
	// Check if page_views table is empty
	var count int64
	if err := db.Model(&models.PageView{}).Count(&count).Error; err != nil {
		return err
	}

	// If we already have data, skip seeding
	if count > 0 {
		log.Printf("Analytics data already exists, skipping seed")
		return nil
	}

	// Get sample data
	visitorIDs := getSampleVisitorIDs()
	userAgents := getSampleUserAgents()
	locations := getSampleLocations()
	pages := getSamplePages()

	// Generate random page views for the last 90 days
	now := time.Now()
	var pageViews []*models.PageView

	// Generate page views for each day
	for i := 90; i >= 0; i-- {
		date := now.AddDate(0, 0, -i)
		pageViews = append(pageViews, generateDailyViews(date, i, visitorIDs, pages, userAgents, locations)...)
	}

	// Save the generated views to the database
	if err := savePageViews(db, pageViews); err != nil {
		return err
	}

	log.Printf("Seeded %d sample page views", len(pageViews))
	return nil
}

// Generate page views for a single day
func generateDailyViews(date time.Time, daysAgo int, visitorIDs []string, pages []string, userAgents []string, locations []struct {
	Country string
	City    string
}) []*models.PageView {
	var dailyViews []*models.PageView

	// More views on weekdays, fewer on weekends
	viewsToGenerate := 5
	weekday := date.Weekday()
	if weekday != time.Saturday && weekday != time.Sunday {
		viewsToGenerate = 10
	}

	// Recent dates have more views
	if daysAgo < 30 {
		viewsToGenerate *= 2
	}

	// Generate views for this day
	for j := 0; j < viewsToGenerate; j++ {
		pageView := generatePageView(date, visitorIDs, pages, userAgents, locations)
		dailyViews = append(dailyViews, pageView)
	}

	return dailyViews
}

// Save page views to the database in batches
func savePageViews(db *gorm.DB, pageViews []*models.PageView) error {
	// Create analytics data in batches to avoid potential DB issues
	batchSize := 100
	for i := 0; i < len(pageViews); i += batchSize {
		end := i + batchSize
		if end > len(pageViews) {
			end = len(pageViews)
		}
		batch := pageViews[i:end]

		// Use Create for each record in the batch
		for _, view := range batch {
			if err := db.Create(view).Error; err != nil {
				return err
			}
		}
	}

	return nil
}
