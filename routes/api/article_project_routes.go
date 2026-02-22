package api

import (
	"web-porto-backend/internal/auth"
	"web-porto-backend/internal/handlers"
	"web-porto-backend/middleware"

	"github.com/gin-gonic/gin"
)

// SetupArticleProjectRoutes configures routes for articles and projects
func SetupArticleProjectRoutes(router *gin.Engine, handlerRegistry *handlers.HandlerRegistry, authService *auth.AuthService) {
	// API v1 group
	v1 := router.Group("/api/v1")

	// Public article routes
	articles := v1.Group("/articles")
	{
		articles.GET("", handlerRegistry.ArticleHandler.GetAll)
		articles.GET("/published", handlerRegistry.ArticleHandler.GetPublished)
		articles.GET("/:id", handlerRegistry.ArticleHandler.GetByID)
		articles.GET("/slug/:slug", handlerRegistry.ArticleHandler.GetBySlug)
		articles.GET("/category/:slug", handlerRegistry.ArticleHandler.GetByCategory)
		articles.GET("/tag/:name", handlerRegistry.ArticleHandler.GetByTag)
	}

	// Public project routes
	projects := v1.Group("/projects")
	{
		projects.GET("", handlerRegistry.ProjectHandler.GetAll)
		projects.GET("/published", handlerRegistry.ProjectHandler.GetPublished)
		projects.GET("/:id", handlerRegistry.ProjectHandler.GetByID)
		projects.GET("/slug/:slug", handlerRegistry.ProjectHandler.GetBySlug)
		projects.GET("/category/:slug", handlerRegistry.ProjectHandler.GetByCategory)
		projects.GET("/technology/:name", handlerRegistry.ProjectHandler.GetByTechnology)
	}

	// Public experience routes
	experiences := v1.Group("/experiences")
	{
		experiences.GET("", handlerRegistry.ExperienceHandler.GetAll)
		experiences.GET("/current", handlerRegistry.ExperienceHandler.GetCurrent)
		experiences.GET("/:id", handlerRegistry.ExperienceHandler.GetByID)
	}

	// Protected routes group
	protected := v1.Group("")
	protected.Use(middleware.JWTAuth(authService))

	// Protected article routes
	protectedArticles := protected.Group("/articles")
	{
		protectedArticles.POST("", handlerRegistry.ArticleHandler.Create)
		protectedArticles.PUT("/:id", handlerRegistry.ArticleHandler.Update)
		protectedArticles.PATCH("/:id", handlerRegistry.ArticleHandler.Patch)
		protectedArticles.DELETE("/:id", handlerRegistry.ArticleHandler.Delete)
		protectedArticles.POST("/:id/images", handlerRegistry.ArticleHandler.AddImage)
		protectedArticles.POST("/:id/videos", handlerRegistry.ArticleHandler.AddVideo)
		protectedArticles.DELETE("/:id/images/:imageId", handlerRegistry.ArticleHandler.DeleteImage)
		protectedArticles.DELETE("/:id/videos/:videoId", handlerRegistry.ArticleHandler.DeleteVideo)
	}

	// Protected project routes
	protectedProjects := protected.Group("/projects")
	{
		protectedProjects.POST("", handlerRegistry.ProjectHandler.Create)
		protectedProjects.PUT("/:id", handlerRegistry.ProjectHandler.Update)
		protectedProjects.PATCH("/:id", handlerRegistry.ProjectHandler.Patch)
		protectedProjects.DELETE("/:id", handlerRegistry.ProjectHandler.Delete)
		protectedProjects.POST("/:id/images", handlerRegistry.ProjectHandler.AddImage)
		protectedProjects.POST("/:id/videos", handlerRegistry.ProjectHandler.AddVideo)
		protectedProjects.DELETE("/:id/images/:imageId", handlerRegistry.ProjectHandler.DeleteImage)
		protectedProjects.DELETE("/:id/videos/:videoId", handlerRegistry.ProjectHandler.DeleteVideo)
		protectedProjects.POST("/:id/technologies", handlerRegistry.ProjectHandler.AddTechnology)
		protectedProjects.DELETE("/:id/technologies/:techId", handlerRegistry.ProjectHandler.RemoveTechnology)
	}

	// Protected experience routes
	protectedExperiences := protected.Group("/experiences")
	{
		protectedExperiences.POST("", handlerRegistry.ExperienceHandler.Create)
		protectedExperiences.PUT("/:id", handlerRegistry.ExperienceHandler.Update)
		protectedExperiences.PATCH("/:id", handlerRegistry.ExperienceHandler.Patch)
		protectedExperiences.DELETE("/:id", handlerRegistry.ExperienceHandler.Delete)
	}

	// Tags management
	tags := protected.Group("/tags")
	{
		tags.GET("", handlerRegistry.TagHandler.GetAll)
		tags.POST("", handlerRegistry.TagHandler.Create)
		tags.PUT("/:id", handlerRegistry.TagHandler.Update)
		tags.DELETE("/:id", handlerRegistry.TagHandler.Delete)
	}

	// Media upload routes
	if handlerRegistry.MediaHandler != nil {
		// Public: list media
		media := v1.Group("/media")
		{
			media.GET("", handlerRegistry.MediaHandler.GetAll)
		}

		// Protected: upload and delete
		protectedMedia := protected.Group("/media")
		{
			protectedMedia.POST("/upload", handlerRegistry.MediaHandler.Upload)
			protectedMedia.DELETE("/:id", handlerRegistry.MediaHandler.Delete)
		}
	}
}
