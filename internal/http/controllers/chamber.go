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

type ChamberController struct {
	baseCtx    context.Context
	chamberSvc model.ChamberUseCase
	doctorRepo model.DoctorRepo
}

func NewChamberController(
	baseCtx context.Context,
	chamberSvc model.ChamberUseCase,
	doctorRepo model.DoctorRepo,
) *ChamberController {
	return &ChamberController{
		baseCtx:    baseCtx,
		chamberSvc: chamberSvc,
		doctorRepo: doctorRepo,
	}
}

func (controller *ChamberController) authorizeDoctor(c echo.Context, doctorID int) error {
	user, err := contextutil.GetUserFromContext(c)
	if err != nil {
		return response.Unauthorized(c, "Unauthorized")
	}

	if user.Role != consts.RoleAdmin {
		doctor, err := controller.doctorRepo.GetDoctorByUserID(user.ID)
		if err != nil || doctor.ID != doctorID {
			return response.Unauthorized(c, "Unauthorized access to this doctor's chambers")
		}
	} else {
		_, err := controller.doctorRepo.GetDoctorByID(doctorID)
		if err != nil {
			return response.BadRequest(c, "Doctor not found")
		}
	}

	return nil
}

func (controller *ChamberController) Create(c echo.Context) error {
	ctx := c.Request().Context()

	var req types.ChamberReq
	if err := c.Bind(&req); err != nil {
		return response.BadRequest(c, err.Error())
	}

	if err := req.Validate(); err != nil {
		return response.BadRequest(c, err.Error())
	}

	if err := controller.authorizeDoctor(c, req.DoctorID); err != nil {
		return err
	}

	resp, err := controller.chamberSvc.Create(ctx, req)
	if err != nil {
		return response.InternalServerError(c, err.Error())
	}

	return response.Success(c, "chamber created successfully", resp)
}

func (controller *ChamberController) Update(c echo.Context) error {
	ctx := c.Request().Context()

	var req types.ChamberUpdateReq
	if err := c.Bind(&req); err != nil {
		return response.BadRequest(c, err.Error())
	}

	if err := req.Validate(); err != nil {
		return response.BadRequest(c, err.Error())
	}

	// 1. Get existing chamber to verify association
	filter := types.ChamberFilter{ID: req.ID}
	existing, err := controller.chamberSvc.Get(ctx, filter)
	if err != nil {
		return response.InternalServerError(c, "failed to retrieve chamber")
	}

	user, err := contextutil.GetUserFromContext(c)
	if err != nil {
		return response.Unauthorized(c, "Unauthorized")
	}

	// 2. Verify Authorization (Owner or Admin)
	if user.Role != consts.RoleAdmin {
		doctor, err := controller.doctorRepo.GetDoctorByUserID(user.ID)
		if err != nil || existing.DoctorID != doctor.ID {
			return response.Unauthorized(c, "unauthorized to update this chamber")
		}
	}

	resp, err := controller.chamberSvc.Update(ctx, req)
	if err != nil {
		return response.InternalServerError(c, err.Error())
	}

	return response.Success(c, "chamber updated successfully", resp)
}

func (controller *ChamberController) Get(c echo.Context) error {
	ctx := c.Request().Context()

	var filter types.ChamberFilter
	if err := c.Bind(&filter); err != nil {
		return response.BadRequest(c, "invalid chamber id")
	}

	if err := filter.Validate(); err != nil {
		return response.BadRequest(c, err.Error())
	}

	chamber, err := controller.chamberSvc.Get(ctx, filter)
	if err != nil {
		return response.InternalServerError(c, "failed to get chamber")
	}

	user, err := contextutil.GetUserFromContext(c)
	if err != nil {
		return response.Unauthorized(c, "Unauthorized")
	}

	// Verify Ownership
	if user.Role != consts.RoleAdmin {
		doctor, err := controller.doctorRepo.GetDoctorByUserID(user.ID)
		if err != nil || chamber.DoctorID != doctor.ID {
			return response.Unauthorized(c, "unauthorized to view this chamber profile")
		}
	}

	return response.Success(c, "chamber fetched successfully", chamber)
}

func (controller *ChamberController) List(c echo.Context) error {
	ctx := c.Request().Context()

	var req types.ChamberListReq
	if err := c.Bind(&req); err != nil {
		return response.BadRequest(c, err.Error())
	}

	if req.DoctorID == 0 {
		return response.BadRequest(c, "doctor id is required")
	}

	if err := controller.authorizeDoctor(c, req.DoctorID); err != nil {
		return err
	}

	// Set defaults
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 {
		req.Limit = 10
	}

	resp, err := controller.chamberSvc.List(ctx, req)
	if err != nil {
		return response.InternalServerError(c, err.Error())
	}

	return response.Success(c, "chambers fetched successfully", resp)
}
