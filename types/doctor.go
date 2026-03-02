package types

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type DoctorReq struct {
	FullName       string   `json:"full_name"`
	Degree         []string `json:"degree"`
	Specialization []string `json:"specialization"`
	Phone          string   `json:"phone"`
	Bio            string   `json:"bio"`
	SignatureURL   string   `json:"signature_url"`
}

func (req DoctorReq) Validate() error {
	return validation.ValidateStruct(&req,
		validation.Field(&req.FullName, validation.Required, validation.Length(2, 150)),
		validation.Field(&req.Degree, validation.Required),
		validation.Field(&req.Specialization, validation.Required),
	)
}

type DoctorResp struct {
	ID             int        `json:"id"`
	UserID         int        `json:"user_id"`
	Email          string     `json:"email"`
	FullName       string     `json:"full_name"`
	Degree         []string   `json:"degree"`
	Specialization []string   `json:"specialization"`
	Phone          string     `json:"phone"`
	Bio            string     `json:"bio"`
	SignatureURL   string     `json:"signature_url"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      *time.Time `json:"updated_at"`
}

type DoctorListReq struct {
	Pagination
}

type DoctorFilter struct {
	ID int `json:"id" query:"id" param:"id"`
}

func (f DoctorFilter) Validate() error {
	return validation.ValidateStruct(&f,
		validation.Field(&f.ID, validation.Required),
	)
}

type PaginatedDoctorResp struct {
	Pagination
	Records []DoctorResp `json:"records"`
}
