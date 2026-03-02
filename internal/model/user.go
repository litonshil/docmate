package model

import (
	"context"
	"docmate/types"
	"time"
)

type User struct {
	ID        int        `json:"id"`
	UserName  string     `json:"user_name"`
	Email     string     `json:"email"`
	Password  string     `json:"password"`
	Role      string     `json:"role"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

type UserUseCase interface {
	Create(ctx context.Context, req types.UserReq) (types.UserResp, error)
	Get(ctx context.Context, userID int) (types.UserResp, error)
	List(ctx context.Context, req types.UserListReq) (types.PaginatedResponse, error)
	Login(ctx context.Context, req types.LoginReq) (types.LoginResp, error)
}

type UserRepo interface {
	Create(req User) (User, error)
	Get(userID int) (User, error)
	GetByEmail(email string) (User, error)
	List(offset, limit int) ([]User, int, error)
}
