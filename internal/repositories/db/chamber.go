package db

import (
	"docmate/internal/model"
)

func (r *Repository) CreateChamber(chamber model.Chamber) (model.Chamber, error) {
	err := r.client.Create(&chamber).Error

	return chamber, err
}

func (r *Repository) UpdateChamber(chamber model.Chamber) (model.Chamber, error) {
	// Omit CreatedAt and let GORM handle UpdatedAt automatically
	err := r.client.Model(&model.Chamber{}).Where("id = ?", chamber.ID).Omit("created_at").Updates(&chamber).Error

	return chamber, err
}

func (r *Repository) GetChamberByID(id int) (model.Chamber, error) {
	var chamber model.Chamber
	if err := r.dbClient(nil).Model(&model.Chamber{}).Where("id = ?", id).First(&chamber).Error; err != nil {
		return model.Chamber{}, err
	}

	return chamber, nil
}

func (r *Repository) ListChambers(offset, limit, doctorID int) ([]model.Chamber, int64, error) {
	var chambers []model.Chamber
	var total int64

	query := r.client.Model(&model.Chamber{})

	if doctorID > 0 {
		query = query.Where("doctor_id = ?", doctorID)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Offset(offset).Limit(limit).Find(&chambers).Error; err != nil {
		return nil, 0, err
	}

	return chambers, total, nil
}
