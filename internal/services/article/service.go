package article

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
	"web-porto-backend/internal/domain/dto"
	"web-porto-backend/internal/domain/models"
	articleRepo "web-porto-backend/internal/repositories/article"
	categoryRepo "web-porto-backend/internal/repositories/category"
	tagRepo "web-porto-backend/internal/repositories/tag"
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

// Service handles business logic for articles
type Service struct {
	articleRepo  articleRepo.Repository
	categoryRepo categoryRepo.Repository
	tagRepo      tagRepo.Repository
	userService  userService.Service
}

// NewService creates a new article service
func NewService(
	articleRepo articleRepo.Repository,
	categoryRepo categoryRepo.Repository,
	tagRepo tagRepo.Repository,
	userService userService.Service,
) *Service {
	return &Service{
		articleRepo:  articleRepo,
		categoryRepo: categoryRepo,
		tagRepo:      tagRepo,
		userService:  userService,
	}
}

// convertStringIDsToInts converts string IDs to integer IDs, handling both numeric IDs and tag names
func (s *Service) convertStringIDsToInts(stringIDs []string, isForTags bool) ([]int, error) {
	var intIDs []int
	for _, strID := range stringIDs {
		if strID == "" {
			continue
		}

		// Try to convert to integer first
		if intID, err := strconv.Atoi(strID); err == nil {
			// It's a numeric ID, use it directly
			intIDs = append(intIDs, intID)
		} else if isForTags {
			// It's a tag name, try to find or create the tag
			tag, err := s.tagRepo.GetByName(strID)
			if err != nil {
				// Tag doesn't exist, create it
				newTag := &models.Tag{
					Name: strID,
					Slug: slug.Make(strID),
				}
				createdTag, err := s.tagRepo.Create(newTag)
				if err != nil {
					return nil, fmt.Errorf("failed to create tag '%s': %v", strID, err)
				}
				intIDs = append(intIDs, createdTag.ID)
			} else {
				intIDs = append(intIDs, tag.ID)
			}
		} else {
			return nil, fmt.Errorf("failed to convert ID '%s' to integer: %v", strID, err)
		}
	}
	return intIDs, nil
}

// resolveCategoryIDs resolves category IDs from either int or string arrays
func (s *Service) resolveCategoryIDs(categories []int, categoryIds []int, categoryIdStrs []string) ([]int, error) {
	// Priority: categories -> categoryIds -> categoryIdStrs
	if len(categories) > 0 {
		return categories, nil
	}
	if len(categoryIds) > 0 {
		return categoryIds, nil
	}
	if len(categoryIdStrs) > 0 {
		return s.convertStringIDsToInts(categoryIdStrs, false)
	}
	return []int{}, nil
}

// resolveTagIDs resolves tag IDs from either int or string arrays
func (s *Service) resolveTagIDs(tags []int, tagIds []int, tagIdStrs []string) ([]int, error) {
	// Priority: tags -> tagIds -> tagIdStrs
	if len(tags) > 0 {
		return tags, nil
	}
	if len(tagIds) > 0 {
		return tagIds, nil
	}
	if len(tagIdStrs) > 0 {
		return s.convertStringIDsToInts(tagIdStrs, true)
	}
	return []int{}, nil
}

// generateUniqueSlug memastikan slug unik di database; append -2, -3, dst jika sudah ada.
func (s *Service) generateUniqueSlug(base string) string {
	candidate := base
	for i := 2; i <= 100; i++ {
		existing, err := s.articleRepo.GetBySlug(candidate)
		if err != nil || existing == nil {
			return candidate
		}
		candidate = fmt.Sprintf("%s-%d", base, i)
	}
	return fmt.Sprintf("%s-%s", base, uuid.New().String()[:8])
}

