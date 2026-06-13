package types

import (
	"regexp"
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
		validation.Field(&req.Phone, validation.Required, validation.Match(regexp.MustCompile(`^\d{11}$`)).Error("phone must be exactly 11 digits")),
		validation.Field(&req.Email, is.Email),
		validation.Field(&req.BloodGroup, validation.In("", "A+", "A-", "B+", "B-", "AB+", "AB-", "O+", "O-")),
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
		validation.Field(&req.FullName, validation.Required, validation.Length(2, 150)),
		validation.Field(&req.Gender, validation.Required, validation.In("male", "female", "other")),
		validation.Field(&req.Age, validation.Required, validation.Min(0), validation.Max(150)),
		validation.Field(&req.Phone, validation.Required, validation.Match(regexp.MustCompile(`^\d{11}$`)).Error("phone must be exactly 11 digits")),
		validation.Field(&req.Email, is.Email),
		validation.Field(&req.BloodGroup, validation.In("", "A+", "A-", "B+", "B-", "AB+", "AB-", "O+", "O-")),
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
	Search   string `json:"search" query:"search"`
	DoctorID int    `json:"doctor_id" query:"doctor_id"`
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
