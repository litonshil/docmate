package controllers

import (
	"context"
	"docmate/internal/consts"
	"docmate/internal/model"
	"docmate/response"
	"docmate/types"
	"docmate/utils/contextutil"

	"github.com/labstack/echo/v4"
)

type PatientController struct {
	baseCtx    context.Context
	patientSvc model.PatientUseCase
	doctorRepo model.DoctorRepo
}

func NewPatientController(
	baseCtx context.Context,
	patientSvc model.PatientUseCase,
	doctorRepo model.DoctorRepo,
) *PatientController {
	return &PatientController{
		baseCtx:    baseCtx,
		patientSvc: patientSvc,
		doctorRepo: doctorRepo,
	}
}

func (controller *PatientController) Create(c echo.Context) error {
	ctx := c.Request().Context()
	user, err := contextutil.GetUserFromContext(c)
	if err != nil {
		return response.Unauthorized(c, "Unauthorized")
	}

	doctor, err := controller.doctorRepo.GetDoctorByUserID(user.ID)
	if err != nil {
		return response.BadRequest(c, "Doctor profile not found for user")
	}

	var req types.PatientReq
	if err := c.Bind(&req); err != nil {
		return response.BadRequest(c, err.Error())
	}

	if err := req.Validate(); err != nil {
		return response.BadRequest(c, err.Error())
	}

	resp, err := controller.patientSvc.Create(ctx, req, doctor.ID)
	if err != nil {
		return response.InternalServerError(c, err.Error())
	}

	return response.Success(c, "patient created successfully", resp)
}

func (controller *PatientController) Update(c echo.Context) error {
	ctx := c.Request().Context()
	user, err := contextutil.GetUserFromContext(c)
	if err != nil {
		return response.Unauthorized(c, "Unauthorized")
	}

	doctor, err := controller.doctorRepo.GetDoctorByUserID(user.ID)
	if err != nil {
		return response.BadRequest(c, "Doctor profile not found for user")
	}

	var req types.PatientUpdateReq
	if err := c.Bind(&req); err != nil {
		return response.BadRequest(c, err.Error())
	}

	if err := req.Validate(); err != nil {
		return response.BadRequest(c, err.Error())
	}

	// 1. Get existing patient to verify authorization
	filter := types.PatientFilter{ID: req.ID}
	existing, err := controller.patientSvc.Get(ctx, filter)
	if err != nil {
		return response.InternalServerError(c, "failed to retrieve patient profile")
	}

	// 2. Verify Authorization (Owner or Admin)
	if existing.DoctorID != doctor.ID && user.Role != consts.RoleAdmin {
		return response.Unauthorized(c, "unauthorized to update this patient profile")
	}

	resp, err := controller.patientSvc.Update(ctx, req)
	if err != nil {
		return response.InternalServerError(c, err.Error())
	}

	return response.Success(c, "patient updated successfully", resp)
}

func (controller *PatientController) Get(c echo.Context) error {
	ctx := c.Request().Context()
	user, err := contextutil.GetUserFromContext(c)
	if err != nil {
		return response.Unauthorized(c, "Unauthorized")
	}

	doctor, err := controller.doctorRepo.GetDoctorByUserID(user.ID)
	if err != nil {
		return response.BadRequest(c, "Doctor profile not found for user")
	}

	var filter types.PatientFilter
	if err := c.Bind(&filter); err != nil {
		return response.BadRequest(c, "invalid patient id")
	}

	if err := filter.Validate(); err != nil {
		return response.BadRequest(c, err.Error())
	}

	patient, err := controller.patientSvc.Get(ctx, filter)
	if err != nil {
		return response.InternalServerError(c, "failed to get patient")
	}

	// Verify Ownership
	if patient.DoctorID != doctor.ID && user.Role != consts.RoleAdmin {
		return response.Unauthorized(c, "unauthorized to view this patient profile")
	}

	return response.Success(c, "patient fetched successfully", patient)
}

func (controller *PatientController) List(c echo.Context) error {
	ctx := c.Request().Context()
	user, err := contextutil.GetUserFromContext(c)
	if err != nil {
		return response.Unauthorized(c, "Unauthorized")
	}

	doctor, err := controller.doctorRepo.GetDoctorByUserID(user.ID)
	if err != nil {
		return response.BadRequest(c, "Doctor profile not found for user")
	}

	var req types.PatientListReq
	if err := c.Bind(&req); err != nil {
		return response.BadRequest(c, err.Error())
	}

	// Set defaults
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 {
		req.Limit = 10
	}

	resp, err := controller.patientSvc.List(ctx, req, doctor.ID)
	if err != nil {
		return response.InternalServerError(c, err.Error())
	}

	return response.Success(c, "patients fetched successfully", resp)
}
