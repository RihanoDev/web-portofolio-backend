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
			continue
		}

		if id, err := strconv.Atoi(name); err == nil {
			technologyIDs = append(technologyIDs, id)
			continue
		}

		tag, err := s.tagService.GetByName(name)
		if err != nil {
			createReq := &dto.CreateTagRequest{Name: name}
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

func (s *Service) deduplicateIDs(ids []int) []int {
	keys := make(map[int]bool)
	list := []int{}
	for _, entry := range ids {
		if _, value := keys[entry]; !value && entry > 0 {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func (s *Service) convertCategoryNamesToIDs(categoryNames []string) ([]int, error) {
	var categoryIDs []int
	for _, name := range categoryNames {
		if name == "" {
			continue
		}
		if id, err := strconv.Atoi(name); err == nil {
			categoryIDs = append(categoryIDs, id)
			continue
		}
		cat, err := s.categoryRepo.FindByName(name)
		if err != nil {
			newCat := &models.Category{Name: name, Slug: slug.Make(name)}
			err := s.categoryRepo.Create(newCat)
			if err != nil {
				return nil, fmt.Errorf("failed to create category '%s': %v", name, err)
			}
			categoryIDs = append(categoryIDs, newCat.ID)
		} else {
			categoryIDs = append(categoryIDs, cat.ID)
		}
	}
	return categoryIDs, nil
}

func (s *Service) resolveCategoryIDs(categoryID *int, categories []int, categoryIds []int, categoryIdStrs []string) ([]int, error) {
	var allIDs []int
	if categoryID != nil && *categoryID > 0 {
		allIDs = append(allIDs, *categoryID)
	}
	if len(categories) > 0 {
		allIDs = append(allIDs, categories...)
	}
	if len(categoryIds) > 0 {
		allIDs = append(allIDs, categoryIds...)
	}
	if len(categoryIdStrs) > 0 {
		ids, err := s.convertCategoryNamesToIDs(categoryIdStrs)
		if err == nil {
			allIDs = append(allIDs, ids...)
		}
	}
	return s.deduplicateIDs(allIDs), nil
}

func (s *Service) resolveTechnologies(technologies []int, technologyNames []string) ([]int, error) {
	var allIDs []int
	if len(technologies) > 0 {
		allIDs = append(allIDs, technologies...)
	}
	if len(technologyNames) > 0 {
		ids, err := s.convertTechnologyNamesToIDs(technologyNames)
		if err == nil {
			allIDs = append(allIDs, ids...)
		}
	}
	return s.deduplicateIDs(allIDs), nil
}

func (s *Service) resolveTags(tags []int, tagIds []int, tagIdStrs []string, tagNames []string) ([]int, error) {
	var allIDs []int
	if len(tags) > 0 {
		allIDs = append(allIDs, tags...)
	}
	if len(tagIds) > 0 {
		allIDs = append(allIDs, tagIds...)
	}
	if len(tagIdStrs) > 0 {
		ids, err := s.convertTechnologyNamesToIDs(tagIdStrs)
		if err == nil {
			allIDs = append(allIDs, ids...)
		}
	}
	if len(tagNames) > 0 {
		ids, err := s.convertTechnologyNamesToIDs(tagNames)
		if err == nil {
			allIDs = append(allIDs, ids...)
		}
	}
	return s.deduplicateIDs(allIDs), nil
}

func (s *Service) generateUniqueSlug(base string) string {
	candidate := base
	for i := 2; i <= 100; i++ {
		existing, err := s.projectRepo.GetBySlug(candidate)
		if err != nil || existing == nil {
			return candidate
		}
		candidate = fmt.Sprintf("%s-%d", base, i)
	}
	return fmt.Sprintf("%s-%s", base, uuid.New().String()[:8])
}

func (s *Service) CreateProject(req dto.CreateProjectRequest) (*dto.ProjectResponse, error) {
	baseSlug := req.Slug
	if baseSlug == "" {
		baseSlug = slug.Make(req.Title)
	}
	req.Slug = s.generateUniqueSlug(baseSlug)

	authorID := req.AuthorID
	if authorID == 0 {
		defaultUser, _ := s.userService.GetDefaultAdmin()
		if defaultUser != nil {
			authorID = int(defaultUser.ID)
		}
	}

	metadata := req.Metadata
	if metadata == nil {
		metadata = make(map[string]interface{})
	}
	metadata["githubUrl"] = req.GitHubURL
	metadata["liveDemoUrl"] = req.LiveDemoURL
	metadataJSON, _ := json.Marshal(metadata)

	project := &models.Project{
		ID:           uuid.New().String(),
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

	if err := s.projectRepo.Create(project); err != nil {
		return nil, err
	}

	techIDs, _ := s.resolveTechnologies(req.Technologies, req.TechnologyNames)
	if len(techIDs) > 0 {
		s.projectRepo.UpdateProjectTechnologies(project.ID, techIDs)
	}

	tagIDs, _ := s.resolveTags(req.Tags, req.TagIds, req.TagIdStrs, req.TagNames)
	if len(tagIDs) > 0 {
		s.projectRepo.UpdateProjectTags(project.ID, tagIDs)
	}

	catIDs, _ := s.resolveCategoryIDs(req.CategoryID, req.Categories, req.CategoryIds, req.CategoryIdStrs)
	if len(catIDs) > 0 {
		s.projectRepo.UpdateProjectCategories(project.ID, catIDs)
		firstID := catIDs[0]
		project.CategoryID = &firstID
		s.projectRepo.Update(project)
	}

	if len(req.Images) > 0 {
		var images []models.ProjectImage
		for _, img := range req.Images {
			images = append(images, models.ProjectImage{
				ProjectID: project.ID,
				URL:       img.URL,
				Caption:   img.Caption,
				SortOrder: img.SortOrder,
			})
		}
		s.projectRepo.UpdateProjectImages(project.ID, images)
	}

	if len(req.Videos) > 0 {
		var videos []models.ProjectVideo
		for _, vid := range req.Videos {
			videos = append(videos, models.ProjectVideo{
				ProjectID: project.ID,
				URL:       vid.URL,
				Caption:   vid.Caption,
				SortOrder: vid.SortOrder,
			})
		}
		s.projectRepo.UpdateProjectVideos(project.ID, videos)
	}

	return s.GetProjectByID(project.ID)
}

func (s *Service) GetProjectByID(id string) (*dto.ProjectResponse, error) {
	project, err := s.projectRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	return s.mapToResponse(project), nil
}

func (s *Service) GetProjectBySlug(slug string) (*dto.ProjectResponse, error) {
	project, err := s.projectRepo.GetBySlug(slug)
	if err != nil {
		return nil, err
	}
	return s.mapToResponse(project), nil
}

func (s *Service) UpdateProject(id string, req dto.UpdateProjectRequest) (*dto.ProjectResponse, error) {
	if len(id) > 5 && id[:5] == "temp-" {
		createReq := dto.CreateProjectRequest{
			CategoryID:      req.CategoryID,
			CategoryIds:     req.CategoryIds,
			CategoryIdStrs:  req.CategoryIdStrs,
			Categories:      req.Categories,
			Technologies:    req.Technologies,
			TechnologyNames: req.TechnologyNames,
			Tags:            req.Tags,
			TagIds:          req.TagIds,
			TagIdStrs:       req.TagIdStrs,
			TagNames:        req.TagNames,
			Metadata:        req.Metadata,
			Images:          req.Images,
			Videos:          req.Videos,
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
	if req.Status != nil {
		project.Status = *req.Status
	}
	if req.GitHubURL != nil {
		project.GitHubURL = *req.GitHubURL
	}
	if req.LiveDemoURL != nil {
		project.LiveDemoURL = *req.LiveDemoURL
	}

	if req.Slug != nil && *req.Slug != "" {
		project.Slug = *req.Slug
	} else if req.Title != nil && project.Slug == "" {
		project.Slug = s.generateUniqueSlug(slug.Make(*req.Title))
	}

	metaMap := make(map[string]interface{})
	if project.Metadata != "" {
		_ = json.Unmarshal([]byte(project.Metadata), &metaMap)
	}
	for k, v := range req.Metadata {
		metaMap[k] = v
	}
	metaMap["githubUrl"] = project.GitHubURL
	metaMap["liveDemoUrl"] = project.LiveDemoURL
	metadataJSON, _ := json.Marshal(metaMap)
	project.Metadata = string(metadataJSON)

	if err := s.projectRepo.Update(project); err != nil {
		return nil, err
	}

	if req.Technologies != nil || req.TechnologyNames != nil {
		techIDs, err := s.resolveTechnologies(req.Technologies, req.TechnologyNames)
		if err == nil {
			s.projectRepo.UpdateProjectTechnologies(project.ID, techIDs)
		}
	}

	if req.Tags != nil || req.TagIds != nil || req.TagIdStrs != nil || req.TagNames != nil {
		tagIDs, err := s.resolveTags(req.Tags, req.TagIds, req.TagIdStrs, req.TagNames)
		if err == nil {
			s.projectRepo.UpdateProjectTags(project.ID, tagIDs)
		}
	}

	if req.CategoryID != nil || req.Categories != nil || req.CategoryIds != nil || req.CategoryIdStrs != nil {
		catIDs, err := s.resolveCategoryIDs(req.CategoryID, req.Categories, req.CategoryIds, req.CategoryIdStrs)
		if err == nil {
			s.projectRepo.UpdateProjectCategories(project.ID, catIDs)
			if len(catIDs) > 0 {
				if req.CategoryID != nil && *req.CategoryID > 0 {
					project.CategoryID = req.CategoryID
				} else {
					firstID := catIDs[0]
					project.CategoryID = &firstID
				}
			} else {
				project.CategoryID = nil
			}
			s.projectRepo.Update(project)
		}
	}

	if len(req.Images) > 0 {
		var images []models.ProjectImage
		for _, img := range req.Images {
			images = append(images, models.ProjectImage{
				ProjectID: project.ID,
				URL:       img.URL,
				Caption:   img.Caption,
				SortOrder: img.SortOrder,
			})
		}
		s.projectRepo.UpdateProjectImages(project.ID, images)
	}

	if len(req.Videos) > 0 {
		var videos []models.ProjectVideo
		for _, vid := range req.Videos {
			videos = append(videos, models.ProjectVideo{
				ProjectID: project.ID,
				URL:       vid.URL,
				Caption:   vid.Caption,
				SortOrder: vid.SortOrder,
			})
		}
		s.projectRepo.UpdateProjectVideos(project.ID, videos)
	}

	return s.GetProjectByID(project.ID)
}

func (s *Service) DeleteProject(id string) error {
	return s.projectRepo.Delete(id)
}

func (s *Service) ListProjects(page, size int) (*dto.PaginatedResponse, error) {
	offset := (page - 1) * size
	projects, total, err := s.projectRepo.GetAll(size, offset)
	if err != nil {
		return nil, err
	}

	projectList := make([]dto.ProjectListResponse, 0)
	for _, project := range projects {
		projectList = append(projectList, s.mapToListResponse(project))
	}

	return &dto.PaginatedResponse{
		Data: projectList,
		Pagination: dto.PaginationResponse{
			TotalCount:  total,
			CurrentPage: page,
			PageSize:    size,
			TotalPages:  int((total + int64(size) - 1) / int64(size)),
			HasNext:     int64(page*size) < total,
			HasPrevious: page > 1,
		},
	}, nil
}

func (s *Service) GetProjectsByCategorySlug(slug string, page, size int) (*dto.PaginatedResponse, error) {
	offset := (page - 1) * size
	projects, total, err := s.projectRepo.GetByCategorySlug(slug, size, offset)
	if err != nil {
		return nil, err
	}

	projectList := make([]dto.ProjectListResponse, 0)
	for _, project := range projects {
		projectList = append(projectList, s.mapToListResponse(project))
	}

	return &dto.PaginatedResponse{
		Data: projectList,
		Pagination: dto.PaginationResponse{
			TotalCount:  total,
			CurrentPage: page,
			PageSize:    size,
		},
	}, nil
}

func (s *Service) mapToResponse(project *models.Project) *dto.ProjectResponse {
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
		Categories:   []dto.CategoryResponse{},
		Metadata:     make(map[string]interface{}),
	}
	if response.Author.Username == "" {
		response.Author.Username = "user"
	}

	for _, cat := range project.Categories {
		response.Categories = append(response.Categories, dto.CategoryResponse{
			ID: cat.ID, Name: cat.Name, Slug: cat.Slug,
		})
	}
	if len(response.Categories) > 0 {
		response.Category = &response.Categories[0]
	} else if project.Category.ID != 0 {
		response.Category = &dto.CategoryResponse{
			ID: project.Category.ID, Name: project.Category.Name, Slug: project.Category.Slug,
		}
		response.Categories = append(response.Categories, *response.Category)
	}

	for _, img := range project.Images {
		response.Images = append(response.Images, dto.ProjectImageResponse{
			ID: img.ID, URL: img.URL, Caption: img.Caption, SortOrder: img.SortOrder,
		})
	}
	for _, vid := range project.Videos {
		response.Videos = append(response.Videos, dto.ProjectVideoResponse{
			ID: vid.ID, URL: vid.URL, Caption: vid.Caption, SortOrder: vid.SortOrder,
		})
	}
	for _, tech := range project.Technologies {
		response.Technologies = append(response.Technologies, dto.TagResponse{
			ID: tech.ID, Name: tech.Name, Slug: tech.Slug,
		})
	}
	for _, tag := range project.Tags {
		response.Tags = append(response.Tags, dto.TagResponse{
			ID: tag.ID, Name: tag.Name, Slug: tag.Slug,
		})
	}

	if project.Metadata != "" {
		_ = json.Unmarshal([]byte(project.Metadata), &response.Metadata)
	}

	return response
}

func (s *Service) mapToListResponse(project *models.Project) dto.ProjectListResponse {
	response := dto.ProjectListResponse{
		ID:             project.ID,
		Title:          project.Title,
		Slug:           project.Slug,
		Description:    project.Description,
		Content:        project.Content,
		ThumbnailURL:   project.ThumbnailURL,
		Status:         project.Status,
		CategoryID:     project.CategoryID,
		AuthorName:     project.Author.Username,
		GitHubURL:      project.GitHubURL,
		LiveDemoURL:    project.LiveDemoURL,
		Tags:           []dto.TagResponse{},
		Technologies:   []dto.TagResponse{},
		TagStrs:        []string{},
		TechnologyStrs: []string{},
		Categories:     []string{},
		CategoryModels: []dto.CategoryResponse{},
		Metadata:       make(map[string]interface{}),
		CreatedAt:      project.CreatedAt,
	}
	if response.AuthorName == "" {
		response.AuthorName = "user"
	}

	if len(project.Categories) > 0 {
		response.Category = project.Categories[0].Name
		for _, cat := range project.Categories {
			response.Categories = append(response.Categories, cat.Name)
			response.CategoryModels = append(response.CategoryModels, dto.CategoryResponse{
				ID: cat.ID, Name: cat.Name, Slug: cat.Slug,
			})
		}
	} else if project.Category.ID != 0 {
		response.Category = project.Category.Name
		catRes := dto.CategoryResponse{
			ID: project.Category.ID, Name: project.Category.Name, Slug: project.Category.Slug,
		}
		response.Categories = append(response.Categories, catRes.Name)
		response.CategoryModels = append(response.CategoryModels, catRes)
	}

	for _, tech := range project.Technologies {
		techRes := dto.TagResponse{
			ID: tech.ID, Name: tech.Name, Slug: tech.Slug,
		}
		response.Technologies = append(response.Technologies, techRes)
		response.TechnologyStrs = append(response.TechnologyStrs, tech.Name)
	}

	for _, tag := range project.Tags {
		tagRes := dto.TagResponse{
			ID: tag.ID, Name: tag.Name, Slug: tag.Slug,
		}
		response.Tags = append(response.Tags, tagRes)
		response.TagStrs = append(response.TagStrs, tag.Name)
	}

	if project.Metadata != "" {
		_ = json.Unmarshal([]byte(project.Metadata), &response.Metadata)
	}

	return response
}

func (s *Service) AddProjectImage(projectID string, imageData dto.ProjectImageData) (*dto.ProjectImageResponse, error) {
	return &dto.ProjectImageResponse{
		ID: uuid.New().String(), URL: imageData.URL, Caption: imageData.Caption,
	}, nil
}

func (s *Service) AddProjectVideo(projectID string, videoData dto.ProjectVideoData) (*dto.ProjectVideoResponse, error) {
	return &dto.ProjectVideoResponse{
		ID: uuid.New().String(), URL: videoData.URL, Caption: videoData.Caption,
	}, nil
}
