package middleware

import (
	"net/http"
	"strings"

	"web-porto-backend/internal/auth"

	"github.com/gin-gonic/gin"
)

func JWTAuth(authService *auth.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Debugging: log request headers
		reqPath := c.FullPath()
		reqMethod := c.Request.Method

		// Get auth header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":  "Authorization header required",
				"path":   reqPath,
				"method": reqMethod,
			})
			c.Abort()
			return
		}

		// Log the auth header for debugging (partially redacted)
		headerLen := len(authHeader)
		headerPrefix := authHeader
		if headerLen > 15 {
			headerPrefix = authHeader[:15] + "..."
		}

		// Check for "Bearer " prefix
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":         "Invalid Authorization format, must be 'Bearer <token>'",
				"header_format": headerPrefix,
			})
			c.Abort()
			return
		}

		// Extract token
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Validate token
		claims, err := authService.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":  "Invalid token: " + err.Error(),
				"path":   reqPath,
				"method": reqMethod,
			})
			c.Abort()
			return
		}

		// Store claims in context
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)
		c.Next()
	}
}

func RequireRole(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		if userRole != role && userRole != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
			c.Abort()
			return
		}

		c.Next()
	}
}
