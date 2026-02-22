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
	UpdateProjectTags(projectID string, tagIDs []int) error
	UpdateProjectCategories(projectID string, categoryIDs []int) error
	UpdateProjectImages(projectID string, images []models.ProjectImage) error
	UpdateProjectVideos(projectID string, videos []models.ProjectVideo) error
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
	err := r.db.Preload("Author").
		Preload("Category").
		Preload("Categories").
		Preload("Technologies").
		Preload("Tags").
		Preload("Images").
		Preload("Videos").
		Where("id = ?", id).First(&project).Error
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
	err := r.db.Preload("Author").
		Preload("Category").
		Preload("Categories").
		Preload("Technologies").
		Preload("Tags").
		Preload("Images").
		Preload("Videos").
		Limit(limit).Offset(offset).Find(&projects).Error

	return projects, total, err
}

func (r *repository) Update(project *models.Project) error {
	return r.db.Save(project).Error
}

func (r *repository) Delete(id string) error {
	if len(id) > 0 && id[:5] == "temp-" {
		return nil
	}

	var project models.Project
	if err := r.db.First(&project, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil
		}
		return err
	}

	// Clean up many-to-many relationships
	_ = r.db.Model(&project).Association("Categories").Clear()
	_ = r.db.Model(&project).Association("Technologies").Clear()
	_ = r.db.Model(&project).Association("Tags").Clear()

	return r.db.Delete(&project).Error
}

func (r *repository) GetBySlug(slug string) (*models.Project, error) {
	var project models.Project
	err := r.db.Preload("Author").
		Preload("Category").
		Preload("Categories").
		Preload("Technologies").
		Preload("Tags").
		Preload("Images").
		Preload("Videos").
		Where("slug = ?", slug).First(&project).Error
	if err != nil {
		return nil, err
	}
	return &project, nil
}

func (r *repository) GetByCategorySlug(slug string, limit, offset int) ([]*models.Project, int64, error) {
	var projects []*models.Project
	var total int64

	query := r.db.Model(&models.Project{}).
		Joins("JOIN project_categories pc ON pc.project_id = projects.id").
		Joins("JOIN categories c ON c.id = pc.category_id").
		Where("c.slug = ?", slug)

	query.Count(&total)

	err := query.Preload("Author").Preload("Categories").
		Limit(limit).Offset(offset).Find(&projects).Error

	return projects, total, err
}

func (r *repository) UpdateProjectTechnologies(projectID string, technologyIDs []int) error {
	var project models.Project
	if err := r.db.First(&project, "id = ?", projectID).Error; err != nil {
		return err
	}
	var tags []models.Tag
	if len(technologyIDs) > 0 {
		r.db.Where("id IN ?", technologyIDs).Find(&tags)
	}
	return r.db.Model(&project).Association("Technologies").Replace(tags)
}

func (r *repository) UpdateProjectTags(projectID string, tagIDs []int) error {
	var project models.Project
	if err := r.db.First(&project, "id = ?", projectID).Error; err != nil {
		return err
	}
	var tags []models.Tag
	if len(tagIDs) > 0 {
		r.db.Where("id IN ?", tagIDs).Find(&tags)
	}
	return r.db.Model(&project).Association("Tags").Replace(tags)
}

func (r *repository) UpdateProjectCategories(projectID string, categoryIDs []int) error {
	var project models.Project
	if err := r.db.First(&project, "id = ?", projectID).Error; err != nil {
		return err
	}
	var categories []models.Category
	if len(categoryIDs) > 0 {
		r.db.Where("id IN ?", categoryIDs).Find(&categories)
	}
	return r.db.Model(&project).Association("Categories").Replace(categories)
}

func (r *repository) UpdateProjectImages(projectID string, images []models.ProjectImage) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("project_id = ?", projectID).Delete(&models.ProjectImage{}).Error; err != nil {
			return err
		}
		if len(images) > 0 {
			if err := tx.Create(&images).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *repository) UpdateProjectVideos(projectID string, videos []models.ProjectVideo) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("project_id = ?", projectID).Delete(&models.ProjectVideo{}).Error; err != nil {
			return err
		}
		if len(videos) > 0 {
			if err := tx.Create(&videos).Error; err != nil {
				return err
			}
		}
		return nil
	})
}
