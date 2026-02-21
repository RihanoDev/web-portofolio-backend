package tag

import (
	"net/http"
	"strconv"
	"web-porto-backend/common/response"
	httpAdapter "web-porto-backend/internal/adapters/http"
	"web-porto-backend/internal/domain/dto"

	"github.com/gin-gonic/gin"
)

// Handler handles HTTP requests for tag endpoints
type Handler struct {
	tagService interface {
		GetAll() ([]dto.TagResponse, error)
		GetByID(id int) (*dto.TagResponse, error)
		GetByName(name string) (*dto.TagResponse, error)
		Create(tag *dto.CreateTagRequest) (*dto.TagResponse, error)
		Update(id int, tag *dto.UpdateTagRequest) (*dto.TagResponse, error)
		Delete(id int) error
	}
	httpAdapter *httpAdapter.HTTPAdapter
}

// NewHandler creates a new tag handler
func NewHandler(tagService interface {
	GetAll() ([]dto.TagResponse, error)
	GetByID(id int) (*dto.TagResponse, error)
	GetByName(name string) (*dto.TagResponse, error)
	Create(tag *dto.CreateTagRequest) (*dto.TagResponse, error)
	Update(id int, tag *dto.UpdateTagRequest) (*dto.TagResponse, error)
	Delete(id int) error
}, httpAdapter *httpAdapter.HTTPAdapter) *Handler {
	return &Handler{
		tagService:  tagService,
		httpAdapter: httpAdapter,
	}
}

// GetAll gets all tags
func (h *Handler) GetAll(c *gin.Context) {
	tags, err := h.tagService.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to retrieve tags", err.Error()))
		return
	}

	h.httpAdapter.SendSuccessResponse(c, http.StatusOK, tags, "Tags retrieved successfully")
}

// Create creates a new tag
func (h *Handler) Create(c *gin.Context) {
	var tagRequest dto.CreateTagRequest

	if err := c.ShouldBindJSON(&tagRequest); err != nil {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid tag data", err.Error()))
		return
	}

	tag, err := h.tagService.Create(&tagRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to create tag", err.Error()))
		return
	}

	h.httpAdapter.SendSuccessResponse(c, http.StatusCreated, tag, "Tag created successfully")
}

// Update updates a tag
func (h *Handler) Update(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("ID is required"))
		return
	}

	// Parse ID to int
	tagID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid tag ID", err.Error()))
		return
	}

	var tagRequest dto.UpdateTagRequest

	if err := c.ShouldBindJSON(&tagRequest); err != nil {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid tag data", err.Error()))
		return
	}

	tag, err := h.tagService.Update(tagID, &tagRequest)
	if err != nil {
		c.JSON(http.StatusNotFound, response.NewErrorResponse("Failed to update tag", err.Error()))
		return
	}

	h.httpAdapter.SendSuccessResponse(c, http.StatusOK, tag, "Tag updated successfully")
}

// Delete deletes a tag
func (h *Handler) Delete(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("ID is required"))
		return
	}

	// Parse ID to int
	tagID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid tag ID", err.Error()))
		return
	}

	if err := h.tagService.Delete(tagID); err != nil {
		c.JSON(http.StatusNotFound, response.NewErrorResponse("Failed to delete tag", err.Error()))
		return
	}

	h.httpAdapter.SendSuccessResponse(c, http.StatusOK, nil, "Tag deleted successfully")
}
