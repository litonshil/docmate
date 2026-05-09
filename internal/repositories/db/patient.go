package db

import (
	"docmate/internal/model"
)

func (r *Repository) CreatePatient(patient model.Patient) (model.Patient, error) {
	err := r.client.Create(&patient).Error

	return patient, err
}

func (r *Repository) UpdatePatient(patient model.Patient) (model.Patient, error) {
	// Omit CreatedAt and let GORM handle UpdatedAt automatically
	err := r.client.Model(&model.Patient{}).Where("id = ?", patient.ID).Omit("created_at").Updates(&patient).Error

	return patient, err
}

func (r *Repository) GetPatientByID(id int) (model.Patient, error) {
	var patient model.Patient
	if err := r.dbClient(nil).Model(&model.Patient{}).Where("id = ?", id).First(&patient).Error; err != nil {
		return model.Patient{}, err
	}

	return patient, nil
}

func (r *Repository) ListPatients(offset, limit, doctorID int, search string) ([]model.Patient, int64, error) {
	var patients []model.Patient
	var total int64

	query := r.client.Model(&model.Patient{})

	if doctorID > 0 {
		query = query.Where("doctor_id = ?", doctorID)
	}

	if search != "" {
		searchTerm := "%" + search + "%"
		query = query.Where("full_name ILIKE ? OR phone ILIKE ?", searchTerm, searchTerm)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Offset(offset).Limit(limit).Find(&patients).Error; err != nil {
		return nil, 0, err
	}

	return patients, total, nil
}
