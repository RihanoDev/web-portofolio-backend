package project

import (
	"web-porto-backend/internal/domain/models"

	"gorm.io/gorm"
)

type Repository interface {
	Create(project *models.Project) error
	GetByID(id string) (*models.Project, error)
	GetAll(limit, offset int) ([]*models.Project, int64, error)
	Update(project *models.Project) error
	Delete(id string) error
	GetBySlug(slug string) (*models.Project, error)
	GetByCategorySlug(slug string, limit, offset int) ([]*models.Project, int64, error)
	UpdateProjectTechnologies(projectID string, technologyIDs []int) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(project *models.Project) error {
	return r.db.Create(project).Error
}

func (r *repository) GetByID(id string) (*models.Project, error) {
	// Check if ID is a valid UUID
	if len(id) > 0 && id[:5] == "temp-" {
		return nil, gorm.ErrRecordNotFound
	}

	var project models.Project
	err := r.db.Preload("Author").Preload("Category").Preload("Tags").Where("id = ?", id).First(&project).Error
	if err != nil {
		return nil, err
	}
	return &project, nil
}

func (r *repository) GetAll(limit, offset int) ([]*models.Project, int64, error) {
	var projects []*models.Project
	var total int64

	// Count total records
	r.db.Model(&models.Project{}).Count(&total)

	// Get paginated results with preloaded associations
	err := r.db.Preload("Author").Preload("Category").Preload("Tags").
		Limit(limit).Offset(offset).Find(&projects).Error

	return projects, total, err
}

func (r *repository) Update(project *models.Project) error {
	return r.db.Save(project).Error
}

func (r *repository) Delete(id string) error {
	// Check if ID is a valid UUID
	if len(id) > 0 && id[:5] == "temp-" {
		// Nothing to delete since it was a temporary ID
		return nil
	}
	return r.db.Where("id = ?", id).Delete(&models.Project{}).Error
}

func (r *repository) GetBySlug(slug string) (*models.Project, error) {
	var project models.Project
	err := r.db.Preload("Author").Preload("Category").Preload("Tags").
		Where("slug = ?", slug).First(&project).Error
	if err != nil {
		return nil, err
	}
	return &project, nil
}

func (r *repository) GetByCategorySlug(slug string, limit, offset int) ([]*models.Project, int64, error) {
	var projects []*models.Project
	var total int64
	var categoryID int

	// First, get the category ID by slug
	if err := r.db.Model(&models.Category{}).Where("slug = ?", slug).Select("id").Scan(&categoryID).Error; err != nil {
		return nil, 0, err
	}

	// Then, get projects with that category ID
	query := r.db.Model(&models.Project{}).Where("category_id = ?", categoryID)
	query.Count(&total)

	err := query.Preload("Author").Preload("Category").
		Limit(limit).Offset(offset).Find(&projects).Error

	return projects, total, err
}

func (r *repository) UpdateProjectTechnologies(projectID string, technologyIDs []int) error {
	// Find the project
	var project models.Project
	if err := r.db.Where("id = ?", projectID).First(&project).Error; err != nil {
		return err
	}

	// Get tags by IDs
	var tags []models.Tag
	if len(technologyIDs) > 0 {
		if err := r.db.Where("id IN ?", technologyIDs).Find(&tags).Error; err != nil {
			return err
		}
	}

	// Replace associations using GORM
	if err := r.db.Model(&project).Association("Tags").Replace(tags); err != nil {
		return err
	}

	return nil
}
