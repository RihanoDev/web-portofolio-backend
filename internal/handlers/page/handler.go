package page

import (
	"net/http"
	"web-porto-backend/common/response"
	httpAdapter "web-porto-backend/internal/adapters/http"
	"web-porto-backend/internal/domain/models"
	"web-porto-backend/internal/services/page"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service     page.Service
	httpAdapter *httpAdapter.HTTPAdapter
}

func NewHandler(service page.Service, httpAdapter *httpAdapter.HTTPAdapter) *Handler {
	return &Handler{
		service:     service,
		httpAdapter: httpAdapter,
	}
}

type CreatePageRequest struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
	Status  string `json:"status" binding:"required,oneof=draft published"`
}

type UpdatePageRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	Status  string `json:"status" binding:"oneof=draft published"`
}

func (h *Handler) Create(c *gin.Context) {
	var req CreatePageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request data", err.Error()))
		return
	}

	page := &models.Page{
		Title:   req.Title,
		Content: req.Content,
		Status:  req.Status,
	}

	if err := h.service.Create(page); err != nil {
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to create page", err.Error()))
		return
	}

	h.httpAdapter.SendSuccessResponse(c, http.StatusCreated, page, "Page created successfully")
}

func (h *Handler) GetAll(c *gin.Context) {
	pagination := h.httpAdapter.GetPaginationFromQuery(c)

	pages, paginationInfo, err := h.service.GetAll(pagination.Page, pagination.Limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to get pages", err.Error()))
		return
	}

	responseData := response.NewPaginatedResponse(pages, pagination.Page, pagination.Limit, paginationInfo.Total, "Pages retrieved successfully")
	c.JSON(http.StatusOK, responseData)
}

func (h *Handler) GetByID(c *gin.Context) {
	id, err := h.httpAdapter.ParseIDParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid page ID", err.Error()))
		return
	}

	page, err := h.service.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, response.NewErrorResponse("Page not found", err.Error()))
		return
	}

	h.httpAdapter.SendSuccessResponse(c, http.StatusOK, page, "Page retrieved successfully")
}

func (h *Handler) Update(c *gin.Context) {
	id, err := h.httpAdapter.ParseIDParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid page ID", err.Error()))
		return
	}

	var req UpdatePageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request data", err.Error()))
		return
	}

	page := &models.Page{
		Title:   req.Title,
		Content: req.Content,
		Status:  req.Status,
	}

	if err := h.service.Update(id, page); err != nil {
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to update page", err.Error()))
		return
	}

	h.httpAdapter.SendSuccessResponse(c, http.StatusOK, nil, "Page updated successfully")
}

func (h *Handler) Delete(c *gin.Context) {
	id, err := h.httpAdapter.ParseIDParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid page ID", err.Error()))
		return
	}

	if err := h.service.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to delete page", err.Error()))
		return
	}

	h.httpAdapter.SendSuccessResponse(c, http.StatusOK, nil, "Page deleted successfully")
}

func (h *Handler) GetBySlug(c *gin.Context) {
	slug := c.Param("slug")
	if slug == "" {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("Slug is required"))
		return
	}

	page, err := h.service.GetBySlug(slug)
	if err != nil {
		c.JSON(http.StatusNotFound, response.NewErrorResponse("Page not found", err.Error()))
		return
	}

	h.httpAdapter.SendSuccessResponse(c, http.StatusOK, page, "Page retrieved successfully")
}

func (h *Handler) GetPublished(c *gin.Context) {
	pagination := h.httpAdapter.GetPaginationFromQuery(c)

	pages, paginationInfo, err := h.service.GetPublished(pagination.Page, pagination.Limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to get published pages", err.Error()))
		return
	}

	responseData := response.NewPaginatedResponse(pages, pagination.Page, pagination.Limit, paginationInfo.Total, "Published pages retrieved successfully")
	c.JSON(http.StatusOK, responseData)
}
