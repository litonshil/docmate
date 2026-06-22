package db

import (
	"docmate/internal/model"
)

func (r *Repository) CreateAssistant(assistant model.Assistant, chamberIDs []int) (model.Assistant, error) {
	err := r.client.Create(&assistant).Error
	if err != nil {
		return model.Assistant{}, err
	}

	for _, cid := range chamberIDs {
		err = r.client.Exec("INSERT INTO assistant_chambers (assistant_id, chamber_id) VALUES (?, ?)", assistant.ID, cid).Error
		if err != nil {
			return model.Assistant{}, err
		}
	}

	return assistant, nil
}

func (r *Repository) UpdateAssistant(assistant model.Assistant, chamberIDs []int) (model.Assistant, error) {
	err := r.client.Save(&assistant).Error
	if err != nil {
		return model.Assistant{}, err
	}

	err = r.client.Exec("DELETE FROM assistant_chambers WHERE assistant_id = ?", assistant.ID).Error
	if err != nil {
		return model.Assistant{}, err
	}

	for _, cid := range chamberIDs {
		err = r.client.Exec("INSERT INTO assistant_chambers (assistant_id, chamber_id) VALUES (?, ?)", assistant.ID, cid).Error
		if err != nil {
			return model.Assistant{}, err
		}
	}

	return assistant, nil
}

func (r *Repository) GetAssistantByID(id int) (model.Assistant, error) {
	var assistant model.Assistant
	err := r.client.Preload("User").First(&assistant, id).Error

	return assistant, err
}

func (r *Repository) GetAssistantByUserID(userID int) (model.Assistant, error) {
	var assistant model.Assistant
	err := r.client.Preload("User").Where("user_id = ?", userID).First(&assistant).Error

	return assistant, err
}

func (r *Repository) ListAssistantsByDoctorID(doctorID int) ([]model.Assistant, error) {
	var assistants []model.Assistant
	err := r.client.Preload("User").Where("doctor_id = ?", doctorID).Order("created_at desc").Find(&assistants).Error

	return assistants, err
}

func (r *Repository) GetChambersByAssistantID(assistantID int) ([]model.Chamber, error) {
	var chambers []model.Chamber
	err := r.client.Table("chambers").
		Select("chambers.*").
		Joins("JOIN assistant_chambers ON assistant_chambers.chamber_id = chambers.id").
		Where("assistant_chambers.assistant_id = ?", assistantID).
		Find(&chambers).Error

	return chambers, err
}
