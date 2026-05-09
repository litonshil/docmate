package db

import (
	"docmate/internal/model"

	"gorm.io/gorm/clause"
)

func (r *Repository) Upsert(setting model.PrescriptionSetting) (model.PrescriptionSetting, error) {
	err := r.client.Model(&model.PrescriptionSetting{}).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "doctor_id"}, {Name: "chamber_id"}},
			DoUpdates: clause.AssignmentColumns([]string{"header_left_bangla", "header_right_english", "chamber_info", "visiting_hour", "template_type", "updated_at"}),
		}).Create(&setting).Error

	return setting, err
}

func (r *Repository) GetByChamber(doctorID, chamberID int) (model.PrescriptionSetting, error) {
	var setting model.PrescriptionSetting
	err := r.client.Model(&model.PrescriptionSetting{}).
		Where("doctor_id = ? AND chamber_id = ?", doctorID, chamberID).
		First(&setting).Error

	return setting, err
}
