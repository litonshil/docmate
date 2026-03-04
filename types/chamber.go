package types

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type VisitingSlot struct {
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
}

type VisitingDay struct {
	Day   string         `json:"day"`
	Slots []VisitingSlot `json:"slots"`
}

type ChamberReq struct {
	DoctorID      int           `json:"-" param:"doctor_id"`
	Name          string        `json:"name"`
	Address       string        `json:"address"`
	Area          string        `json:"area"`
	City          string        `json:"city"`
	Country       string        `json:"country"`
	Phone         string        `json:"phone"`
	Email         string        `json:"email"`
	Fee           float64       `json:"fee"`
	FollowUpFee   *float64      `json:"follow_up_fee"`
	VisitingHours []VisitingDay `json:"visiting_hours"`
	IsActive      *bool         `json:"is_active"`
}

func (req ChamberReq) Validate() error {
	return validation.ValidateStruct(&req,
		validation.Field(&req.DoctorID, validation.Required),
		validation.Field(&req.Name, validation.Required, validation.Length(2, 255)),
		validation.Field(&req.Address, validation.Required, validation.Length(2, 500)),
		validation.Field(&req.City, validation.Required, validation.Length(2, 100)),
		validation.Field(&req.Country, validation.Length(2, 100)),
		validation.Field(&req.Phone, validation.Length(5, 20)),
		validation.Field(&req.Email, is.Email),
		validation.Field(&req.Fee, validation.Required, validation.Min(0.0)),
	)
}

type ChamberUpdateReq struct {
	ID            int           `json:"-" param:"id"`
	DoctorID      int           `json:"-" param:"doctor_id"`
	Name          string        `json:"name"`
	Address       string        `json:"address"`
	Area          string        `json:"area"`
	City          string        `json:"city"`
	Country       string        `json:"country"`
	Phone         string        `json:"phone"`
	Email         string        `json:"email"`
	Fee           float64       `json:"fee"`
	FollowUpFee   *float64      `json:"follow_up_fee"`
	VisitingHours []VisitingDay `json:"visiting_hours"`
	IsActive      *bool         `json:"is_active"`
}

func (req ChamberUpdateReq) Validate() error {
	return validation.ValidateStruct(&req,
		validation.Field(&req.ID, validation.Required),
		validation.Field(&req.DoctorID, validation.Required),
		validation.Field(&req.Name, validation.Length(2, 255)),
		validation.Field(&req.Address, validation.Length(2, 500)),
		validation.Field(&req.City, validation.Length(2, 100)),
		validation.Field(&req.Country, validation.Length(2, 100)),
		validation.Field(&req.Phone, validation.Length(5, 20)),
		validation.Field(&req.Email, is.Email),
		validation.Field(&req.Fee, validation.Min(0.0)),
	)
}

type ChamberFilter struct {
	ID       int `json:"-" param:"id"`
	DoctorID int `json:"-" param:"doctor_id"`
}

func (f ChamberFilter) Validate() error {
	return validation.ValidateStruct(&f,
		validation.Field(&f.ID, validation.Required),
		validation.Field(&f.DoctorID, validation.Required),
	)
}

type ChamberListReq struct {
	Pagination
	DoctorID int `json:"-" param:"doctor_id"`
}

type ChamberResp struct {
	ID            int           `json:"id"`
	DoctorID      int           `json:"doctor_id"`
	Name          string        `json:"name"`
	Address       string        `json:"address"`
	Area          string        `json:"area"`
	City          string        `json:"city"`
	Country       string        `json:"country"`
	Phone         string        `json:"phone"`
	Email         string        `json:"email"`
	Fee           float64       `json:"fee"`
	FollowUpFee   *float64      `json:"follow_up_fee"`
	VisitingHours []VisitingDay `json:"visiting_hours"`
	IsActive      bool          `json:"is_active"`
	CreatedAt     time.Time     `json:"created_at"`
	UpdatedAt     *time.Time    `json:"updated_at"`
}

type PaginatedChamberResp struct {
	Pagination
	Records []ChamberResp `json:"records"`
}
