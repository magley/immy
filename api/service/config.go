package service

import (
	"immy-api/model"
	"immy-api/repo"
)

type ConfigService struct {
	ConfigRepo *repo.ConfigRepo
}

func (s *ConfigService) GetConfig() (*model.Config, error) {
	return s.ConfigRepo.GetConfig()
}

func (s *ConfigService) SetPostingEnabled(enabled bool) (*model.Config, error) {
	config, err := s.ConfigRepo.GetConfig()
	if err != nil {
		return nil, err
	}
	return s.ConfigRepo.SetPostingEnabled(config, enabled)
}
