package model

import (
	"context"
	"docmate/types"
	"time"

	"gorm.io/datatypes"
)

type Prescription struct {
	ID              int            `json:"id"`
	DoctorID        int            `json:"doctor_id"`
	PatientID       int            `json:"patient_id"`
	PatientName     string         `json:"patient_name" gorm:"->"`
	ChamberID       int            `json:"chamber_id"`
	Vitals          datatypes.JSON `json:"vitals" gorm:"type:jsonb;default:'{}'"`
	ChiefComplaints datatypes.JSON `json:"chief_complaints" gorm:"type:jsonb"`
	Diagnosis       datatypes.JSON `json:"diagnosis" gorm:"type:jsonb"`
	Medications     datatypes.JSON `json:"medications" gorm:"type:jsonb"`
	Investigations  datatypes.JSON `json:"investigations" gorm:"type:jsonb"`
	Advice          string         `json:"advice"`
	Status          string         `json:"status"`
	FollowUpDate    *time.Time     `json:"follow_up_date"`
	FileUrl         *string        `json:"file_url"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       *time.Time     `json:"updated_at"`
	DeletedAt       *time.Time     `json:"deleted_at"`
}

type PrescriptionUseCase interface {
	Create(ctx context.Context, req types.PrescriptionReq) (types.PrescriptionResp, error)
	Get(ctx context.Context, id int, doctorID int) (types.PrescriptionResp, error)
	Update(ctx context.Context, id int, req types.PrescriptionReq) (types.PrescriptionResp, error)
	List(ctx context.Context, req types.PrescriptionListReq) (types.PaginatedPrescriptionResp, error)
}

type PrescriptionRepo interface {
	CreatePrescription(p Prescription) (Prescription, error)
	GetPrescriptionByID(id int) (Prescription, error)
	UpdatePrescription(p Prescription) (Prescription, error)
	ListPrescriptions(doctorID int, limit, offset int, search string) ([]Prescription, int, error)
	ListPrescriptionsByPatient(doctorID, patientID int, limit, offset int, search string) ([]Prescription, int, error)
}
