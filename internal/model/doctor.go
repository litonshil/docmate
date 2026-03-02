package model

import (
	"context"
	"docmate/types"
	"time"
)

type Doctor struct {
	ID             int        `json:"id"`
	UserID         int        `json:"user_id"`
	Email          string     `json:"email"`
	FullName       string     `json:"full_name"`
	Degree         []string   `json:"degree" gorm:"type:degree_type[]"`
	Specialization []string   `json:"specialization" gorm:"type:specialization_type[]"`
	Phone          string     `json:"phone"`
	Bio            string     `json:"bio"`
	SignatureURL   string     `json:"signature_url"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      *time.Time `json:"updated_at"`
	DeletedAt      *time.Time `json:"deleted_at"`
}

type DoctorUseCase interface {
	Create(ctx context.Context, req types.DoctorReq) (types.DoctorResp, error)
	Get(ctx context.Context, id int) (types.DoctorResp, error)
	Update(ctx context.Context, id int, req types.DoctorReq) (types.DoctorResp, error)
	List(ctx context.Context, req types.DoctorListReq) (types.PaginatedDoctorResp, error)
}

type DoctorRepo interface {
	CreateDoctor(doctor Doctor) (Doctor, error)
	UpdateDoctor(doctor Doctor) (Doctor, error)
	GetDoctorByID(id int) (Doctor, error)
	ListDoctors(offset, limit int) ([]Doctor, int, error)

	GetDoctorByUserID(userID int) (Doctor, error)
	UpsertDoctor(doctor Doctor) (Doctor, error)
}
