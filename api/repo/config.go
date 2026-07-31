package repo

import (
	"gorm.io/gorm"

	model "immy-api/model"
)

type ConfigRepo struct {
	DB *gorm.DB
}

func (r *ConfigRepo) GetConfig() (*model.Config, error) {
	var config *model.Config
	result := r.DB.Where("id = ?", 1).Find(&config)
	return config, result.Error
}

func (r *ConfigRepo) SetPostingEnabled(config *model.Config, enabled bool) (*model.Config, error) {
	config.PostingEnabled = enabled
	result := r.DB.Save(&config)
	return config, result.Error
}
