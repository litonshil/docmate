package db

import (
	"docmate/internal/model"
)

func (r *Repository) CreatePrescription(p model.Prescription) (model.Prescription, error) {
	err := r.client.Create(&p).Error

	return p, err
}

func (r *Repository) UpdatePrescription(p model.Prescription) (model.Prescription, error) {
	err := r.client.Save(&p).Error

	return p, err
}

func (r *Repository) GetPrescriptionByID(id int) (model.Prescription, error) {
	var p model.Prescription
	err := r.client.Table("prescriptions").
		Select("prescriptions.*, patients.full_name as patient_name").
		Joins("left join patients on patients.id = prescriptions.patient_id").
		Where("prescriptions.id = ?", id).
		Scan(&p).Error

	return p, err
}

func (r *Repository) ListPrescriptions(doctorID int, limit, offset int, search string) ([]model.Prescription, int64, error) {
	var prescriptions []model.Prescription
	var count int64

	query := r.client.Table("prescriptions").
		Select("prescriptions.*, patients.full_name as patient_name").
		Joins("left join patients on patients.id = prescriptions.patient_id").
		Where("prescriptions.doctor_id = ?", doctorID)

	if search != "" {
		if idSearch := extractID(search); idSearch != "" {
			query = query.Where("(prescriptions.id = ? OR patients.full_name ILIKE ?)", idSearch, "%"+search+"%")
		} else {
			query = query.Where("patients.full_name ILIKE ?", "%"+search+"%")
		}
	}

	err := query.Count(&count).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Order("prescriptions.created_at desc").Limit(limit).Offset(offset).Scan(&prescriptions).Error

	return prescriptions, count, err
}

func extractID(s string) string {
	if len(s) > 3 && s[:3] == "PR-" {
		return s[3:]
	}

	return ""
}

func (r *Repository) ListPrescriptionsByPatient(doctorID, patientID int, limit, offset int, search string) ([]model.Prescription, int64, error) {
	var prescriptions []model.Prescription
	var count int64

	query := r.client.Table("prescriptions").
		Select("prescriptions.*, patients.full_name as patient_name").
		Joins("left join patients on patients.id = prescriptions.patient_id").
		Where("prescriptions.doctor_id = ? AND prescriptions.patient_id = ?", doctorID, patientID)

	if search != "" {
		if idSearch := extractID(search); idSearch != "" {
			query = query.Where("(prescriptions.id = ? OR patients.full_name ILIKE ?)", idSearch, "%"+search+"%")
		} else {
			query = query.Where("patients.full_name ILIKE ?", "%"+search+"%")
		}
	}

	err := query.Count(&count).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Order("prescriptions.created_at desc").Limit(limit).Offset(offset).Scan(&prescriptions).Error

	return prescriptions, count, err
}
