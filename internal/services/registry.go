package services

import (
	"web-porto-backend/internal/repositories"
	analyticsSrvc "web-porto-backend/internal/services/analytics"
	articleSrvc "web-porto-backend/internal/services/article"
	categorySrvc "web-porto-backend/internal/services/category"
	commentSrvc "web-porto-backend/internal/services/comment"
	experienceSrvc "web-porto-backend/internal/services/experience"
	pageSrvc "web-porto-backend/internal/services/page"
	projectSrvc "web-porto-backend/internal/services/project"
	tagSrvc "web-porto-backend/internal/services/tag"
	userSrvc "web-porto-backend/internal/services/user"
)

type ServiceRegistry struct {
	AnalyticsService  analyticsSrvc.Service
	ArticleService    *articleSrvc.Service
	CategoryService   categorySrvc.Service
	CommentService    commentSrvc.Service
	ExperienceService *experienceSrvc.Service
	UserService       userSrvc.Service
	PageService       pageSrvc.Service
	ProjectService    *projectSrvc.Service
	TagService        tagSrvc.Service
}

func NewServiceRegistry(repo *repositories.RepositoryRegistry) *ServiceRegistry {
	// Create user service first
	userService := userSrvc.NewService(repo.UserRepository)

	// Create tag service
	tagService := tagSrvc.NewService(repo.TagRepository)

	return &ServiceRegistry{
		AnalyticsService: analyticsSrvc.NewService(repo.AnalyticsRepository),
		ArticleService: articleSrvc.NewService(
			repo.ArticleRepository,
			repo.CategoryRepository,
			repo.TagRepository,
			userService,
		),
		CategoryService: categorySrvc.NewService(repo.CategoryRepository),
		CommentService:  commentSrvc.NewService(repo.CommentRepository),
		ExperienceService: experienceSrvc.NewService(
			repo.ExperienceRepository,
			tagService,
		),
		UserService: userService,
		PageService: pageSrvc.NewService(repo.PageRepository),
		ProjectService: projectSrvc.NewService(
			repo.ProjectRepository,
			repo.CategoryRepository,
			userService,
			tagService,
		),
		TagService: tagService,
	}
}
