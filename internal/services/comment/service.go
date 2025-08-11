package comment

import (
	"web-porto-backend/internal/domain/models"
	"web-porto-backend/internal/repositories/comment"
)

type Service interface {
	GetAll() ([]models.Comment, error)
	GetByID(id int) (*models.Comment, error)
	GetByPostID(postID int) ([]models.Comment, error)
	Create(comment *models.Comment) error
	Update(comment *models.Comment) error
	Delete(id int) error
}

type service struct {
	repo comment.Repository
}

func NewService(repo comment.Repository) Service {
	return &service{repo}
}

func (s *service) GetAll() ([]models.Comment, error) {
	return s.repo.FindAll()
}

func (s *service) GetByID(id int) (*models.Comment, error) {
	return s.repo.FindByID(id)
}
func (s *service) GetByPostID(postID int) ([]models.Comment, error) {
	return s.repo.FindByPostID(postID)
}

func (s *service) Create(comment *models.Comment) error {
	return s.repo.Create(comment)
}

func (s *service) Update(comment *models.Comment) error {
	return s.repo.Update(comment)
}

func (s *service) Delete(id int) error {
	return s.repo.Delete(id)
}
