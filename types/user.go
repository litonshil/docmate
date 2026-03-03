package types

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type User struct {
	ID       int    `json:"id"`
	UserName string `json:"user_name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type AuthUser struct {
	ID       int    `json:"id"`
	UserName string `json:"user_name"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}

type UserReq struct {
	ID       int    `json:"id"`
	UserName string `json:"user_name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type UserResp struct {
	ID        int        `json:"id"`
	UserName  string     `json:"user_name"`
	Email     string     `json:"email"`
	Password  string     `json:"-"`
	Role      string     `json:"role"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

type LoginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (r LoginReq) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Email, validation.Required),
		validation.Field(&r.Password, validation.Required),
	)
}

type LoginResp struct {
	Token string   `json:"token"`
	User  UserResp `json:"user"`
}

type UserFilter struct {
	ID int `json:"id" query:"id" param:"id"`
}

func (u UserFilter) Validate() error {
	return validation.ValidateStruct(&u,
		validation.Field(&u.ID, validation.Required),
	)
}

// Validate method for AgentReq.
func (ur UserReq) Validate() error {
	return validation.ValidateStruct(&ur,
		// ID is mandatory
		validation.Field(&ur.UserName, validation.Required),
		validation.Field(&ur.Email, validation.Required),
		validation.Field(&ur.Password, validation.Required),
	)
}

type Pagination struct {
	Page     int `json:"page" query:"page"`
	Limit    int `json:"limit" query:"limit"`
	Total    int `json:"total"`
	LastPage int `json:"last_page"`
}

type UserListReq struct {
	Pagination
}

type PaginatedResponse struct {
	Pagination
	Records []UserResp `json:"records"`
}
