package types

import (
	"regexp"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type AppointmentReq struct {
	PatientID       int           `json:"patient_id"`
	ChamberID       int           `json:"chamber_id"`
	AppointmentDate string        `json:"appointment_date"` // YYYY-MM-DD
	StartTime       string        `json:"start_time"`
	EndTime         string        `json:"end_time"`
	Reason          string        `json:"reason"`
	Notes           string        `json:"notes"`
	QuickPatient    *QuickPatient `json:"quick_patient"` // Optional for Quick Add
}

func (req AppointmentReq) Validate() error {
	return validation.ValidateStruct(&req,
		validation.Field(&req.ChamberID, validation.Required),
		validation.Field(&req.AppointmentDate, validation.Required, validation.Match(regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`))),
		validation.Field(&req.StartTime, validation.Required),
		validation.Field(&req.PatientID, validation.Required.When(req.QuickPatient == nil).Error("patient_id is required when quick_patient is not provided")),
		validation.Field(&req.QuickPatient),
	)
}

type QuickPatient struct {
	FullName string `json:"full_name"`
	Phone    string `json:"phone"`
	Gender   string `json:"gender"`
	Age      int16  `json:"age"`
}

func (q QuickPatient) Validate() error {
	return validation.ValidateStruct(&q,
		validation.Field(&q.FullName, validation.Required, validation.Length(2, 100)),
		validation.Field(&q.Phone, validation.Required, validation.Match(regexp.MustCompile(`^\d{11}$`)).Error("phone must be exactly 11 digits")),
		validation.Field(&q.Gender, validation.Required, validation.In("male", "female", "other")),
		validation.Field(&q.Age, validation.Required, validation.Min(0), validation.Max(120)),
	)
}

type AppointmentResp struct {
	ID              int          `json:"id"`
	DoctorID        int          `json:"doctor_id"`
	PatientID       int          `json:"patient_id"`
	ChamberID       int          `json:"chamber_id"`
	AppointmentDate time.Time    `json:"appointment_date"`
	StartTime       string       `json:"start_time"`
	EndTime         string       `json:"end_time"`
	Status          string       `json:"status"`
	Reason          string       `json:"reason"`
	Notes           string       `json:"notes"`
	CreatedAt       time.Time    `json:"created_at"`
	UpdatedAt       time.Time    `json:"updated_at"`
	Patient         *PatientResp `json:"patient,omitempty"`
	Chamber         *ChamberResp `json:"chamber,omitempty"`
}

type UpdateAppointmentStatusReq struct {
	Status string `json:"status"`
}
