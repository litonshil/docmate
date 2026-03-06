package db

import (
	"docmate/internal/model"
	"fmt"
)

func (repo *Repository) CreateMedicine(medicine model.Medicine) (model.Medicine, error) {
	err := repo.client.Create(&medicine).Error

	return medicine, err
}

func (repo *Repository) UpdateMedicine(medicine model.Medicine) (model.Medicine, error) {
	err := repo.client.Model(&model.Medicine{}).Where("id = ?", medicine.ID).Updates(&medicine).Error

	return medicine, err
}

func (repo *Repository) GetMedicineByID(id int) (model.Medicine, error) {
	var medicine model.Medicine
	if err := repo.client.Where("id = ?", id).First(&medicine).Error; err != nil {
		return model.Medicine{}, err
	}

	return medicine, nil
}

func (repo *Repository) DeleteMedicine(id int) error {
	return repo.client.Delete(&model.Medicine{}, id).Error
}

func (repo *Repository) ListMedicines(offset, limit int, search string) ([]model.Medicine, int, error) {
	var medicines []model.Medicine
	var total int64

	db := repo.client.Model(&model.Medicine{})

	if search != "" {
		searchQuery := fmt.Sprintf("%%%s%%", search)
		db = db.Where("brand_name ILIKE ? OR generic_name ILIKE ?", searchQuery, searchQuery)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := db.Offset(offset).Limit(limit).Find(&medicines).Error; err != nil {
		return nil, 0, err
	}

	return medicines, int(total), nil
}
