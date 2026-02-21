package tag

import (
	"net/http"
	"strconv"
	"web-porto-backend/common/response"

	"github.com/gin-gonic/gin"
)

// GetByID gets a tag by ID
func (h *Handler) GetByID(c *gin.Context) {
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

	tag, err := h.tagService.GetByID(tagID)
	if err != nil {
		c.JSON(http.StatusNotFound, response.NewErrorResponse("Tag not found", err.Error()))
		return
	}

	h.httpAdapter.SendSuccessResponse(c, http.StatusOK, tag, "Tag retrieved successfully")
}

// GetByName gets a tag by name
func (h *Handler) GetByName(c *gin.Context) {
	name := c.Param("name")
	if name == "" {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("Name is required"))
		return
	}

	tag, err := h.tagService.GetByName(name)
	if err != nil {
		c.JSON(http.StatusNotFound, response.NewErrorResponse("Tag not found", err.Error()))
		return
	}

	h.httpAdapter.SendSuccessResponse(c, http.StatusOK, tag, "Tag retrieved successfully")
}
