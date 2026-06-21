package db

import (
	"docmate/internal/model"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (r *Repository) GetDoctorSettings(doctorID int) ([]model.DoctorAISetting, error) {
	var settings []model.DoctorAISetting
	err := r.client.Model(&model.DoctorAISetting{}).
		Where("doctor_id = ?", doctorID).
		Find(&settings).Error

	return settings, err
}

func (r *Repository) UpsertDoctorSetting(setting model.DoctorAISetting) (model.DoctorAISetting, error) {
	err := r.client.Model(&model.DoctorAISetting{}).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "doctor_id"}, {Name: "ai_providers_id"}},
			DoUpdates: clause.AssignmentColumns([]string{"individual_api_key", "is_active", "use_individual_key", "is_ai_enabled", "allow_global_api", "ai_request_status", "updated_at"}),
		}).Create(&setting).Error

	return setting, err
}

func (r *Repository) DeactivateAllDoctorSettings(doctorID int) error {
	return r.client.Model(&model.DoctorAISetting{}).
		Where("doctor_id = ?", doctorID).
		Update("is_active", false).Error
}

func (r *Repository) AdminUpdateDoctorSettings(doctorID int, isAIEnabled bool, allowGlobalAPI bool, useIndividualKey bool, requestStatus *string) error {
	// First check if any rows exist for this doctor
	var count int64
	err := r.client.Model(&model.DoctorAISetting{}).Where("doctor_id = ?", doctorID).Count(&count).Error
	if err != nil {
		return err
	}

	var statusVal *string = requestStatus
	if statusVal == nil {
		// Default behavior based on isAIEnabled
		if isAIEnabled {
			approvedStr := "approved"
			statusVal = &approvedStr
		} else {
			statusVal = nil
		}
	} else if *statusVal == "none" || *statusVal == "" {
		statusVal = nil
	}

	if count == 0 {
		// If no settings exist yet, create a default row for the doctor (using Gemini provider)
		var gemini model.AIProvider
		err = r.client.Model(&model.AIProvider{}).Where("slug = ?", "gemini").First(&gemini).Error
		if err == nil {
			defaultSetting := model.DoctorAISetting{
				DoctorID:         doctorID,
				AIProvidersID:    gemini.ID,
				IsActive:         true,
				IsAIEnabled:      isAIEnabled,
				AllowGlobalAPI:   allowGlobalAPI,
				UseIndividualKey: useIndividualKey,
				AIRequestStatus:  statusVal,
				CreatedAt:        time.Now(),
				UpdatedAt:        time.Now(),
			}

			return r.client.Create(&defaultSetting).Error
		}
	}

	// Update all existing rows
	updates := map[string]interface{}{
		"is_ai_enabled":      isAIEnabled,
		"allow_global_api":   allowGlobalAPI,
		"use_individual_key": useIndividualKey,
		"ai_request_status":  statusVal,
		"updated_at":         time.Now(),
	}

	return r.client.Model(&model.DoctorAISetting{}).
		Where("doctor_id = ?", doctorID).
		Updates(updates).Error
}

func (r *Repository) UpdateAIRequestStatus(doctorID int, status string) error {
	var count int64
	err := r.client.Model(&model.DoctorAISetting{}).Where("doctor_id = ?", doctorID).Count(&count).Error
	if err != nil {
		return err
	}

	var statusVal *string
	if status != "" && status != "none" {
		statusVal = &status
	}

	if count == 0 {
		var gemini model.AIProvider
		err = r.client.Model(&model.AIProvider{}).Where("slug = ?", "gemini").First(&gemini).Error
		if err == nil {
			defaultSetting := model.DoctorAISetting{
				DoctorID:         doctorID,
				AIProvidersID:    gemini.ID,
				IsActive:         true,
				IsAIEnabled:      false,
				AllowGlobalAPI:   false,
				UseIndividualKey: false,
				AIRequestStatus:  statusVal,
				CreatedAt:        time.Now(),
				UpdatedAt:        time.Now(),
			}

			return r.client.Create(&defaultSetting).Error
		}
	}

	return r.client.Model(&model.DoctorAISetting{}).
		Where("doctor_id = ?", doctorID).
		Updates(map[string]interface{}{
			"ai_request_status": statusVal,
			"updated_at":        time.Now(),
		}).Error
}

func (r *Repository) GetProviders() ([]model.AIProvider, error) {
	var providers []model.AIProvider
	err := r.client.Model(&model.AIProvider{}).Order("id asc").Find(&providers).Error

	return providers, err
}

func (r *Repository) UpdateProviders(providers []model.AIProvider) error {
	return r.client.Transaction(func(tx *gorm.DB) error {
		for _, p := range providers {
			p.UpdatedAt = time.Now()
			err := tx.Model(&model.AIProvider{}).
				Clauses(clause.OnConflict{
					Columns:   []clause.Column{{Name: "slug"}},
					DoUpdates: clause.AssignmentColumns([]string{"name", "api_key", "model", "is_active", "updated_at"}),
				}).Create(&p).Error
			if err != nil {
				return err
			}
		}

		return nil
	})
}

func (r *Repository) GetProviderByID(id int) (model.AIProvider, error) {
	var provider model.AIProvider
	err := r.client.Model(&model.AIProvider{}).Where("id = ?", id).First(&provider).Error

	return provider, err
}

func (r *Repository) GetProviderBySlug(slug string) (model.AIProvider, error) {
	var provider model.AIProvider
	err := r.client.Model(&model.AIProvider{}).Where("slug = ?", slug).First(&provider).Error

	return provider, err
}
