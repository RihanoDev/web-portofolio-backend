package routes

import (
	"web-porto-backend/internal/adapters/websocket"
	"web-porto-backend/internal/auth"
	"web-porto-backend/internal/handlers"
	"web-porto-backend/routes/api"
	ws "web-porto-backend/routes/websocket"

	"github.com/gin-gonic/gin"
)

// SetupRoutes configures all application routes
func SetupRoutes(router *gin.Engine, handlerRegistry *handlers.HandlerRegistry, authService *auth.AuthService) {
	// Setup API routes using the API package
	api.SetupAPIRoutes(router, handlerRegistry, authService)
}

// SetupWebSocketRoutes configures all WebSocket routes
func SetupWebSocketRoutes(router *gin.Engine, wsManager *websocket.Manager) {
	// Setup WebSocket routes
	ws.SetupWebSocketRoutes(router, wsManager)
}
