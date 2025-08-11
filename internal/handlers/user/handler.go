package user

import (
"net/http"
"web-porto-backend/common/response"
httpAdapter "web-porto-backend/internal/adapters/http"
"web-porto-backend/internal/domain/models"
"web-porto-backend/internal/services/user"

"github.com/gin-gonic/gin"
)

type Handler struct {
service     user.Service
httpAdapter *httpAdapter.HTTPAdapter
}

func NewHandler(service user.Service, httpAdapter *httpAdapter.HTTPAdapter) *Handler {
return &Handler{
service:     service,
httpAdapter: httpAdapter,
}
}

type CreateUserRequest struct {
Username string `json:"username" binding:"required"`
Email    string `json:"email" binding:"required,email"`
Password string `json:"password" binding:"required,min=6"`
Role     string `json:"role" binding:"oneof=admin user"`
}

type UpdateUserRequest struct {
Username string `json:"username"`
Email    string `json:"email" binding:"omitempty,email"`
Role     string `json:"role" binding:"omitempty,oneof=admin user"`
}

func (h *Handler) CreateUser(c *gin.Context) {
var req CreateUserRequest
if err := c.ShouldBindJSON(&req); err != nil {
c.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request data", err.Error()))
return
}

user := &models.User{
Username:     req.Username,
Email:        req.Email,
PasswordHash: req.Password, // Service will hash this
Role:         req.Role,
}

if err := h.service.Create(user); err != nil {
c.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to create user", err.Error()))
return
}

// Don't return password in response
user.PasswordHash = ""
h.httpAdapter.SendSuccessResponse(c, http.StatusCreated, user, "User created successfully")
}

func (h *Handler) GetUser(c *gin.Context) {
id, err := h.httpAdapter.ParseIDParam(c, "id")
if err != nil {
c.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid user ID", err.Error()))
return
}

user, err := h.service.GetByID(id)
if err != nil {
c.JSON(http.StatusNotFound, response.NewErrorResponse("User not found", err.Error()))
return
}

// Don't return password
user.PasswordHash = ""
h.httpAdapter.SendSuccessResponse(c, http.StatusOK, user, "User retrieved successfully")
}

func (h *Handler) GetUsers(c *gin.Context) {
pagination := h.httpAdapter.GetPaginationFromQuery(c)

users, paginationInfo, err := h.service.GetAll(pagination.Page, pagination.Limit)
if err != nil {
c.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to get users", err.Error()))
return
}

// Remove passwords from response
for _, user := range users {
user.PasswordHash = ""
}

responseData := response.NewPaginatedResponse(users, pagination.Page, pagination.Limit, paginationInfo.Total, "Users retrieved successfully")
c.JSON(http.StatusOK, responseData)
}

func (h *Handler) UpdateUser(c *gin.Context) {
id, err := h.httpAdapter.ParseIDParam(c, "id")
if err != nil {
c.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid user ID", err.Error()))
return
}

var req UpdateUserRequest
if err := c.ShouldBindJSON(&req); err != nil {
c.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request data", err.Error()))
return
}

user := &models.User{
Username: req.Username,
Email:    req.Email,
Role:     req.Role,
}

if err := h.service.Update(id, user); err != nil {
c.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to update user", err.Error()))
return
}

h.httpAdapter.SendSuccessResponse(c, http.StatusOK, nil, "User updated successfully")
}

func (h *Handler) DeleteUser(c *gin.Context) {
id, err := h.httpAdapter.ParseIDParam(c, "id")
if err != nil {
c.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid user ID", err.Error()))
return
}

if err := h.service.Delete(id); err != nil {
c.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to delete user", err.Error()))
return
}

h.httpAdapter.SendSuccessResponse(c, http.StatusOK, nil, "User deleted successfully")
}

func (h *Handler) GetUserByEmail(c *gin.Context) {
email := c.Query("email")
if email == "" {
c.JSON(http.StatusBadRequest, response.NewErrorResponse("Email parameter is required"))
return
}

user, err := h.service.GetByEmail(email)
if err != nil {
c.JSON(http.StatusNotFound, response.NewErrorResponse("User not found", err.Error()))
return
}

// Don't return password
user.PasswordHash = ""
h.httpAdapter.SendSuccessResponse(c, http.StatusOK, user, "User retrieved successfully")
}
