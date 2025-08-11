package services

import (
	"web-porto-backend/internal/repositories"
	analyticsSrvc "web-porto-backend/internal/services/analytics"
	categorySrvc "web-porto-backend/internal/services/category"
	commentSrvc "web-porto-backend/internal/services/comment"
	pageSrvc "web-porto-backend/internal/services/page"
	postSrvc "web-porto-backend/internal/services/post"
	userSrvc "web-porto-backend/internal/services/user"
)

type ServiceRegistry struct {
	AnalyticsService analyticsSrvc.Service
	CategoryService  categorySrvc.Service
	CommentService   commentSrvc.Service
	UserService      userSrvc.Service
	PostService      postSrvc.Service
	PageService      pageSrvc.Service
}

func NewServiceRegistry(repo *repositories.RepositoryRegistry) *ServiceRegistry {
	return &ServiceRegistry{
		AnalyticsService: analyticsSrvc.NewService(repo.AnalyticsRepository),
		CategoryService:  categorySrvc.NewService(repo.CategoryRepository),
		CommentService:   commentSrvc.NewService(repo.CommentRepository),
		UserService:      userSrvc.NewService(repo.UserRepository),
		PostService:      postSrvc.NewService(repo.PostRepository),
		PageService:      pageSrvc.NewService(repo.PageRepository),
	}
}
