package types

import "time"

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

type QuickPatient struct {
	FullName string `json:"full_name"`
	Phone    string `json:"phone"`
	Gender   string `json:"gender"`
	Age      int16  `json:"age"`
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
