package migrations

import (
	"log"

	"gorm.io/gorm"

	"web-porto-backend/internal/domain/models"
)

// SeedCategories creates initial categories if the categories table is empty
func SeedCategories(db *gorm.DB) error {
	// Ensure categories table exists
	if err := db.AutoMigrate(&models.Category{}); err != nil {
		return err
	}

	// Check if categories exist
	var count int64
	if err := db.Model(&models.Category{}).Count(&count).Error; err != nil {
		return err
	}

	// If categories already exist, do nothing
	if count > 0 {
		return nil
	}

	// Create initial categories
	categories := []models.Category{
		{Name: "Web Development", Slug: "web-development", Description: "Projects related to web development"},
		{Name: "Mobile Development", Slug: "mobile-development", Description: "Projects related to mobile app development"},
		{Name: "Data Science", Slug: "data-science", Description: "Projects related to data science and analytics"},
		{Name: "DevOps", Slug: "devops", Description: "Projects related to DevOps and infrastructure"},
		{Name: "Machine Learning", Slug: "machine-learning", Description: "Projects related to machine learning and AI"},
	}

	if err := db.Create(&categories).Error; err != nil {
		return err
	}

	log.Printf("Seeded %d initial categories", len(categories))
	return nil
}
