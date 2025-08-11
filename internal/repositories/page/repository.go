package page

import (
	"web-porto-backend/internal/domain/models"

	"gorm.io/gorm"
)

type Repository interface {
	Create(page *models.Page) error
	GetByID(id uint) (*models.Page, error)
	GetAll(limit, offset int) ([]*models.Page, int64, error)
	Update(page *models.Page) error
	Delete(id uint) error
	GetBySlug(slug string) (*models.Page, error)
	GetPublished(limit, offset int) ([]*models.Page, int64, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(page *models.Page) error {
	return r.db.Create(page).Error
}

func (r *repository) GetByID(id uint) (*models.Page, error) {
	var page models.Page
	err := r.db.First(&page, id).Error
	if err != nil {
		return nil, err
	}
	return &page, nil
}

func (r *repository) GetAll(limit, offset int) ([]*models.Page, int64, error) {
	var pages []*models.Page
	var total int64

	// Count total records
	r.db.Model(&models.Page{}).Count(&total)

	// Get paginated results
	err := r.db.Limit(limit).Offset(offset).Find(&pages).Error

	return pages, total, err
}

func (r *repository) Update(page *models.Page) error {
	return r.db.Save(page).Error
}

func (r *repository) Delete(id uint) error {
	return r.db.Delete(&models.Page{}, id).Error
}

func (r *repository) GetBySlug(slug string) (*models.Page, error) {
	var page models.Page
	err := r.db.Where("slug = ?", slug).First(&page).Error
	if err != nil {
		return nil, err
	}
	return &page, nil
}

func (r *repository) GetPublished(limit, offset int) ([]*models.Page, int64, error) {
	var pages []*models.Page
	var total int64

	query := r.db.Model(&models.Page{}).Where("status = ?", "published")
	query.Count(&total)

	err := query.Limit(limit).Offset(offset).Find(&pages).Error

	return pages, total, err
}
