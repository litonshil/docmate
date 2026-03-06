package medicine

import (
	"context"
	"docmate/internal/model"
	"docmate/types"
	"log/slog"
)

type Service struct {
	repo model.MedicineRepo
}

func NewService(repo model.MedicineRepo) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, req types.MedicineReq) (types.MedicineResp, error) {
	medicine := model.Medicine{
		CreatedBy:    req.CreatedBy,
		BrandName:    req.BrandName,
		GenericName:  req.GenericName,
		Form:         model.MedicineFormType(req.Form),
		Strength:     req.Strength,
		Manufacturer: req.Manufacturer,
		Description:  req.Description,
		IsActive:     true,
	}

	created, err := s.repo.CreateMedicine(medicine)
	if err != nil {
		slog.Error("failed to create medicine", "error", err)

		return types.MedicineResp{}, err
	}

	return mapToResponse(created), nil
}

func (s *Service) Get(ctx context.Context, id int) (types.MedicineResp, error) {
	medicine, err := s.repo.GetMedicineByID(id)
	if err != nil {
		return types.MedicineResp{}, err
	}

	return mapToResponse(medicine), nil
}

func (s *Service) Update(ctx context.Context, req types.MedicineUpdateReq) (types.MedicineResp, error) {
	medicine := model.Medicine{
		ID:           req.ID,
		BrandName:    req.BrandName,
		GenericName:  req.GenericName,
		Form:         model.MedicineFormType(req.Form),
		Strength:     req.Strength,
		Manufacturer: req.Manufacturer,
		Description:  req.Description,
		IsActive:     req.IsActive,
	}

	updated, err := s.repo.UpdateMedicine(medicine)
	if err != nil {
		return types.MedicineResp{}, err
	}

	return mapToResponse(updated), nil
}

func (s *Service) Delete(ctx context.Context, id int) error {
	return s.repo.DeleteMedicine(id)
}

func (s *Service) List(ctx context.Context, req types.MedicineListReq) (types.PaginatedMedicineResp, error) {
	if req.Limit == 0 {
		req.Limit = 10
	}
	if req.Page == 0 {
		req.Page = 1
	}

	offset := (req.Page - 1) * req.Limit
	medicines, total, err := s.repo.ListMedicines(offset, req.Limit, req.Search)
	if err != nil {
		return types.PaginatedMedicineResp{}, err
	}

	var records []types.MedicineResp
	for _, m := range medicines {
		records = append(records, mapToResponse(m))
	}

	lastPage := total / req.Limit
	if total%req.Limit > 0 {
		lastPage++
	}

	return types.PaginatedMedicineResp{
		Pagination: types.Pagination{
			Page:     req.Page,
			Limit:    req.Limit,
			Total:    total,
			LastPage: lastPage,
		},
		Records: records,
	}, nil
}

func mapToResponse(m model.Medicine) types.MedicineResp {
	return types.MedicineResp{
		ID:           m.ID,
		CreatedBy:    m.CreatedBy,
		BrandName:    m.BrandName,
		GenericName:  m.GenericName,
		Form:         string(m.Form),
		Strength:     m.Strength,
		Manufacturer: m.Manufacturer,
		Description:  m.Description,
		IsActive:     m.IsActive,
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
	}
}
