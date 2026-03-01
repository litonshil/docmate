package model

import (
	"context"
	"docmate/types"
	"time"
)

type User struct {
	ID        int        `json:"id"`
	Type      string     `json:"type"`
	UserName  string     `json:"user_name"`
	Password  string     `json:"password"`
	FirstName string     `json:"first_name"`
	LastName  string     `json:"last_name"`
	Email     string     `json:"email"`
	Phone     string     `json:"phone"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

type UserResp struct {
	ID        int        `json:"id"`
	Type      string     `json:"type"`
	UserName  string     `json:"user_name"`
	Password  string     `json:"password"`
	FirstName string     `json:"first_name"`
	LastName  string     `json:"last_name"`
	Email     string     `json:"email"`
	Phone     string     `json:"phone"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

type IntermediateUserPermissions struct {
	ID             int    `json:"id"`
	UserName       string `json:"user_name"`
	Email          string `json:"email"`
	RoleName       string `json:"role_name"`
	Password       string `json:"-"`
	PermissionName string `json:"permission_name"`
}

type IntermediateUserRoles struct {
	UserID          int    `json:"user_id"`
	RoleID          int    `json:"role_id"`
	RoleName        string `json:"role_name"`
	RoleDescription string `json:"role_description"`
}

type DeleteUserRoleReq struct {
	UserID int `json:"user_id"`
	RoleID int `json:"role_id"`
}

type ResetPasswordReq struct {
	ID       int    `json:"id"`
	Password string `json:"password"`
}

type UserUseCase interface {
	CreateUser(ctx context.Context, req types.UserReq) (types.UserResp, error)
	GetUser(ctx context.Context, userID int) (types.UserResp, error)
	ListUsers(ctx context.Context, req types.UserListReq) (types.PaginatedResponse, error)
	Login(ctx context.Context, req types.LoginReq) (types.LoginResp, error)
}

type UserRepo interface {
	CreateUser(req User) (User, error)
	GetUser(userID int) (UserResp, error)
	GetUserByEmail(email string) (UserResp, error)
	ListUsers(offset, limit int) ([]UserResp, int, error)
}
