package db

import (
	"docmate/internal/model"
	"time"
)

func (r *Repository) CreateAppointment(appointment *model.Appointment) error {
	return r.client.Create(appointment).Error
}

func (r *Repository) GetAppointmentByID(id int) (*model.Appointment, error) {
	var appointment model.Appointment
	err := r.client.Preload("Patient").Preload("Chamber").First(&appointment, id).Error

	return &appointment, err
}

func (r *Repository) UpdateAppointment(appointment *model.Appointment) error {
	return r.client.Save(appointment).Error
}

func (r *Repository) UpdateAppointmentFields(id int, updates map[string]interface{}) error {
	return r.client.Model(&model.Appointment{}).Where("id = ?", id).Updates(updates).Error
}

func (r *Repository) ListAppointments(doctorID int, dateFrom, dateTo *time.Time, status string, search string, page, limit int) ([]model.Appointment, int64, error) {
	var appointments []model.Appointment
	var total int64

	query := r.client.Model(&model.Appointment{})
	if doctorID != 0 {
		query = query.Where("doctor_id = ?", doctorID)
	}

	if search != "" {
		var patientIDs []int
		searchPattern := "%" + search + "%"
		patientQuery := r.client.Model(&model.Patient{})
		if doctorID != 0 {
			patientQuery = patientQuery.Where("doctor_id = ?", doctorID)
		}
		patientQuery.Where("full_name ILIKE ? OR phone ILIKE ?", searchPattern, searchPattern).
			Pluck("id", &patientIDs)

		if len(patientIDs) == 0 {
			// If no patients found, return empty result
			return []model.Appointment{}, 0, nil
		}
		query = query.Where("patient_id IN ?", patientIDs)
	}

	if dateFrom != nil {
		query = query.Where("appointment_date >= ?", dateFrom)
	}
	if dateTo != nil {
		query = query.Where("appointment_date <= ?", dateTo)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	q := query.Preload("Patient").Preload("Chamber")

	if status == "pending" {
		q = q.Order("appointments.appointment_date asc, appointments.start_time asc")
	} else {
		q = q.Order("appointments.created_at desc")
	}

	err = q.Offset(offset).Limit(limit).Find(&appointments).Error

	return appointments, total, err
}

func (r *Repository) GetTodayAppointments(doctorID int) ([]model.Appointment, error) {
	var appointments []model.Appointment
	today := time.Now().Truncate(24 * time.Hour)
	err := r.client.Preload("Patient").
		Where("doctor_id = ? AND appointment_date = ? AND status != ?", doctorID, today, model.AppointmentStatusCancelled).
		Order("start_time asc").
		Find(&appointments).Error

	return appointments, err
}
