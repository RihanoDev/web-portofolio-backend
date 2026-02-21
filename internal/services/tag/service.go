package tag

import (
	"regexp"
	"strings"
	"web-porto-backend/internal/domain/dto"
	"web-porto-backend/internal/domain/models"
	"web-porto-backend/internal/repositories/tag"
)

// slugify creates a URL-friendly version of the input string
func slugify(text string) string {
	// Convert to lowercase
	text = strings.ToLower(text)

	// Replace spaces with hyphens
	text = strings.ReplaceAll(text, " ", "-")

	// Remove any character that's not alphanumeric or hyphen
	reg, _ := regexp.Compile("[^a-z0-9-]+")
	text = reg.ReplaceAllString(text, "")

	// Remove multiple consecutive hyphens
	for strings.Contains(text, "--") {
		text = strings.ReplaceAll(text, "--", "-")
	}

	// Trim hyphens from beginning and end
	text = strings.Trim(text, "-")

	return text
}

// Service defines tag business logic operations
type Service interface {
	GetAll() ([]dto.TagResponse, error)
	GetByID(id int) (*dto.TagResponse, error)
	GetByName(name string) (*dto.TagResponse, error)
	GetBySlug(slug string) (*dto.TagResponse, error)
	Create(tag *dto.CreateTagRequest) (*dto.TagResponse, error)
	Update(id int, tag *dto.UpdateTagRequest) (*dto.TagResponse, error)
	Delete(id int) error
}

type service struct {
	repo tag.Repository
}

// NewService creates a new tag service
func NewService(repo tag.Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetAll() ([]dto.TagResponse, error) {
	tags, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}

	var responses []dto.TagResponse
	for _, tag := range tags {
		responses = append(responses, dto.TagResponse{
			ID:   tag.ID,
			Name: tag.Name,
			Slug: tag.Slug,
		})
	}

	return responses, nil
}

func (s *service) GetByID(id int) (*dto.TagResponse, error) {
	tag, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	return &dto.TagResponse{
		ID:   tag.ID,
		Name: tag.Name,
		Slug: tag.Slug,
	}, nil
}

func (s *service) GetByName(name string) (*dto.TagResponse, error) {
	tag, err := s.repo.GetByName(name)
	if err != nil {
		return nil, err
	}

	return &dto.TagResponse{
		ID:   tag.ID,
		Name: tag.Name,
		Slug: tag.Slug,
	}, nil
}

func (s *service) GetBySlug(slug string) (*dto.TagResponse, error) {
	tag, err := s.repo.GetBySlug(slug)
	if err != nil {
		return nil, err
	}

	return &dto.TagResponse{
		ID:   tag.ID,
		Name: tag.Name,
		Slug: tag.Slug,
	}, nil
}

func (s *service) Create(request *dto.CreateTagRequest) (*dto.TagResponse, error) {
	// Generate slug if not provided
	slug := request.Slug
	if slug == "" {
		// In real implementation, we'd use a proper slug generator function
		slug = slugify(request.Name)
	}

	tag := &models.Tag{
		Name: request.Name,
		Slug: slug,
	}

	createdTag, err := s.repo.Create(tag)
	if err != nil {
		return nil, err
	}

	return &dto.TagResponse{
		ID:   createdTag.ID,
		Name: createdTag.Name,
		Slug: createdTag.Slug,
	}, nil
}

func (s *service) Update(id int, request *dto.UpdateTagRequest) (*dto.TagResponse, error) {
	existingTag, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Update name
	existingTag.Name = request.Name

	// Update slug if provided
	if request.Slug != "" {
		existingTag.Slug = request.Slug
	} else if existingTag.Name != request.Name {
		// Generate new slug if name changed and no slug provided
		existingTag.Slug = slugify(request.Name)
	}

	updatedTag, err := s.repo.Update(existingTag)
	if err != nil {
		return nil, err
	}

	return &dto.TagResponse{
		ID:   updatedTag.ID,
		Name: updatedTag.Name,
		Slug: updatedTag.Slug,
	}, nil
}

func (s *service) Delete(id int) error {
	return s.repo.Delete(id)
}
