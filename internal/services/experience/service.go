package experience

import (
	"encoding/json"
	"fmt"
	"time"
	"web-porto-backend/internal/domain/dto"
	"web-porto-backend/internal/domain/models"
	experienceRepo "web-porto-backend/internal/repositories/experience"
	tagService "web-porto-backend/internal/services/tag"
)

// Helper functions for parsing metadata
func getString(data map[string]interface{}, key string) string {
	if val, ok := data[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}

func getStringArray(data map[string]interface{}, key string) []string {
	if val, ok := data[key]; ok {
		if arr, ok := val.([]interface{}); ok {
			result := make([]string, 0, len(arr))
			for _, item := range arr {
				if str, ok := item.(string); ok {
					result = append(result, str)
				}
			}
			return result
		}
	}
	return []string{}
}

type Service struct {
	experienceRepo experienceRepo.Repository
	tagService     tagService.Service
}

func NewService(
	experienceRepo experienceRepo.Repository,
	tagService tagService.Service,
) *Service {
	return &Service{
		experienceRepo: experienceRepo,
		tagService:     tagService,
	}
}

// convertTechnologyNamesToIDs converts technology names to their corresponding IDs
func (s *Service) convertTechnologyNamesToIDs(technologyNames []string) ([]int, error) {
	var technologyIDs []int

	for _, name := range technologyNames {
		if name == "" {
			continue // Skip empty names
		}

		tag, err := s.tagService.GetByName(name)
		if err != nil {
			// If tag doesn't exist, create it
			createReq := &dto.CreateTagRequest{
				Name: name,
			}
			newTag, err := s.tagService.Create(createReq)
			if err != nil {
				return nil, fmt.Errorf("failed to create tag '%s': %v", name, err)
			}
			technologyIDs = append(technologyIDs, newTag.ID)
		} else {
			technologyIDs = append(technologyIDs, tag.ID)
		}
	}

	return technologyIDs, nil
}

// resolveTechnologies resolves technology IDs from either IDs or names
func (s *Service) resolveTechnologies(technologies []int, technologyNames []string) ([]int, error) {
	// If we have technology IDs, use them
	if len(technologies) > 0 {
		return technologies, nil
	}

	// If we have technology names, convert them to IDs
	if len(technologyNames) > 0 {
		return s.convertTechnologyNamesToIDs(technologyNames)
	}

	// No technologies provided
	return []int{}, nil
}

// CreateExperience creates a new work experience entry
func (s *Service) CreateExperience(req dto.CreateExperienceRequest) (*dto.ExperienceResponse, error) {
	// Parse dates
	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		return nil, err
	}

	var endDate *time.Time
	if req.EndDate != "" && !req.Current {
		parsedEndDate, err := time.Parse("2006-01-02", req.EndDate)
		if err != nil {
			return nil, err
		}
		endDate = &parsedEndDate
	}

	// Create metadata JSON
	metadata := map[string]interface{}{
		"originalId":  "",
		"lastUpdated": time.Now().Format(time.RFC3339),
		"version":     "1.0",
	}

	// Add any additional metadata from request
	if req.Metadata != nil {
		for k, v := range req.Metadata {
			metadata[k] = v
		}
	}

	metadataJSON, err := json.Marshal(metadata)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal metadata: %v", err)
	}

	// Create experience model
	experience := &models.Experience{
		Title:            req.Title,
		Company:          req.Company,
		Location:         req.Location,
		StartDate:        startDate,
		EndDate:          endDate,
		Current:          req.Current,
		Description:      req.Description,
		Responsibilities: models.StringArray(req.Responsibilities),
		CompanyURL:       req.CompanyURL,
		LogoURL:          req.LogoURL,
		Metadata:         string(metadataJSON),
	}

	// Create the experience
	if err := s.experienceRepo.Create(experience); err != nil {
		return nil, err
	}

	// Handle technologies if provided
	if len(req.TechnologyIDs) > 0 || len(req.TechnologyNames) > 0 {
		// Resolve technology IDs from either IDs or names
		technologyIDs, err := s.resolveTechnologies(req.TechnologyIDs, req.TechnologyNames)
		if err != nil {
			return nil, fmt.Errorf("failed to resolve technologies: %v", err)
		}

		// Update experience technologies
		if err := s.experienceRepo.UpdateExperienceTechnologies(experience.ID, technologyIDs); err != nil {
			return nil, fmt.Errorf("failed to update experience technologies: %v", err)
		}
	}

	// Return response
	return s.mapToResponse(experience), nil
}

// GetExperienceByID retrieves an experience entry by ID
func (s *Service) GetExperienceByID(id int) (*dto.ExperienceResponse, error) {
	experience, err := s.experienceRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	return s.mapToResponse(experience), nil
}

// ListExperiences retrieves a paginated list of experiences
func (s *Service) ListExperiences(page, size int) (*dto.PaginatedResponse, error) {
	offset := (page - 1) * size
	experiences, total, err := s.experienceRepo.GetAll(size, offset)
	if err != nil {
		return nil, err
	}

	// Convert experiences to response objects
	experienceList := make([]dto.ExperienceResponse, 0)
	for _, exp := range experiences {
		experienceList = append(experienceList, *s.mapToResponse(exp))
	}

	// Create pagination info
	pagination := dto.PaginationResponse{
		TotalCount:  total,
		CurrentPage: page,
		PageSize:    size,
		TotalPages:  int((total + int64(size) - 1) / int64(size)),
		HasNext:     int64(page*size) < total,
		HasPrevious: page > 1,
	}

	return &dto.PaginatedResponse{
		Data:       experienceList,
		Pagination: pagination,
	}, nil
}

