package routes

import (
	"web-porto-backend/internal/auth"
	"web-porto-backend/internal/handlers"
	"web-porto-backend/routes/api"

	"github.com/gin-gonic/gin"
)

// SetupRoutes configures all application routes
func SetupRoutes(router *gin.Engine, handlerRegistry *handlers.HandlerRegistry, authService *auth.AuthService) {
	// Setup API routes using the API package
	api.SetupAPIRoutes(router, handlerRegistry, authService)
}
