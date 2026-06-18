package db

import (
	"docmate/internal/model"
	"time"

	"gorm.io/gorm/clause"
)

func (r *Repository) UpsertAISetting(setting model.AISetting) (model.AISetting, error) {
	err := r.client.Model(&model.AISetting{}).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "doctor_id"}},
			DoUpdates: clause.AssignmentColumns([]string{"provider", "individual_api_key", "use_individual_key", "updated_at"}),
		}).Create(&setting).Error

	return setting, err
}

func (r *Repository) AdminUpdateAISetting(setting model.AISetting) (model.AISetting, error) {
	err := r.client.Model(&model.AISetting{}).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "doctor_id"}},
			DoUpdates: clause.AssignmentColumns([]string{"is_ai_enabled", "allow_global_api", "use_individual_key", "updated_at"}),
		}).Create(&setting).Error

	return setting, err
}

func (r *Repository) GetAISettingByDoctor(doctorID int) (model.AISetting, error) {
	var setting model.AISetting
	err := r.client.Model(&model.AISetting{}).
		Where("doctor_id = ?", doctorID).
		First(&setting).Error

	return setting, err
}

func (r *Repository) GetGlobalSetting(key string) (string, error) {
	var setting model.GlobalSetting
	err := r.client.Model(&model.GlobalSetting{}).Where("key = ?", key).First(&setting).Error

	return setting.Value, err
}

func (r *Repository) SetGlobalSetting(key string, value string) error {
	var setting model.GlobalSetting
	setting.Key = key
	setting.Value = value
	setting.UpdatedAt = time.Now()

	return r.client.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "key"}},
		DoUpdates: clause.AssignmentColumns([]string{"value", "updated_at"}),
	}).Create(&setting).Error
}
