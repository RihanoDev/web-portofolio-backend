package project

import (
	"encoding/json"
	"fmt"
	"web-porto-backend/internal/domain/dto"
	"web-porto-backend/internal/domain/models"
	categoryRepo "web-porto-backend/internal/repositories/category"
	projectRepo "web-porto-backend/internal/repositories/project"
	tagService "web-porto-backend/internal/services/tag"
	userService "web-porto-backend/internal/services/user"

	"github.com/google/uuid"
	"github.com/gosimple/slug"
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

// Service handles business logic for projects
type Service struct {
	projectRepo  projectRepo.Repository
	categoryRepo categoryRepo.Repository
	userService  userService.Service
	tagService   tagService.Service
}

// NewService creates a new project service
func NewService(
	projectRepo projectRepo.Repository,
	categoryRepo categoryRepo.Repository,
	userService userService.Service,
	tagService tagService.Service,
) *Service {
	return &Service{
		projectRepo:  projectRepo,
		categoryRepo: categoryRepo,
		userService:  userService,
		tagService:   tagService,
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

// CreateProject creates a new project
func (s *Service) CreateProject(req dto.CreateProjectRequest) (*dto.ProjectResponse, error) {
	// Generate slug if not provided
	if req.Slug == "" {
		req.Slug = slug.Make(req.Title)
	}

	// Use default authorID if not provided - get from database
	authorID := req.AuthorID
	if authorID == 0 {
		// Get default admin user from database
		defaultUser, err := s.userService.GetDefaultAdmin()
		if err != nil {
			return nil, fmt.Errorf("failed to get default admin user: %v", err)
		}
		authorID = int(defaultUser.ID)
	}

	// Create metadata JSON
	metadata := map[string]interface{}{
		"githubUrl":   req.GitHubURL,
		"liveDemoUrl": req.LiveDemoURL,
	}

	// Add any additional metadata
	for k, v := range req.Metadata {
		metadata[k] = v
	}

	metadataJSON, err := json.Marshal(metadata)
	if err != nil {
		return nil, err
	}

	// Create project model
	project := &models.Project{
		ID:           uuid.New().String(), // Generate a proper UUID here
		Title:        req.Title,
		Slug:         req.Slug,
		Description:  req.Description,
		Content:      req.Content,
		ThumbnailURL: req.ThumbnailURL,
		Status:       req.Status,
		CategoryID:   req.CategoryID,
		AuthorID:     authorID,
		GitHubURL:    req.GitHubURL,
		LiveDemoURL:  req.LiveDemoURL,
		Metadata:     string(metadataJSON),
	}

	// Create the project
	if err := s.projectRepo.Create(project); err != nil {
		return nil, err
	}

	// Handle technologies if provided
	if len(req.Technologies) > 0 || len(req.TechnologyNames) > 0 {
		// Resolve technology IDs from either IDs or names
		technologyIDs, err := s.resolveTechnologies(req.Technologies, req.TechnologyNames)
		if err != nil {
			return nil, fmt.Errorf("failed to resolve technologies: %v", err)
		}

		// Update project technologies
		if err := s.projectRepo.UpdateProjectTechnologies(project.ID, technologyIDs); err != nil {
			return nil, fmt.Errorf("failed to update project technologies: %v", err)
		}
	}

	// Return project response
	return &dto.ProjectResponse{
		ID:           project.ID,
		Title:        project.Title,
		Slug:         project.Slug,
		Description:  project.Description,
		Content:      project.Content,
		ThumbnailURL: project.ThumbnailURL,
		Status:       project.Status,
		GitHubURL:    project.GitHubURL,
		LiveDemoURL:  project.LiveDemoURL,
		CreatedAt:    project.CreatedAt,
		UpdatedAt:    project.UpdatedAt,
		Author: dto.AuthorResponse{
			ID:       project.AuthorID,
			Username: "user", // Default username
		},
	}, nil
}

// GetProjectByID retrieves a project by ID
func (s *Service) GetProjectByID(id string) (*dto.ProjectResponse, error) {
	project, err := s.projectRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Create response
	response := &dto.ProjectResponse{
		ID:           project.ID,
		Title:        project.Title,
		Slug:         project.Slug,
		Description:  project.Description,
		Content:      project.Content,
		ThumbnailURL: project.ThumbnailURL,
		Status:       project.Status,
		GitHubURL:    project.GitHubURL,
		LiveDemoURL:  project.LiveDemoURL,
		CreatedAt:    project.CreatedAt,
		UpdatedAt:    project.UpdatedAt,
		Author: dto.AuthorResponse{
			ID:       project.AuthorID,
			Username: "user", // Default username
		},
		Images:       []dto.ProjectImageResponse{},
		Videos:       []dto.ProjectVideoResponse{},
		Technologies: []dto.TagResponse{},
		Metadata:     make(map[string]interface{}),
	}

	// Parse metadata and restore rich data
	if project.Metadata != "" {
		var metadata map[string]interface{}
		if err := json.Unmarshal([]byte(project.Metadata), &metadata); err == nil {
			response.Metadata = metadata

			// Restore images from metadata
			if imagesData, ok := metadata["images"].([]interface{}); ok {
				for _, imgData := range imagesData {
					if imgMap, ok := imgData.(map[string]interface{}); ok {
						image := dto.ProjectImageResponse{
							ID:        getString(imgMap, "id"),
							URL:       getString(imgMap, "url"),
							Caption:   getString(imgMap, "caption"),
							SortOrder: getInt(imgMap, "sortOrder"),
						}
						response.Images = append(response.Images, image)
					}
				}
			}

			// Restore videos from metadata
			if videosData, ok := metadata["videos"].([]interface{}); ok {
				for _, vidData := range videosData {
					if vidMap, ok := vidData.(map[string]interface{}); ok {
						video := dto.ProjectVideoResponse{
							ID:        getString(vidMap, "id"),
							URL:       getString(vidMap, "url"),
							Caption:   getString(vidMap, "caption"),
							SortOrder: getInt(vidMap, "sortOrder"),
						}
						response.Videos = append(response.Videos, video)
					}
				}
			}

			// Restore additional URLs from metadata
			if githubURL := getString(metadata, "githubUrl"); githubURL != "" {
				response.GitHubURL = githubURL
			}
			if liveDemoURL := getString(metadata, "liveDemoUrl"); liveDemoURL != "" {
				response.LiveDemoURL = liveDemoURL
			}
			if demoURL := getString(metadata, "demoUrl"); demoURL != "" && response.LiveDemoURL == "" {
				response.LiveDemoURL = demoURL
			}
			if featuredImageURL := getString(metadata, "featuredImageUrl"); featuredImageURL != "" {
				response.ThumbnailURL = featuredImageURL
			}
			if thumbnailURL := getString(metadata, "thumbnailUrl"); thumbnailURL != "" && response.ThumbnailURL == "" {
				response.ThumbnailURL = thumbnailURL
			}

			// Restore technologies from metadata
			if technologies := getStringArray(metadata, "technologies"); len(technologies) > 0 {
				for _, tech := range technologies {
					response.Technologies = append(response.Technologies, dto.TagResponse{
						ID:   0, // Technology tags might not have IDs in metadata
						Name: tech,
					})
				}
			}
		}
	}

	// Add technologies from relational tags (primary source)
	for _, tag := range project.Tags {
		response.Technologies = append(response.Technologies, dto.TagResponse{
			ID:   tag.ID,
			Name: tag.Name,
			Slug: tag.Slug,
		})
	}

	// Add category if available
	if project.Category != nil {
		response.Category = &dto.CategoryResponse{
			ID:   *project.CategoryID,
			Name: project.Category.Name,
			Slug: project.Category.Slug,
		}
	}

	return response, nil
}

// GetProjectBySlug retrieves a project by slug
func (s *Service) GetProjectBySlug(slug string) (*dto.ProjectResponse, error) {
	project, err := s.projectRepo.GetBySlug(slug)
	if err != nil {
		return nil, err
	}

	// Create response (similar to GetProjectByID)
	response := &dto.ProjectResponse{
		ID:           project.ID,
		Title:        project.Title,
		Slug:         project.Slug,
		Description:  project.Description,
		Content:      project.Content,
		ThumbnailURL: project.ThumbnailURL,
		Status:       project.Status,
		GitHubURL:    project.GitHubURL,
		LiveDemoURL:  project.LiveDemoURL,
		CreatedAt:    project.CreatedAt,
		UpdatedAt:    project.UpdatedAt,
		Author: dto.AuthorResponse{
			ID:       project.AuthorID,
			Username: "user", // Default username
		},
		Images:       []dto.ProjectImageResponse{},
		Videos:       []dto.ProjectVideoResponse{},
		Technologies: []dto.TagResponse{},
		Metadata:     make(map[string]interface{}),
	}

	// Parse metadata and restore rich data
	if project.Metadata != "" {
		var metadata map[string]interface{}
		if err := json.Unmarshal([]byte(project.Metadata), &metadata); err == nil {
			response.Metadata = metadata

			// Restore images from metadata
			if imagesData, ok := metadata["images"].([]interface{}); ok {
				for _, imgData := range imagesData {
					if imgMap, ok := imgData.(map[string]interface{}); ok {
						image := dto.ProjectImageResponse{
							ID:        getString(imgMap, "id"),
							URL:       getString(imgMap, "url"),
							Caption:   getString(imgMap, "caption"),
							SortOrder: getInt(imgMap, "sortOrder"),
						}
						response.Images = append(response.Images, image)
					}
				}
			}

			// Restore videos from metadata
			if videosData, ok := metadata["videos"].([]interface{}); ok {
				for _, vidData := range videosData {
					if vidMap, ok := vidData.(map[string]interface{}); ok {
						video := dto.ProjectVideoResponse{
							ID:        getString(vidMap, "id"),
							URL:       getString(vidMap, "url"),
							Caption:   getString(vidMap, "caption"),
							SortOrder: getInt(vidMap, "sortOrder"),
						}
						response.Videos = append(response.Videos, video)
					}
				}
			}

			// Restore additional URLs from metadata
			if githubURL := getString(metadata, "githubUrl"); githubURL != "" {
				response.GitHubURL = githubURL
			}
			if liveDemoURL := getString(metadata, "liveDemoUrl"); liveDemoURL != "" {
				response.LiveDemoURL = liveDemoURL
			}
			if demoURL := getString(metadata, "demoUrl"); demoURL != "" && response.LiveDemoURL == "" {
				response.LiveDemoURL = demoURL
			}
			if featuredImageURL := getString(metadata, "featuredImageUrl"); featuredImageURL != "" {
				response.ThumbnailURL = featuredImageURL
			}
			if thumbnailURL := getString(metadata, "thumbnailUrl"); thumbnailURL != "" && response.ThumbnailURL == "" {
				response.ThumbnailURL = thumbnailURL
			}

			// Restore technologies from metadata
			if technologies := getStringArray(metadata, "technologies"); len(technologies) > 0 {
				for _, tech := range technologies {
					response.Technologies = append(response.Technologies, dto.TagResponse{
						ID:   0, // Technology tags might not have IDs in metadata
						Name: tech,
					})
				}
			}
		}
	}

	// Add technologies from relational tags (primary source)
	for _, tag := range project.Tags {
		response.Technologies = append(response.Technologies, dto.TagResponse{
			ID:   tag.ID,
			Name: tag.Name,
			Slug: tag.Slug,
		})
	}

	// Add category if available
	if project.Category != nil {
		response.Category = &dto.CategoryResponse{
			ID:   *project.CategoryID,
			Name: project.Category.Name,
			Slug: project.Category.Slug,
		}
	}

	return response, nil
}

// GetProjectsByCategorySlug retrieves projects by category slug (simplified)
func (s *Service) GetProjectsByCategorySlug(slug string, page, size int) (*dto.PaginatedResponse, error) {
	// This is a simplified implementation
	return &dto.PaginatedResponse{
		Data:       []interface{}{},
		Pagination: dto.PaginationResponse{},
	}, nil
}

// ListProjects retrieves a paginated list of projects
func (s *Service) ListProjects(page, size int) (*dto.PaginatedResponse, error) {
	offset := (page - 1) * size
	projects, total, err := s.projectRepo.GetAll(size, offset)
	if err != nil {
		return nil, err
	}

	// Convert projects to project list responses
	projectList := make([]dto.ProjectListResponse, 0)
	for _, project := range projects {
		listResponse := dto.ProjectListResponse{
			ID:           project.ID,
			Title:        project.Title,
			Slug:         project.Slug,
			Description:  project.Description,
			ThumbnailURL: project.ThumbnailURL,
			Status:       project.Status,
			AuthorName:   "user", // Default username
			GitHubURL:    project.GitHubURL,
			LiveDemoURL:  project.LiveDemoURL,
			Technologies: []string{}, // Initialize empty array
			CreatedAt:    project.CreatedAt,
		}

		// Add category if available
		if project.Category != nil {
			listResponse.Category = project.Category.Name
		}

		// Parse metadata for list view
		if project.Metadata != "" {
			var metadata map[string]interface{}
			if err := json.Unmarshal([]byte(project.Metadata), &metadata); err == nil {
				// Restore URLs from metadata if not set
				if listResponse.GitHubURL == "" {
					if githubURL := getString(metadata, "githubUrl"); githubURL != "" {
						listResponse.GitHubURL = githubURL
					}
				}
				if listResponse.LiveDemoURL == "" {
					if liveDemoURL := getString(metadata, "liveDemoUrl"); liveDemoURL != "" {
						listResponse.LiveDemoURL = liveDemoURL
					}
					if demoURL := getString(metadata, "demoUrl"); demoURL != "" && listResponse.LiveDemoURL == "" {
						listResponse.LiveDemoURL = demoURL
					}
				}
				if listResponse.ThumbnailURL == "" {
					if featuredImageURL := getString(metadata, "featuredImageUrl"); featuredImageURL != "" {
						listResponse.ThumbnailURL = featuredImageURL
					}
					if thumbnailURL := getString(metadata, "thumbnailUrl"); thumbnailURL != "" && listResponse.ThumbnailURL == "" {
						listResponse.ThumbnailURL = thumbnailURL
					}
				}

				// Keep metadata technologies as additional info only
			}
		}

		// Add technologies from relational tags (primary source)
		for _, tag := range project.Tags {
			listResponse.Technologies = append(listResponse.Technologies, tag.Name)
		}

		projectList = append(projectList, listResponse)
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
		Data:       projectList,
		Pagination: pagination,
	}, nil
}

// UpdateProject updates an existing project
func (s *Service) UpdateProject(id string, req dto.UpdateProjectRequest) (*dto.ProjectResponse, error) {
	// Handle temporary IDs from frontend
	if len(id) > 0 && id[:5] == "temp-" {
		// Create a new project instead
		createReq := dto.CreateProjectRequest{
			Title:           req.Title,
			Slug:            req.Slug,
			Description:     req.Description,
			Content:         req.Content,
			ThumbnailURL:    req.ThumbnailURL,
			Status:          req.Status,
			CategoryID:      req.CategoryID,
			GitHubURL:       req.GitHubURL,
			LiveDemoURL:     req.LiveDemoURL,
			Technologies:    req.Technologies,
			TechnologyNames: req.TechnologyNames,
			Metadata:        req.Metadata,
		}
		return s.CreateProject(createReq)
	}

	// Get existing project
	project, err := s.projectRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Update ALL fields from request (allow empty values to be set)
	project.Title = req.Title
	project.Description = req.Description
	project.Content = req.Content
	project.ThumbnailURL = req.ThumbnailURL
	project.GitHubURL = req.GitHubURL
	project.LiveDemoURL = req.LiveDemoURL

	// Update slug if provided, otherwise generate from title
	if req.Slug != "" {
		project.Slug = req.Slug
	} else if req.Title != "" {
		project.Slug = slug.Make(req.Title)
	}

	// Update status and category if provided
	if req.Status != "" {
		project.Status = req.Status
	}
	if req.CategoryID != nil {
		project.CategoryID = req.CategoryID
	}

	// Update metadata - always set to ensure it's updated
	if req.Metadata != nil {
		metadataJSON, err := json.Marshal(req.Metadata)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal metadata: %v", err)
		}
		project.Metadata = string(metadataJSON)
	} else {
		// Set empty JSON object if no metadata provided
		project.Metadata = "{}"
	}

	// Handle technologies update - always update if provided
	if len(req.Technologies) > 0 || len(req.TechnologyNames) > 0 {
		// Resolve technology IDs from either IDs or names
		technologyIDs, err := s.resolveTechnologies(req.Technologies, req.TechnologyNames)
		if err != nil {
			return nil, fmt.Errorf("failed to resolve technologies: %v", err)
		}

		// Update project technologies
		if err := s.projectRepo.UpdateProjectTechnologies(project.ID, technologyIDs); err != nil {
			return nil, fmt.Errorf("failed to update project technologies: %v", err)
		}
	} else {
		// If no technologies provided, clear existing ones
		if err := s.projectRepo.UpdateProjectTechnologies(project.ID, []int{}); err != nil {
			return nil, fmt.Errorf("failed to clear project technologies: %v", err)
		}
	}

	// Update the project
	if err := s.projectRepo.Update(project); err != nil {
		return nil, err
	}

	// Create response
	response := &dto.ProjectResponse{
		ID:           project.ID,
		Title:        project.Title,
		Slug:         project.Slug,
		Description:  project.Description,
		Content:      project.Content,
		ThumbnailURL: project.ThumbnailURL,
		Status:       project.Status,
		GitHubURL:    project.GitHubURL,
		LiveDemoURL:  project.LiveDemoURL,
		CreatedAt:    project.CreatedAt,
		UpdatedAt:    project.UpdatedAt,
		Author: dto.AuthorResponse{
			ID:       project.AuthorID,
			Username: "user", // Default username
		},
	}

	// Add category if available
	if project.Category != nil {
		response.Category = &dto.CategoryResponse{
			ID:   *project.CategoryID,
			Name: project.Category.Name,
			Slug: project.Category.Slug,
		}
	}

	return response, nil
}

// DeleteProject deletes a project by ID
func (s *Service) DeleteProject(id string) error {
	return s.projectRepo.Delete(id)
}

// AddProjectImage adds a new image to a project (simplified stub)
func (s *Service) AddProjectImage(projectID string, imageData dto.ProjectImageData) (*dto.ProjectImageResponse, error) {
	// This is a simplified implementation
	return &dto.ProjectImageResponse{
		ID:        uuid.New().String(),
		URL:       imageData.URL,
		Caption:   imageData.Caption,
		SortOrder: imageData.SortOrder,
	}, nil
}

// AddProjectVideo adds a new video to a project (simplified stub)
func (s *Service) AddProjectVideo(projectID string, videoData dto.ProjectVideoData) (*dto.ProjectVideoResponse, error) {
	// This is a simplified implementation
	return &dto.ProjectVideoResponse{
		ID:        uuid.New().String(),
		URL:       videoData.URL,
		Caption:   videoData.Caption,
		SortOrder: videoData.SortOrder,
	}, nil
}
