package category

import (
	"web-porto-backend/internal/domain/models"
	"web-porto-backend/internal/repositories/category"
)

type Service interface {
	GetAll() ([]models.Category, error)
	GetByID(id int) (*models.Category, error)
	Create(category *models.Category) error
	Update(category *models.Category) error
	Delete(id int) error
}

type service struct {
	repo category.Repository
}

func NewService(repo category.Repository) Service {
	return &service{repo}
}

func (s *service) GetAll() ([]models.Category, error) {
	return s.repo.FindAll()
}

func (s *service) GetByID(id int) (*models.Category, error) {
	return s.repo.FindByID(id)
}

func (s *service) Create(category *models.Category) error {
	return s.repo.Create(category)
}

func (s *service) Update(category *models.Category) error {
	return s.repo.Update(category)
}

func (s *service) Delete(id int) error {
	return s.repo.Delete(id)
}