// CreateArticle creates a new article
func (s *Service) CreateArticle(req dto.CreateArticleRequest) (*dto.ArticleResponse, error) {
	// Generate slug unik
	baseSlug := req.Slug
	if baseSlug == "" {
		baseSlug = slug.Make(req.Title)
	}
	req.Slug = s.generateUniqueSlug(baseSlug)

	// Calculate read time (words per minute: 200)
	readTime := len(strings.Fields(req.Content)) / 200
	if readTime < 1 {
		readTime = 1
	}

	// Set default authorID if not provided - get from database
	authorID := req.AuthorID
	if authorID == 0 {
		// Get default admin user from database
		defaultUser, err := s.userService.GetDefaultAdmin()
		if err != nil {
			return nil, fmt.Errorf("failed to get default admin user: %v", err)
		}
		authorID = int(defaultUser.ID)
	}

	// Create article
	article := &models.Article{
		Title:            req.Title,
		Slug:             req.Slug,
		Excerpt:          req.Excerpt,
		Content:          req.Content,
		FeaturedImageURL: req.FeaturedImageURL,
		Status:           req.Status,
		AuthorID:         authorID,
		ReadTime:         readTime,
	}

	if req.PublishAt != nil && !req.PublishAt.IsZero() {
		article.PublishedAt = req.PublishAt
	} else if req.Status == "published" {
		now := time.Now()
		article.PublishedAt = &now
	}

	// Handle metadata
	metadata := map[string]interface{}{}

	// Add any additional metadata from request
	if req.Metadata != nil {
		for k, v := range req.Metadata {
			metadata[k] = v
		}
	}

	// Convert to JSON string
	metadataJSON, err := json.Marshal(metadata)
	if err != nil {
		return nil, err
	}
	article.Metadata = string(metadataJSON)

	// Create the article
	if err := s.articleRepo.Create(article); err != nil {
		return nil, err
	}

	// Add categories if any
	categoryIDs, err := s.resolveCategoryIDs(req.Categories, req.CategoryIds, req.CategoryIdStrs)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve category IDs: %v", err)
	}
	if len(categoryIDs) > 0 {
		for _, categoryID := range categoryIDs {
			category, err := s.categoryRepo.FindByID(categoryID)
			if err != nil {
				// Log error but continue
				fmt.Printf("Error adding category: %v\n", err)
				continue
			}
			article.Categories = append(article.Categories, *category)
		}
		// Update article with categories
		if err := s.articleRepo.Update(article); err != nil {
			return nil, err
		}
	}

	// Add tags if any
	tagIDs, err := s.resolveTagIDs(req.Tags, req.TagIds, req.TagIdStrs)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve tag IDs: %v", err)
	}
	if len(tagIDs) > 0 {
		for _, tagID := range tagIDs {
			tag, err := s.tagRepo.GetByID(tagID)
			if err != nil {
				fmt.Printf("Error adding tag: %v\n", err)
				continue
			}
			article.Tags = append(article.Tags, *tag)
		}
		// Update article with tags
		if err := s.articleRepo.Update(article); err != nil {
			return nil, err
		}
	}

	// Convert article to response
	return s.mapArticleToResponse(article), nil
}

// GetArticleByID retrieves an article by ID
func (s *Service) GetArticleByID(id string) (*dto.ArticleResponse, error) {
	article, err := s.articleRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	return s.mapArticleToResponse(article), nil
}

// GetArticleBySlug retrieves an article by slug
func (s *Service) GetArticleBySlug(slug string) (*dto.ArticleResponse, error) {
	article, err := s.articleRepo.GetBySlug(slug)
	if err != nil {
		return nil, err
	}

	// Update view count
	article.ViewCount++
	if err := s.articleRepo.Update(article); err != nil {
		// Log error but don't fail the request
		fmt.Printf("Error updating view count: %v\n", err)
	}

	return s.mapArticleToResponse(article), nil
}

// GetArticlesByCategorySlug retrieves articles by category slug
func (s *Service) GetArticlesByCategorySlug(slug string, page, size int) (*dto.PaginatedResponse, error) {
	// Get category by slug first
	category, err := s.categoryRepo.FindBySlug(slug)
	if err != nil {
		return nil, fmt.Errorf("category not found: %w", err)
	}

	// Get articles by category ID
	offset := (page - 1) * size
	articles, total, err := s.articleRepo.GetByCategory(category.ID, size, offset)
	if err != nil {
		return nil, err
	}

	// Map articles to response objects
	articlesResponse := make([]interface{}, 0, len(articles))
	for _, article := range articles {
		articlesResponse = append(articlesResponse, s.mapArticleToListResponse(article))
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
		Data:       articlesResponse,
		Pagination: pagination,
	}, nil
}

