package types

import "time"

type MedicineReq struct {
	CreatedBy    int    `json:"created_by"`
	BrandName    string `json:"brand_name" validate:"required"`
	GenericName  string `json:"generic_name" validate:"required"`
	Form         string `json:"form" validate:"required"`
	Strength     string `json:"strength"`
	Manufacturer string `json:"manufacturer"`
	Description  string `json:"description"`
}

type MedicineUpdateReq struct {
	ID           int    `json:"id" validate:"required"`
	BrandName    string `json:"brand_name" validate:"required"`
	GenericName  string `json:"generic_name" validate:"required"`
	Form         string `json:"form" validate:"required"`
	Strength     string `json:"strength"`
	Manufacturer string `json:"manufacturer"`
	Description  string `json:"description"`
	IsActive     bool   `json:"is_active"`
}

type MedicineResp struct {
	ID           int       `json:"id"`
	CreatedBy    int       `json:"created_by"`
	BrandName    string    `json:"brand_name"`
	GenericName  string    `json:"generic_name"`
	Form         string    `json:"form"`
	Strength     string    `json:"strength"`
	Manufacturer string    `json:"manufacturer"`
	Description  string    `json:"description"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type MedicineListReq struct {
	Page   int    `query:"page"`
	Limit  int    `query:"limit"`
	Search string `query:"search"`
}

type PaginatedMedicineResp struct {
	Pagination Pagination     `json:"pagination"`
	Records    []MedicineResp `json:"records"`
}
