package doctor

import (
	"context"
	"docmate/internal/model"
	"docmate/types"
	"encoding/json"
	"fmt"
	"log/slog"
)

type Service struct {
	doctorRepo model.DoctorRepo
}

func NewService(doctorRepo model.DoctorRepo) *Service {
	return &Service{
		doctorRepo: doctorRepo,
	}
}

func (service *Service) Create(ctx context.Context, req types.DoctorReq) (types.DoctorResp, error) {
	// Check if already exists so we don't duplicate
	existing, err := service.doctorRepo.GetDoctorByUserID(req.UserID)
	if err == nil && existing.ID != 0 {
		return types.DoctorResp{}, fmt.Errorf("doctor profile already exists for this user")
	}

	degreeJson, _ := json.Marshal(req.Degree)
	specializationJson, _ := json.Marshal(req.Specialization)

	payload := model.Doctor{
		UserID:         req.UserID,
		Email:          req.Email, // Link to the user's primary email automatically
		FullName:       req.FullName,
		Degree:         degreeJson,
		Specialization: specializationJson,
		Phone:          req.Phone,
		Bio:            req.Bio,
		SignatureURL:   req.SignatureURL,
	}

	doctor, err := service.doctorRepo.CreateDoctor(payload)
	if err != nil {
		slog.Error("failed to create doctor profile", "error", err.Error())

		return types.DoctorResp{}, err
	}

	return mapToDoctorResponse(doctor), nil
}

func (service *Service) Update(ctx context.Context, req types.DoctorUpdateReq) (types.DoctorResp, error) {
	// 1. Get the existing profile to verify properties that should not be updated (Email, UserID)
	existing, err := service.doctorRepo.GetDoctorByID(req.ID)
	if err != nil {
		return types.DoctorResp{}, fmt.Errorf("failed to retrieve doctor profile: %w", err)
	}

	degreeJson, _ := json.Marshal(req.Degree)
	specializationJson, _ := json.Marshal(req.Specialization)

	payload := model.Doctor{
		ID:             req.ID,
		UserID:         existing.UserID,
		Email:          existing.Email,
		FullName:       req.FullName,
		Degree:         degreeJson,
		Specialization: specializationJson,
		Phone:          req.Phone,
		Bio:            req.Bio,
		SignatureURL:   req.SignatureURL,
	}

	doctor, err := service.doctorRepo.UpdateDoctor(payload)
	if err != nil {
		slog.Error("failed to update doctor profile", "error", err.Error())

		return types.DoctorResp{}, err
	}

	return mapToDoctorResponse(doctor), nil
}

func (service *Service) Get(ctx context.Context, filter types.DoctorFilter) (types.DoctorResp, error) {
	var doctor model.Doctor
	var err error

	if filter.ID != 0 {
		doctor, err = service.doctorRepo.GetDoctorByID(filter.ID)
	} else {
		doctor, err = service.doctorRepo.GetDoctorByUserID(filter.UserID)
	}

	if err != nil {
		slog.Error("failed to get doctor profile", "error", err.Error())

		return types.DoctorResp{}, fmt.Errorf("failed to get doctor profile: %w", err)
	}

	return mapToDoctorResponse(doctor), nil
}

func (service *Service) List(ctx context.Context, req types.DoctorListReq) (types.PaginatedDoctorResp, error) {
	offset := (req.Page - 1) * req.Limit
	doctors, total, err := service.doctorRepo.ListDoctors(offset, req.Limit)
	if err != nil {
		slog.Error("failed to list doctors list", "error", err.Error())

		return types.PaginatedDoctorResp{}, err
	}

	var records []types.DoctorResp
	for _, doctor := range doctors {
		records = append(records, mapToDoctorResponse(doctor))
	}

	lastPage := total / req.Limit
	if total%req.Limit > 0 {
		lastPage++
	}

	return types.PaginatedDoctorResp{
		Pagination: types.Pagination{
			Page:     req.Page,
			Limit:    req.Limit,
			Total:    total,
			LastPage: lastPage,
		},
		Records: records,
	}, nil
}

func mapToDoctorResponse(doc model.Doctor) types.DoctorResp {
	var degree []string
	var specialization []string
	_ = json.Unmarshal(doc.Degree, &degree)
	_ = json.Unmarshal(doc.Specialization, &specialization)

	return types.DoctorResp{
		ID:             doc.ID,
		UserID:         doc.UserID,
		Email:          doc.Email,
		FullName:       doc.FullName,
		Degree:         degree,
		Specialization: specialization,
		Phone:          doc.Phone,
		Bio:            doc.Bio,
		SignatureURL:   doc.SignatureURL,
		CreatedAt:      doc.CreatedAt,
		UpdatedAt:      doc.UpdatedAt,
	}
}