// ListArticles retrieves a paginated list of articles
func (s *Service) ListArticles(page, size int) (*dto.PaginatedResponse, error) {
	offset := (page - 1) * size
	articles, total, err := s.articleRepo.GetAll(size, offset)
	if err != nil {
		return nil, err
	}

	// Map articles to response objects
	articlesResponse := make([]interface{}, 0, len(articles))
	for _, article := range articles {
		articlesResponse = append(articlesResponse, s.mapArticleToListResponse(article))
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
		Data:       articlesResponse,
		Pagination: pagination,
	}, nil
}

// UpdateArticle updates an existing article
func (s *Service) UpdateArticle(id string, req dto.UpdateArticleRequest) (*dto.ArticleResponse, error) {
	// Check if this is a new article (using temporary ID)
	if len(id) > 0 && id[:5] == "temp-" {
		createReq := dto.CreateArticleRequest{
			AuthorID:       1, // Default to admin user
			Categories:     req.Categories,
			CategoryIds:    req.CategoryIds,
			CategoryIdStrs: req.CategoryIdStrs,
			Tags:           req.Tags,
			TagIds:         req.TagIds,
			TagIdStrs:      req.TagIdStrs,
			PublishAt:      req.PublishAt,
			Metadata:       req.Metadata,
		}
		if req.Title != nil {
			createReq.Title = *req.Title
		}
		if req.Slug != nil {
			createReq.Slug = *req.Slug
		}
		if req.Excerpt != nil {
			createReq.Excerpt = *req.Excerpt
		}
		if req.Content != nil {
			createReq.Content = *req.Content
		}
		if req.FeaturedImageURL != nil {
			createReq.FeaturedImageURL = *req.FeaturedImageURL
		}
		if req.Status != nil {
			createReq.Status = *req.Status
		}
		return s.CreateArticle(createReq)
	}

	// Get existing article
	article, err := s.articleRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if req.Title != nil {
		article.Title = *req.Title
	}
	if req.Excerpt != nil {
		article.Excerpt = *req.Excerpt
	}
	if req.Content != nil {
		article.Content = *req.Content
	}
	if req.FeaturedImageURL != nil {
		article.FeaturedImageURL = *req.FeaturedImageURL
	}

	// Update slug jika disediakan, atau generate dari title dengan unique check
	if req.Slug != nil && *req.Slug != "" {
		article.Slug = *req.Slug
	} else if req.Title != nil && *req.Title != "" && *req.Title != article.Title {
		// Hanya re-generate slug jika title benar-benar berubah
		article.Slug = s.generateUniqueSlug(slug.Make(*req.Title))
	}

	// Update status if provided
	if req.Status != nil && *req.Status != "" {
		article.Status = *req.Status
		// Update published date if status changes to published
		if *req.Status == "published" && article.PublishedAt == nil {
			now := time.Now()
			article.PublishedAt = &now
		}
	}

	// Update publish date if provided
	if req.PublishAt != nil {
		article.PublishedAt = req.PublishAt
	}

	// Recalculate read time if content changed
	if req.Content != nil && *req.Content != "" {
		readTime := len(strings.Fields(*req.Content)) / 200
		if readTime < 1 {
			readTime = 1
		}
		article.ReadTime = readTime
	}

	// Update categories - always update if provided
	categoryIDs, err := s.resolveCategoryIDs(req.Categories, req.CategoryIds, req.CategoryIdStrs)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve category IDs: %v", err)
	}

	// Always clear and update categories (even if empty)
	var newCategories []models.Category
	for _, categoryID := range categoryIDs {
		category, err := s.categoryRepo.FindByID(categoryID)
		if err != nil {
			// Log error but continue
			fmt.Printf("Error adding category: %v\n", err)
			continue
		}
		newCategories = append(newCategories, *category)
	}
	article.Categories = newCategories

	// Update tags - always update if provided
	tagIDs, err := s.resolveTagIDs(req.Tags, req.TagIds, req.TagIdStrs)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve tag IDs: %v", err)
	}

	// Always clear and update tags (even if empty)
	var newTags []models.Tag
	for _, tagID := range tagIDs {
		tag, err := s.tagRepo.GetByID(tagID)
		if err != nil {
			fmt.Printf("Error adding tag: %v\n", err)
			continue
		}
		newTags = append(newTags, *tag)
	}
	article.Tags = newTags // Update metadata - always set to ensure it's updated
	if req.Metadata != nil {
		metadataJSON, err := json.Marshal(req.Metadata)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal metadata: %v", err)
		}
		article.Metadata = string(metadataJSON)
	} else {
		// Set empty JSON object if no metadata provided
		article.Metadata = "{}"
	}

	// Update the article
	if err := s.articleRepo.Update(article); err != nil {
		return nil, err
	}

	// Reload dari DB agar mendapatkan Categories dan Tags yang terbaru
	updatedArticle, err := s.articleRepo.GetByID(id)
	if err != nil {
		// Jika gagal reload, kembalikan response dengan data yang ada
		return s.mapArticleToResponse(article), nil
	}

	return s.mapArticleToResponse(updatedArticle), nil
}

