package patient

import (
	"context"
	"docmate/internal/model"
	"docmate/types"
	"fmt"
	"log/slog"
)

type Service struct {
	patientRepo model.PatientRepo
}

func NewService(patientRepo model.PatientRepo) *Service {
	return &Service{
		patientRepo: patientRepo,
	}
}

func (service *Service) Create(ctx context.Context, req types.PatientReq, doctorID int) (types.PatientResp, error) {
	payload := model.Patient{
		DoctorID:       doctorID,
		FullName:       req.FullName,
		Gender:         req.Gender,
		Age:            req.Age,
		Phone:          req.Phone,
		Email:          req.Email,
		BloodGroup:     req.BloodGroup,
		Allergies:      req.Allergies,
		MedicalHistory: req.MedicalHistory,
	}

	patient, err := service.patientRepo.CreatePatient(payload)
	if err != nil {
		slog.Error("failed to create patient", "error", err.Error())

		return types.PatientResp{}, err
	}

	return mapToPatientResponse(patient), nil
}

func (service *Service) Update(ctx context.Context, req types.PatientUpdateReq) (types.PatientResp, error) {
	// 1. Get existing patient to verify non-mutable properties
	existing, err := service.patientRepo.GetPatientByID(req.ID)
	if err != nil {
		return types.PatientResp{}, fmt.Errorf("failed to retrieve patient: %w", err)
	}

	payload := model.Patient{
		ID:             req.ID,
		DoctorID:       existing.DoctorID,
		FullName:       req.FullName,
		Gender:         req.Gender,
		Age:            req.Age,
		Phone:          req.Phone,
		Email:          req.Email,
		BloodGroup:     req.BloodGroup,
		Allergies:      req.Allergies,
		MedicalHistory: req.MedicalHistory,
	}

	patient, err := service.patientRepo.UpdatePatient(payload)
	if err != nil {
		slog.Error("failed to update patient", "error", err.Error())

		return types.PatientResp{}, err
	}

	return mapToPatientResponse(patient), nil
}

func (service *Service) Get(ctx context.Context, filter types.PatientFilter) (types.PatientResp, error) {
	patient, err := service.patientRepo.GetPatientByID(filter.ID)
	if err != nil {
		slog.Error("failed to get patient by ID", "error", err.Error())

		return types.PatientResp{}, err
	}

	return mapToPatientResponse(patient), nil
}

func (service *Service) List(ctx context.Context, req types.PatientListReq, doctorID int) (types.PaginatedPatientResp, error) {
	offset := (req.Page - 1) * req.Limit
	patients, total, err := service.patientRepo.ListPatients(offset, req.Limit, doctorID, req.Name, req.Phone)
	if err != nil {
		slog.Error("failed to list patients", "error", err.Error())

		return types.PaginatedPatientResp{}, err
	}

	var records []types.PatientResp
	for _, patient := range patients {
		records = append(records, mapToPatientResponse(patient))
	}

	lastPage := total / req.Limit
	if total%req.Limit > 0 {
		lastPage++
	}

	return types.PaginatedPatientResp{
		Pagination: types.Pagination{
			Page:     req.Page,
			Limit:    req.Limit,
			Total:    total,
			LastPage: lastPage,
		},
		Records: records,
	}, nil
}

func mapToPatientResponse(patient model.Patient) types.PatientResp {
	return types.PatientResp{
		ID:             patient.ID,
		DoctorID:       patient.DoctorID,
		FullName:       patient.FullName,
		Gender:         patient.Gender,
		Age:            patient.Age,
		Phone:          patient.Phone,
		Email:          patient.Email,
		BloodGroup:     patient.BloodGroup,
		Allergies:      patient.Allergies,
		MedicalHistory: patient.MedicalHistory,
		CreatedAt:      patient.CreatedAt,
		UpdatedAt:      patient.UpdatedAt,
	}
}
