package category

import (
	"web-porto-backend/internal/domain/models"

	"gorm.io/gorm"
)

type Repository interface {
	FindAll() ([]models.Category, error)
	FindByID(id int) (*models.Category, error)
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

func (r *repository) Create(category *models.Category) error {
	return r.db.Create(category).Error
}

func (r *repository) Update(category *models.Category) error {
	return r.db.Save(category).Error
}

func (r *repository) Delete(id int) error {
	return r.db.Delete(&models.Category{}, id).Error
}
