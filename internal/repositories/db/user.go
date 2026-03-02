package db

import (
	"docmate/internal/model"
)

func (repo *Repository) Create(user model.User) (model.User, error) {
	err := repo.client.Create(&user).Error

	return user, err
}

func (repo *Repository) Get(userID int) (model.User, error) {
	var user model.User
	if err := repo.client.Model(&model.User{}).Where("id = ?", userID).First(&user).Error; err != nil {
		return model.User{}, err
	}

	return user, nil
}

func (repo *Repository) GetByEmail(email string) (model.User, error) {
	var user model.User
	if err := repo.client.Model(&model.User{}).Where("email = ?", email).First(&user).Error; err != nil {
		return model.User{}, err
	}

	return user, nil
}

func (repo *Repository) List(offset, limit int) ([]model.User, int, error) {
	var users []model.User
	var total int64

	if err := repo.client.Model(&model.User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := repo.client.Model(&model.User{}).Offset(offset).Limit(limit).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, int(total), nil
}
