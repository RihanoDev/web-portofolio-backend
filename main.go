package main

import (
	"fmt"
	"log"
	"strconv"
	"time"
	"web-porto-backend/config"
	"web-porto-backend/internal/adapters/http"
	"web-porto-backend/internal/auth"
	"web-porto-backend/internal/domain/models"
	"web-porto-backend/internal/handlers"
	"web-porto-backend/internal/migrations"
	"web-porto-backend/internal/repositories"
	"web-porto-backend/internal/services"
	"web-porto-backend/middleware"
	"web-porto-backend/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Database connection
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Jakarta",
		cfg.Database.Host,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Name,
		cfg.Database.Port,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto-migrate analytics table (optional safety)
	if err := db.AutoMigrate(&models.PageView{}); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// Apply SQL migrations in database_schema
	log.Println("Running SQL migrations (database_schema)...")
	if err := migrations.RunMigrations(db); err != nil {
		log.Fatal("Failed running migrations:", err)
	}

	// Initialize layers
	repositoryRegistry := repositories.NewRepositoryRegistry(db)
	serviceRegistry := services.NewServiceRegistry(repositoryRegistry)
	authService := auth.NewAuthService(cfg.JWT.Secret)

	// Initialize HTTP adapter
	httpAdapter := http.NewHTTPAdapter()

	handlerRegistry := handlers.NewHandlerRegistry(serviceRegistry, authService, httpAdapter) // Setup Gin router
	router := gin.Default()

	// Add CORS middleware
	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:3000",
			"http://localhost:8080",
			"http://localhost:5173",
			"http://localhost:5174",
			"https://api.rihanodev.com",
			"https://cms.rihanodev.com",
			"https://rihanodev.com",
			"https://www.rihanodev.com",
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-API-Key"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Add logging middleware
	logger := logrus.New()
	router.Use(middleware.Logger(logger))

	// Register routes using the new routes structure
	routes.SetupRoutes(router, handlerRegistry, authService)

	// Register analytics track route with API key + rate limiting if API key configured
	if cfg.Analytics.APIKey != "" {
		router.POST(
			"/api/v1/analytics/track",
			middleware.APIKeyAuth(cfg),
			middleware.RateLimit(10, 10*time.Second),
			handlerRegistry.AnalyticsHandler.Track,
		)
	} else {
		// Without API key, still apply a rate limit to be safe
		router.POST(
			"/api/v1/analytics/track",
			middleware.RateLimit(10, 10*time.Second),
			handlerRegistry.AnalyticsHandler.Track,
		)
	}

	// Start server
	port := strconv.Itoa(cfg.Server.Port)
	log.Printf("Server starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
