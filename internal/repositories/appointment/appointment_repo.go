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

func (r *appointmentRepo) UpdateAppointmentFields(id int, updates map[string]interface{}) error {
	return r.db.Model(&model.Appointment{}).Where("id = ?", id).Updates(updates).Error
}

func (r *appointmentRepo) ListAppointments(doctorID int, dateFrom, dateTo *time.Time, status string, search string, page, limit int) ([]model.Appointment, int64, error) {
	var appointments []model.Appointment
	var total int64
	query := r.db.Model(&model.Appointment{}).Where("doctor_id = ?", doctorID)

	if search != "" {
		var patientIDs []int
		searchPattern := "%" + search + "%"
		r.db.Model(&model.Patient{}).
			Where("doctor_id = ? AND (full_name ILIKE ? OR phone ILIKE ?)", doctorID, searchPattern, searchPattern).
			Pluck("id", &patientIDs)

		if len(patientIDs) == 0 {
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

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	q := r.db.Model(&model.Appointment{}).
		Joins("JOIN patients ON patients.id = appointments.patient_id").
		Where("appointments.doctor_id = ?", doctorID).
		Preload("Patient").Preload("Chamber")

	if search != "" {
		searchPattern := "%" + search + "%"
		q = q.Where("(patients.full_name ILIKE ? OR patients.phone LIKE ?)", searchPattern, searchPattern)
	}

	if status == "pending" {
		q = q.Order("appointments.appointment_date asc, appointments.start_time asc")
	} else {
		q = q.Order("appointments.created_at desc")
	}

	err := q.Offset(offset).Limit(limit).Find(&appointments).Error

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
