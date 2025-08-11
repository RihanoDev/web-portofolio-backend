package category

import (
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
	categories, err := h.service.GetAll()
	if err != nil {
		h.httpAdapter.SendErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
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

	if err := h.service.Create(&category); err != nil {
		h.httpAdapter.SendErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	h.httpAdapter.SendSuccessResponse(c, http.StatusCreated, category, "Category created successfully")
}

func (h *Handler) Update(c *gin.Context) {
	id, err := h.httpAdapter.ParseIntIDParam(c, "id")
	if err != nil {
		h.httpAdapter.SendErrorResponse(c, http.StatusBadRequest, "invalid id")
		return
	}

	var category models.Category
	if err := h.httpAdapter.BindJSON(c, &category); err != nil {
		h.httpAdapter.SendValidationErrorResponse(c, err.Error())
		return
	}

	category.ID = id
	if err := h.service.Update(&category); err != nil {
		h.httpAdapter.SendErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

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
