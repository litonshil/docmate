package model

import (
	"context"
	"docmate/types"
	"time"
)

type Patient struct {
	ID             int        `json:"id"`
	DoctorID       int        `json:"doctor_id"`
	FullName       string     `json:"full_name"`
	Gender         string     `json:"gender"`
	Age            int16      `json:"age"`
	Phone          string     `json:"phone"`
	Email          string     `json:"email"`
	BloodGroup     string     `json:"blood_group"`
	Allergies      []string   `json:"allergies" gorm:"type:text[]"`
	MedicalHistory string     `json:"medical_history"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      *time.Time `json:"updated_at"`
	DeletedAt      *time.Time `json:"deleted_at"`
}

type PatientUseCase interface {
	Create(ctx context.Context, req types.PatientReq, doctorID int) (types.PatientResp, error)
	Get(ctx context.Context, filter types.PatientFilter) (types.PatientResp, error)
	Update(ctx context.Context, req types.PatientUpdateReq) (types.PatientResp, error)
	List(ctx context.Context, req types.PatientListReq, doctorID int) (types.PaginatedPatientResp, error)
}

type PatientRepo interface {
	CreatePatient(patient Patient) (Patient, error)
	UpdatePatient(patient Patient) (Patient, error)
	GetPatientByID(id int) (Patient, error)
	ListPatients(offset, limit, doctorID int, name, phone string) ([]Patient, int, error)
}
