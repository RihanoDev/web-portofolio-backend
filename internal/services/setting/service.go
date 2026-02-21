package setting

import (
	"web-porto-backend/internal/domain/models"
	"web-porto-backend/internal/repositories/setting"
)

type Service interface {
	GetSetting(key string) (string, error)
	GetSettings(keys []string) (map[string]string, error)
	GetAllSettings() (map[string]string, error)
	SaveSettings(updates map[string]string) error
}

type service struct {
	repo setting.Repository
}

func NewService(repo setting.Repository) Service {
	return &service{repo}
}

func (s *service) GetSetting(key string) (string, error) {
	setting, err := s.repo.GetByKey(key)
	if err != nil {
		return "", err
	}
	return setting.Value, nil
}

func (s *service) GetSettings(keys []string) (map[string]string, error) {
	settings, err := s.repo.GetByKeys(keys)
	if err != nil {
		return nil, err
	}

	result := make(map[string]string)
	for _, setting := range settings {
		result[setting.Key] = setting.Value
	}
	return result, nil
}

func (s *service) GetAllSettings() (map[string]string, error) {
	settings, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}

	result := make(map[string]string)
	for _, setting := range settings {
		result[setting.Key] = setting.Value
	}
	return result, nil
}

func (s *service) SaveSettings(updates map[string]string) error {
	var settings []models.Setting
	for k, v := range updates {
		settings = append(settings, models.Setting{Key: k, Value: v})
	}
	return s.repo.SaveMany(settings)
}
