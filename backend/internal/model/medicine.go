package model

import (
	"context"
	"docmate/types"
	"time"
)

type MedicineFormType string

const (
	MedicineFormTablet      MedicineFormType = "tablet"
	MedicineFormCapsule     MedicineFormType = "capsule"
	MedicineFormSyrup       MedicineFormType = "syrup"
	MedicineFormSuspension  MedicineFormType = "suspension"
	MedicineFormInjection   MedicineFormType = "injection"
	MedicineFormInhaler     MedicineFormType = "inhaler"
	MedicineFormDrops       MedicineFormType = "drops"
	MedicineFormCream       MedicineFormType = "cream"
	MedicineFormOintment    MedicineFormType = "ointment"
	MedicineFormGel         MedicineFormType = "gel"
	MedicineFormPatch       MedicineFormType = "patch"
	MedicineFormSuppository MedicineFormType = "suppository"
	MedicineFormPowder      MedicineFormType = "powder"
	MedicineFormSachet      MedicineFormType = "sachet"
	MedicineFormOther       MedicineFormType = "other"
)

type Medicine struct {
	ID           int              `json:"id" gorm:"primaryKey"`
	CreatedBy    int              `json:"created_by"`
	BrandName    string           `json:"brand_name"`
	GenericName  string           `json:"generic_name"`
	Form         MedicineFormType `json:"form" gorm:"type:medicine_form_type"`
	Strength     string           `json:"strength"`
	Manufacturer string           `json:"manufacturer"`
	Description  string           `json:"description"`
	IsActive     bool             `json:"is_active" gorm:"default:true"`
	CreatedAt    time.Time        `json:"created_at"`
	UpdatedAt    time.Time        `json:"updated_at"`
	DeletedAt    *time.Time       `json:"deleted_at" gorm:"index"`
}

type MedicineUseCase interface {
	Create(ctx context.Context, req types.MedicineReq) (types.MedicineResp, error)
	Get(ctx context.Context, id int) (types.MedicineResp, error)
	Update(ctx context.Context, req types.MedicineUpdateReq) (types.MedicineResp, error)
	Delete(ctx context.Context, id int) error
	List(ctx context.Context, req types.MedicineListReq) (types.PaginatedMedicineResp, error)
}

type MedicineRepo interface {
	CreateMedicine(medicine Medicine) (Medicine, error)
	UpdateMedicine(medicine Medicine) (Medicine, error)
	GetMedicineByID(id int) (Medicine, error)
	DeleteMedicine(id int) error
	ListMedicines(offset, limit int, search string) ([]Medicine, int, error)
}
