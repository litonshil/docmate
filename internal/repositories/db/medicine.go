package db

import (
	"docmate/internal/model"
	"fmt"
)

func (r *Repository) CreateMedicine(medicine model.Medicine) (model.Medicine, error) {
	err := r.client.Create(&medicine).Error

	return medicine, err
}

func (r *Repository) UpdateMedicine(medicine model.Medicine) (model.Medicine, error) {
	err := r.client.Model(&model.Medicine{}).Where("id = ?", medicine.ID).Updates(&medicine).Error

	return medicine, err
}

func (r *Repository) GetMedicineByID(id int) (model.Medicine, error) {
	var medicine model.Medicine
	if err := r.client.Where("id = ?", id).First(&medicine).Error; err != nil {
		return model.Medicine{}, err
	}

	return medicine, nil
}

func (r *Repository) DeleteMedicine(id int) error {
	return r.client.Delete(&model.Medicine{}, id).Error
}

func (r *Repository) ListMedicines(offset, limit int, search string) ([]model.Medicine, int64, error) {
	var medicines []model.Medicine
	var total int64

	db := r.client.Model(&model.Medicine{})

	if search != "" {
		searchQuery := fmt.Sprintf("%s%%", search)
		db = db.Where("brand_name ILIKE ? OR generic_name ILIKE ?", searchQuery, searchQuery)
		db = db.Order("brand_name ASC")
	} else {
		db = db.Order("brand_name ASC")
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := db.Offset(offset).Limit(limit).Find(&medicines).Error; err != nil {
		return nil, 0, err
	}

	return medicines, total, nil
}
