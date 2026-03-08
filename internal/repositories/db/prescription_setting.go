package db

import (
	"docmate/internal/model"

	"gorm.io/gorm/clause"
)

func (repo *Repository) Upsert(setting model.PrescriptionSetting) (model.PrescriptionSetting, error) {
	err := repo.client.Model(&model.PrescriptionSetting{}).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "doctor_id"}, {Name: "chamber_id"}},
			DoUpdates: clause.AssignmentColumns([]string{"header_left_bangla", "header_right_english", "footer_info_bangla", "footer_info_english", "template_type", "updated_at"}),
		}).Create(&setting).Error

	return setting, err
}

func (repo *Repository) GetByChamber(doctorID, chamberID int) (model.PrescriptionSetting, error) {
	var setting model.PrescriptionSetting
	err := repo.client.Model(&model.PrescriptionSetting{}).
		Where("doctor_id = ? AND chamber_id = ?", doctorID, chamberID).
		First(&setting).Error

	return setting, err
}
