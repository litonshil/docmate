package model

import (
	"context"
	"docmate/types"
	"time"
)

type PrescriptionSetting struct {
	ID                 int        `json:"id"`
	DoctorID           int        `json:"doctor_id"`
	ChamberID          int        `json:"chamber_id"`
	HeaderLeftBangla   string     `json:"header_left_bangla"`
	HeaderRightEnglish string     `json:"header_right_english"`
	FooterInfoBangla   string     `json:"footer_info_bangla"`
	FooterInfoEnglish  string     `json:"footer_info_english"`
	TemplateType       string     `json:"template_type"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
	DeletedAt          *time.Time `json:"deleted_at"`
}

type PrescriptionSettingUseCase interface {
	Upsert(ctx context.Context, req types.PrescriptionSettingReq) (types.PrescriptionSettingResp, error)
	GetByChamber(ctx context.Context, doctorID, chamberID int) (types.PrescriptionSettingResp, error)
}

type PrescriptionSettingRepo interface {
	Upsert(setting PrescriptionSetting) (PrescriptionSetting, error)
	GetByChamber(doctorID, chamberID int) (PrescriptionSetting, error)
}
