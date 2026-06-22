package types

import (
	"regexp"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type AssistantReq struct {
	Name       string `json:"name"`
	Phone      string `json:"phone"`
	Email      string `json:"email"`
	Password   string `json:"password"`
	ChamberIDs []int  `json:"chamber_ids"`
}

func (req AssistantReq) Validate() error {
	return validation.ValidateStruct(&req,
		validation.Field(&req.Name, validation.Required, validation.Length(2, 100)),
		validation.Field(&req.Phone, validation.Required, validation.Match(regexp.MustCompile(`^[0-9]{11}$`)).Error("phone must be exactly 11 digits")),
		validation.Field(&req.Email, validation.Required, is.Email),
		validation.Field(&req.Password, validation.Required, validation.Length(8, 100).Error("password must be at least 8 characters")),
	)
}

type AssistantUpdateReq struct {
	Name       string `json:"name"`
	Phone      string `json:"phone"`
	Email      string `json:"email"`
	Password   string `json:"password,omitempty"`
	IsActive   *bool  `json:"is_active"`
	ChamberIDs []int  `json:"chamber_ids"`
}

func (req AssistantUpdateReq) Validate() error {
	return validation.ValidateStruct(&req,
		validation.Field(&req.Name, validation.Required, validation.Length(2, 100)),
		validation.Field(&req.Phone, validation.Required, validation.Match(regexp.MustCompile(`^[0-9]{11}$`)).Error("phone must be exactly 11 digits")),
		validation.Field(&req.Email, validation.Required, is.Email),
		validation.Field(&req.Password, validation.By(func(value interface{}) error {
			s, ok := value.(string)
			if !ok || s == "" {
				return nil
			}

			return validation.Validate(s, validation.Length(8, 100).Error("password must be at least 8 characters"))
		})),
		validation.Field(&req.IsActive, validation.Required),
	)
}

type AssistantResp struct {
	ID        int           `json:"id"`
	UserID    int           `json:"user_id"`
	DoctorID  int           `json:"doctor_id"`
	Name      string        `json:"name"`
	Phone     string        `json:"phone"`
	Email     string        `json:"email"`
	IsActive  bool          `json:"is_active"`
	Chambers  []ChamberResp `json:"chambers,omitempty"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
}
