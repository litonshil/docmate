package ai_setting

import (
	"context"
	"docmate/config"
	"docmate/internal/llm"
	"docmate/internal/model"
	"docmate/types"
	"errors"
	"log/slog"

	"gorm.io/gorm"
)

type Service struct {
	repo       model.AISettingRepo
	llmFactory *llm.Factory
}

func NewService(repo model.AISettingRepo, llmFactory *llm.Factory) *Service {
	return &Service{
		repo:       repo,
		llmFactory: llmFactory,
	}
}

func (s *Service) Upsert(ctx context.Context, req types.AISettingReq) (types.AISettingResp, error) {
	payload := model.AISetting{
		DoctorID:         req.DoctorID,
		Provider:         req.Provider,
		IndividualAPIKey: req.IndividualAPIKey,
		UseIndividualKey: req.UseIndividualKey,
	}

	setting, err := s.repo.UpsertAISetting(payload)
	if err != nil {
		slog.Error("failed to upsert ai setting", "error", err.Error())

		return types.AISettingResp{}, err
	}

	return mapToResponse(setting), nil
}

func (s *Service) AdminUpdate(ctx context.Context, req types.AdminAISettingUpdateReq) (types.AISettingResp, error) {
	payload := model.AISetting{
		DoctorID:         req.DoctorID,
		IsAIEnabled:      req.IsAIEnabled,
		AllowGlobalAPI:   req.AllowGlobalAPI,
		UseIndividualKey: req.UseIndividualKey,
	}

	setting, err := s.repo.AdminUpdateAISetting(payload)
	if err != nil {
		slog.Error("failed to admin update ai setting", "error", err.Error())

		return types.AISettingResp{}, err
	}

	return mapToResponse(setting), nil
}

func (s *Service) GetByDoctor(ctx context.Context, doctorID int) (types.AISettingResp, error) {
	setting, err := s.repo.GetAISettingByDoctor(doctorID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			defaultSetting := model.AISetting{
				DoctorID:         doctorID,
				IsAIEnabled:      false,
				AllowGlobalAPI:   false,
				UseIndividualKey: false,
			}

			return mapToResponse(defaultSetting), nil
		}
		slog.Error("failed to get ai setting", "doctor_id", doctorID, "error", err.Error())

		return types.AISettingResp{}, err
	}

	return mapToResponse(setting), nil
}

func (s *Service) GetSuggestions(ctx context.Context, doctorID int, complaints []string) (*types.AISuggestionResp, error) {
	setting, err := s.repo.GetAISettingByDoctor(doctorID)
	if err != nil {
		slog.Warn("AI assistance lookup failed", "doctor_id", doctorID, "error", err.Error())

		return nil, errors.New("AI assistance is not enabled for this doctor")
	}

	if !setting.IsAIEnabled {
		return nil, errors.New("AI assistance is locked for this account")
	}

	var apiKey string
	var providerName string

	// Resolution logic
	if setting.UseIndividualKey && setting.IndividualAPIKey != "" {
		slog.Info("using individual API key", "doctor_id", doctorID)
		apiKey = setting.IndividualAPIKey
		providerName = setting.Provider
	} else if setting.AllowGlobalAPI {
		slog.Info("using global system API key", "doctor_id", doctorID)

		dbKey, keyErr := s.repo.GetGlobalSetting("ai_global_api_key")
		dbProvider, providerErr := s.repo.GetGlobalSetting("ai_global_provider")

		if keyErr == nil && dbKey != "" {
			apiKey = dbKey
		} else {
			apiKey = config.AI().GeminiAPIKey
		}

		if providerErr == nil && dbProvider != "" {
			providerName = dbProvider
		} else {
			providerName = config.AI().Provider
		}
	} else {
		slog.Warn("no credentials available", "doctor_id", doctorID)

		return nil, errors.New("individual API key required but missing")
	}

	if apiKey == "" {
		return nil, errors.New("no valid API key found for AI suggestions")
	}

	provider, err := s.llmFactory.GetProvider(providerName)
	if err != nil {
		slog.Error("failed to get provider", "provider", providerName, "error", err.Error())

		return nil, err
	}

	slog.Info("requesting AI suggestions", "doctor_id", doctorID, "complaints_count", len(complaints))

	return provider.GenerateSuggestions(ctx, apiKey, complaints)
}

func mapToResponse(setting model.AISetting) types.AISettingResp {
	return types.AISettingResp{
		ID:               setting.ID,
		DoctorID:         setting.DoctorID,
		IsAIEnabled:      setting.IsAIEnabled,
		AllowGlobalAPI:   setting.AllowGlobalAPI,
		Provider:         setting.Provider,
		IndividualAPIKey: setting.IndividualAPIKey,
		UseIndividualKey: setting.UseIndividualKey,
		CreatedAt:        setting.CreatedAt,
		UpdatedAt:        setting.UpdatedAt,
	}
}

func (s *Service) GetGlobalSetting(ctx context.Context, key string) (string, error) {
	val, err := s.repo.GetGlobalSetting(key)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", nil
		}

		return "", err
	}

	return val, nil
}

func (s *Service) SetGlobalSetting(ctx context.Context, key string, value string) error {
	return s.repo.SetGlobalSetting(key, value)
}
