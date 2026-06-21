package model

import (
	"context"
	"docmate/types"
	"time"
)

type AIProvider struct {
	ID        int       `json:"id" gorm:"primaryKey;column:id"`
	Name      string    `json:"name" gorm:"column:name"`
	Slug      string    `json:"slug" gorm:"column:slug;unique"`
	APIKey    string    `json:"api_key" gorm:"column:api_key"`
	Model     string    `json:"model" gorm:"column:model"`
	IsActive  bool      `json:"is_active" gorm:"column:is_active"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at"`
}

func (AIProvider) TableName() string {
	return "ai_providers"
}

type DoctorAISetting struct {
	ID               int       `json:"id" gorm:"primaryKey;column:id"`
	DoctorID         int       `json:"doctor_id" gorm:"column:doctor_id;uniqueIndex:idx_doc_provider"`
	AIProvidersID    int       `json:"ai_providers_id" gorm:"column:ai_providers_id;uniqueIndex:idx_doc_provider"`
	IndividualAPIKey string    `json:"individual_api_key" gorm:"column:individual_api_key"`
	IsActive         bool      `json:"is_active" gorm:"column:is_active"`
	IsAIEnabled      bool      `json:"is_ai_enabled" gorm:"column:is_ai_enabled"`
	AllowGlobalAPI   bool      `json:"allow_global_api" gorm:"column:allow_global_api"`
	UseIndividualKey bool      `json:"use_individual_key" gorm:"column:use_individual_key"`
	AIRequestStatus  *string   `json:"ai_request_status" gorm:"column:ai_request_status;default:null"`
	CreatedAt        time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt        time.Time `json:"updated_at" gorm:"column:updated_at"`
}

func (DoctorAISetting) TableName() string {
	return "doctor_ai_settings"
}

type AISettingUseCase interface {
	Upsert(ctx context.Context, req types.AISettingReq) (types.AISettingResp, error)
	AdminUpdate(ctx context.Context, req types.AdminAISettingUpdateReq) (types.AISettingResp, error)
	GetByDoctor(ctx context.Context, doctorID int) (types.AISettingResp, error)
	GetSuggestions(ctx context.Context, doctorID int, complaints []string) (*types.AISuggestionResp, error)
	RequestActivation(ctx context.Context, doctorID int) (types.AISettingResp, error)

	GetProviders(ctx context.Context) ([]types.AIProviderConfig, error)
	UpdateProviders(ctx context.Context, req []types.AIProviderConfig) error
	GetActiveProviders(ctx context.Context) ([]string, error)
}

type AISettingRepo interface {
	GetDoctorSettings(doctorID int) ([]DoctorAISetting, error)
	UpsertDoctorSetting(setting DoctorAISetting) (DoctorAISetting, error)
	DeactivateAllDoctorSettings(doctorID int) error
	AdminUpdateDoctorSettings(doctorID int, isAIEnabled bool, allowGlobalAPI bool, useIndividualKey bool, requestStatus *string) error
	UpdateAIRequestStatus(doctorID int, status string) error

	GetProviders() ([]AIProvider, error)
	UpdateProviders(providers []AIProvider) error
	GetProviderByID(id int) (AIProvider, error)
	GetProviderBySlug(slug string) (AIProvider, error)
}
