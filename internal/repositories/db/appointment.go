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

func (r *Repository) ListAppointments(doctorID int, dateFrom, dateTo *time.Time, status string, page, limit int) ([]model.Appointment, int64, error) {
	var appointments []model.Appointment
	var total int64

	query := r.client.Model(&model.Appointment{}).Where("doctor_id = ?", doctorID)

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
	err = query.Preload("Patient").Preload("Chamber").
		Order("appointment_date desc, start_time desc").
		Offset(offset).Limit(limit).
		Find(&appointments).Error

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
