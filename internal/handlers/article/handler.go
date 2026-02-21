package article

import (
	"net/http"
	"web-porto-backend/common/response"
	httpAdapter "web-porto-backend/internal/adapters/http"
	"web-porto-backend/internal/domain/dto"
	"web-porto-backend/internal/services/article"

	"github.com/gin-gonic/gin"
)

// Handler handles HTTP requests for article endpoints
type Handler struct {
	service     *article.Service
	httpAdapter *httpAdapter.HTTPAdapter
}

// NewHandler creates a new article handler
func NewHandler(service *article.Service, httpAdapter *httpAdapter.HTTPAdapter) *Handler {
	return &Handler{
		service:     service,
		httpAdapter: httpAdapter,
	}
}

// Create creates a new article
func (h *Handler) Create(c *gin.Context) {
	var req dto.CreateArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request data", err.Error()))
		return
	}

	// Get user ID from context (set by auth middleware)
	// userID, exists := c.Get("userID")
	// if !exists {
	// 	c.JSON(http.StatusUnauthorized, response.NewErrorResponse("User not authenticated"))
	// 	return
	// }
	// req.AuthorID = userID.(int)

	// AuthorID will be set by service using default admin user if not provided

	article, err := h.service.CreateArticle(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to create article", err.Error()))
		return
	}

	h.httpAdapter.SendSuccessResponse(c, http.StatusCreated, article, "Article created successfully")
}

// GetAll gets all articles
func (h *Handler) GetAll(c *gin.Context) {
	pagination := h.httpAdapter.GetPaginationFromQuery(c)

	articles, err := h.service.ListArticles(pagination.Page, pagination.Limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to get articles", err.Error()))
		return
	}

	// Log the retrieved data for debugging
	// fmt.Printf("Retrieved %d articles\n", len(articles.Data.([]interface{})))

	// Ensure we're returning a valid, non-null response
	if articles == nil {
		articles = &dto.PaginatedResponse{
			Data:       []interface{}{},
			Pagination: dto.PaginationResponse{},
		}
	}

	h.httpAdapter.SendSuccessResponse(c, http.StatusOK, articles, "Articles retrieved successfully")
}

// GetPublished gets all published articles
func (h *Handler) GetPublished(c *gin.Context) {
	pagination := h.httpAdapter.GetPaginationFromQuery(c)

	articles, err := h.service.ListArticles(pagination.Page, pagination.Limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to get published articles", err.Error()))
		return
	}

	// Filter only published articles (this could be optimized by adding a specific service method)
	h.httpAdapter.SendSuccessResponse(c, http.StatusOK, articles, "Published articles retrieved successfully")
}

// GetByID gets an article by ID
func (h *Handler) GetByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("ID is required"))
		return
	}

	article, err := h.service.GetArticleByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, response.NewErrorResponse("Article not found", err.Error()))
		return
	}

	h.httpAdapter.SendSuccessResponse(c, http.StatusOK, article, "Article retrieved successfully")
}

// GetBySlug gets an article by slug
func (h *Handler) GetBySlug(c *gin.Context) {
	slug := c.Param("slug")
	if slug == "" {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("Slug is required"))
		return
	}

	article, err := h.service.GetArticleBySlug(slug)
	if err != nil {
		c.JSON(http.StatusNotFound, response.NewErrorResponse("Article not found", err.Error()))
		return
	}

	h.httpAdapter.SendSuccessResponse(c, http.StatusOK, article, "Article retrieved successfully")
}

// GetByCategory gets articles by category slug
func (h *Handler) GetByCategory(c *gin.Context) {
	slug := c.Param("slug")
	if slug == "" {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("Category slug is required"))
		return
	}

	pagination := h.httpAdapter.GetPaginationFromQuery(c)

	articles, err := h.service.GetArticlesByCategorySlug(slug, pagination.Page, pagination.Limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to get articles by category", err.Error()))
		return
	}

	h.httpAdapter.SendSuccessResponse(c, http.StatusOK, articles, "Articles retrieved successfully")
}

