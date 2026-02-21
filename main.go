package main

import (
	"fmt"
	"log"
	"strconv"
	"time"
	appLogger "web-porto-backend/common/logger"
	"web-porto-backend/config"
	httpAdapter "web-porto-backend/internal/adapters/http"
	"web-porto-backend/internal/adapters/websocket"
	"web-porto-backend/internal/auth"
	"web-porto-backend/internal/domain/models"
	"web-porto-backend/internal/handlers"
	"web-porto-backend/internal/migrations"
	"web-porto-backend/internal/repositories"
	"web-porto-backend/internal/services"
	analyticsSrvc "web-porto-backend/internal/services/analytics"
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

	// Setup application logger
	var baseLogger appLogger.Logger
	if cfg.App.Debug {
		baseLogger = appLogger.NewDevelopmentLogger()
	} else {
		baseLogger = appLogger.NewProductionLogger()
	}
	appLogger.SetDefaultLogger(baseLogger)

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

	// Seed initial data (admin user) for development
	if err := migrations.Seed(db); err != nil {
		log.Fatal("Failed seeding data:", err)
	}

	// Initialize layers
	repositoryRegistry := repositories.NewRepositoryRegistry(db)
	serviceRegistry := services.NewServiceRegistry(repositoryRegistry)
	authService := auth.NewAuthService(cfg.JWT.Secret)

	// Initialize WebSocket manager
	wsManager := websocket.NewManager()
	go wsManager.Start() // Start WebSocket manager in a goroutine

	// Initialize HTTP adapter
	httpAdpt := httpAdapter.NewHTTPAdapter()

	// Connect WebSocket manager to analytics service
	analyticsService, ok := serviceRegistry.AnalyticsService.(analyticsSrvc.Service)
	if ok {
		analyticsService.SetWebsocketManager(wsManager)
	} else {
		log.Println("Warning: Could not connect WebSocket manager to analytics service")
	}

	handlerRegistry := handlers.NewHandlerRegistry(serviceRegistry, authService, httpAdpt) // Setup Gin router
	router := gin.Default()

	// Add CORS middleware with specific origins for development and production
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // Development URLs
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-API-Key", "X-Requested-With", "Sec-WebSocket-Protocol", "Sec-WebSocket-Version", "Sec-WebSocket-Key", "Upgrade", "Connection"},
		ExposeHeaders:    []string{"Content-Length", "Content-Type", "Upgrade", "Connection"},
		AllowCredentials: true, // Enable credentials for auth
		MaxAge:           12 * time.Hour,
	}))

	// Add logging middleware (use logrus for HTTP access logs)
	logger := logrus.New()
	router.Use(middleware.Logger(logger))

	// Register API routes
	routes.SetupRoutes(router, handlerRegistry, authService)

	// Register WebSocket routes
	routes.SetupWebSocketRoutes(router, wsManager)

	// Start server
	port := strconv.Itoa(cfg.Server.Port)
	log.Printf("Server starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
