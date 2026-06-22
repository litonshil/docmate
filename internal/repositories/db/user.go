package db

import (
	"docmate/internal/model"
)

func (r *Repository) Create(user model.User) (model.User, error) {
	err := r.client.Create(&user).Error

	return user, err
}

func (r *Repository) Get(userID int) (model.User, error) {
	var user model.User
	if err := r.client.Model(&model.User{}).Where("id = ?", userID).First(&user).Error; err != nil {
		return model.User{}, err
	}

	return user, nil
}

func (r *Repository) GetByEmail(email string) (model.User, error) {
	var user model.User
	if err := r.client.Model(&model.User{}).Where("email = ?", email).First(&user).Error; err != nil {
		return model.User{}, err
	}

	return user, nil
}

func (r *Repository) List(offset, limit int) ([]model.User, int64, error) {
	var users []model.User
	var total int64

	if err := r.client.Model(&model.User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.client.Model(&model.User{}).Offset(offset).Limit(limit).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (r *Repository) Update(user model.User) (model.User, error) {
	err := r.client.Save(&user).Error

	return user, err
}
