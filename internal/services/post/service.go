package post

import (
"web-porto-backend/common/utils"
"web-porto-backend/internal/domain/models"
"web-porto-backend/internal/repositories/post"
)

// PaginationInfo represents pagination metadata
type PaginationInfo struct {
Page       int   `json:"page"`
Limit      int   `json:"limit"`
Total      int64 `json:"total"`
TotalPages int   `json:"total_pages"`
}

type Service interface {
Create(postData *models.Post) error
GetByID(id uint) (*models.Post, error)
GetAll(page, limit int) ([]*models.Post, *PaginationInfo, error)
Update(id uint, postData *models.Post) error
Delete(id uint) error
GetBySlug(slug string) (*models.Post, error)
GetByAuthorID(authorID int, page, limit int) ([]*models.Post, *PaginationInfo, error)
GetPublished(page, limit int) ([]*models.Post, *PaginationInfo, error)
}

type service struct {
repo post.Repository
}

func NewService(repo post.Repository) Service {
return &service{repo: repo}
}

func (s *service) Create(postData *models.Post) error {
// Generate slug if not provided
if postData.Slug == "" {
postData.Slug = utils.StringToSlug(postData.Title)
}
return s.repo.Create(postData)
}

func (s *service) GetByID(id uint) (*models.Post, error) {
return s.repo.GetByID(id)
}

func (s *service) GetAll(page, limit int) ([]*models.Post, *PaginationInfo, error) {
// Validate pagination parameters
page, limit = utils.ValidatePageAndLimit(page, limit)
offset := (page - 1) * limit

posts, total, err := s.repo.GetAll(limit, offset)
if err != nil {
return nil, nil, err
}

pagination := s.calculatePagination(int(total), page, limit)
return posts, pagination, nil
}

func (s *service) Update(id uint, postData *models.Post) error {
existingPost, err := s.repo.GetByID(id)
if err != nil {
return err
}

// Update fields according to actual model
existingPost.Title = postData.Title
existingPost.Content = postData.Content
existingPost.Status = postData.Status
existingPost.AuthorID = postData.AuthorID
existingPost.PublishedAt = postData.PublishedAt

// Update slug if title changed
if existingPost.Title != postData.Title {
existingPost.Slug = utils.StringToSlug(postData.Title)
}

return s.repo.Update(existingPost)
}

func (s *service) Delete(id uint) error {
return s.repo.Delete(id)
}

func (s *service) GetBySlug(slug string) (*models.Post, error) {
return s.repo.GetBySlug(slug)
}

func (s *service) GetByAuthorID(authorID int, page, limit int) ([]*models.Post, *PaginationInfo, error) {
// Validate pagination parameters
page, limit = utils.ValidatePageAndLimit(page, limit)
offset := (page - 1) * limit

posts, total, err := s.repo.GetByAuthorID(authorID, limit, offset)
if err != nil {
return nil, nil, err
}

pagination := s.calculatePagination(int(total), page, limit)
return posts, pagination, nil
}

func (s *service) GetPublished(page, limit int) ([]*models.Post, *PaginationInfo, error) {
// Validate pagination parameters
page, limit = utils.ValidatePageAndLimit(page, limit)
offset := (page - 1) * limit

posts, total, err := s.repo.GetPublished(limit, offset)
if err != nil {
return nil, nil, err
}

pagination := s.calculatePagination(int(total), page, limit)
return posts, pagination, nil
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
