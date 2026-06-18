package model

import "time"

type GlobalSetting struct {
	Key       string    `gorm:"primaryKey;column:key"`
	Value     string    `gorm:"column:value"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (GlobalSetting) TableName() string {
	return "global_settings"
}
