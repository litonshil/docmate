package db

import (
	"docmate/internal/model"
)

func (repo *Repository) CreateUser(user model.User) (model.User, error) {
	err := repo.client.Create(&user).Error
	return user, err
}

func (repo *Repository) GetByID(id int) (*model.User, error) {
	var user model.User
	err := repo.client.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (repo *Repository) GetUser(userID int) (model.UserResp, error) {
	var user model.UserResp
	if err := repo.client.Model(&model.User{}).Where("id = ?", userID).First(&user).Error; err != nil {
		return model.UserResp{}, err
	}
	return user, nil
}

func (repo *Repository) ListUsers(offset, limit int) ([]model.UserResp, int, error) {
	var users []model.UserResp
	var total int64

	if err := repo.client.Model(&model.User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := repo.client.Model(&model.User{}).Offset(offset).Limit(limit).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, int(total), nil
}
