package db

import (
	"docmate/internal/model"

	"gorm.io/gorm/clause"
)

func (repo *Repository) UpsertAISetting(setting model.AISetting) (model.AISetting, error) {
	err := repo.client.Model(&model.AISetting{}).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "doctor_id"}},
			DoUpdates: clause.AssignmentColumns([]string{"provider", "individual_api_key", "use_individual_key", "updated_at"}),
		}).Create(&setting).Error

	return setting, err
}

func (repo *Repository) AdminUpdateAISetting(setting model.AISetting) (model.AISetting, error) {
	err := repo.client.Model(&model.AISetting{}).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "doctor_id"}},
			DoUpdates: clause.AssignmentColumns([]string{"is_ai_enabled", "allow_global_api", "updated_at"}),
		}).Create(&setting).Error

	return setting, err
}

func (repo *Repository) GetAISettingByDoctor(doctorID int) (model.AISetting, error) {
	var setting model.AISetting
	err := repo.client.Model(&model.AISetting{}).
		Where("doctor_id = ?", doctorID).
		First(&setting).Error

	return setting, err
}
