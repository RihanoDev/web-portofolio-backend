package response

// APIResponse represents a standard API response structure
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// PaginatedResponse represents a paginated API response
type PaginatedResponse struct {
	Success    bool        `json:"success"`
	Message    string      `json:"message,omitempty"`
	Data       interface{} `json:"data,omitempty"`
	Error      string      `json:"error,omitempty"`
	Pagination Pagination  `json:"pagination"`
}

// Pagination metadata for paginated responses
type Pagination struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

// Response constants
const (
	ErrInvalidID      = "invalid id"
	ErrInvalidInput   = "invalid input"
	ErrNotFound       = "resource not found"
	ErrUnauthorized   = "unauthorized"
	ErrForbidden      = "forbidden"
	ErrInternalServer = "internal server error"
	ErrValidation     = "validation failed"
	ErrDuplicateEntry = "duplicate entry"

	MsgSuccess = "success"
	MsgCreated = "resource created successfully"
	MsgUpdated = "resource updated successfully"
	MsgDeleted = "resource deleted successfully"
)

// NewSuccessResponse creates a success response
func NewSuccessResponse(data interface{}, message string) APIResponse {
	if message == "" {
		message = MsgSuccess
	}
	return APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	}
}

// NewErrorResponse creates a new error response
func NewErrorResponse(message string, details ...string) *APIResponse {
	resp := &APIResponse{
		Success: false,
		Message: message,
	}

	// If additional details provided, combine them
	if len(details) > 0 && details[0] != "" {
		resp.Error = details[0]
	}

	return resp
}

// NewPaginatedResponse creates a paginated response
func NewPaginatedResponse(data interface{}, page, limit int, total int64, message string) PaginatedResponse {
	if message == "" {
		message = MsgSuccess
	}

	totalPages := int(total) / limit
	if int(total)%limit != 0 {
		totalPages++
	}

	return PaginatedResponse{
		Success: true,
		Message: message,
		Data:    data,
		Pagination: Pagination{
			Page:       page,
			Limit:      limit,
			Total:      total,
			TotalPages: totalPages,
		},
	}
}

// NewValidationErrorResponse creates a validation error response
func NewValidationErrorResponse(message string) APIResponse {
	if message == "" {
		message = ErrValidation
	}
	return APIResponse{
		Success: false,
		Error:   message,
	}
}
