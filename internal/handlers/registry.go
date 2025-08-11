package handlers

import (
	httpAdapter "web-porto-backend/internal/adapters/http"
	"web-porto-backend/internal/auth"
	analyticsHandler "web-porto-backend/internal/handlers/analytics"
	authHandler "web-porto-backend/internal/handlers/auth"
	categoryHandler "web-porto-backend/internal/handlers/category"
	commentHandler "web-porto-backend/internal/handlers/comment"
	pageHandler "web-porto-backend/internal/handlers/page"
	postHandler "web-porto-backend/internal/handlers/post"
	userHandler "web-porto-backend/internal/handlers/user"
	"web-porto-backend/internal/services"
)

type HandlerRegistry struct {
	AnalyticsHandler *analyticsHandler.Handler
	CategoryHandler  *categoryHandler.Handler
	CommentHandler   *commentHandler.Handler
	AuthHandler      *authHandler.Handler
	PostHandler      *postHandler.Handler
	PageHandler      *pageHandler.Handler
	UserHandler      *userHandler.Handler
}

func NewHandlerRegistry(svc *services.ServiceRegistry, authService *auth.AuthService, httpAdapter *httpAdapter.HTTPAdapter) *HandlerRegistry {
	return &HandlerRegistry{
		AnalyticsHandler: analyticsHandler.NewHandler(svc.AnalyticsService, httpAdapter),
		// CategoryHandler: categoryHandler.NewHandler(svc.CategoryService, httpAdapter),
		// CommentHandler:  commentHandler.NewHandler(svc.CommentService, httpAdapter),
		AuthHandler: authHandler.NewHandler(svc.UserService, authService, httpAdapter),
		PostHandler: postHandler.NewHandler(svc.PostService, httpAdapter),
		PageHandler: pageHandler.NewHandler(svc.PageService, httpAdapter),
		UserHandler: userHandler.NewHandler(svc.UserService, httpAdapter),
	}
}
