package appointment

import (
	"docmate/internal/model"
	"time"

	"gorm.io/gorm"
)

type appointmentRepo struct {
	db *gorm.DB
}

func NewAppointmentRepo(db *gorm.DB) model.AppointmentRepo {
	return &appointmentRepo{db: db}
}

func (r *appointmentRepo) CreateAppointment(appointment *model.Appointment) error {
	return r.db.Create(appointment).Error
}

func (r *appointmentRepo) GetAppointmentByID(id int) (*model.Appointment, error) {
	var appointment model.Appointment
	err := r.db.Preload("Patient").Preload("Chamber").First(&appointment, id).Error

	return &appointment, err
}

func (r *appointmentRepo) UpdateAppointment(appointment *model.Appointment) error {
	return r.db.Save(appointment).Error
}

func (r *appointmentRepo) ListAppointments(doctorID int, dateFrom, dateTo *time.Time, status string, page, limit int) ([]model.Appointment, int64, error) {
	var appointments []model.Appointment
	var total int64
	query := r.db.Model(&model.Appointment{}).Where("doctor_id = ?", doctorID)

	if dateFrom != nil {
		query = query.Where("appointment_date >= ?", dateFrom)
	}
	if dateTo != nil {
		query = query.Where("appointment_date <= ?", dateTo)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	err := query.Preload("Patient").Preload("Chamber").
		Order("appointment_date asc, start_time asc").
		Offset(offset).Limit(limit).
		Find(&appointments).Error

	return appointments, total, err
}

func (r *appointmentRepo) GetTodayAppointments(doctorID int) ([]model.Appointment, error) {
	var appointments []model.Appointment
	today := time.Now().Truncate(24 * time.Hour)
	err := r.db.Preload("Patient").
		Where("doctor_id = ? AND appointment_date = ? AND status != ?", doctorID, today, model.AppointmentStatusCancelled).
		Order("start_time asc").
		Find(&appointments).Error

	return appointments, err
}
