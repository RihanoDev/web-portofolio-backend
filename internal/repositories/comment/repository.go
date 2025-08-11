package comment

import (
	"web-porto-backend/internal/domain/models"

	"gorm.io/gorm"
)

type Repository interface {
	FindAll() ([]models.Comment, error)
	FindByID(id int) (*models.Comment, error)
	FindByPostID(postID int) ([]models.Comment, error)
	Create(comment *models.Comment) error
	Update(comment *models.Comment) error
	Delete(id int) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) FindAll() ([]models.Comment, error) {
	var comments []models.Comment
	err := r.db.Find(&comments).Error
	return comments, err
}

func (r *repository) FindByID(id int) (*models.Comment, error) {
	var comment models.Comment
	err := r.db.First(&comment, id).Error
	if err != nil {
		return nil, err
	}
	return &comment, nil
}
func (r *repository) FindByPostID(postID int) ([]models.Comment, error) {
	var comments []models.Comment
	err := r.db.Where("post_id = ?", postID).Find(&comments).Error
	return comments, err
}

func (r *repository) Create(comment *models.Comment) error {
	return r.db.Create(comment).Error
}

func (r *repository) Update(comment *models.Comment) error {
	return r.db.Save(comment).Error
}

func (r *repository) Delete(id int) error {
	return r.db.Delete(&models.Comment{}, id).Error
}
