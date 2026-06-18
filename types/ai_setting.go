package types

import "time"

type AISettingReq struct {
	DoctorID         int    `json:"doctor_id"`
	Provider         string `json:"provider"`
	IndividualAPIKey string `json:"individual_api_key"`
	UseIndividualKey bool   `json:"use_individual_key"`
}

type AISettingResp struct {
	ID               int       `json:"id"`
	DoctorID         int       `json:"doctor_id"`
	IsAIEnabled      bool      `json:"is_ai_enabled"`
	AllowGlobalAPI   bool      `json:"allow_global_api"`
	Provider         string    `json:"provider"`
	IndividualAPIKey string    `json:"individual_api_key,omitempty"`
	UseIndividualKey bool      `json:"use_individual_key"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
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
	DoctorID         int  `json:"doctor_id"`
	IsAIEnabled      bool `json:"is_ai_enabled"`
	AllowGlobalAPI   bool `json:"allow_global_api"`
	UseIndividualKey bool `json:"use_individual_key"`
}

type GlobalSettingReq struct {
	Value string `json:"value"`
}

type GlobalSettingResp struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
