package websocket

import (
	"web-porto-backend/internal/adapters/websocket"
	"web-porto-backend/middleware"

	"github.com/gin-gonic/gin"
)

// SetupWebSocketRoutes configures WebSocket routes
func SetupWebSocketRoutes(router *gin.Engine, wsManager *websocket.Manager) {
	// WebSocket routes group with special middleware
	ws := router.Group("/ws")
	ws.Use(middleware.WebSocketMiddleware())
	{
		// WebSocket endpoint for real-time analytics
		ws.GET("/analytics", wsManager.ServeWs)
	}
}
