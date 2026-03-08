package model

import (
	"context"
	"docmate/types"
	"time"

	"gorm.io/datatypes"
)

type Chamber struct {
	ID            int            `json:"id"`
	DoctorID      int            `json:"doctor_id"`
	Name          string         `json:"name"`
	Address       string         `json:"address"`
	Area          string         `json:"area"`
	City          string         `json:"city"`
	Country       string         `json:"country"`
	Phone         string         `json:"phone"`
	Email         string         `json:"email"`
	Fee           float64        `json:"fee"`
	FollowUpFee   *float64       `json:"follow_up_fee"`
	VisitingHours datatypes.JSON `json:"visiting_hours"`
	IsActive      bool           `json:"is_active"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     *time.Time     `json:"updated_at"`
	DeletedAt     *time.Time     `json:"deleted_at"`
}

type ChamberUseCase interface {
	Create(ctx context.Context, req types.ChamberReq) (types.ChamberResp, error)
	Get(ctx context.Context, filter types.ChamberFilter) (types.ChamberResp, error)
	Update(ctx context.Context, req types.ChamberUpdateReq) (types.ChamberResp, error)
	List(ctx context.Context, req types.ChamberListReq) (types.PaginatedChamberResp, error)
}

type ChamberRepo interface {
	CreateChamber(chamber Chamber) (Chamber, error)
	UpdateChamber(chamber Chamber) (Chamber, error)
	GetChamberByID(id int) (Chamber, error)
	ListChambers(offset, limit, doctorID int) ([]Chamber, int, error)
}
