package page

import (
	"web-porto-backend/common/utils"
	"web-porto-backend/internal/domain/models"
	"web-porto-backend/internal/repositories/page"
)

// PaginationInfo represents pagination metadata
type PaginationInfo struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

type Service interface {
	Create(pageData *models.Page) error
	GetByID(id uint) (*models.Page, error)
	GetAll(page, limit int) ([]*models.Page, *PaginationInfo, error)
	Update(id uint, pageData *models.Page) error
	Delete(id uint) error
	GetBySlug(slug string) (*models.Page, error)
	GetPublished(page, limit int) ([]*models.Page, *PaginationInfo, error)
}

type service struct {
	repo page.Repository
}

func NewService(repo page.Repository) Service {
	return &service{repo: repo}
}

func (s *service) Create(pageData *models.Page) error {
	// Generate slug if not provided
	if pageData.Slug == "" {
		pageData.Slug = utils.StringToSlug(pageData.Title)
	}
	return s.repo.Create(pageData)
}

func (s *service) GetByID(id uint) (*models.Page, error) {
	return s.repo.GetByID(id)
}

func (s *service) GetAll(page, limit int) ([]*models.Page, *PaginationInfo, error) {
	// Validate pagination parameters
	page, limit = utils.ValidatePageAndLimit(page, limit)
	offset := (page - 1) * limit

	pages, total, err := s.repo.GetAll(limit, offset)
	if err != nil {
		return nil, nil, err
	}

	pagination := s.calculatePagination(int(total), page, limit)
	return pages, pagination, nil
}

func (s *service) Update(id uint, pageData *models.Page) error {
	existingPage, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	// Update fields according to actual model
	existingPage.Title = pageData.Title
	existingPage.Content = pageData.Content
	existingPage.Status = pageData.Status // Update slug if title changed
	if existingPage.Title != pageData.Title {
		existingPage.Slug = utils.StringToSlug(pageData.Title)
	}

	return s.repo.Update(existingPage)
}

func (s *service) Delete(id uint) error {
	return s.repo.Delete(id)
}

func (s *service) GetBySlug(slug string) (*models.Page, error) {
	return s.repo.GetBySlug(slug)
}

func (s *service) GetPublished(page, limit int) ([]*models.Page, *PaginationInfo, error) {
	// Validate pagination parameters
	page, limit = utils.ValidatePageAndLimit(page, limit)
	offset := (page - 1) * limit

	pages, total, err := s.repo.GetPublished(limit, offset)
	if err != nil {
		return nil, nil, err
	}

	pagination := s.calculatePagination(int(total), page, limit)
	return pages, pagination, nil
}

// calculatePagination calculates pagination metadata
func (s *service) calculatePagination(total, page, limit int) *PaginationInfo {
	totalPages := total / limit
	if total%limit != 0 {
		totalPages++
	}

	return &PaginationInfo{
		Page:       page,
		Limit:      limit,
		Total:      int64(total),
		TotalPages: totalPages,
	}
}
