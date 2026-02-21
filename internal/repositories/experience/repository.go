package experience

import (
	"web-porto-backend/internal/domain/models"

	"gorm.io/gorm"
)

type Repository interface {
	Create(experience *models.Experience) error
	GetByID(id int) (*models.Experience, error)
	GetAll(limit, offset int) ([]*models.Experience, int64, error)
	Update(experience *models.Experience) error
	Delete(id int) error
	GetCurrent() ([]*models.Experience, error)
	UpdateExperienceTechnologies(experienceID int, technologyIDs []int) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(experience *models.Experience) error {
	return r.db.Create(experience).Error
}

func (r *repository) GetByID(id int) (*models.Experience, error) {
	var experience models.Experience
	err := r.db.Preload("Technologies").Where("id = ?", id).First(&experience).Error
	if err != nil {
		return nil, err
	}
	return &experience, nil
}

func (r *repository) GetAll(limit, offset int) ([]*models.Experience, int64, error) {
	var experiences []*models.Experience
	var total int64

	// Count total records
	r.db.Model(&models.Experience{}).Count(&total)

	// Get paginated results ordered by start_date desc with preloaded technologies
	err := r.db.Preload("Technologies").Order("start_date DESC").
		Limit(limit).Offset(offset).Find(&experiences).Error

	return experiences, total, err
}

func (r *repository) Update(experience *models.Experience) error {
	return r.db.Save(experience).Error
}

func (r *repository) Delete(id int) error {
	return r.db.Delete(&models.Experience{}, id).Error
}

func (r *repository) GetCurrent() ([]*models.Experience, error) {
	var experiences []*models.Experience
	err := r.db.Preload("Technologies").Where("current = ?", true).
		Order("start_date DESC").
		Find(&experiences).Error

	return experiences, err
}

func (r *repository) UpdateExperienceTechnologies(experienceID int, technologyIDs []int) error {
	// Get the experience instance
	var experience models.Experience
	if err := r.db.First(&experience, experienceID).Error; err != nil {
		return err
	}

	// Get the tag records
	var tags []models.Tag
	if len(technologyIDs) > 0 {
		if err := r.db.Where("id IN ?", technologyIDs).Find(&tags).Error; err != nil {
			return err
		}
	}

	// Replace the associations
	return r.db.Model(&experience).Association("Technologies").Replace(tags)
}
