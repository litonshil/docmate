package types

import "time"

// Vitals represents patient vital signs
type Vitals struct {
	WeightKg      *float64 `json:"weight_kg,omitempty"`
	HeightCm      *float64 `json:"height_cm,omitempty"`
	BloodPressure string   `json:"blood_pressure,omitempty"`
	TemperatureF  *float64 `json:"temperature_f,omitempty"`
	PulseBpm      *int     `json:"pulse_bpm,omitempty"`
	SpO2Percent   *int     `json:"spo2_percent,omitempty"`
}

// Medication defines a single prescribed medicine
type Medication struct {
	MedicineID   *int   `json:"medicine_id,omitempty"`
	MedicineName string `json:"medicine_name" validate:"required"`
	GenericName  string `json:"generic_name"`
	Form         string `json:"form"`
	Strength     string `json:"strength"`
	Dosage       string `json:"dosage" validate:"required"`
	Frequency    string `json:"frequency" validate:"required"`
	Timing       string `json:"timing"`
	Duration     string `json:"duration" validate:"required"`
	Instructions string `json:"instructions"`
	SortOrder    int    `json:"sort_order"`
}

type PrescriptionReq struct {
	DoctorID        int          `json:"-"` // Set by auth token
	PatientID       int          `json:"patient_id" validate:"required"`
	ChamberID       int          `json:"chamber_id" validate:"required"`
	Vitals          Vitals       `json:"vitals"`
	ChiefComplaints []string     `json:"chief_complaints"`
	Diagnosis       []string     `json:"diagnosis"`
	Medications     []Medication `json:"medications"`
	Investigations  []string     `json:"investigations"`
	Advice          string       `json:"advice"`
	Status          string       `json:"status"`
	FollowUpDate    *time.Time   `json:"follow_up_date"`
}

type PrescriptionResp struct {
	ID              int          `json:"id"`
	DoctorID        int          `json:"doctor_id"`
	PatientID       int          `json:"patient_id"`
	ChamberID       int          `json:"chamber_id"`
	Vitals          Vitals       `json:"vitals"`
	ChiefComplaints []string     `json:"chief_complaints"`
	Diagnosis       []string     `json:"diagnosis"`
	Medications     []Medication `json:"medications"`
	Investigations  []string     `json:"investigations"`
	Advice          string       `json:"advice"`
	Status          string       `json:"status"`
	FollowUpDate    *time.Time   `json:"follow_up_date,omitempty"`
	CreatedAt       time.Time    `json:"created_at"`
	UpdatedAt       time.Time    `json:"updated_at"`
}

type PrescriptionListReq struct {
	Page      int `query:"page"`
	Limit     int `query:"limit"`
	DoctorID  int // Internal
	PatientID int `query:"patient_id"`
}

type PaginatedPrescriptionResp struct {
	Pagination Pagination         `json:"pagination"`
	Records    []PrescriptionResp `json:"records"`
}
