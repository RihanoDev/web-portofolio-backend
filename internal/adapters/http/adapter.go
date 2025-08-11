package http

import (
	"net/http"
	"strconv"
	"web-porto-backend/common/response"
	"web-porto-backend/common/utils"

	"github.com/gin-gonic/gin"
)

// HTTPAdapter handles HTTP-specific concerns
type HTTPAdapter struct{}

// NewHTTPAdapter creates a new HTTP adapter
func NewHTTPAdapter() *HTTPAdapter {
	return &HTTPAdapter{}
}

// SendSuccessResponse sends a successful response
func (h *HTTPAdapter) SendSuccessResponse(c *gin.Context, statusCode int, data interface{}, message string) {
	resp := response.NewSuccessResponse(data, message)
	c.JSON(statusCode, resp)
}

// SendErrorResponse sends an error response
func (h *HTTPAdapter) SendErrorResponse(c *gin.Context, statusCode int, error string) {
	resp := response.NewErrorResponse(error)
	c.JSON(statusCode, resp)
}

// SendPaginatedResponse sends a paginated response
func (h *HTTPAdapter) SendPaginatedResponse(c *gin.Context, data interface{}, page, limit int, total int64, message string) {
	resp := response.NewPaginatedResponse(data, page, limit, total, message)
	c.JSON(http.StatusOK, resp)
}

// SendValidationErrorResponse sends a validation error response
func (h *HTTPAdapter) SendValidationErrorResponse(c *gin.Context, message string) {
	resp := response.NewValidationErrorResponse(message)
	c.JSON(http.StatusBadRequest, resp)
}

// GetPaginationFromQuery extracts pagination parameters from query string
func (h *HTTPAdapter) GetPaginationFromQuery(c *gin.Context) utils.Pagination {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	return utils.NewPagination(page, limit)
}

// ParseIDParam parses ID parameter from URL
func (h *HTTPAdapter) ParseIDParam(c *gin.Context, paramName string) (uint, error) {
	return utils.ParseID(c.Param(paramName))
}

// ParseIntIDParam parses int ID parameter from URL
func (h *HTTPAdapter) ParseIntIDParam(c *gin.Context, paramName string) (int, error) {
	return utils.ParseIntID(c.Param(paramName))
}

// BindJSON binds JSON request to struct
func (h *HTTPAdapter) BindJSON(c *gin.Context, obj interface{}) error {
	return c.ShouldBindJSON(obj)
}

// GetUserContext extracts user context from request
func (h *HTTPAdapter) GetUserContext(c *gin.Context) (userID interface{}, username string, role string) {
	userID, _ = c.Get("user_id")
	username = c.GetString("username")
	role = c.GetString("role")
	return
}
