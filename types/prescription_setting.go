package types

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type PrescriptionSettingReq struct {
	DoctorID           int    `json:"-" param:"doctor_id"`
	ChamberID          int    `json:"chamber_id"`
	HeaderLeftBangla   string `json:"header_left_bangla"`
	HeaderRightEnglish string `json:"header_right_english"`
	FooterInfoBangla   string `json:"footer_info_bangla"`
	FooterInfoEnglish  string `json:"footer_info_english"`
	TemplateType       string `json:"template_type"`
}

func (req PrescriptionSettingReq) Validate() error {
	return validation.ValidateStruct(&req,
		validation.Field(&req.DoctorID, validation.Required),
		validation.Field(&req.ChamberID, validation.Required),
		validation.Field(&req.TemplateType, validation.Required, validation.In("standard", "modern")),
	)
}

type PrescriptionSettingResp struct {
	ID                 int       `json:"id"`
	DoctorID           int       `json:"doctor_id"`
	ChamberID          int       `json:"chamber_id"`
	HeaderLeftBangla   string    `json:"header_left_bangla"`
	HeaderRightEnglish string    `json:"header_right_english"`
	FooterInfoBangla   string    `json:"footer_info_bangla"`
	FooterInfoEnglish  string    `json:"footer_info_english"`
	TemplateType       string    `json:"template_type"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}
