package tag

import (
	"web-porto-backend/internal/domain/models"

	"gorm.io/gorm"
)

// Repository handles tag database operations
type Repository interface {
	GetAll() ([]models.Tag, error)
	GetByID(id int) (*models.Tag, error)
	GetByName(name string) (*models.Tag, error)
	Create(tag *models.Tag) (*models.Tag, error)
	Update(tag *models.Tag) (*models.Tag, error)
	Delete(id int) error
}

type repository struct {
	db *gorm.DB
}

// NewRepository creates a new tag repository
func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) GetAll() ([]models.Tag, error) {
	var tags []models.Tag
	if err := r.db.Find(&tags).Error; err != nil {
		return nil, err
	}
	return tags, nil
}

func (r *repository) GetByID(id int) (*models.Tag, error) {
	var tag models.Tag
	if err := r.db.Where("id = ?", id).First(&tag).Error; err != nil {
		return nil, err
	}
	return &tag, nil
}

func (r *repository) GetByName(name string) (*models.Tag, error) {
	var tag models.Tag
	if err := r.db.Where("name = ?", name).First(&tag).Error; err != nil {
		return nil, err
	}
	return &tag, nil
}

func (r *repository) Create(tag *models.Tag) (*models.Tag, error) {
	if err := r.db.Create(tag).Error; err != nil {
		return nil, err
	}
	return tag, nil
}

func (r *repository) Update(tag *models.Tag) (*models.Tag, error) {
	if err := r.db.Save(tag).Error; err != nil {
		return nil, err
	}
	return tag, nil
}

func (r *repository) Delete(id int) error {
	return r.db.Delete(&models.Tag{}, id).Error
}
