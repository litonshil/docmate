package ai_setting

import (
	"context"
	"docmate/internal/llm"
	"docmate/internal/model"
	"docmate/types"
	"errors"
	"log/slog"
	"time"

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
	// First fetch existing settings to preserve admin flags
	existing, err := s.repo.GetDoctorSettings(req.DoctorID)
	var isAIEnabled bool = false
	var allowGlobalAPI bool = false
	var useIndividualKey bool = false

	if err == nil && len(existing) > 0 {
		// Use flags from the first existing row (as they are identical across provider rows)
		isAIEnabled = existing[0].IsAIEnabled
		allowGlobalAPI = existing[0].AllowGlobalAPI
		useIndividualKey = existing[0].UseIndividualKey
	}

	// Fetch all providers to validate IDs and map slugs
	providers, err := s.repo.GetProviders()
	if err != nil {
		slog.Error("failed to get active providers list", "error", err.Error())

		return types.AISettingResp{}, err
	}

	// Deactivate active status on all rows for this doctor before setting the new active provider
	err = s.repo.DeactivateAllDoctorSettings(req.DoctorID)
	if err != nil {
		slog.Error("failed to deactivate doctor settings", "error", err.Error())

		return types.AISettingResp{}, err
	}

	// Upsert a setting row for each globally active provider
	for _, prov := range providers {
		if !prov.IsActive {
			continue
		}

		apiKey := req.ProviderKeys[prov.Slug]
		// If key not in request, but exists in previous DB configuration, preserve it
		if apiKey == "" {
			for _, ext := range existing {
				if ext.AIProvidersID == prov.ID {
					apiKey = ext.IndividualAPIKey

					break
				}
			}
		}

		isActive := prov.ID == req.AIProviderID

		row := model.DoctorAISetting{
			DoctorID:         req.DoctorID,
			AIProvidersID:    prov.ID,
			IndividualAPIKey: apiKey,
			IsActive:         isActive,
			IsAIEnabled:      isAIEnabled,
			AllowGlobalAPI:   allowGlobalAPI,
			UseIndividualKey: useIndividualKey,
			CreatedAt:        time.Now(),
			UpdatedAt:        time.Now(),
		}

		_, err = s.repo.UpsertDoctorSetting(row)
		if err != nil {
			slog.Error("failed to upsert doctor setting", "provider", prov.Slug, "error", err.Error())

			return types.AISettingResp{}, err
		}
	}

	return s.GetByDoctor(ctx, req.DoctorID)
}

func (s *Service) AdminUpdate(ctx context.Context, req types.AdminAISettingUpdateReq) (types.AISettingResp, error) {
	err := s.repo.AdminUpdateDoctorSettings(req.DoctorID, req.IsAIEnabled, req.AllowGlobalAPI, req.UseIndividualKey)
	if err != nil {
		slog.Error("failed to admin update doctor settings", "error", err.Error())

		return types.AISettingResp{}, err
	}

	return s.GetByDoctor(ctx, req.DoctorID)
}

func (s *Service) GetByDoctor(ctx context.Context, doctorID int) (types.AISettingResp, error) {
	settings, err := s.repo.GetDoctorSettings(doctorID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		slog.Error("failed to get doctor settings", "doctor_id", doctorID, "error", err.Error())

		return types.AISettingResp{}, err
	}

	providers, err := s.repo.GetProviders()
	if err != nil {
		slog.Error("failed to fetch providers", "error", err.Error())

		return types.AISettingResp{}, err
	}

	// Map providers for easy lookup
	provMap := make(map[int]model.AIProvider)
	var activeProviders []types.ActiveProviderInfo
	var defaultProviderID int
	var defaultProviderSlug string

	for _, p := range providers {
		provMap[p.ID] = p
		if p.IsActive {
			activeProviders = append(activeProviders, types.ActiveProviderInfo{
				ID:    p.ID,
				Name:  p.Name,
				Slug:  p.Slug,
				Model: p.Model,
			})
			if defaultProviderID == 0 {
				defaultProviderID = p.ID
				defaultProviderSlug = p.Slug
			}
		}
	}

	// Build individual keys dictionary
	keys := make(map[string]string)
	var activeSetting *model.DoctorAISetting

	for _, s := range settings {
		p, exists := provMap[s.AIProvidersID]
		if !exists || !p.IsActive {
			continue
		}
		keys[p.Slug] = s.IndividualAPIKey
		if s.IsActive {
			// Save reference to active setting row
			activeSetting = &s
		}
	}

	// If doctor has no settings yet
	if len(settings) == 0 || activeSetting == nil {
		return types.AISettingResp{
			DoctorID:         doctorID,
			IsAIEnabled:      false,
			AllowGlobalAPI:   false,
			UseIndividualKey: false,
			AIProviderID:     defaultProviderID,
			ProviderSlug:     defaultProviderSlug,
			ProviderKeys:     keys,
			ActiveProviders:  activeProviders,
		}, nil
	}

	pSlug := ""
	if p, ok := provMap[activeSetting.AIProvidersID]; ok {
		pSlug = p.Slug
	}

	return types.AISettingResp{
		ID:               activeSetting.ID,
		DoctorID:         activeSetting.DoctorID,
		IsAIEnabled:      activeSetting.IsAIEnabled,
		AllowGlobalAPI:   activeSetting.AllowGlobalAPI,
		AIProviderID:     activeSetting.AIProvidersID,
		ProviderSlug:     pSlug,
		UseIndividualKey: activeSetting.UseIndividualKey,
		ProviderKeys:     keys,
		ActiveProviders:  activeProviders,
		CreatedAt:        activeSetting.CreatedAt,
		UpdatedAt:        activeSetting.UpdatedAt,
	}, nil
}

