package types

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type User struct {
	ID        int    `json:"id"`
	Type      string `json:"type"`
	UserName  string `json:"user_name"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
}

type UserReq struct {
	ID        int    `json:"id"`
	Type      string `json:"type"`
	UserName  string `json:"user_name"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
}

type UserUpdateReq struct {
	ID        int    `json:"id" param:"id" query:"id"`
	UserName  string `json:"user_name"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Phone     string `json:"phone"`
}

type LoggedInUser struct {
	ID          int      `json:"user_id"`
	AccessUuid  string   `json:"access_uuid"`
	RefreshUuid string   `json:"refresh_uuid"`
	Role        string   `json:"role"`
	Permissions []string `json:"permissions"`
}

type UserResp struct {
	ID        int        `json:"id"`
	Type      string     `json:"type"`
	UserName  string     `json:"user_name"`
	Password  string     `json:"-"`
	FirstName string     `json:"first_name"`
	LastName  string     `json:"last_name"`
	Email     string     `json:"email"`
	Phone     string     `json:"phone"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

type UserRolePermissionsInfo struct {
	ID          int      `json:"id"`
	UserName    string   `json:"user_name"`
	Email       string   `json:"email"`
	Role        string   `json:"role"`
	Password    string   `json:"-"`
	Permissions []string `json:"permissions"`
}

type CurrentUser struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	UserName string `json:"user_name"`
}

type UserFilter struct {
	ID int `json:"id" query:"id" param:"id"`
}

type UserWithParamsResp struct {
	CurrentUser
	Role        string   `json:"role,omitempty"`
	Permissions []string `json:"permissions,omitempty"`
}

type AssignRoleRequest struct {
	ID     int `json:"id" query:"id" param:"id"`
	RoleID int `json:"role_id"`
}

type UserRoleParam struct {
	ID     int `json:"id" query:"id" param:"id"`
	RoleID int `json:"role_id" query:"role_id" param:"role_id"`
}

type DeleteUserRoleReq struct {
	ID     int `json:"id"`
	RoleID int `json:"role_id"`
}

type ResetPasswordReq struct {
	ID       int    `json:"id"`
	Password string `json:"password"`
}

func (rp ResetPasswordReq) Validate() error {
	return validation.ValidateStruct(&rp,
		validation.Field(&rp.Password, validation.Required),
	)
}

func (uur UserUpdateReq) Validate() error {
	return validation.ValidateStruct(&uur,
		validation.Field(&uur.ID, validation.Required),
	)
}

func (u UserFilter) Validate() error {
	return validation.ValidateStruct(&u,
		validation.Field(&u.ID, validation.Required),
	)
}

func (r AssignRoleRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.ID, validation.Required),
		validation.Field(&r.RoleID, validation.Required),
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

func (u UserRoleParam) Validate() error {
	return validation.ValidateStruct(&u,
		validation.Field(&u.ID, validation.Required),
		validation.Field(&u.RoleID, validation.Required),
	)
}

type DeleteUserReq struct {
	ID int `json:"id" query:"id" param:"id"`
}

func (r DeleteUserReq) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.ID, validation.Required),
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
