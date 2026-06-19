package controllers

import (
	"context"
	"docmate/internal/consts"
	"docmate/internal/model"
	"docmate/response"
	"docmate/types"
	"docmate/utils/contextutil"
	"log/slog"
	"strconv"

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

	if user.Role == consts.RoleAdmin {
		var req types.AISettingReq
		if err := echoCtx.Bind(&req); err != nil {
			return response.BadRequest(echoCtx, err.Error())
		}
		resp := types.AISettingResp{
			IsAIEnabled:      true,
			AllowGlobalAPI:   true,
			AIProviderID:     req.AIProviderID,
			UseIndividualKey: req.UseIndividualKey,
		}

		return response.Success(echoCtx, "ai settings updated successfully", resp)
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

func (c *AISuggestionController) AdminGetSettings(echoCtx echo.Context) error {
	ctx := echoCtx.Request().Context()

	doctorID, err := strconv.Atoi(echoCtx.Param("id"))
	if err != nil {
		return response.BadRequest(echoCtx, "invalid doctor id")
	}

	resp, err := c.svc.GetByDoctor(ctx, doctorID)
	if err != nil {
		return response.InternalServerError(echoCtx, err.Error())
	}

	return response.Success(echoCtx, "ai settings fetched successfully", resp)
}

func (c *AISuggestionController) AdminUpdateSettings(echoCtx echo.Context) error {
	ctx := echoCtx.Request().Context()

	doctorID, err := strconv.Atoi(echoCtx.Param("id"))
	if err != nil {
		return response.BadRequest(echoCtx, "invalid doctor id")
	}

	var req types.AdminAISettingUpdateReq
	if err := echoCtx.Bind(&req); err != nil {
		return response.BadRequest(echoCtx, err.Error())
	}
	req.DoctorID = doctorID

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

	if user.Role == consts.RoleAdmin {
		// Mock response for admin so they can access the AI settings page
		resp := types.AISettingResp{
			IsAIEnabled:      true,
			AllowGlobalAPI:   true,
			AIProviderID:     1,
			ProviderSlug:     "gemini",
			UseIndividualKey: false,
		}

		return response.Success(echoCtx, "ai settings fetched successfully", resp)
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

func (c *AISuggestionController) GetGlobalSettings(echoCtx echo.Context) error {
	ctx := echoCtx.Request().Context()

	providers, err := c.svc.GetProviders(ctx)
	if err != nil {
		return response.InternalServerError(echoCtx, err.Error())
	}

	return response.Success(echoCtx, "global settings fetched successfully", providers)
}

func (c *AISuggestionController) UpdateGlobalSettings(echoCtx echo.Context) error {
	ctx := echoCtx.Request().Context()

	var req []types.AIProviderConfig
	if err := echoCtx.Bind(&req); err != nil {
		return response.BadRequest(echoCtx, err.Error())
	}

	err := c.svc.UpdateProviders(ctx, req)
	if err != nil {
		return response.InternalServerError(echoCtx, err.Error())
	}

	return response.Success(echoCtx, "global settings updated successfully", nil)
}
