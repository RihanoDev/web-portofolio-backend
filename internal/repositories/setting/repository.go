package setting

import (
	"web-porto-backend/internal/domain/models"

	"gorm.io/gorm"
)

type Repository interface {
	GetByKey(key string) (*models.Setting, error)
	GetByKeys(keys []string) ([]models.Setting, error)
	GetAll() ([]models.Setting, error)
	Save(setting *models.Setting) error
	SaveMany(settings []models.Setting) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) GetByKey(key string) (*models.Setting, error) {
	var setting models.Setting
	if err := r.db.Where("key = ?", key).First(&setting).Error; err != nil {
		return nil, err
	}
	return &setting, nil
}

func (r *repository) GetByKeys(keys []string) ([]models.Setting, error) {
	var settings []models.Setting
	if err := r.db.Where("key IN ?", keys).Find(&settings).Error; err != nil {
		return nil, err
	}
	return settings, nil
}

func (r *repository) GetAll() ([]models.Setting, error) {
	var settings []models.Setting
	if err := r.db.Find(&settings).Error; err != nil {
		return nil, err
	}
	return settings, nil
}

func (r *repository) Save(setting *models.Setting) error {
	return r.db.Save(setting).Error
}

func (r *repository) SaveMany(settings []models.Setting) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		for _, s := range settings {
			if err := tx.Save(&s).Error; err != nil {
				return err
			}
		}
		return nil
	})
}
