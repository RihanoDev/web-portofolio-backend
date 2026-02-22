package experience

import (
	"encoding/json"
	"fmt"
	"time"
	"web-porto-backend/common/utils"
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

func getInt(data map[string]interface{}, key string) int {
	if val, ok := data[key]; ok {
		if num, ok := val.(float64); ok {
			return int(num)
		}
		if num, ok := val.(int); ok {
			return num
		}
	}
	return 0
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

// parseFlexibleDate tries multiple date formats
func parseFlexibleDate(dateStr string) (time.Time, error) {
	formats := []string{
		"2006-01-02",
		"2006-01",
		"2006/01/02",
		"2006/01",
	}

	for _, format := range formats {
		if t, err := time.Parse(format, dateStr); err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("invalid date format: %s", dateStr)
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

		// Try to find by slug first as it's the unique identifier most likely to conflict
		slug := utils.StringToSlug(name)
		var existingTag *dto.TagResponse

		// Attempt to find by name exactly first
		tag, err := s.tagService.GetBySlug(slug)
		if err == nil {
			existingTag = tag
		} else {
			// If not found by slug, it's safe to try creating it
			fmt.Printf("[ExperienceService] Tag slug '%s' not found, will attempt to create\n", slug)
		}

		if existingTag != nil {
			technologyIDs = append(technologyIDs, existingTag.ID)
		} else {
			// If tag doesn't exist, create it
			createReq := &dto.CreateTagRequest{
				Name: name,
				Slug: slug,
			}
			newTag, err := s.tagService.Create(createReq)
			if err != nil {
				// Final attempt: maybe it was created by another process in the meantime?
				// Or check if it exists by name just in case slugify logic differs
				tag, errRetry := s.tagService.GetByName(name)
				if errRetry == nil {
					technologyIDs = append(technologyIDs, tag.ID)
					continue
				}

				fmt.Printf("[ExperienceService] Failed to create tag '%s' (slug: %s): %v\n", name, slug, err)
				return nil, fmt.Errorf("failed to handle technology '%s': %v", name, err)
			}
			technologyIDs = append(technologyIDs, newTag.ID)
		}
	}

	return technologyIDs, nil
}

// resolveTechnologies resolves technology IDs from all sources
func (s *Service) resolveTechnologies(technologies []int, technologyNames []string) ([]int, error) {
	var allIDs []int

	// Add int IDs
	if len(technologies) > 0 {
		allIDs = append(allIDs, technologies...)
	}

	// Convert and add technology names
	if len(technologyNames) > 0 {
		ids, err := s.convertTechnologyNamesToIDs(technologyNames)
		if err != nil {
			return nil, err
		}
		allIDs = append(allIDs, ids...)
	}

	// Deduplicate
	return s.deduplicateIDs(allIDs), nil
}

// deduplicateIDs removes duplicate integer IDs
func (s *Service) deduplicateIDs(ids []int) []int {
	uniqueMap := make(map[int]bool)
	var result []int
	for _, id := range ids {
		if id > 0 && !uniqueMap[id] {
			uniqueMap[id] = true
			result = append(result, id)
		}
	}
	return result
}

// CreateExperience creates a new work experience entry
func (s *Service) CreateExperience(req dto.CreateExperienceRequest) (*dto.ExperienceResponse, error) {
	// Validate required fields
	if req.Title == "" {
		return nil, fmt.Errorf("title is required")
	}
	if req.Company == "" {
		return nil, fmt.Errorf("company is required")
	}
	if req.StartDate == "" {
		return nil, fmt.Errorf("startDate is required")
	}

	// Parse dates
	startDate, err := parseFlexibleDate(req.StartDate)
	if err != nil {
		return nil, fmt.Errorf("invalid startDate format: %v", err)
	}

	var endDate *time.Time
	if req.EndDate != "" && !req.Current {
		parsedEndDate, err := parseFlexibleDate(req.EndDate)
		if err != nil {
			return nil, fmt.Errorf("invalid endDate format: %v", err)
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

	// Force responsibilities into metadata for rich support
	metadata["responsibilities"] = req.Responsibilities

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
		fmt.Printf("[ExperienceService.Create] Repository error: %v\n", err)
		return nil, fmt.Errorf("database error: %v", err)
	}

	// Handle technologies if provided
	if len(req.TechnologyIDs) > 0 || len(req.TechnologyNames) > 0 {
		fmt.Printf("[ExperienceService.Create] Resolving technologies: IDs=%v, Names=%v\n", req.TechnologyIDs, req.TechnologyNames)
		// Resolve technology IDs from either IDs or names
		technologyIDs, err := s.resolveTechnologies(req.TechnologyIDs, req.TechnologyNames)
		if err != nil {
			fmt.Printf("[ExperienceService.Create] Resolution error: %v\n", err)
			return nil, fmt.Errorf("failed to resolve technologies: %v", err)
		}

		// Update experience technologies
		if err := s.experienceRepo.UpdateExperienceTechnologies(experience.ID, technologyIDs); err != nil {
			fmt.Printf("[ExperienceService.Create] Association error: %v\n", err)
			return nil, fmt.Errorf("failed to update experience technologies: %v", err)
		}
	}

	// Handle images if provided
	if len(req.Images) > 0 {
		var expImages []models.ExperienceImage
		for _, img := range req.Images {
			expImages = append(expImages, models.ExperienceImage{
				ExperienceID: experience.ID,
				URL:          img.URL,
				Caption:      img.Caption,
				SortOrder:    img.SortOrder,
			})
		}
		if err := s.experienceRepo.UpdateExperienceImages(experience.ID, expImages); err != nil {
			return nil, fmt.Errorf("failed to add experience images: %v", err)
		}
	}

	// Return response with fresh data
	return s.GetExperienceByID(experience.ID)
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
	fmt.Printf("[UpdateExperience] Updating id=%d title=%v startDate=%v endDate=%v current=%v techs(ids)=%v techs(names)=%v\n",
		id, req.Title, req.StartDate, req.EndDate, req.Current, req.TechnologyIDs, req.TechnologyNames)

	// Get existing experience
	experience, err := s.experienceRepo.GetByID(id)
	if err != nil {
		fmt.Printf("[UpdateExperience] GetByID error: %v\n", err)
		return nil, err
	}

	// Update fields if provided
	if req.Title != nil && *req.Title != "" {
		experience.Title = *req.Title
	}
	if req.Company != nil && *req.Company != "" {
		experience.Company = *req.Company
	}
	if req.Location != nil && *req.Location != "" {
		experience.Location = *req.Location
	}
	if req.StartDate != nil && *req.StartDate != "" {
		startDate, err := parseFlexibleDate(*req.StartDate)
		if err != nil {
			fmt.Printf("[UpdateExperience] startDate parse error: %v (input: %q)\n", err, *req.StartDate)
			return nil, fmt.Errorf("invalid startDate format: %v", err)
		}
		experience.StartDate = startDate
	}

	// Handle end date and current status
	if req.Current != nil {
		experience.Current = *req.Current

		// If marked as current, set end_date to null
		if *req.Current {
			experience.EndDate = nil
		} else if req.EndDate != nil && *req.EndDate != "" {
			// If not current and end_date provided, update it
			endDate, err := parseFlexibleDate(*req.EndDate)
			if err != nil {
				return nil, fmt.Errorf("invalid endDate format: %v", err)
			}
			experience.EndDate = &endDate
		}
	} else if req.EndDate != nil && *req.EndDate != "" {
		// If only end_date provided (current flag not changed)
		endDate, err := parseFlexibleDate(*req.EndDate)
		if err != nil {
			return nil, fmt.Errorf("invalid endDate format: %v", err)
		}
		experience.EndDate = &endDate
		experience.Current = false
	}

	if req.Description != nil && *req.Description != "" {
		experience.Description = *req.Description
	}
	if req.Responsibilities != nil {
		experience.Responsibilities = models.StringArray(*req.Responsibilities)
	}

	if req.CompanyURL != nil && *req.CompanyURL != "" {
		experience.CompanyURL = *req.CompanyURL
	}
	if req.LogoURL != nil && *req.LogoURL != "" {
		experience.LogoURL = *req.LogoURL
	}

	// Generate comprehensive metadata JSON - always update
	metadataContent := make(map[string]interface{})
	if experience.Metadata != "" && experience.Metadata != "{}" {
		_ = json.Unmarshal([]byte(experience.Metadata), &metadataContent)
	}
	metadataContent["originalId"] = experience.ID
	metadataContent["lastUpdated"] = time.Now().Format(time.RFC3339)
	metadataContent["version"] = 2
	metadataContent["companyUrl"] = experience.CompanyURL
	metadataContent["logoUrl"] = experience.LogoURL

	// Ensure responsibilities are always in sync in metadata
	resps := []string(experience.Responsibilities)
	if resps == nil {
		resps = []string{}
	}
	metadataContent["responsibilities"] = resps

	// Merge user metadata jika ada
	if req.Metadata != nil {
		reservedKeys := map[string]bool{
			"responsibilities": true,
			"originalId":       true,
			"lastUpdated":      true,
			"version":          true,
			"companyUrl":       true,
			"logoUrl":          true,
		}
		for k, v := range req.Metadata {
			if !reservedKeys[k] {
				metadataContent[k] = v
			}
		}
	}
	metadataBytes, err := json.Marshal(metadataContent)
	if err == nil {
		experience.Metadata = string(metadataBytes)
	}

	// Update the experience in DB first
	if err := s.experienceRepo.Update(experience); err != nil {
		return nil, fmt.Errorf("failed to update experience: %v", err)
	}

	// Handle technologies update - after main update
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

	// Update images if provided
	if req.Images != nil {
		var expImages []models.ExperienceImage
		for _, img := range req.Images {
			expImages = append(expImages, models.ExperienceImage{
				ExperienceID: experience.ID,
				URL:          img.URL,
				Caption:      img.Caption,
				SortOrder:    img.SortOrder,
			})
		}
		if err := s.experienceRepo.UpdateExperienceImages(experience.ID, expImages); err != nil {
			return nil, fmt.Errorf("failed to update experience images: %v", err)
		}
	}

	// Reload dari DB agar mendapatkan Technologies yang terbaru
	updatedExp, err := s.experienceRepo.GetByID(id)
	if err != nil {
		// Jika gagal reload, kembalikan hasil mapToResponse dengan data yang ada
		return s.mapToResponse(experience), nil
	}

	return s.mapToResponse(updatedExp), nil
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
		Images:           []dto.ExperienceImageResponse{},
		Metadata:         make(map[string]interface{}),
		CreatedAt:        experience.CreatedAt,
		UpdatedAt:        experience.UpdatedAt,
	}

	// Add images from database (primary source)
	for _, img := range experience.Images {
		response.Images = append(response.Images, dto.ExperienceImageResponse{
			ID:        img.ID,
			URL:       img.URL,
			Caption:   img.Caption,
			SortOrder: img.SortOrder,
		})
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
	if experience.Metadata != "" && experience.Metadata != "{}" {
		var metadata map[string]interface{}
		if err := json.Unmarshal([]byte(experience.Metadata), &metadata); err == nil {
			response.Metadata = metadata

			// Restore additional data from metadata if needed
			// IMPORTANT: If responsibilities exist in metadata, they are the source of truth for rich data
			if val, ok := metadata["responsibilities"]; ok {
				if arr, ok := val.([]interface{}); ok {
					resps := make([]string, 0, len(arr))
					for _, item := range arr {
						if str, ok := item.(string); ok {
							resps = append(resps, str)
						}
					}
					response.Responsibilities = resps
				} else if arr, ok := val.([]string); ok {
					response.Responsibilities = arr
				}
			}

			// Restore images from metadata if DB images are empty (backward compatibility)
			if len(response.Images) == 0 {
				if imagesData, ok := metadata["images"].([]interface{}); ok {
					for _, imgData := range imagesData {
						if imgMap, ok := imgData.(map[string]interface{}); ok {
							image := dto.ExperienceImageResponse{
								ID:        getString(imgMap, "id"),
								URL:       getString(imgMap, "url"),
								Caption:   getString(imgMap, "caption"),
								SortOrder: getInt(imgMap, "sortOrder"),
							}
							response.Images = append(response.Images, image)
						}
					}
				}
			}

			if companyURL := getString(metadata, "companyUrl"); companyURL != "" {
				response.CompanyURL = companyURL
			}
			if logoURL := getString(metadata, "logoUrl"); logoURL != "" {
				response.LogoURL = logoURL
			}

			// Restore images from metadata if DB images are empty (backward compatibility)
			if len(response.Images) == 0 {
				if imagesData, ok := metadata["images"].([]interface{}); ok {
					for _, imgData := range imagesData {
						if imgMap, ok := imgData.(map[string]interface{}); ok {
							image := dto.ExperienceImageResponse{
								ID:        getString(imgMap, "id"),
								URL:       getString(imgMap, "url"),
								Caption:   getString(imgMap, "caption"),
								SortOrder: getInt(imgMap, "sortOrder"),
							}
							response.Images = append(response.Images, image)
						}
					}
				}
			}
		}
	}

	return response
}
