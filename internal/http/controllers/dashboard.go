package controllers

import (
	"context"
	"docmate/internal/consts"
	"docmate/internal/model"
	"docmate/response"
	"docmate/utils/contextutil"

	"github.com/labstack/echo/v4"
)

type DashboardController struct {
	baseCtx      context.Context
	dashboardSvc model.DashboardUseCase
	doctorRepo   model.DoctorRepo
}

func NewDashboardController(
	baseCtx context.Context,
	dashboardSvc model.DashboardUseCase,
	doctorRepo model.DoctorRepo,
) *DashboardController {
	return &DashboardController{
		baseCtx:      baseCtx,
		dashboardSvc: dashboardSvc,
		doctorRepo:   doctorRepo,
	}
}

func (controller *DashboardController) GetSummary(c echo.Context) error {
	ctx := c.Request().Context()

	// Extract authenticated User
	user, err := contextutil.GetUserFromContext(c)
	if err != nil {
		return response.Unauthorized(c, "Unauthorized")
	}

	var doctorID int
	if user.Role == consts.RoleDoctor {
		doctor, err := controller.doctorRepo.GetDoctorByUserID(user.ID)
		if err != nil {
			return response.BadRequest(c, "Doctor profile not found for user")
		}
		doctorID = doctor.ID
	} else if user.Role == consts.RoleAdmin {
		doctorID = 0
	} else {
		return response.Forbidden(c, "Forbidden")
	}

	// Fetch Summary
	resp, err := controller.dashboardSvc.GetSummary(ctx, doctorID)
	if err != nil {
		return response.InternalServerError(c, "Failed to load dashboard summary")
	}

	return response.Success(c, "Dashboard summary fetched successfully", resp)
}
