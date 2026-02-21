package repositories

import (
	analyticsRepo "web-porto-backend/internal/repositories/analytics"
	articleRepo "web-porto-backend/internal/repositories/article"
	categoryRepo "web-porto-backend/internal/repositories/category"
	commentRepo "web-porto-backend/internal/repositories/comment"
	experienceRepo "web-porto-backend/internal/repositories/experience"
	pageRepo "web-porto-backend/internal/repositories/page"
	projectRepo "web-porto-backend/internal/repositories/project"
	tagRepo "web-porto-backend/internal/repositories/tag"
	userRepo "web-porto-backend/internal/repositories/user"

	"gorm.io/gorm"
)

type RepositoryRegistry struct {
	AnalyticsRepository  analyticsRepo.Repository
	ArticleRepository    articleRepo.Repository
	CategoryRepository   categoryRepo.Repository
	CommentRepository    commentRepo.Repository
	ExperienceRepository experienceRepo.Repository
	UserRepository       userRepo.Repository
	PageRepository       pageRepo.Repository
	ProjectRepository    projectRepo.Repository
	TagRepository        tagRepo.Repository
}

func NewRepositoryRegistry(db *gorm.DB) *RepositoryRegistry {
	return &RepositoryRegistry{
		AnalyticsRepository:  analyticsRepo.NewRepository(db),
		ArticleRepository:    articleRepo.NewRepository(db),
		CategoryRepository:   categoryRepo.NewRepository(db),
		CommentRepository:    commentRepo.NewRepository(db),
		ExperienceRepository: experienceRepo.NewRepository(db),
		UserRepository:       userRepo.NewRepository(db),
		PageRepository:       pageRepo.NewRepository(db),
		ProjectRepository:    projectRepo.NewRepository(db),
		TagRepository:        tagRepo.NewRepository(db),
	}
}
