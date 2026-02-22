package project

import (
	"encoding/json"
	"fmt"
	"strconv"
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

// resolveTechnologies resolves technology IDs from all sources
func (s *Service) resolveTechnologies(technologies []int, technologyNames []string) ([]int, error) {
	var allIDs []int

	// If we have technology IDs, add them
	if len(technologies) > 0 {
		allIDs = append(allIDs, technologies...)
	}

	// If we have technology names, convert them to IDs and add them
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

// resolveCategoryIDs resolves category IDs from all sources
func (s *Service) resolveCategoryIDs(categories []int, categoryIds []int, categoryIdStrs []string) ([]int, error) {
	var ids []int
	if len(categories) > 0 {
		ids = append(ids, categories...)
	}
	if len(categoryIds) > 0 {
		ids = append(ids, categoryIds...)
	}
	for _, idStr := range categoryIdStrs {
		if id, err := strconv.Atoi(idStr); err == nil {
			ids = append(ids, id)
		} else if idStr != "" {
			// Try to find or create category by name
			cat, err := s.categoryRepo.FindByName(idStr)
			if err != nil {
				// Category doesn't exist, create it
				newCat := &models.Category{
					Name: idStr,
					Slug: slug.Make(idStr),
				}
				if err := s.categoryRepo.Create(newCat); err != nil {
					fmt.Printf("Error creating category '%s': %v\n", idStr, err)
					continue
				}
				ids = append(ids, newCat.ID)
			} else {
				ids = append(ids, cat.ID)
			}
		}
	}
	return s.deduplicateIDs(ids), nil
}

// generateUniqueSlug memastikan slug yang dihasilkan unik di database.
// Jika slug sudah ada, tambahkan suffix -2, -3, dst.
func (s *Service) generateUniqueSlug(base string) string {
	candidate := base
	for i := 2; i <= 100; i++ {
		existing, err := s.projectRepo.GetBySlug(candidate)
		if err != nil || existing == nil {
			// Slug belum dipakai — bisa digunakan
			return candidate
		}
		// Sudah ada, coba suffix berikutnya
		candidate = fmt.Sprintf("%s-%d", base, i)
	}
	// Fallback: tambah UUID pendek agar tetap unik
	return fmt.Sprintf("%s-%s", base, uuid.New().String()[:8])
}

// CreateProject creates a new project
func (s *Service) CreateProject(req dto.CreateProjectRequest) (*dto.ProjectResponse, error) {
	// Generate slug unik jika tidak disediakan
	baseSlug := req.Slug
	if baseSlug == "" {
		baseSlug = slug.Make(req.Title)
	}
	req.Slug = s.generateUniqueSlug(baseSlug)

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
		AuthorID:     authorID,
		GitHubURL:    req.GitHubURL,
		LiveDemoURL:  req.LiveDemoURL,
		Metadata:     string(metadataJSON),
	}

	// Create the project
	if err := s.projectRepo.Create(project); err != nil {
		return nil, err
	}

	// Handle categories if provided
	categoryIDs, err := s.resolveCategoryIDs(req.Categories, req.CategoryIds, req.CategoryIdStrs)
	if err == nil && len(categoryIDs) > 0 {
		if err := s.projectRepo.UpdateProjectCategories(project.ID, categoryIDs); err != nil {
			fmt.Printf("Error adding category: %v\n", err)
		}
	} else if req.CategoryID != nil && *req.CategoryID > 0 {
		// Backward compatibility for single categoryId
		if err := s.projectRepo.UpdateProjectCategories(project.ID, []int{*req.CategoryID}); err != nil {
			fmt.Printf("Error adding single category: %v\n", err)
		}
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

	// Handle images if provided
	if len(req.Images) > 0 {
		var projectImages []models.ProjectImage
		for _, img := range req.Images {
			projectImages = append(projectImages, models.ProjectImage{
				ProjectID: project.ID,
				URL:       img.URL,
				Caption:   img.Caption,
				SortOrder: img.SortOrder,
			})
		}
		if err := s.projectRepo.UpdateProjectImages(project.ID, projectImages); err != nil {
			return nil, fmt.Errorf("failed to add project images: %v", err)
		}
	}

	// Handle videos if provided
	if len(req.Videos) > 0 {
		var projectVideos []models.ProjectVideo
		for _, vid := range req.Videos {
			projectVideos = append(projectVideos, models.ProjectVideo{
				ProjectID: project.ID,
				URL:       vid.URL,
				Caption:   vid.Caption,
				SortOrder: vid.SortOrder,
			})
		}
		if err := s.projectRepo.UpdateProjectVideos(project.ID, projectVideos); err != nil {
			return nil, fmt.Errorf("failed to add project videos: %v", err)
		}
	}

	// Return project response with full data using GetProjectByID
	return s.GetProjectByID(project.ID)
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
			Username: project.Author.Username,
		},
		Images:       []dto.ProjectImageResponse{},
		Videos:       []dto.ProjectVideoResponse{},
		Technologies: []dto.TagResponse{},
		Metadata:     make(map[string]interface{}),
	}

	// Add images from database (primary source)
	for _, img := range project.Images {
		response.Images = append(response.Images, dto.ProjectImageResponse{
			ID:        img.ID,
			URL:       img.URL,
			Caption:   img.Caption,
			SortOrder: img.SortOrder,
		})
	}

	// Add videos from database (primary source)
	for _, vid := range project.Videos {
		response.Videos = append(response.Videos, dto.ProjectVideoResponse{
			ID:        vid.ID,
			URL:       vid.URL,
			Caption:   vid.Caption,
			SortOrder: vid.SortOrder,
		})
	}

	// Parse metadata and restore rich data
	if project.Metadata != "" {
		var metadata map[string]interface{}
		if err := json.Unmarshal([]byte(project.Metadata), &metadata); err == nil {
			response.Metadata = metadata

			// Restore images from metadata if DB images are empty (backward compatibility)
			if len(response.Images) == 0 {
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
			}

			// Restore videos from metadata if DB videos are empty (backward compatibility)
			if len(response.Videos) == 0 {
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

	// Add categories
	if len(project.Categories) > 0 {
		response.Categories = make([]dto.CategoryResponse, 0)
		for _, cat := range project.Categories {
			catResp := dto.CategoryResponse{
				ID:   cat.ID,
				Name: cat.Name,
				Slug: cat.Slug,
			}
			response.Categories = append(response.Categories, catResp)
		}
		// For backward compatibility
		response.Category = &response.Categories[0]
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

	// Add categories
	if len(project.Categories) > 0 {
		response.Categories = make([]dto.CategoryResponse, 0)
		for _, cat := range project.Categories {
			catResp := dto.CategoryResponse{
				ID:   cat.ID,
				Name: cat.Name,
				Slug: cat.Slug,
			}
			response.Categories = append(response.Categories, catResp)
		}
		// For backward compatibility
		response.Category = &response.Categories[0]
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
			Content:      project.Content,
			ThumbnailURL: project.ThumbnailURL,
			Status:       project.Status,
			AuthorName:   "user", // Default username
			GitHubURL:    project.GitHubURL,
			LiveDemoURL:  project.LiveDemoURL,
			Tags:         make([]dto.TagResponse, 0), // Initialize empty array
			Technologies: []string{},                 // Initialize empty array
			Images:       []dto.ProjectImageResponse{},
			Videos:       []dto.ProjectVideoResponse{},
			Metadata:     make(map[string]interface{}),
			CreatedAt:    project.CreatedAt,
		}

		// Add category if available
		// Add categories
		catModels := make([]dto.CategoryResponse, 0)
		for _, cat := range project.Categories {
			catModels = append(catModels, dto.CategoryResponse{
				ID:   cat.ID,
				Name: cat.Name,
				Slug: cat.Slug,
			})
		}
		listResponse.Categories = catModels

		if len(catModels) > 0 {
			listResponse.Category = catModels[0].Name
			listResponse.CategoryID = &catModels[0].ID
		}

		// Parse metadata for list view
		if project.Metadata != "" {
			var metadata map[string]interface{}
			if err := json.Unmarshal([]byte(project.Metadata), &metadata); err == nil {
				listResponse.Metadata = metadata

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
							listResponse.Images = append(listResponse.Images, image)
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
							listResponse.Videos = append(listResponse.Videos, video)
						}
					}
				}

				// Keep metadata technologies as additional info only
			}
		}

		// Add technologies from relational tags (primary source)
		for _, tag := range project.Tags {
			listResponse.Technologies = append(listResponse.Technologies, tag.Name)
			listResponse.Tags = append(listResponse.Tags, dto.TagResponse{
				ID:   tag.ID,
				Name: tag.Name,
				Slug: tag.Slug,
			})
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
			Categories:      req.Categories,
			CategoryIds:     req.CategoryIds,
			CategoryIdStrs:  req.CategoryIdStrs,
			CategoryID:      req.CategoryID,
			Technologies:    req.Technologies,
			TechnologyNames: req.TechnologyNames,
			Metadata:        req.Metadata,
		}
		if req.Title != nil {
			createReq.Title = *req.Title
		}
		if req.Slug != nil {
			createReq.Slug = *req.Slug
		}
		if req.Description != nil {
			createReq.Description = *req.Description
		}
		if req.Content != nil {
			createReq.Content = *req.Content
		}
		if req.ThumbnailURL != nil {
			createReq.ThumbnailURL = *req.ThumbnailURL
		}
		if req.Status != nil {
			createReq.Status = *req.Status
		}
		if req.GitHubURL != nil {
			createReq.GitHubURL = *req.GitHubURL
		}
		if req.LiveDemoURL != nil {
			createReq.LiveDemoURL = *req.LiveDemoURL
		}
		return s.CreateProject(createReq)
	}

	// Get existing project
	project, err := s.projectRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if req.Title != nil {
		project.Title = *req.Title
	}
	if req.Description != nil {
		project.Description = *req.Description
	}
	if req.Content != nil {
		project.Content = *req.Content
	}
	if req.ThumbnailURL != nil {
		project.ThumbnailURL = *req.ThumbnailURL
	}
	if req.GitHubURL != nil {
		project.GitHubURL = *req.GitHubURL
	}
	if req.LiveDemoURL != nil {
		project.LiveDemoURL = *req.LiveDemoURL
	}

	// Update slug jika disediakan, atau generate dari title dengan unique check
	if req.Slug != nil && *req.Slug != "" {
		project.Slug = *req.Slug
	} else if req.Title != nil && *req.Title != "" && *req.Title != project.Title {
		// Hanya re-generate slug jika title benar-benar berubah
		project.Slug = s.generateUniqueSlug(slug.Make(*req.Title))
	}

	// Update status and category if provided
	if req.Status != nil && *req.Status != "" {
		project.Status = *req.Status
	}
	// Update categories - always update if provided
	categoryIDs, err := s.resolveCategoryIDs(req.Categories, req.CategoryIds, req.CategoryIdStrs)
	if err == nil && (len(categoryIDs) > 0 || len(req.CategoryIdStrs) > 0) {
		if err := s.projectRepo.UpdateProjectCategories(project.ID, categoryIDs); err != nil {
			fmt.Printf("Error updating categories: %v\n", err)
		}
	} else if req.CategoryID != nil {
		if *req.CategoryID > 0 {
			if err := s.projectRepo.UpdateProjectCategories(project.ID, []int{*req.CategoryID}); err != nil {
				fmt.Printf("Error updating single category: %v\n", err)
			}
		} else {
			// Clear categories if CategoryID=0 or null
			if err := s.projectRepo.UpdateProjectCategories(project.ID, []int{}); err != nil {
				fmt.Printf("Error clearing categories: %v\n", err)
			}
		}
	}

	// Update metadata - always set to ensure it's updated
	metaMap := make(map[string]interface{})
	if project.Metadata != "" {
		_ = json.Unmarshal([]byte(project.Metadata), &metaMap)
	}

	metaMap["githubUrl"] = project.GitHubURL
	metaMap["liveDemoUrl"] = project.LiveDemoURL
	if req.Metadata != nil {
		for k, v := range req.Metadata {
			metaMap[k] = v
		}
	}
	metadataJSON, err := json.Marshal(metaMap)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal metadata: %v", err)
	}
	project.Metadata = string(metadataJSON)

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

	// Update images if provided
	if req.Images != nil {
		var projectImages []models.ProjectImage
		for _, img := range req.Images {
			projectImages = append(projectImages, models.ProjectImage{
				ProjectID: project.ID,
				URL:       img.URL,
				Caption:   img.Caption,
				SortOrder: img.SortOrder,
			})
		}
		if err := s.projectRepo.UpdateProjectImages(project.ID, projectImages); err != nil {
			return nil, fmt.Errorf("failed to update project images: %v", err)
		}
	}

	// Update videos if provided
	if req.Videos != nil {
		var projectVideos []models.ProjectVideo
		for _, vid := range req.Videos {
			projectVideos = append(projectVideos, models.ProjectVideo{
				ProjectID: project.ID,
				URL:       vid.URL,
				Caption:   vid.Caption,
				SortOrder: vid.SortOrder,
			})
		}
		if err := s.projectRepo.UpdateProjectVideos(project.ID, projectVideos); err != nil {
			return nil, fmt.Errorf("failed to update project videos: %v", err)
		}
	}

	// Update the project
	if err := s.projectRepo.Update(project); err != nil {
		return nil, err
	}

	// Reload project dari DB agar mendapatkan Tags (technologies) yang terbaru
	updatedProject, err := s.projectRepo.GetByID(project.ID)
	if err != nil {
		// Jika gagal reload, kembalikan response tanpa technologies
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
				Username: "user",
			},
		}, nil
	}

	// Gunakan GetProjectByID untuk membangun response lengkap dengan technologies
	return s.GetProjectByID(updatedProject.ID)
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