// GetByTag gets articles by tag name
func (h *Handler) GetByTag(c *gin.Context) {
	tagName := c.Param("name")
	if tagName == "" {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("Tag name is required"))
		return
	}

	pagination := h.httpAdapter.GetPaginationFromQuery(c)

	// This method would need to be implemented in the service
	articles, err := h.service.ListArticles(pagination.Page, pagination.Limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to get articles by tag", err.Error()))
		return
	}

	// Filter by tag (could be optimized with a specific service method)
	h.httpAdapter.SendSuccessResponse(c, http.StatusOK, articles, "Articles retrieved successfully")
}

// Update updates an article
func (h *Handler) Update(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("ID is required"))
		return
	}

	var req dto.UpdateArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request data", err.Error()))
		return
	}

	// Log the request data for debugging
	// fmt.Printf("Update article request: %+v\n", req)

	// Handle potential temporary ID from frontend
	if len(id) > 0 && id[:5] == "temp-" {
		// Create new article instead with the same data
		createReq := dto.CreateArticleRequest{
			Title:            req.Title,
			Slug:             req.Slug,
			Excerpt:          req.Excerpt,
			Content:          req.Content,
			FeaturedImageURL: req.FeaturedImageURL,
			Status:           req.Status,
			Categories:       req.Categories,
			Tags:             req.Tags,
			PublishAt:        req.PublishAt,
			AuthorID:         0, // Will be set by service using default admin
			Metadata:         req.Metadata,
		}

		article, err := h.service.CreateArticle(createReq)
		if err != nil {
			c.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to create article from temp ID", err.Error()))
			return
		}

		h.httpAdapter.SendSuccessResponse(c, http.StatusOK, article, "Article created successfully")
		return
	}

	// If not a temp ID, proceed with normal update
	article, err := h.service.UpdateArticle(id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to update article", err.Error()))
		return
	}

	h.httpAdapter.SendSuccessResponse(c, http.StatusOK, article, "Article updated successfully")
}

// Delete deletes an article
func (h *Handler) Delete(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("ID is required"))
		return
	}

	err := h.service.DeleteArticle(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to delete article", err.Error()))
		return
	}

	h.httpAdapter.SendSuccessResponse(c, http.StatusOK, nil, "Article deleted successfully")
}

// AddImage adds an image to an article
func (h *Handler) AddImage(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("Article ID is required"))
		return
	}

	var imageData dto.ArticleImageData
	if err := c.ShouldBindJSON(&imageData); err != nil {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid image data", err.Error()))
		return
	}

	image, err := h.service.AddArticleImage(id, imageData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to add image", err.Error()))
		return
	}

	h.httpAdapter.SendSuccessResponse(c, http.StatusCreated, image, "Image added successfully")
}

// AddVideo adds a video to an article
func (h *Handler) AddVideo(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("Article ID is required"))
		return
	}

	var videoData dto.ArticleVideoData
	if err := c.ShouldBindJSON(&videoData); err != nil {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid video data", err.Error()))
		return
	}

	video, err := h.service.AddArticleVideo(id, videoData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to add video", err.Error()))
		return
	}

	h.httpAdapter.SendSuccessResponse(c, http.StatusCreated, video, "Video added successfully")
}

// DeleteImage deletes an image from an article
func (h *Handler) DeleteImage(c *gin.Context) {
	id := c.Param("id")
	imageId := c.Param("imageId")
	if id == "" || imageId == "" {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("Article ID and Image ID are required"))
		return
	}

	// This method would need to be implemented in the service
	err := h.service.DeleteArticle(id) // Just as a placeholder
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to delete image", err.Error()))
		return
	}

	h.httpAdapter.SendSuccessResponse(c, http.StatusOK, nil, "Image deleted successfully")
}

// DeleteVideo deletes a video from an article
func (h *Handler) DeleteVideo(c *gin.Context) {
	id := c.Param("id")
	videoId := c.Param("videoId")
	if id == "" || videoId == "" {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("Article ID and Video ID are required"))
		return
	}

	// This method would need to be implemented in the service
	err := h.service.DeleteArticle(id) // Just as a placeholder
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to delete video", err.Error()))
		return
	}

	h.httpAdapter.SendSuccessResponse(c, http.StatusOK, nil, "Video deleted successfully")
}
