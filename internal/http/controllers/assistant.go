package controllers

import (
	"context"
	"docmate/internal/model"
	"docmate/response"
	"docmate/types"
	"docmate/utils/contextutil"
	"strconv"

	"github.com/labstack/echo/v4"
)

type AssistantController struct {
	baseCtx      context.Context
	assistantSvc model.AssistantUseCase
	doctorRepo   model.DoctorRepo
}

func NewAssistantController(
	baseCtx context.Context,
	assistantSvc model.AssistantUseCase,
	doctorRepo model.DoctorRepo,
) *AssistantController {
	return &AssistantController{
		baseCtx:      baseCtx,
		assistantSvc: assistantSvc,
		doctorRepo:   doctorRepo,
	}
}

func (ctrl *AssistantController) Create(c echo.Context) error {
	ctx := c.Request().Context()
	user, err := contextutil.GetUserFromContext(c)
	if err != nil {
		return response.Unauthorized(c, "Unauthorized")
	}

	doctor, err := ctrl.doctorRepo.GetDoctorByUserID(user.ID)
	if err != nil {
		return response.BadRequest(c, "Doctor profile not found for user")
	}

	var req types.AssistantReq
	if err := c.Bind(&req); err != nil {
		return response.BadRequest(c, err.Error())
	}

	resp, err := ctrl.assistantSvc.Create(ctx, req, doctor.ID)
	if err != nil {
		return response.InternalServerError(c, err.Error())
	}

	return response.Created(c, "assistant created successfully", resp)
}

func (ctrl *AssistantController) Update(c echo.Context) error {
	ctx := c.Request().Context()
	id, _ := strconv.Atoi(c.Param("id"))
	if id == 0 {
		return response.BadRequest(c, "Invalid assistant ID")
	}

	user, err := contextutil.GetUserFromContext(c)
	if err != nil {
		return response.Unauthorized(c, "Unauthorized")
	}

	doctor, err := ctrl.doctorRepo.GetDoctorByUserID(user.ID)
	if err != nil {
		return response.BadRequest(c, "Doctor profile not found for user")
	}

	var req types.AssistantUpdateReq
	if err := c.Bind(&req); err != nil {
		return response.BadRequest(c, err.Error())
	}

	resp, err := ctrl.assistantSvc.Update(ctx, id, req, doctor.ID)
	if err != nil {
		return response.InternalServerError(c, err.Error())
	}

	return response.Success(c, "assistant updated successfully", resp)
}

func (ctrl *AssistantController) Get(c echo.Context) error {
	ctx := c.Request().Context()
	id, _ := strconv.Atoi(c.Param("id"))
	if id == 0 {
		return response.BadRequest(c, "Invalid assistant ID")
	}

	user, err := contextutil.GetUserFromContext(c)
	if err != nil {
		return response.Unauthorized(c, "Unauthorized")
	}

	doctor, err := ctrl.doctorRepo.GetDoctorByUserID(user.ID)
	if err != nil {
		return response.BadRequest(c, "Doctor profile not found for user")
	}

	resp, err := ctrl.assistantSvc.Get(ctx, id, doctor.ID)
	if err != nil {
		return response.InternalServerError(c, err.Error())
	}

	return response.Success(c, "assistant details fetched successfully", resp)
}

func (ctrl *AssistantController) List(c echo.Context) error {
	ctx := c.Request().Context()
	user, err := contextutil.GetUserFromContext(c)
	if err != nil {
		return response.Unauthorized(c, "Unauthorized")
	}

	doctor, err := ctrl.doctorRepo.GetDoctorByUserID(user.ID)
	if err != nil {
		return response.BadRequest(c, "Doctor profile not found for user")
	}

	resp, err := ctrl.assistantSvc.List(ctx, doctor.ID)
	if err != nil {
		return response.InternalServerError(c, err.Error())
	}

	return response.Success(c, "assistants list fetched successfully", resp)
}
