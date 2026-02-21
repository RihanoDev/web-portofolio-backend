package handlers

import (
	httpAdapter "web-porto-backend/internal/adapters/http"
	"web-porto-backend/internal/auth"
	analyticsHandler "web-porto-backend/internal/handlers/analytics"
	articleHandler "web-porto-backend/internal/handlers/article"
	authHandler "web-porto-backend/internal/handlers/auth"
	categoryHandler "web-porto-backend/internal/handlers/category"
	commentHandler "web-porto-backend/internal/handlers/comment"
	experienceHandler "web-porto-backend/internal/handlers/experience"
	mediaHandler "web-porto-backend/internal/handlers/media"
	pageHandler "web-porto-backend/internal/handlers/page"
	postHandler "web-porto-backend/internal/handlers/post"
	projectHandler "web-porto-backend/internal/handlers/project"
	settingHandler "web-porto-backend/internal/handlers/setting"
	tagHandler "web-porto-backend/internal/handlers/tag"
	userHandler "web-porto-backend/internal/handlers/user"
	"web-porto-backend/internal/services"

	"gorm.io/gorm"
)

type HandlerRegistry struct {
	AnalyticsHandler  *analyticsHandler.Handler
	ArticleHandler    *articleHandler.Handler
	CategoryHandler   *categoryHandler.Handler
	CommentHandler    *commentHandler.Handler
	AuthHandler       *authHandler.Handler
	ExperienceHandler *experienceHandler.Handler
	MediaHandler      *mediaHandler.Handler
	PostHandler       *postHandler.Handler
	PageHandler       *pageHandler.Handler
	ProjectHandler    *projectHandler.Handler
	SettingHandler    *settingHandler.Handler
	TagHandler        *tagHandler.Handler
	UserHandler       *userHandler.Handler
}

func NewHandlerRegistry(svc *services.ServiceRegistry, authService *auth.AuthService, httpAdapter *httpAdapter.HTTPAdapter) *HandlerRegistry {
	return NewHandlerRegistryWithDB(svc, authService, httpAdapter, nil, "", "")
}

func NewHandlerRegistryWithDB(svc *services.ServiceRegistry, authService *auth.AuthService, httpAdapter *httpAdapter.HTTPAdapter, db *gorm.DB, uploadDir string, baseURL string) *HandlerRegistry {
	// Create adapter from article service to post service for backward compatibility
	postServiceAdapter := postHandler.NewPostServiceAdapter(svc.ArticleService)

	var mHandler *mediaHandler.Handler
	if db != nil {
		mHandler = mediaHandler.NewHandler(db, uploadDir, baseURL)
	}

	return &HandlerRegistry{
		AnalyticsHandler:  analyticsHandler.NewHandler(svc.AnalyticsService, httpAdapter),
		ArticleHandler:    articleHandler.NewHandler(svc.ArticleService, httpAdapter),
		CategoryHandler:   categoryHandler.NewHandler(svc.CategoryService), // Using existing constructor
		CommentHandler:    commentHandler.NewHandler(svc.CommentService),   // Using existing constructor
		AuthHandler:       authHandler.NewHandler(svc.UserService, authService, httpAdapter),
		ExperienceHandler: experienceHandler.NewHandler(svc.ExperienceService, httpAdapter),
		MediaHandler:      mHandler,
		PostHandler:       postHandler.NewHandler(postServiceAdapter, httpAdapter),
		PageHandler:       pageHandler.NewHandler(svc.PageService, httpAdapter),
		ProjectHandler:    projectHandler.NewHandler(svc.ProjectService, httpAdapter),
		SettingHandler:    settingHandler.NewHandler(svc.SettingService),
		TagHandler:        tagHandler.NewHandler(svc.TagService, httpAdapter),
		UserHandler:       userHandler.NewHandler(svc.UserService, httpAdapter),
	}
}
