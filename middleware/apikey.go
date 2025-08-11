package middleware

import (
	"net/http"
	"web-porto-backend/config"

	"github.com/gin-gonic/gin"
)

// APIKeyAuth validates X-API-Key header against config.analytics.api_key
func APIKeyAuth(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		expected := cfg.Analytics.APIKey
		if expected == "" {
			c.Next()
			return
		}
		provided := c.GetHeader("X-API-Key")
		if provided == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "X-API-Key header required"})
			c.Abort()
			return
		}
		if provided != expected {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid API key"})
			c.Abort()
			return
		}
		c.Next()
	}
}
