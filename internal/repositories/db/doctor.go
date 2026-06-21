package db

import (
	"docmate/internal/model"

	"gorm.io/gorm/clause"
)

func (r *Repository) CreateDoctor(doctor model.Doctor) (model.Doctor, error) {
	err := r.client.Create(&doctor).Error

	return doctor, err
}

func (r *Repository) UpdateDoctor(doctor model.Doctor) (model.Doctor, error) {
	// Omit CreatedAt and let GORM handle UpdatedAt automatically
	err := r.client.Model(&model.Doctor{}).Where("id = ?", doctor.ID).Omit("created_at").Updates(&doctor).Error

	return doctor, err
}

func (r *Repository) GetDoctorByID(id int) (model.Doctor, error) {
	var doctor model.Doctor
	if err := r.client.Model(&model.Doctor{}).Where("id = ?", id).First(&doctor).Error; err != nil {
		return model.Doctor{}, err
	}

	return doctor, nil
}

func (r *Repository) ListDoctors(offset, limit int) ([]model.Doctor, int64, error) {
	var doctors []model.Doctor
	var total int64

	if err := r.client.Model(&model.Doctor{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := r.client.Model(&model.Doctor{}).
		Select("doctors.*, ai.ai_request_status, COALESCE(ai.is_ai_enabled, false) as is_ai_enabled").
		Joins("LEFT JOIN (SELECT doctor_id, MAX(ai_request_status) as ai_request_status, MAX(CASE WHEN is_ai_enabled = true THEN 1 ELSE 0 END) = 1 as is_ai_enabled FROM doctor_ai_settings GROUP BY doctor_id) ai ON ai.doctor_id = doctors.id").
		Offset(offset).
		Limit(limit).
		Order("doctors.id desc").
		Find(&doctors).Error

	if err != nil {
		return nil, 0, err
	}

	return doctors, total, nil
}

func (r *Repository) GetDoctorByUserID(userID int) (model.Doctor, error) {
	var doctor model.Doctor
	if err := r.client.Model(&model.Doctor{}).Where("user_id = ?", userID).First(&doctor).Error; err != nil {
		return model.Doctor{}, err
	}

	return doctor, nil
}

func (r *Repository) UpsertDoctor(doctor model.Doctor) (model.Doctor, error) {
	// Upsert profile based on user_id conflict
	err := r.client.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "user_id"}}, // conflict target
		DoUpdates: clause.AssignmentColumns([]string{
			"full_name",
			"degree",
			"specialization",
			"phone",
			"bio",
			"signature_url",
			"updated_at",
		}),
	}).Create(&doctor).Error

	return doctor, err
}
