package repositories

import (
	analyticsRepo "web-porto-backend/internal/repositories/analytics"
	categoryRepo "web-porto-backend/internal/repositories/category"
	commentRepo "web-porto-backend/internal/repositories/comment"
	pageRepo "web-porto-backend/internal/repositories/page"
	postRepo "web-porto-backend/internal/repositories/post"
	userRepo "web-porto-backend/internal/repositories/user"

	"gorm.io/gorm"
)

type RepositoryRegistry struct {
	AnalyticsRepository analyticsRepo.Repository
	CategoryRepository  categoryRepo.Repository
	CommentRepository   commentRepo.Repository
	UserRepository      userRepo.Repository
	PostRepository      postRepo.Repository
	PageRepository      pageRepo.Repository
}

func NewRepositoryRegistry(db *gorm.DB) *RepositoryRegistry {
	return &RepositoryRegistry{
		AnalyticsRepository: analyticsRepo.NewRepository(db),
		CategoryRepository:  categoryRepo.NewRepository(db),
		CommentRepository:   commentRepo.NewRepository(db),
		UserRepository:      userRepo.NewRepository(db),
		PostRepository:      postRepo.NewRepository(db),
		PageRepository:      pageRepo.NewRepository(db),
	}
}
