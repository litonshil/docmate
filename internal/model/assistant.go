package model

import (
	"context"
	"docmate/types"
	"time"
)

type Assistant struct {
	ID        int       `json:"id" gorm:"primaryKey"`
	UserID    int       `json:"user_id" gorm:"not null;uniqueIndex"`
	DoctorID  int       `json:"doctor_id" gorm:"not null;index"`
	Name      string    `json:"name" gorm:"not null"`
	Phone     string    `json:"phone" gorm:"not null"`
	IsActive  bool      `json:"is_active" gorm:"not null;default:true"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Associations
	User   User   `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Doctor Doctor `json:"doctor,omitempty" gorm:"foreignKey:DoctorID"`
}

type AssistantUseCase interface {
	Create(ctx context.Context, req types.AssistantReq, doctorID int) (types.AssistantResp, error)
	Update(ctx context.Context, id int, req types.AssistantUpdateReq, doctorID int) (types.AssistantResp, error)
	Get(ctx context.Context, id int, doctorID int) (types.AssistantResp, error)
	List(ctx context.Context, doctorID int) ([]types.AssistantResp, error)
}

type AssistantRepo interface {
	CreateAssistant(assistant Assistant, chamberIDs []int) (Assistant, error)
	UpdateAssistant(assistant Assistant, chamberIDs []int) (Assistant, error)
	GetAssistantByID(id int) (Assistant, error)
	GetAssistantByUserID(userID int) (Assistant, error)
	ListAssistantsByDoctorID(doctorID int) ([]Assistant, error)
	GetChambersByAssistantID(assistantID int) ([]Chamber, error)
}
