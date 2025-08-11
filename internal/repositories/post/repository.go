package post

import (
	"web-porto-backend/internal/domain/models"

	"gorm.io/gorm"
)

type Repository interface {
	Create(post *models.Post) error
	GetByID(id uint) (*models.Post, error)
	GetAll(limit, offset int) ([]*models.Post, int64, error)
	Update(post *models.Post) error
	Delete(id uint) error
	GetBySlug(slug string) (*models.Post, error)
	GetByAuthorID(authorID int, limit, offset int) ([]*models.Post, int64, error)
	GetPublished(limit, offset int) ([]*models.Post, int64, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(post *models.Post) error {
	return r.db.Create(post).Error
}

func (r *repository) GetByID(id uint) (*models.Post, error) {
	var post models.Post
	err := r.db.Preload("Author").First(&post, id).Error
	if err != nil {
		return nil, err
	}
	return &post, nil
}

func (r *repository) GetAll(limit, offset int) ([]*models.Post, int64, error) {
	var posts []*models.Post
	var total int64

	// Count total records
	r.db.Model(&models.Post{}).Count(&total)

	// Get paginated results with preloaded associations
	err := r.db.Preload("Author").
		Limit(limit).Offset(offset).Find(&posts).Error

	return posts, total, err
}

func (r *repository) Update(post *models.Post) error {
	return r.db.Save(post).Error
}

func (r *repository) Delete(id uint) error {
	return r.db.Delete(&models.Post{}, id).Error
}

func (r *repository) GetBySlug(slug string) (*models.Post, error) {
	var post models.Post
	err := r.db.Preload("Author").
		Where("slug = ?", slug).First(&post).Error
	if err != nil {
		return nil, err
	}
	return &post, nil
}

func (r *repository) GetByAuthorID(authorID int, limit, offset int) ([]*models.Post, int64, error) {
	var posts []*models.Post
	var total int64

	query := r.db.Model(&models.Post{}).Where("author_id = ?", authorID)
	query.Count(&total)

	err := query.Preload("Author").
		Limit(limit).Offset(offset).Find(&posts).Error

	return posts, total, err
}

func (r *repository) GetPublished(limit, offset int) ([]*models.Post, int64, error) {
	var posts []*models.Post
	var total int64

	query := r.db.Model(&models.Post{}).Where("status = ?", "published")
	query.Count(&total)

	err := query.Preload("Author").
		Limit(limit).Offset(offset).Find(&posts).Error

	return posts, total, err
}
