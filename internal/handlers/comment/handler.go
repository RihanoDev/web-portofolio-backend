package comment

import (
	"net/http"
	httpAdapter "web-porto-backend/internal/adapters/http"
	"web-porto-backend/internal/domain/models"
	"web-porto-backend/internal/services/comment"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service     comment.Service
	httpAdapter *httpAdapter.HTTPAdapter
}

func NewHandler(service comment.Service) *Handler {
	return &Handler{
		service:     service,
		httpAdapter: httpAdapter.NewHTTPAdapter(),
	}
}

func (h *Handler) GetAll(c *gin.Context) {
	comments, err := h.service.GetAll()
	if err != nil {
		h.httpAdapter.SendErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	h.httpAdapter.SendSuccessResponse(c, http.StatusOK, comments, "")
}

func (h *Handler) GetByID(c *gin.Context) {
	id, err := h.httpAdapter.ParseIntIDParam(c, "id")
	if err != nil {
		h.httpAdapter.SendErrorResponse(c, http.StatusBadRequest, "invalid id")
		return
	}

	comment, err := h.service.GetByID(id)
	if err != nil {
		h.httpAdapter.SendErrorResponse(c, http.StatusNotFound, err.Error())
		return
	}

	h.httpAdapter.SendSuccessResponse(c, http.StatusOK, comment, "")
}

func (h *Handler) GetByPost(c *gin.Context) {
	postID, err := h.httpAdapter.ParseIntIDParam(c, "postId")
	if err != nil {
		h.httpAdapter.SendErrorResponse(c, http.StatusBadRequest, "invalid post id")
		return
	}

	comments, err := h.service.GetByPostID(postID)
	if err != nil {
		h.httpAdapter.SendErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	h.httpAdapter.SendSuccessResponse(c, http.StatusOK, comments, "")
}

func (h *Handler) Create(c *gin.Context) {
	var comment models.Comment
	if err := h.httpAdapter.BindJSON(c, &comment); err != nil {
		h.httpAdapter.SendValidationErrorResponse(c, err.Error())
		return
	}

	if err := h.service.Create(&comment); err != nil {
		h.httpAdapter.SendErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	h.httpAdapter.SendSuccessResponse(c, http.StatusCreated, comment, "Comment created successfully")
}

func (h *Handler) Update(c *gin.Context) {
	id, err := h.httpAdapter.ParseIntIDParam(c, "id")
	if err != nil {
		h.httpAdapter.SendErrorResponse(c, http.StatusBadRequest, "invalid id")
		return
	}

	var comment models.Comment
	if err := h.httpAdapter.BindJSON(c, &comment); err != nil {
		h.httpAdapter.SendValidationErrorResponse(c, err.Error())
		return
	}

	comment.ID = id
	if err := h.service.Update(&comment); err != nil {
		h.httpAdapter.SendErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	h.httpAdapter.SendSuccessResponse(c, http.StatusOK, comment, "Comment updated successfully")
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

	h.httpAdapter.SendSuccessResponse(c, http.StatusOK, nil, "Comment deleted successfully")
}
