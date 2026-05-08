package controllers

import (
	"context"
	"docmate/internal/model"
	"docmate/response"
	"docmate/types"
	"docmate/utils/contextutil"
	"log/slog"

	"github.com/labstack/echo/v4"
)

type AISuggestionController struct {
	baseCtx context.Context
	svc     model.AISettingUseCase
	docRepo model.DoctorRepo
}

func NewAISuggestionController(baseCtx context.Context, svc model.AISettingUseCase, docRepo model.DoctorRepo) *AISuggestionController {
	return &AISuggestionController{
		baseCtx: baseCtx,
		svc:     svc,
		docRepo: docRepo,
	}
}

func (c *AISuggestionController) GetSuggestions(echoCtx echo.Context) error {
	ctx := echoCtx.Request().Context()
	user, err := contextutil.GetUserFromContext(echoCtx)
	if err != nil {
		return response.Unauthorized(echoCtx, "Unauthorized")
	}

	doc, err := c.docRepo.GetDoctorByUserID(user.ID)
	if err != nil {
		return response.BadRequest(echoCtx, "doctor profile not found")
	}

	var req types.AISuggestionReq
	if err := echoCtx.Bind(&req); err != nil {
		return response.BadRequest(echoCtx, err.Error())
	}

	resp, err := c.svc.GetSuggestions(ctx, doc.ID, req.Complaints)
	if err != nil {
		slog.Error("failed to get suggestions", "error", err.Error(), "doctor_id", doc.ID)

		return response.InternalServerError(echoCtx, err.Error())
	}

	return response.Success(echoCtx, "suggestions generated successfully", resp)
}

func (c *AISuggestionController) UpsertSettings(echoCtx echo.Context) error {
	ctx := echoCtx.Request().Context()
	user, err := contextutil.GetUserFromContext(echoCtx)
	if err != nil {
		return response.Unauthorized(echoCtx, "Unauthorized")
	}

	doc, err := c.docRepo.GetDoctorByUserID(user.ID)
	if err != nil {
		return response.BadRequest(echoCtx, "doctor profile not found")
	}

	var req types.AISettingReq
	if err := echoCtx.Bind(&req); err != nil {
		return response.BadRequest(echoCtx, err.Error())
	}
	req.DoctorID = doc.ID

	resp, err := c.svc.Upsert(ctx, req)
	if err != nil {
		return response.InternalServerError(echoCtx, err.Error())
	}

	return response.Success(echoCtx, "ai settings updated successfully", resp)
}

func (c *AISuggestionController) AdminUpdateSettings(echoCtx echo.Context) error {
	ctx := echoCtx.Request().Context()

	var req types.AdminAISettingUpdateReq
	if err := echoCtx.Bind(&req); err != nil {
		return response.BadRequest(echoCtx, err.Error())
	}

	resp, err := c.svc.AdminUpdate(ctx, req)
	if err != nil {
		return response.InternalServerError(echoCtx, err.Error())
	}

	return response.Success(echoCtx, "ai settings updated by admin", resp)
}

func (c *AISuggestionController) GetSettings(echoCtx echo.Context) error {
	ctx := echoCtx.Request().Context()
	user, err := contextutil.GetUserFromContext(echoCtx)
	if err != nil {
		return response.Unauthorized(echoCtx, "Unauthorized")
	}

	doc, err := c.docRepo.GetDoctorByUserID(user.ID)
	if err != nil {
		return response.BadRequest(echoCtx, "doctor profile not found")
	}

	resp, err := c.svc.GetByDoctor(ctx, doc.ID)
	if err != nil {
		return response.InternalServerError(echoCtx, err.Error())
	}

	return response.Success(echoCtx, "ai settings fetched successfully", resp)
}
