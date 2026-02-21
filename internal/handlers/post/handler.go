package post

import (
	"net/http"
	"web-porto-backend/common/response"
	"web-porto-backend/common/utils"
	httpAdapter "web-porto-backend/internal/adapters/http"
	"web-porto-backend/internal/domain/models"
	"web-porto-backend/internal/services/post"

	"github.com/gin-gonic/gin"
)

// Constants for messages
const (
	// Error messages
	msgInvalidPostID      = "Invalid post ID"
	msgPostNotFound       = "Post not found"
	msgInvalidRequestData = "Invalid request data"
	msgEmptyIDParameter   = "ID parameter is empty"

	// Success messages
	msgPostCreated    = "Post created successfully"
	msgPostUpdated    = "Post updated successfully"
	msgPostDeleted    = "Post deleted successfully"
	msgPostRetrieved  = "Post retrieved successfully"
	msgPostsRetrieved = "Posts retrieved successfully"
)

type Handler struct {
	service     post.Service
	httpAdapter *httpAdapter.HTTPAdapter
}

func NewHandler(service post.Service, httpAdapter *httpAdapter.HTTPAdapter) *Handler {
	return &Handler{
		service:     service,
		httpAdapter: httpAdapter,
	}
}

type CreatePostRequest struct {
	Title    string `json:"title" binding:"required"`
	Content  string `json:"content" binding:"required"`
	Status   string `json:"status" binding:"required,oneof=draft published"`
	AuthorID int    `json:"author_id" binding:"required"`
}

type UpdatePostRequest struct {
	Title    string `json:"title"`
	Content  string `json:"content"`
	Status   string `json:"status" binding:"oneof=draft published"`
	AuthorID int    `json:"author_id"`
}

func (h *Handler) Create(c *gin.Context) {
	var req CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request data", err.Error()))
		return
	}

	post := &models.Post{
		Title:    req.Title,
		Content:  req.Content,
		Status:   req.Status,
		AuthorID: req.AuthorID,
	}

	if err := h.service.Create(post); err != nil {
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to create post", err.Error()))
		return
	}

	h.httpAdapter.SendSuccessResponse(c, http.StatusCreated, post, msgPostCreated)
}

func (h *Handler) GetAll(c *gin.Context) {
	pagination := h.httpAdapter.GetPaginationFromQuery(c)

	posts, paginationInfo, err := h.service.GetAll(pagination.Page, pagination.Limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to get posts", err.Error()))
		return
	}

	responseData := response.NewPaginatedResponse(posts, pagination.Page, pagination.Limit, paginationInfo.Total, msgPostsRetrieved)
	c.JSON(http.StatusOK, responseData)
}

func (h *Handler) GetByID(c *gin.Context) {
	// Get ID directly as string since the service expects string
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse(msgInvalidPostID, msgEmptyIDParameter))
		return
	}

	post, err := h.service.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, response.NewErrorResponse(msgPostNotFound, err.Error()))
		return
	}

	h.httpAdapter.SendSuccessResponse(c, http.StatusOK, post, msgPostRetrieved)
}

func (h *Handler) Update(c *gin.Context) {
	// Get ID directly as string since the service expects string
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse(msgInvalidPostID, msgEmptyIDParameter))
		return
	}

	var req UpdatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse(msgInvalidRequestData, err.Error()))
		return
	}

	post := &models.Post{
		Title:    req.Title,
		Content:  req.Content,
		Status:   req.Status,
		AuthorID: req.AuthorID,
	}

	if err := h.service.Update(id, post); err != nil {
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to update post", err.Error()))
		return
	}

	h.httpAdapter.SendSuccessResponse(c, http.StatusOK, nil, msgPostUpdated)
}

func (h *Handler) Delete(c *gin.Context) {
	// Get ID directly as string since the service expects string
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse(msgInvalidPostID, msgEmptyIDParameter))
		return
	}

	if err := h.service.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to delete post", err.Error()))
		return
	}

	h.httpAdapter.SendSuccessResponse(c, http.StatusOK, nil, msgPostDeleted)
}

func (h *Handler) GetBySlug(c *gin.Context) {
	slug := c.Param("slug")
	if slug == "" {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("Slug is required"))
		return
	}

	post, err := h.service.GetBySlug(slug)
	if err != nil {
		c.JSON(http.StatusNotFound, response.NewErrorResponse(msgPostNotFound, err.Error()))
		return
	}

	h.httpAdapter.SendSuccessResponse(c, http.StatusOK, post, msgPostRetrieved)
}

func (h *Handler) GetByAuthor(c *gin.Context) {
	authorIDStr := c.Param("authorId")
	authorID, err := utils.ParseIntID(authorIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid author ID", err.Error()))
		return
	}

	pagination := h.httpAdapter.GetPaginationFromQuery(c)

	posts, paginationInfo, err := h.service.GetByAuthorID(authorID, pagination.Page, pagination.Limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to get posts", err.Error()))
		return
	}

	responseData := response.NewPaginatedResponse(posts, pagination.Page, pagination.Limit, paginationInfo.Total, "Posts retrieved successfully")
	c.JSON(http.StatusOK, responseData)
}

func (h *Handler) GetPublished(c *gin.Context) {
	pagination := h.httpAdapter.GetPaginationFromQuery(c)

	posts, paginationInfo, err := h.service.GetPublished(pagination.Page, pagination.Limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to get published posts", err.Error()))
		return
	}

	responseData := response.NewPaginatedResponse(posts, pagination.Page, pagination.Limit, paginationInfo.Total, "Published posts retrieved successfully")
	c.JSON(http.StatusOK, responseData)
}
