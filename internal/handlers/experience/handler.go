package experience

import (
	"net/http"

	httpAdapter "web-porto-backend/internal/adapters/http"
	"web-porto-backend/internal/domain/dto"
	experienceService "web-porto-backend/internal/services/experience"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service     *experienceService.Service
	httpAdapter *httpAdapter.HTTPAdapter
}

func NewHandler(service *experienceService.Service, httpAdapter *httpAdapter.HTTPAdapter) *Handler {
	return &Handler{
		service:     service,
		httpAdapter: httpAdapter,
	}
}

// GetAll retrieves all experiences
func (h *Handler) GetAll(c *gin.Context) {
	// Get pagination parameters
	pagination := h.httpAdapter.GetPaginationFromQuery(c)

	result, err := h.service.ListExperiences(pagination.Page, pagination.Limit)
	if err != nil {
		h.httpAdapter.SendErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve experiences")
		return
	}

	h.httpAdapter.SendSuccessResponse(c, http.StatusOK, result, "Experiences retrieved successfully")
}

// GetByID retrieves an experience by ID
func (h *Handler) GetByID(c *gin.Context) {
	id, err := h.httpAdapter.ParseIntIDParam(c, "id")
	if err != nil {
		h.httpAdapter.SendErrorResponse(c, http.StatusBadRequest, "Invalid experience ID")
		return
	}

	experience, err := h.service.GetExperienceByID(id)
	if err != nil {
		h.httpAdapter.SendErrorResponse(c, http.StatusNotFound, "Experience not found")
		return
	}

	h.httpAdapter.SendSuccessResponse(c, http.StatusOK, experience, "Experience retrieved successfully")
}

// Create creates a new experience
func (h *Handler) Create(c *gin.Context) {
	var req dto.CreateExperienceRequest
	if err := h.httpAdapter.BindJSON(c, &req); err != nil {
		h.httpAdapter.SendValidationErrorResponse(c, "Invalid request body")
		return
	}

	experience, err := h.service.CreateExperience(req)
	if err != nil {
		h.httpAdapter.SendErrorResponse(c, http.StatusInternalServerError, "Failed to create experience")
		return
	}

	h.httpAdapter.SendSuccessResponse(c, http.StatusCreated, experience, "Experience created successfully")
}

// Update updates an existing experience
func (h *Handler) Update(c *gin.Context) {
	id, err := h.httpAdapter.ParseIntIDParam(c, "id")
	if err != nil {
		h.httpAdapter.SendErrorResponse(c, http.StatusBadRequest, "Invalid experience ID")
		return
	}

	var req dto.UpdateExperienceRequest
	if err := h.httpAdapter.BindJSON(c, &req); err != nil {
		h.httpAdapter.SendValidationErrorResponse(c, "Invalid request body")
		return
	}

	experience, err := h.service.UpdateExperience(id, req)
	if err != nil {
		h.httpAdapter.SendErrorResponse(c, http.StatusInternalServerError, "Failed to update experience")
		return
	}

	h.httpAdapter.SendSuccessResponse(c, http.StatusOK, experience, "Experience updated successfully")
}

// Delete deletes an experience by ID
func (h *Handler) Delete(c *gin.Context) {
	id, err := h.httpAdapter.ParseIntIDParam(c, "id")
	if err != nil {
		h.httpAdapter.SendErrorResponse(c, http.StatusBadRequest, "Invalid experience ID")
		return
	}

	err = h.service.DeleteExperience(id)
	if err != nil {
		h.httpAdapter.SendErrorResponse(c, http.StatusInternalServerError, "Failed to delete experience")
		return
	}

	h.httpAdapter.SendSuccessResponse(c, http.StatusOK, nil, "Experience deleted successfully")
}

// GetCurrent retrieves current active experiences
func (h *Handler) GetCurrent(c *gin.Context) {
	experiences, err := h.service.GetCurrentExperiences()
	if err != nil {
		h.httpAdapter.SendErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve current experiences")
		return
	}

	h.httpAdapter.SendSuccessResponse(c, http.StatusOK, experiences, "Current experiences retrieved successfully")
}
