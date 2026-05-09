package model

import (
	"context"
	"docmate/types"
	"time"

	"gorm.io/gorm"
)

type AppointmentStatus string

const (
	AppointmentStatusPending   AppointmentStatus = "pending"
	AppointmentStatusConfirmed AppointmentStatus = "confirmed"
	AppointmentStatusCancelled AppointmentStatus = "cancelled"
	AppointmentStatusCompleted AppointmentStatus = "completed"
)

type Appointment struct {
	ID              int               `json:"id" gorm:"primaryKey"`
	DoctorID        int               `json:"doctor_id" gorm:"not null;index"`
	PatientID       int               `json:"patient_id" gorm:"not null;index"`
	ChamberID       int               `json:"chamber_id" gorm:"not null;index"`
	AppointmentDate time.Time         `json:"appointment_date" gorm:"not null;type:date;index"`
	StartTime       string            `json:"start_time" gorm:"not null;type:varchar(10)"` // HH:MM AM/PM
	EndTime         string            `json:"end_time" gorm:"type:varchar(10)"`
	Status          AppointmentStatus `json:"status" gorm:"type:varchar(20);default:'pending';index"`
	Reason          string            `json:"reason" gorm:"type:text"`
	Notes           string            `json:"notes" gorm:"type:text"`
	CreatedAt       time.Time         `json:"created_at"`
	UpdatedAt       time.Time         `json:"updated_at"`
	DeletedAt       gorm.DeletedAt    `json:"-" gorm:"index"`

	// Associations
	Doctor  Doctor  `json:"doctor,omitempty" gorm:"foreignKey:DoctorID"`
	Patient Patient `json:"patient,omitempty" gorm:"foreignKey:PatientID"`
	Chamber Chamber `json:"chamber,omitempty" gorm:"foreignKey:ChamberID"`
}

type AppointmentRepo interface {
	CreateAppointment(appointment *Appointment) error
	GetAppointmentByID(id int) (*Appointment, error)
	UpdateAppointment(appointment *Appointment) error
	ListAppointments(doctorID int, dateFrom, dateTo *time.Time, status string, page, limit int) ([]Appointment, int64, error)
	GetTodayAppointments(doctorID int) ([]Appointment, error)
}

type AppointmentUseCase interface {
	BookAppointment(ctx context.Context, req types.AppointmentReq, doctorID int) (types.AppointmentResp, error)
	UpdateStatus(ctx context.Context, id int, status AppointmentStatus) error
	GetAppointment(ctx context.Context, id int) (types.AppointmentResp, error)
	ListAppointments(ctx context.Context, doctorID int, dateFrom, dateTo string, status string, page, limit int) (types.PaginatedResponse[types.AppointmentResp], error)
}
