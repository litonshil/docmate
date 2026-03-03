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

type DoctorController struct {
	baseCtx   context.Context
	doctorSvc model.DoctorUseCase
}

func NewDoctorController(
	baseCtx context.Context,
	doctorSvc model.DoctorUseCase,
) *DoctorController {
	return &DoctorController{
		baseCtx:   baseCtx,
		doctorSvc: doctorSvc,
	}
}

func (controller *DoctorController) Create(c echo.Context) error {
	ctx := c.Request().Context()
	user, err := contextutil.GetUserFromContext(c)
	if err != nil {
		return response.Unauthorized(c, "Unauthorized")
	}

	var req types.DoctorReq
	if err := c.Bind(&req); err != nil {
		return response.BadRequest(c, err.Error())
	}

	if err := req.Validate(); err != nil {
		return response.BadRequest(c, err.Error())
	}

	req.UserID = user.ID
	req.Email = user.Email

	resp, err := controller.doctorSvc.Create(ctx, req)
	if err != nil {
		return response.InternalServerError(c, err.Error())
	}

	return response.Success(c, "doctor created successfully", resp)
}

func (controller *DoctorController) Update(c echo.Context) error {
	ctx := c.Request().Context()
	user, err := contextutil.GetUserFromContext(c)
	if err != nil {
		return response.Unauthorized(c, "Unauthorized")
	}

	var req types.DoctorUpdateReq
	if err := c.Bind(&req); err != nil {
		return response.BadRequest(c, err.Error())
	}

	if err := req.Validate(); err != nil {
		return response.BadRequest(c, err.Error())
	}

	req.UserID = user.ID

	// 1. Get existing profile to verify authorization
	filter := types.DoctorFilter{ID: req.ID}
	existing, err := controller.doctorSvc.Get(ctx, filter)
	if err != nil {
		return response.InternalServerError(c, "failed to retrieve doctor profile")
	}

	// 2. Verify Authorization (Owner or Admin)
	if existing.UserID != user.ID && user.Role != consts.RoleAdmin {
		return response.Unauthorized(c, "unauthorized to update this doctor profile")
	}

	resp, err := controller.doctorSvc.Update(ctx, req)
	if err != nil {
		return response.InternalServerError(c, err.Error())
	}

	return response.Success(c, "doctor updated successfully", resp)
}

func (controller *DoctorController) Get(c echo.Context) error {
	ctx := c.Request().Context()
	user, err := contextutil.GetUserFromContext(c)
	if err != nil {
		return response.Unauthorized(c, "Unauthorized")
	}

	var filter types.DoctorFilter
	if err := c.Bind(&filter); err != nil {
		return response.BadRequest(c, "invalid doctor id")
	}

	if err := filter.Validate(); err != nil {
		return response.BadRequest(c, err.Error())
	}

	filter.UserID = user.ID

	doctor, err := controller.doctorSvc.Get(ctx, filter)
	if err != nil {
		return response.InternalServerError(c, "failed to get doctor")
	}

	return response.Success(c, "doctor fetched successfully", doctor)
}

func (controller *DoctorController) List(c echo.Context) error {
	ctx := c.Request().Context()
	user, err := contextutil.GetUserFromContext(c)
	if err != nil || user.Role != consts.RoleAdmin {
		return response.Unauthorized(c, "Unauthorized")
	}

	var req types.DoctorListReq
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

	resp, err := controller.doctorSvc.List(ctx, req)
	if err != nil {
		return response.InternalServerError(c, err.Error())
	}

	return response.Success(c, "doctors fetched successfully", resp)
}
