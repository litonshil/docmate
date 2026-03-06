package db

import (
	"docmate/internal/model"
)

func (repo *Repository) CreatePrescription(p model.Prescription) (model.Prescription, error) {
	err := repo.client.Create(&p).Error
	return p, err
}

func (repo *Repository) UpdatePrescription(p model.Prescription) (model.Prescription, error) {
	err := repo.client.Save(&p).Error
	return p, err
}

func (repo *Repository) GetPrescriptionByID(id int) (model.Prescription, error) {
	var p model.Prescription
	err := repo.client.First(&p, id).Error
	return p, err
}

func (repo *Repository) ListPrescriptions(doctorID int, limit, offset int) ([]model.Prescription, int, error) {
	var prescriptions []model.Prescription
	var count int64

	query := repo.client.Model(&model.Prescription{}).Where("doctor_id = ?", doctorID)

	err := query.Count(&count).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Order("created_at desc").Limit(limit).Offset(offset).Find(&prescriptions).Error
	return prescriptions, int(count), err
}

func (repo *Repository) ListPrescriptionsByPatient(doctorID, patientID int, limit, offset int) ([]model.Prescription, int, error) {
	var prescriptions []model.Prescription
	var count int64

	query := repo.client.Model(&model.Prescription{}).Where("doctor_id = ? AND patient_id = ?", doctorID, patientID)

	err := query.Count(&count).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Order("created_at desc").Limit(limit).Offset(offset).Find(&prescriptions).Error
	return prescriptions, int(count), err
}
