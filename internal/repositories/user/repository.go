package user

import (
	"web-porto-backend/internal/domain/models"

	"gorm.io/gorm"
)

type Repository interface {
	FindAll() ([]models.User, error)
	FindByID(id int) (*models.User, error)
	FindByEmail(email string) (*models.User, error)
	Create(user *models.User) error
	Update(user *models.User) error
	Delete(id int) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) FindAll() ([]models.User, error) {
	var users []models.User
	err := r.db.Find(&users).Error
	return users, err
}

func (r *repository) FindByID(id int) (*models.User, error) {
	var user models.User
	err := r.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *repository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *repository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *repository) Update(user *models.User) error {
	return r.db.Save(user).Error
}

func (r *repository) Delete(id int) error {
	return r.db.Delete(&models.User{}, id).Error
}
