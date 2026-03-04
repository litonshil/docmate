package chamber

import (
	"context"
	"docmate/internal/model"
	"docmate/types"
	"encoding/json"
	"fmt"
	"log/slog"

	"gorm.io/datatypes"
)

type Service struct {
	chamberRepo model.ChamberRepo
}

func NewService(chamberRepo model.ChamberRepo) *Service {
	return &Service{
		chamberRepo: chamberRepo,
	}
}

func (service *Service) Create(ctx context.Context, req types.ChamberReq) (types.ChamberResp, error) {
	visitingHoursJSON, err := json.Marshal(req.VisitingHours)
	if err != nil {
		return types.ChamberResp{}, fmt.Errorf("invalid visiting hours format")
	}

	payload := model.Chamber{
		DoctorID:      req.DoctorID,
		Name:          req.Name,
		Address:       req.Address,
		Area:          req.Area,
		City:          req.City,
		Country:       req.Country,
		Phone:         req.Phone,
		Email:         req.Email,
		Fee:           req.Fee,
		FollowUpFee:   req.FollowUpFee,
		VisitingHours: datatypes.JSON(visitingHoursJSON),
		IsActive:      req.IsActive != nil && *req.IsActive,
	}

	chamber, err := service.chamberRepo.CreateChamber(payload)
	if err != nil {
		slog.Error("failed to create chamber", "error", err.Error())

		return types.ChamberResp{}, err
	}

	return mapToChamberResponse(chamber), nil
}

func (service *Service) Update(ctx context.Context, req types.ChamberUpdateReq) (types.ChamberResp, error) {
	// 1. Get existing chamber to retain non-mutable properties like DoctorID
	existing, err := service.chamberRepo.GetChamberByID(req.ID)
	if err != nil {
		return types.ChamberResp{}, fmt.Errorf("failed to retrieve chamber: %w", err)
	}

	visitingHoursJSON, err := json.Marshal(req.VisitingHours)
	if err != nil {
		return types.ChamberResp{}, fmt.Errorf("invalid visiting hours format")
	}

	isActive := existing.IsActive
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	payload := model.Chamber{
		ID:            req.ID,
		DoctorID:      existing.DoctorID,
		Name:          req.Name,
		Address:       req.Address,
		Area:          req.Area,
		City:          req.City,
		Country:       req.Country,
		Phone:         req.Phone,
		Email:         req.Email,
		Fee:           req.Fee,
		FollowUpFee:   req.FollowUpFee,
		VisitingHours: datatypes.JSON(visitingHoursJSON),
		IsActive:      isActive,
	}

	chamber, err := service.chamberRepo.UpdateChamber(payload)
	if err != nil {
		slog.Error("failed to update chamber", "error", err.Error())

		return types.ChamberResp{}, err
	}

	return mapToChamberResponse(chamber), nil
}

func (service *Service) Get(ctx context.Context, filter types.ChamberFilter) (types.ChamberResp, error) {
	chamber, err := service.chamberRepo.GetChamberByID(filter.ID)
	if err != nil {
		slog.Error("failed to get chamber by ID", "error", err.Error())

		return types.ChamberResp{}, err
	}

	return mapToChamberResponse(chamber), nil
}

func (service *Service) List(ctx context.Context, req types.ChamberListReq) (types.PaginatedChamberResp, error) {
	offset := (req.Page - 1) * req.Limit
	chambers, total, err := service.chamberRepo.ListChambers(offset, req.Limit, req.DoctorID)
	if err != nil {
		slog.Error("failed to list chambers", "error", err.Error())

		return types.PaginatedChamberResp{}, err
	}

	var records []types.ChamberResp
	for _, chamber := range chambers {
		records = append(records, mapToChamberResponse(chamber))
	}

	lastPage := total / req.Limit
	if total%req.Limit > 0 {
		lastPage++
	}

	return types.PaginatedChamberResp{
		Pagination: types.Pagination{
			Page:     req.Page,
			Limit:    req.Limit,
			Total:    total,
			LastPage: lastPage,
		},
		Records: records,
	}, nil
}

func mapToChamberResponse(chamber model.Chamber) types.ChamberResp {
	var visitingHours []types.VisitingDay
	if len(chamber.VisitingHours) > 0 {
		_ = json.Unmarshal(chamber.VisitingHours, &visitingHours)
	}

	return types.ChamberResp{
		ID:            chamber.ID,
		DoctorID:      chamber.DoctorID,
		Name:          chamber.Name,
		Address:       chamber.Address,
		Area:          chamber.Area,
		City:          chamber.City,
		Country:       chamber.Country,
		Phone:         chamber.Phone,
		Email:         chamber.Email,
		Fee:           chamber.Fee,
		FollowUpFee:   chamber.FollowUpFee,
		VisitingHours: visitingHours,
		IsActive:      chamber.IsActive,
		CreatedAt:     chamber.CreatedAt,
		UpdatedAt:     chamber.UpdatedAt,
	}
}