// DeleteArticle deletes an article by ID
func (s *Service) DeleteArticle(id string) error {
	return s.articleRepo.Delete(id)
}

// AddArticleImage adds a new image to an article (simplified stub)
func (s *Service) AddArticleImage(articleID string, imageData dto.ArticleImageData) (*dto.ArticleImageResponse, error) {
	// This is a simplified implementation
	// In a real implementation, you would need to create an actual image record in the database
	return &dto.ArticleImageResponse{
		ID:        uuid.New().String(),
		URL:       imageData.URL,
		Caption:   imageData.Caption,
		AltText:   imageData.AltText,
		SortOrder: imageData.SortOrder,
	}, nil
}

// AddArticleVideo adds a new video to an article (simplified stub)
func (s *Service) AddArticleVideo(articleID string, videoData dto.ArticleVideoData) (*dto.ArticleVideoResponse, error) {
	// This is a simplified implementation
	// In a real implementation, you would need to create an actual video record in the database
	return &dto.ArticleVideoResponse{
		ID:        uuid.New().String(),
		URL:       videoData.URL,
		Caption:   videoData.Caption,
		SortOrder: videoData.SortOrder,
	}, nil
}

// Helper function to map Article to ArticleResponse
func (s *Service) mapArticleToResponse(article *models.Article) *dto.ArticleResponse {
	// Create base response
	response := &dto.ArticleResponse{
		ID:               article.ID,
		Title:            article.Title,
		Slug:             article.Slug,
		Excerpt:          article.Excerpt,
		Content:          article.Content,
		FeaturedImageURL: article.FeaturedImageURL,
		Status:           article.Status,
		ReadTime:         article.ReadTime,
		ViewCount:        article.ViewCount,
		PublishedAt:      article.PublishedAt,
		CreatedAt:        article.CreatedAt,
		UpdatedAt:        article.UpdatedAt,
		Author: dto.AuthorResponse{
			ID: article.AuthorID,
			// In a real implementation, get the username from the User model
			Username: article.Author.Username,
		},
		Images:   []dto.ArticleImageResponse{},
		Videos:   []dto.ArticleVideoResponse{},
		Metadata: make(map[string]interface{}),
	}

	// Add categories
	categories := []dto.CategoryResponse{}
	for _, category := range article.Categories {
		categories = append(categories, dto.CategoryResponse{
			ID:   category.ID,
			Name: category.Name,
			Slug: category.Slug,
		})
	}
	response.Categories = categories

	// Add tags
	tags := []dto.TagResponse{}
	for _, tag := range article.Tags {
		tags = append(tags, dto.TagResponse{
			ID:   tag.ID,
			Name: tag.Name,
		})
	}
	response.Tags = tags

	// Parse the metadata JSON string and restore rich data
	if article.Metadata != "" {
		var metadata map[string]interface{}
		if err := json.Unmarshal([]byte(article.Metadata), &metadata); err == nil {
			response.Metadata = metadata

			// Restore images from metadata
			if imagesData, ok := metadata["images"].([]interface{}); ok {
				for _, imgData := range imagesData {
					if imgMap, ok := imgData.(map[string]interface{}); ok {
						image := dto.ArticleImageResponse{
							ID:        getString(imgMap, "id"),
							URL:       getString(imgMap, "url"),
							Caption:   getString(imgMap, "caption"),
							AltText:   getString(imgMap, "altText"),
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
						video := dto.ArticleVideoResponse{
							ID:        getString(vidMap, "id"),
							URL:       getString(vidMap, "url"),
							Caption:   getString(vidMap, "caption"),
							SortOrder: getInt(vidMap, "sortOrder"),
						}
						response.Videos = append(response.Videos, video)
					}
				}
			}

			// Restore featured image URL from metadata if not set
			if response.FeaturedImageURL == "" {
				if featuredImageURL := getString(metadata, "featuredImageUrl"); featuredImageURL != "" {
					response.FeaturedImageURL = featuredImageURL
				}
			}
		}
	}

	return response
}

// Helper function to map Article to ArticleListResponse
func (s *Service) mapArticleToListResponse(article *models.Article) dto.ArticleListResponse {
	// Create base response
	response := dto.ArticleListResponse{
		ID:               article.ID,
		Title:            article.Title,
		Slug:             article.Slug,
		Excerpt:          article.Excerpt,
		FeaturedImageURL: article.FeaturedImageURL,
		Status:           article.Status,
		AuthorName:       article.Author.Username, // In a real implementation, get the name from the User model
		ReadTime:         article.ReadTime,
		ViewCount:        article.ViewCount,
		PublishedAt:      article.PublishedAt,
		CreatedAt:        article.CreatedAt,
		Content:          article.Content,
		Categories:       []string{},
		CategoryModels:   []dto.CategoryResponse{},
		Tags:             []string{},
		TagModels:        []dto.TagResponse{},
		Images:           []dto.ArticleImageResponse{},
		Videos:           []dto.ArticleVideoResponse{},
		Metadata:         make(map[string]interface{}),
	}

	// Add category names
	for _, category := range article.Categories {
		response.Categories = append(response.Categories, category.Name)
		response.CategoryModels = append(response.CategoryModels, dto.CategoryResponse{
			ID:   category.ID,
			Name: category.Name,
			Slug: category.Slug,
		})
	}

	// Add tag names
	for _, tag := range article.Tags {
		response.Tags = append(response.Tags, tag.Name)
		response.TagModels = append(response.TagModels, dto.TagResponse{
			ID:   tag.ID,
			Name: tag.Name,
			Slug: tag.Slug,
		})
	}

	// Parse the metadata JSON string and restore rich data
	if article.Metadata != "" {
		var metadata map[string]interface{}
		if err := json.Unmarshal([]byte(article.Metadata), &metadata); err == nil {
			response.Metadata = metadata

			// Restore images from metadata
			if imagesData, ok := metadata["images"].([]interface{}); ok {
				for _, imgData := range imagesData {
					if imgMap, ok := imgData.(map[string]interface{}); ok {
						image := dto.ArticleImageResponse{
							ID:        getString(imgMap, "id"),
							URL:       getString(imgMap, "url"),
							Caption:   getString(imgMap, "caption"),
							AltText:   getString(imgMap, "altText"),
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
						video := dto.ArticleVideoResponse{
							ID:        getString(vidMap, "id"),
							URL:       getString(vidMap, "url"),
							Caption:   getString(vidMap, "caption"),
							SortOrder: getInt(vidMap, "sortOrder"),
						}
						response.Videos = append(response.Videos, video)
					}
				}
			}

			// Restore featured image URL from metadata if not set
			if response.FeaturedImageURL == "" {
				if featuredImageURL := getString(metadata, "featuredImageUrl"); featuredImageURL != "" {
					response.FeaturedImageURL = featuredImageURL
				}
			}
		}
	}

	return response
}
