package types

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type PatientReq struct {
	FullName       string   `json:"full_name"`
	Gender         string   `json:"gender"`
	Age            int16    `json:"age"`
	Phone          string   `json:"phone"`
	Email          string   `json:"email"`
	BloodGroup     string   `json:"blood_group"`
	Allergies      []string `json:"allergies"`
	MedicalHistory string   `json:"medical_history"`
}

func (req PatientReq) Validate() error {
	return validation.ValidateStruct(&req,
		validation.Field(&req.FullName, validation.Required, validation.Length(2, 150)),
		validation.Field(&req.Gender, validation.Required, validation.In("male", "female", "other")),
		validation.Field(&req.Age, validation.Required, validation.Min(0), validation.Max(150)),
		validation.Field(&req.Phone, validation.Length(5, 20)),
		validation.Field(&req.Email, is.Email),
		validation.Field(&req.BloodGroup, validation.In("A+", "A-", "B+", "B-", "AB+", "AB-", "O+", "O-")),
	)
}

type PatientUpdateReq struct {
	ID             int      `json:"-" param:"id"`
	FullName       string   `json:"full_name"`
	Gender         string   `json:"gender"`
	Age            int16    `json:"age"`
	Phone          string   `json:"phone"`
	Email          string   `json:"email"`
	BloodGroup     string   `json:"blood_group"`
	Allergies      []string `json:"allergies"`
	MedicalHistory string   `json:"medical_history"`
}

func (req PatientUpdateReq) Validate() error {
	return validation.ValidateStruct(&req,
		validation.Field(&req.ID, validation.Required),
		validation.Field(&req.FullName, validation.Length(2, 150)),
		validation.Field(&req.Gender, validation.In("male", "female", "other")),
		validation.Field(&req.Age, validation.Min(0), validation.Max(150)),
		validation.Field(&req.Phone, validation.Length(5, 20)),
		validation.Field(&req.Email, is.Email),
		validation.Field(&req.BloodGroup, validation.In("A+", "A-", "B+", "B-", "AB+", "AB-", "O+", "O-")),
	)
}

type PatientFilter struct {
	ID int `json:"-" param:"id"`
}

func (f PatientFilter) Validate() error {
	return validation.ValidateStruct(&f,
		validation.Field(&f.ID, validation.Required),
	)
}

type PatientListReq struct {
	Pagination
	Name  string `json:"name" query:"name"`
	Phone string `json:"phone" query:"phone"`
}

type PatientResp struct {
	ID             int        `json:"id"`
	DoctorID       int        `json:"doctor_id"`
	FullName       string     `json:"full_name"`
	Gender         string     `json:"gender"`
	Age            int16      `json:"age"`
	Phone          string     `json:"phone"`
	Email          string     `json:"email"`
	BloodGroup     string     `json:"blood_group"`
	Allergies      []string   `json:"allergies"`
	MedicalHistory string     `json:"medical_history"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      *time.Time `json:"updated_at"`
}

type PaginatedPatientResp struct {
	Pagination
	Records []PatientResp `json:"records"`
}
