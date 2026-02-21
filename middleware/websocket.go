package middleware

import (
	"github.com/gin-gonic/gin"
)

// WebSocketMiddleware handles requests to WebSocket endpoints by adding appropriate headers
func WebSocketMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Allow WebSocket connections from browsers
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization, Sec-WebSocket-Protocol, Sec-WebSocket-Version, Sec-WebSocket-Key, Upgrade, Connection")
		c.Header("Access-Control-Allow-Methods", "GET, OPTIONS")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204) // No content for OPTIONS preflight
			return
		}

		c.Next()
	}
}
