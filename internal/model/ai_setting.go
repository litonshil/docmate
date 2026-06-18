package model

import (
	"context"
	"docmate/types"
	"time"
)

type AISetting struct {
	ID               int        `json:"id" gorm:"primaryKey"`
	DoctorID         int        `json:"doctor_id" gorm:"uniqueIndex"`
	IsAIEnabled      bool       `json:"is_ai_enabled"`
	AllowGlobalAPI   bool       `json:"allow_global_api"`
	Provider         string     `json:"provider"`
	IndividualAPIKey string     `json:"individual_api_key"`
	UseIndividualKey bool       `json:"use_individual_key"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
	DeletedAt        *time.Time `json:"deleted_at"`
}

type AISettingUseCase interface {
	Upsert(ctx context.Context, req types.AISettingReq) (types.AISettingResp, error)
	AdminUpdate(ctx context.Context, req types.AdminAISettingUpdateReq) (types.AISettingResp, error)
	GetByDoctor(ctx context.Context, doctorID int) (types.AISettingResp, error)
	GetSuggestions(ctx context.Context, doctorID int, complaints []string) (*types.AISuggestionResp, error)
	GetGlobalSetting(ctx context.Context, key string) (string, error)
	SetGlobalSetting(ctx context.Context, key string, value string) error
}

type AISettingRepo interface {
	UpsertAISetting(setting AISetting) (AISetting, error)
	AdminUpdateAISetting(setting AISetting) (AISetting, error)
	GetAISettingByDoctor(doctorID int) (AISetting, error)
	GetGlobalSetting(key string) (string, error)
	SetGlobalSetting(key string, value string) error
}
