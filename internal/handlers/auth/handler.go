package auth

import (
	"net/http"
	"web-porto-backend/common/response"
	httpAdapter "web-porto-backend/internal/adapters/http"
	"web-porto-backend/internal/auth"
	"web-porto-backend/internal/domain/models"
	"web-porto-backend/internal/services/user"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	userService user.Service
	jwtService  auth.JWTService
	httpAdapter *httpAdapter.HTTPAdapter
}

func NewHandler(userService user.Service, jwtService auth.JWTService, httpAdapter *httpAdapter.HTTPAdapter) *Handler {
	return &Handler{
		userService: userService,
		jwtService:  jwtService,
		httpAdapter: httpAdapter,
	}
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginResponse struct {
	Token string      `json:"token"`
	User  models.User `json:"user"`
}

func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request data", err.Error()))
		return
	}

	user, err := h.userService.GetByEmail(req.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, response.NewErrorResponse("Invalid credentials"))
		return
	}

	if !h.userService.CheckPassword(user.PasswordHash, req.Password) {
		c.JSON(http.StatusUnauthorized, response.NewErrorResponse("Invalid credentials"))
		return
	}

	token, err := h.jwtService.GenerateToken(user.ID, user.Email, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to generate token", err.Error()))
		return
	}

	// Don't return password hash
	user.PasswordHash = ""

	loginResponse := LoginResponse{
		Token: token,
		User:  *user,
	}

	h.httpAdapter.SendSuccessResponse(c, http.StatusOK, loginResponse, "Login successful")
}

func (h *Handler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request data", err.Error()))
		return
	}

	// Check if user already exists
	existingUser, _ := h.userService.GetByEmail(req.Email)
	if existingUser != nil {
		c.JSON(http.StatusConflict, response.NewErrorResponse("User already exists"))
		return
	}

	user := &models.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: req.Password, // Service should hash this
		Role:         "user",       // Default role
	}

	if err := h.userService.Create(user); err != nil {
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to create user", err.Error()))
		return
	}

	// Generate token for the new user
	token, err := h.jwtService.GenerateToken(user.ID, user.Email, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to generate token", err.Error()))
		return
	}

	// Don't return password hash
	user.PasswordHash = ""

	registerResponse := LoginResponse{
		Token: token,
		User:  *user,
	}

	h.httpAdapter.SendSuccessResponse(c, http.StatusCreated, registerResponse, "Registration successful")
}

func (h *Handler) RefreshToken(c *gin.Context) {
	// Get user from context (set by auth middleware)
	userInterface, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, response.NewErrorResponse("User not found in context"))
		return
	}

	user, ok := userInterface.(*models.User)
	if !ok {
		c.JSON(http.StatusUnauthorized, response.NewErrorResponse("Invalid user data"))
		return
	}

	token, err := h.jwtService.GenerateToken(user.ID, user.Email, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to generate token", err.Error()))
		return
	}

	tokenResponse := map[string]string{"token": token}
	h.httpAdapter.SendSuccessResponse(c, http.StatusOK, tokenResponse, "Token refreshed successfully")
}

func (h *Handler) Logout(c *gin.Context) {
	// In JWT, logout is typically handled client-side by removing the token
	// For enhanced security, you could implement token blacklisting here
	h.httpAdapter.SendSuccessResponse(c, http.StatusOK, nil, "Logout successful")
}

func (h *Handler) GetProfile(c *gin.Context) {
	// Get user from context (set by auth middleware)
	userInterface, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, response.NewErrorResponse("User not found in context"))
		return
	}

	user, ok := userInterface.(*models.User)
	if !ok {
		c.JSON(http.StatusUnauthorized, response.NewErrorResponse("Invalid user data"))
		return
	}

	// Don't return password hash
	user.PasswordHash = ""
	h.httpAdapter.SendSuccessResponse(c, http.StatusOK, user, "Profile retrieved successfully")
}

func (h *Handler) Me(c *gin.Context) {
	// Alias for GetProfile
	h.GetProfile(c)
}
