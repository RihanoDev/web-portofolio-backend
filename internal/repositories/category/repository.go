package category

import (
	"fmt"
	"regexp"
	"strings"
	"time"
	"web-porto-backend/internal/domain/models"

	"gorm.io/gorm"
)

type Repository interface {
	FindAll() ([]models.Category, error)
	FindByID(id int) (*models.Category, error)
	FindBySlug(slug string) (*models.Category, error)
	FindByName(name string) (*models.Category, error)
	Create(category *models.Category) error
	Update(category *models.Category) error
	Delete(id int) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) FindAll() ([]models.Category, error) {
	var categories []models.Category
	err := r.db.Find(&categories).Error
	return categories, err
}

func (r *repository) FindByID(id int) (*models.Category, error) {
	var category models.Category
	err := r.db.First(&category, id).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *repository) FindByName(name string) (*models.Category, error) {
	var category models.Category
	err := r.db.Where("name = ?", name).First(&category).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *repository) FindBySlug(slug string) (*models.Category, error) {
	var category models.Category
	err := r.db.Where("slug = ?", slug).First(&category).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *repository) FindByName(name string) (*models.Category, error) {
	var category models.Category
	err := r.db.Where("name = ?", name).First(&category).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *repository) Create(category *models.Category) error {
	// Ensure slug is not empty
	if category.Slug == "" && category.Name != "" {
		// Generate slug from name
		category.Slug = generateSlug(category.Name)

		// Check if slug already exists, if so add a unique suffix
		var count int64
		r.db.Model(&models.Category{}).Where("slug = ?", category.Slug).Count(&count)
		if count > 0 {
			// Add timestamp as suffix to make unique
			category.Slug = category.Slug + "-" + generateTimestampSuffix()
		}
	}
	return r.db.Create(category).Error
}

func (r *repository) Update(category *models.Category) error {
	// Ensure slug is not empty when updating
	if category.Slug == "" && category.Name != "" {
		category.Slug = generateSlug(category.Name)
	}
	return r.db.Save(category).Error
}

func (r *repository) Delete(id int) error {
	return r.db.Delete(&models.Category{}, id).Error
}

// Helper functions for slug generation

// generateSlug creates a URL-friendly slug from a string
func generateSlug(input string) string {
	// Convert to lowercase
	slug := strings.ToLower(input)

	// Replace spaces with hyphens
	slug = strings.ReplaceAll(slug, " ", "-")

	// Remove special characters
	reg := regexp.MustCompile(`[^a-z0-9\-]`)
	slug = reg.ReplaceAllString(slug, "")

	// Replace multiple hyphens with single hyphen
	reg = regexp.MustCompile(`\-+`)
	slug = reg.ReplaceAllString(slug, "-")

	// Trim hyphens from start and end
	slug = strings.Trim(slug, "-")

	return slug
}

// generateTimestampSuffix creates a timestamp-based suffix for making slugs unique
func generateTimestampSuffix() string {
	return fmt.Sprintf("%d", time.Now().UnixNano()%100000)
}
