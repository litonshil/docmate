package db

import (
	"docmate/internal/model"
)

func (repo *Repository) CreateChamber(chamber model.Chamber) (model.Chamber, error) {
	err := repo.client.Create(&chamber).Error

	return chamber, err
}

func (repo *Repository) UpdateChamber(chamber model.Chamber) (model.Chamber, error) {
	// Omit CreatedAt and let GORM handle UpdatedAt automatically
	err := repo.client.Model(&model.Chamber{}).Where("id = ?", chamber.ID).Omit("created_at").Updates(&chamber).Error

	return chamber, err
}

func (repo *Repository) GetChamberByID(id int) (model.Chamber, error) {
	var chamber model.Chamber
	if err := repo.dbClient(nil).Model(&model.Chamber{}).Where("id = ?", id).First(&chamber).Error; err != nil {
		return model.Chamber{}, err
	}

	return chamber, nil
}

func (repo *Repository) ListChambers(offset, limit, doctorID int) ([]model.Chamber, int, error) {
	var chambers []model.Chamber
	var total int64

	query := repo.client.Model(&model.Chamber{})

	if doctorID > 0 {
		query = query.Where("doctor_id = ?", doctorID)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Offset(offset).Limit(limit).Find(&chambers).Error; err != nil {
		return nil, 0, err
	}

	return chambers, int(total), nil
}
