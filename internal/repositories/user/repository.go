package user

import (
	applog "web-porto-backend/common/logger"
	"web-porto-backend/internal/domain/models"

	"gorm.io/gorm"
)

const (
	msgDBError      = "db error"
	msgFetchedUsers = "fetched users"
	msgFoundUser    = "found user"
	msgCreatedUser  = "created user"
	msgUpdatedUser  = "updated user"
	msgDeletedUser  = "deleted user"
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
	log := applog.GetLogger().WithFields(applog.Fields{"repo": "user", "method": "FindAll"})
	var users []models.User
	const msgDBError = "db error"
	const msgFetchedUsers = "fetched users"
	err := r.db.Find(&users).Error
	if err != nil {
		log.Error(msgDBError, applog.Fields{"error": err.Error()})
	} else {
		log.Info(msgFetchedUsers, applog.Fields{"count": len(users)})
	}
	return users, err
}

func (r *repository) FindByID(id int) (*models.User, error) {
	log := applog.GetLogger().WithFields(applog.Fields{"repo": "user", "method": "FindByID", "id": id})
	var user models.User
	err := r.db.First(&user, id).Error
	if err != nil {
		log.Error(msgDBError, applog.Fields{"error": err.Error()})
		return nil, err
	}
	const msgFoundUser = "found user"
	log.Info(msgFoundUser)
	return &user, nil
}

func (r *repository) FindByEmail(email string) (*models.User, error) {
	log := applog.GetLogger().WithFields(applog.Fields{"repo": "user", "method": "FindByEmail", "email": email})
	var user models.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		log.Error(msgDBError, applog.Fields{"error": err.Error()})
		return nil, err
	}
	log.Info(msgFoundUser)
	return &user, nil
}

func (r *repository) Create(user *models.User) error {
	log := applog.GetLogger().WithFields(applog.Fields{"repo": "user", "method": "Create"})
	const msgCreatedUser = "created user"
	if err := r.db.Create(user).Error; err != nil {
		log.Error(msgDBError, applog.Fields{"error": err.Error()})
		return err
	}
	log.Info(msgCreatedUser, applog.Fields{"id": user.ID})
	return nil
}

func (r *repository) Update(user *models.User) error {
	log := applog.GetLogger().WithFields(applog.Fields{"repo": "user", "method": "Update", "id": user.ID})
	const msgUpdatedUser = "updated user"
	if err := r.db.Save(user).Error; err != nil {
		log.Error(msgDBError, applog.Fields{"error": err.Error()})
		return err
	}
	log.Info(msgUpdatedUser)
	return nil
}

func (r *repository) Delete(id int) error {
	log := applog.GetLogger().WithFields(applog.Fields{"repo": "user", "method": "Delete", "id": id})
	const msgDeletedUser = "deleted user"
	if err := r.db.Delete(&models.User{}, id).Error; err != nil {
		log.Error(msgDBError, applog.Fields{"error": err.Error()})
		return err
	}
	log.Info(msgDeletedUser)
	return nil
}
