package api

import (
	"web-porto-backend/internal/auth"
	"web-porto-backend/internal/handlers"
	"web-porto-backend/middleware"

	"github.com/gin-gonic/gin"
)

// SetupAPIRoutes configures all API routes
func SetupAPIRoutes(router *gin.Engine, handlerRegistry *handlers.HandlerRegistry, authService *auth.AuthService) {
	// Create API v1 group
	v1 := router.Group("/api/v1")

	// Setup routes using handler registry
	setupPublicRoutesWithRegistry(v1, handlerRegistry)
	setupProtectedRoutesWithRegistry(v1, handlerRegistry, authService)

	// Setup article and project routes
	SetupArticleProjectRoutes(router, handlerRegistry, authService)
}

// setupPublicRoutesWithRegistry configures public API routes using handler registry
func setupPublicRoutesWithRegistry(router *gin.RouterGroup, handlerRegistry *handlers.HandlerRegistry) {
	// Analytics routes (public for client-side tracking; can be protected by API key header if needed)
	analytics := router.Group("/analytics")
	{
		analytics.GET("/views", handlerRegistry.AnalyticsHandler.GetViews)
		analytics.GET("", handlerRegistry.AnalyticsHandler.GetAnalytics)
		analytics.GET("/series", handlerRegistry.AnalyticsHandler.GetSeries)
		analytics.POST("/track", handlerRegistry.AnalyticsHandler.Track)
	}

	// Content view tracking routes
	views := router.Group("/views")
	{
		views.POST("/track", handlerRegistry.AnalyticsHandler.TrackContentView)
		views.GET("/count", handlerRegistry.AnalyticsHandler.GetContentViewCount)
		views.GET("/analytics", handlerRegistry.AnalyticsHandler.GetContentViewAnalytics)
	}

	// Authentication routes
	auth := router.Group("/auth")
	{
		auth.POST("/register", handlerRegistry.AuthHandler.Register)
		auth.POST("/login", handlerRegistry.AuthHandler.Login)
	}

	// Public category routes
	categories := router.Group("/categories")
	{
		categories.GET("", handlerRegistry.CategoryHandler.GetAll)
		categories.GET("/:id", handlerRegistry.CategoryHandler.GetByID)
	}

	// Public post routes
	posts := router.Group("/posts")
	{
		posts.GET("", handlerRegistry.PostHandler.GetAll)
		posts.GET("/published", handlerRegistry.PostHandler.GetPublished)
		posts.GET("/:id", handlerRegistry.PostHandler.GetByID)
		posts.GET("/slug/:slug", handlerRegistry.PostHandler.GetBySlug)
		posts.GET("/author/:authorId", handlerRegistry.PostHandler.GetByAuthor)
	}

	// Public page routes
	pages := router.Group("/pages")
	{
		pages.GET("", handlerRegistry.PageHandler.GetAll)
		pages.GET("/published", handlerRegistry.PageHandler.GetPublished)
		pages.GET("/:id", handlerRegistry.PageHandler.GetByID)
		pages.GET("/slug/:slug", handlerRegistry.PageHandler.GetBySlug)
	}

	// Public comment routes (read-only)
	comments := router.Group("/comments")
	{
		comments.GET("", handlerRegistry.CommentHandler.GetAll)
		comments.GET("/:id", handlerRegistry.CommentHandler.GetByID)
		comments.GET("/post/:postId", handlerRegistry.CommentHandler.GetByPost)
	}
}

// setupProtectedRoutesWithRegistry configures protected API routes using handler registry
func setupProtectedRoutesWithRegistry(router *gin.RouterGroup, handlerRegistry *handlers.HandlerRegistry, authService *auth.AuthService) {
	// Protected routes group
	protected := router.Group("")
	protected.Use(middleware.JWTAuth(authService))

	// Auth profile routes
	auth := protected.Group("/auth")
	{
		auth.GET("/me", handlerRegistry.AuthHandler.Me)
	}

	// Protected category routes
	categories := protected.Group("/categories")
	{
		categories.POST("", handlerRegistry.CategoryHandler.Create)
		categories.PUT("/:id", handlerRegistry.CategoryHandler.Update)
		categories.DELETE("/:id", handlerRegistry.CategoryHandler.Delete)
	}

	// Protected post routes
	posts := protected.Group("/posts")
	{
		posts.POST("", handlerRegistry.PostHandler.Create)
		posts.PUT("/:id", handlerRegistry.PostHandler.Update)
		posts.DELETE("/:id", handlerRegistry.PostHandler.Delete)
	}

	// Protected page routes
	pages := protected.Group("/pages")
	{
		pages.POST("", handlerRegistry.PageHandler.Create)
		pages.PUT("/:id", handlerRegistry.PageHandler.Update)
		pages.DELETE("/:id", handlerRegistry.PageHandler.Delete)
	}

	// Protected comment routes
	comments := protected.Group("/comments")
	{
		comments.POST("", handlerRegistry.CommentHandler.Create)
		comments.PUT("/:id", handlerRegistry.CommentHandler.Update)
		comments.DELETE("/:id", handlerRegistry.CommentHandler.Delete)
	}
}
