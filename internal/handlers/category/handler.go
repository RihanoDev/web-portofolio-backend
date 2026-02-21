package category

import (
	"log"
	"net/http"
	httpAdapter "web-porto-backend/internal/adapters/http"
	"web-porto-backend/internal/domain/models"
	"web-porto-backend/internal/services/category"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service     category.Service
	httpAdapter *httpAdapter.HTTPAdapter
}

func NewHandler(service category.Service) *Handler {
	return &Handler{
		service:     service,
		httpAdapter: httpAdapter.NewHTTPAdapter(),
	}
}

func (h *Handler) GetAll(c *gin.Context) {
	// Log the request for debugging
	log.Printf("GetAll categories request received")

	categories, err := h.service.GetAll()
	if err != nil {
		log.Printf("Error fetching categories: %v", err)
		h.httpAdapter.SendErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	// Log the response for debugging
	log.Printf("Categories fetched successfully, count: %d", len(categories))

	h.httpAdapter.SendSuccessResponse(c, http.StatusOK, categories, "")
}

func (h *Handler) GetByID(c *gin.Context) {
	id, err := h.httpAdapter.ParseIntIDParam(c, "id")
	if err != nil {
		h.httpAdapter.SendErrorResponse(c, http.StatusBadRequest, "invalid id")
		return
	}

	category, err := h.service.GetByID(id)
	if err != nil {
		h.httpAdapter.SendErrorResponse(c, http.StatusNotFound, err.Error())
		return
	}

	h.httpAdapter.SendSuccessResponse(c, http.StatusOK, category, "")
}

func (h *Handler) Create(c *gin.Context) {
	var category models.Category
	if err := h.httpAdapter.BindJSON(c, &category); err != nil {
		h.httpAdapter.SendValidationErrorResponse(c, err.Error())
		return
	}

	// Validate required fields
	if category.Name == "" {
		h.httpAdapter.SendValidationErrorResponse(c, "Name is required")
		return
	}

	// Log category data for debugging
	log.Printf("Creating category: %+v", category)

	if err := h.service.Create(&category); err != nil {
		log.Printf("Error creating category: %v", err)
		h.httpAdapter.SendErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	log.Printf("Category created successfully: %+v", category)
	h.httpAdapter.SendSuccessResponse(c, http.StatusCreated, category, "Category created successfully")
}

func (h *Handler) Update(c *gin.Context) {
	id, err := h.httpAdapter.ParseIntIDParam(c, "id")
	if err != nil {
		log.Printf("Invalid ID parameter: %v", err)
		h.httpAdapter.SendErrorResponse(c, http.StatusBadRequest, "invalid id")
		return
	}

	var category models.Category
	if err := h.httpAdapter.BindJSON(c, &category); err != nil {
		h.httpAdapter.SendValidationErrorResponse(c, err.Error())
		return
	}

	// Validate required fields
	if category.Name == "" {
		h.httpAdapter.SendValidationErrorResponse(c, "Name is required")
		return
	}

	// Log category data for debugging
	log.Printf("Updating category %d: %+v", id, category)

	category.ID = id
	if err := h.service.Update(&category); err != nil {
		log.Printf("Error updating category: %v", err)
		h.httpAdapter.SendErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	log.Printf("Category updated successfully: %+v", category)
	h.httpAdapter.SendSuccessResponse(c, http.StatusOK, category, "Category updated successfully")
}

func (h *Handler) Delete(c *gin.Context) {
	id, err := h.httpAdapter.ParseIntIDParam(c, "id")
	if err != nil {
		h.httpAdapter.SendErrorResponse(c, http.StatusBadRequest, "invalid id")
		return
	}

	if err := h.service.Delete(id); err != nil {
		h.httpAdapter.SendErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	h.httpAdapter.SendSuccessResponse(c, http.StatusOK, nil, "Category deleted successfully")
}