func (s *Service) GetSuggestions(ctx context.Context, doctorID int, complaints []string) (*types.AISuggestionResp, error) {
	settings, err := s.repo.GetDoctorSettings(doctorID)
	if err != nil || len(settings) == 0 {
		slog.Warn("AI assistance lookup failed", "doctor_id", doctorID)

		return nil, errors.New("AI assistance is not enabled for this doctor")
	}

	// Find the active row
	var activeSetting *model.DoctorAISetting
	for i := range settings {
		if settings[i].IsActive {
			activeSetting = &settings[i]

			break
		}
	}

	// Fallback to first row if none marked active
	if activeSetting == nil {
		activeSetting = &settings[0]
	}

	if !activeSetting.IsAIEnabled {
		return nil, errors.New("AI assistance is locked for this account")
	}

	// Fetch provider information
	prov, err := s.repo.GetProviderByID(activeSetting.AIProvidersID)
	if err != nil {
		slog.Error("failed to get active provider detail", "id", activeSetting.AIProvidersID, "error", err.Error())

		return nil, errors.New("AI provider lookup failed")
	}

	if !prov.IsActive {
		return nil, errors.New("the selected AI provider is currently disabled by administrator")
	}

	var apiKey string

	if activeSetting.UseIndividualKey {
		if activeSetting.IndividualAPIKey == "" {
			slog.Warn("individual API key required but missing", "doctor_id", doctorID)

			return nil, errors.New("individual API key required but missing")
		}
		apiKey = activeSetting.IndividualAPIKey
	} else if activeSetting.AllowGlobalAPI {
		if prov.APIKey == "" {
			slog.Warn("global API key required but missing for provider", "provider", prov.Slug)

			return nil, errors.New("system global API key is not configured for this provider")
		}
		apiKey = prov.APIKey
	} else {
		return nil, errors.New("AI configuration error: no API access allowed")
	}

	provider, err := s.llmFactory.GetProvider(prov.Slug)
	if err != nil {
		slog.Error("failed to resolve llm provider", "slug", prov.Slug, "error", err.Error())

		return nil, err
	}

	slog.Info("requesting AI suggestions", "doctor_id", doctorID, "provider", prov.Slug, "model", prov.Model)

	return provider.GenerateSuggestions(ctx, apiKey, prov.Model, complaints)
}

func (s *Service) GetProviders(ctx context.Context) ([]types.AIProviderConfig, error) {
	provs, err := s.repo.GetProviders()
	if err != nil {
		return nil, err
	}

	res := make([]types.AIProviderConfig, len(provs))
	for i, p := range provs {
		res[i] = types.AIProviderConfig{
			ID:       p.ID,
			Name:     p.Name,
			Slug:     p.Slug,
			APIKey:   p.APIKey,
			Model:    p.Model,
			IsActive: p.IsActive,
		}
	}

	return res, nil
}

func (s *Service) UpdateProviders(ctx context.Context, req []types.AIProviderConfig) error {
	payload := make([]model.AIProvider, len(req))
	for i, p := range req {
		payload[i] = model.AIProvider{
			ID:       p.ID,
			Name:     p.Name,
			Slug:     p.Slug,
			APIKey:   p.APIKey,
			Model:    p.Model,
			IsActive: p.IsActive,
		}
	}

	return s.repo.UpdateProviders(payload)
}

func (s *Service) GetActiveProviders(ctx context.Context) ([]string, error) {
	provs, err := s.repo.GetProviders()
	if err != nil {
		return nil, err
	}

	var active []string
	for _, p := range provs {
		if p.IsActive {
			active = append(active, p.Slug)
		}
	}

	return active, nil
}
