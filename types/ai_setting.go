package types

import "time"

type AIProviderConfig struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Slug     string `json:"slug"`
	APIKey   string `json:"api_key"`
	Model    string `json:"model"`
	IsActive bool   `json:"is_active"`
}

type ActiveProviderInfo struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Slug  string `json:"slug"`
	Model string `json:"model"`
}

type AISettingReq struct {
	DoctorID         int               `json:"doctor_id"`
	AIProviderID     int               `json:"ai_provider_id"`
	UseIndividualKey bool              `json:"use_individual_key"`
	ProviderKeys     map[string]string `json:"provider_keys"`
}

type AISettingResp struct {
	ID               int                  `json:"id"`
	DoctorID         int                  `json:"doctor_id"`
	IsAIEnabled      bool                 `json:"is_ai_enabled"`
	AllowGlobalAPI   bool                 `json:"allow_global_api"`
	AIProviderID     int                  `json:"ai_provider_id"`
	ProviderSlug     string               `json:"provider_slug"`
	UseIndividualKey bool                 `json:"use_individual_key"`
	AIRequestStatus  *string              `json:"ai_request_status"`
	ProviderKeys     map[string]string    `json:"provider_keys"`
	ActiveProviders  []ActiveProviderInfo `json:"active_providers"`
	CreatedAt        time.Time            `json:"created_at"`
	UpdatedAt        time.Time            `json:"updated_at"`
}

type AISuggestionReq struct {
	Complaints []string `json:"complaints" validate:"required"`
}

type AISuggestionResp struct {
	Diagnoses      []string `json:"diagnoses"`
	Investigations []string `json:"investigations"`
	Disclaimer     string   `json:"disclaimer"`
}

type AdminAISettingUpdateReq struct {
	DoctorID         int     `json:"doctor_id"`
	IsAIEnabled      bool    `json:"is_ai_enabled"`
	AllowGlobalAPI   bool    `json:"allow_global_api"`
	UseIndividualKey bool    `json:"use_individual_key"`
	AIRequestStatus  *string `json:"ai_request_status"`
}
