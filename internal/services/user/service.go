package user

import (
	"web-porto-backend/internal/domain/models"
	"web-porto-backend/internal/repositories/user"

	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	GetAll(page, limit int) ([]*models.User, *PaginationInfo, error)
	GetByID(id uint) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	Create(user *models.User) error
	Update(id uint, user *models.User) error
	Delete(id uint) error
	CheckPassword(hashedPassword, password string) bool
}

// PaginationInfo represents pagination metadata
type PaginationInfo struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

type service struct {
	repo user.Repository
}

func NewService(repo user.Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetAll(page, limit int) ([]*models.User, *PaginationInfo, error) {
	users, err := s.repo.FindAll()
	if err != nil {
		return nil, nil, err
	}

	// Convert []models.User to []*models.User
	userPtrs := make([]*models.User, len(users))
	for i := range users {
		userPtrs[i] = &users[i]
	}

	// Calculate pagination (simplified)
	total := len(users)
	start := (page - 1) * limit
	end := start + limit

	if start >= total {
		userPtrs = []*models.User{}
	} else if end > total {
		userPtrs = userPtrs[start:]
	} else {
		userPtrs = userPtrs[start:end]
	}

	pagination := s.calculatePagination(total, page, limit)
	return userPtrs, pagination, nil
}

func (s *service) GetByID(id uint) (*models.User, error) {
	return s.repo.FindByID(int(id))
}

func (s *service) GetByEmail(email string) (*models.User, error) {
	return s.repo.FindByEmail(email)
}

func (s *service) Create(user *models.User) error {
	// Hash password before saving
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.PasswordHash = string(hashedPassword)

	return s.repo.Create(user)
}

func (s *service) Update(id uint, user *models.User) error {
	existingUser, err := s.repo.FindByID(int(id))
	if err != nil {
		return err
	}

	// Update fields
	existingUser.Username = user.Username
	existingUser.Email = user.Email
	existingUser.Role = user.Role

	// Hash password if provided
	if user.PasswordHash != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.PasswordHash), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		existingUser.PasswordHash = string(hashedPassword)
	}

	return s.repo.Update(existingUser)
}

func (s *service) Delete(id uint) error {
	return s.repo.Delete(int(id))
}

func (s *service) CheckPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

// calculatePagination calculates pagination metadata
func (s *service) calculatePagination(total, page, limit int) *PaginationInfo {
	totalPages := total / limit
	if total%limit != 0 {
		totalPages++
	}

	return &PaginationInfo{
		Page:       page,
		Limit:      limit,
		Total:      int64(total),
		TotalPages: totalPages,
	}
}
