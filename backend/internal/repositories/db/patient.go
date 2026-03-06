package db

import (
	"docmate/internal/model"
)

func (repo *Repository) CreatePatient(patient model.Patient) (model.Patient, error) {
	err := repo.client.Create(&patient).Error

	return patient, err
}

func (repo *Repository) UpdatePatient(patient model.Patient) (model.Patient, error) {
	// Omit CreatedAt and let GORM handle UpdatedAt automatically
	err := repo.client.Model(&model.Patient{}).Where("id = ?", patient.ID).Omit("created_at").Updates(&patient).Error

	return patient, err
}

func (repo *Repository) GetPatientByID(id int) (model.Patient, error) {
	var patient model.Patient
	if err := repo.dbClient(nil).Model(&model.Patient{}).Where("id = ?", id).First(&patient).Error; err != nil {
		return model.Patient{}, err
	}

	return patient, nil
}

func (repo *Repository) ListPatients(offset, limit, doctorID int, name, phone string) ([]model.Patient, int, error) {
	var patients []model.Patient
	var total int64

	query := repo.client.Model(&model.Patient{})

	if doctorID > 0 {
		query = query.Where("doctor_id = ?", doctorID)
	}

	if name != "" {
		query = query.Where("full_name ILIKE ?", "%"+name+"%")
	}

	if phone != "" {
		query = query.Where("phone ILIKE ?", "%"+phone+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Offset(offset).Limit(limit).Find(&patients).Error; err != nil {
		return nil, 0, err
	}

	return patients, int(total), nil
}