// UpdateExperience updates an existing experience entry
func (s *Service) UpdateExperience(id int, req dto.UpdateExperienceRequest) (*dto.ExperienceResponse, error) {
	// Get existing experience
	experience, err := s.experienceRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Update fields if provided
	if req.Title != "" {
		experience.Title = req.Title
	}
	if req.Company != "" {
		experience.Company = req.Company
	}
	if req.Location != "" {
		experience.Location = req.Location
	}
	if req.StartDate != "" {
		startDate, err := time.Parse("2006-01-02", req.StartDate)
		if err != nil {
			return nil, err
		}
		experience.StartDate = startDate
	}

	// Handle end date and current status
	if req.Current != nil {
		experience.Current = *req.Current

		// If marked as current, set end_date to null
		if *req.Current {
			experience.EndDate = nil
		} else if req.EndDate != "" {
			// If not current and end_date provided, update it
			endDate, err := time.Parse("2006-01-02", req.EndDate)
			if err != nil {
				return nil, err
			}
			experience.EndDate = &endDate
		}
	} else if req.EndDate != "" {
		// If only end_date provided (current flag not changed)
		endDate, err := time.Parse("2006-01-02", req.EndDate)
		if err != nil {
			return nil, err
		}
		experience.EndDate = &endDate
		experience.Current = false
	}

	if req.Description != "" {
		experience.Description = req.Description
	}
	if len(req.Responsibilities) > 0 {
		experience.Responsibilities = models.StringArray(req.Responsibilities)
	}
	if req.CompanyURL != "" {
		experience.CompanyURL = req.CompanyURL
	}
	if req.LogoURL != "" {
		experience.LogoURL = req.LogoURL
	}

	// Handle technologies update - always update if provided
	if len(req.TechnologyIDs) > 0 || len(req.TechnologyNames) > 0 {
		// Resolve technology IDs from either IDs or names
		technologyIDs, err := s.resolveTechnologies(req.TechnologyIDs, req.TechnologyNames)
		if err != nil {
			return nil, fmt.Errorf("failed to resolve technologies: %v", err)
		}

		// Update experience technologies
		if err := s.experienceRepo.UpdateExperienceTechnologies(experience.ID, technologyIDs); err != nil {
			return nil, fmt.Errorf("failed to update experience technologies: %v", err)
		}
	}

	// Generate comprehensive metadata JSON
	if req.Metadata != nil {
		metadataBytes, err := json.Marshal(map[string]interface{}{
			"originalId":       experience.ID,
			"lastUpdated":      time.Now(),
			"version":          2,
			"title":            req.Title,
			"company":          req.Company,
			"location":         req.Location,
			"description":      req.Description,
			"companyUrl":       req.CompanyURL,
			"logoUrl":          req.LogoURL,
			"responsibilities": req.Responsibilities,
			"startDate":        req.StartDate,
			"endDate":          req.EndDate,
			"current":          req.Current,
			"userMetadata":     req.Metadata,
		})
		if err == nil {
			experience.Metadata = string(metadataBytes)
		}
	}

	// Update the experience
	if err := s.experienceRepo.Update(experience); err != nil {
		return nil, err
	}

	return s.mapToResponse(experience), nil
}

// DeleteExperience deletes an experience entry
func (s *Service) DeleteExperience(id int) error {
	return s.experienceRepo.Delete(id)
}

// GetCurrentExperiences retrieves currently active experiences
func (s *Service) GetCurrentExperiences() ([]dto.ExperienceResponse, error) {
	experiences, err := s.experienceRepo.GetCurrent()
	if err != nil {
		return nil, err
	}

	result := make([]dto.ExperienceResponse, len(experiences))
	for i, exp := range experiences {
		result[i] = *s.mapToResponse(exp)
	}

	return result, nil
}

// Helper function to map model to response
func (s *Service) mapToResponse(experience *models.Experience) *dto.ExperienceResponse {
	response := &dto.ExperienceResponse{
		ID:               experience.ID,
		Title:            experience.Title,
		Company:          experience.Company,
		Location:         experience.Location,
		StartDate:        experience.StartDate.Format("2006-01-02"),
		Current:          experience.Current,
		Description:      experience.Description,
		Responsibilities: []string(experience.Responsibilities),
		Technologies:     []dto.TagResponse{}, // Initialize empty array
		CompanyURL:       experience.CompanyURL,
		LogoURL:          experience.LogoURL,
		Metadata:         make(map[string]interface{}),
		CreatedAt:        experience.CreatedAt,
		UpdatedAt:        experience.UpdatedAt,
	}

	// Add technologies from relational tags (primary source)
	for _, tag := range experience.Technologies {
		response.Technologies = append(response.Technologies, dto.TagResponse{
			ID:   tag.ID,
			Name: tag.Name,
			Slug: tag.Slug,
		})
	}

	// Add end date if available
	if experience.EndDate != nil {
		endDate := experience.EndDate.Format("2006-01-02")
		response.EndDate = &endDate
	}

	// Parse metadata and restore rich data
	if experience.Metadata != "" {
		var metadata map[string]interface{}
		if err := json.Unmarshal([]byte(experience.Metadata), &metadata); err == nil {
			response.Metadata = metadata

			// Restore additional data from metadata if needed
			// Keep metadata technologies as additional info only
			if responsibilities := getStringArray(metadata, "responsibilities"); len(responsibilities) > 0 {
				response.Responsibilities = responsibilities
			}
			if companyURL := getString(metadata, "companyUrl"); companyURL != "" {
				response.CompanyURL = companyURL
			}
			if logoURL := getString(metadata, "logoUrl"); logoURL != "" {
				response.LogoURL = logoURL
			}
		}
	}

	return response
}
