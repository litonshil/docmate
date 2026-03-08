package controllers

import (
	"context"
	"docmate/internal/consts"
	"docmate/internal/model"
	"docmate/response"
	"docmate/types"
	"docmate/utils/contextutil"
	"strconv"

	"github.com/labstack/echo/v4"
)

type PrescriptionController struct {
	baseCtx context.Context
	svc     model.PrescriptionUseCase
	docRepo model.DoctorRepo
}

func NewPrescriptionController(baseCtx context.Context, svc model.PrescriptionUseCase, docRepo model.DoctorRepo) *PrescriptionController {
	return &PrescriptionController{
		baseCtx: baseCtx,
		svc:     svc,
		docRepo: docRepo,
	}
}

func (c *PrescriptionController) Create(echoCtx echo.Context) error {
	ctx := echoCtx.Request().Context()
	user, err := contextutil.GetUserFromContext(echoCtx)
	if err != nil {
		return response.Unauthorized(echoCtx, "Unauthorized")
	}

	doc, err := c.docRepo.GetDoctorByUserID(user.ID)
	if err != nil {
		return response.BadRequest(echoCtx, "doctor profile not found")
	}

	var req types.PrescriptionReq
	if err := echoCtx.Bind(&req); err != nil {
		return response.BadRequest(echoCtx, err.Error())
	}

	req.DoctorID = doc.ID

	resp, err := c.svc.Create(ctx, req)
	if err != nil {
		return response.InternalServerError(echoCtx, err.Error())
	}

	return response.Success(echoCtx, "prescription created successfully", resp)
}

func (c *PrescriptionController) Update(echoCtx echo.Context) error {
	ctx := echoCtx.Request().Context()
	user, err := contextutil.GetUserFromContext(echoCtx)
	if err != nil {
		return response.Unauthorized(echoCtx, "Unauthorized")
	}

	doc, err := c.docRepo.GetDoctorByUserID(user.ID)
	if err != nil {
		return response.BadRequest(echoCtx, "doctor profile not found")
	}

	idParam := echoCtx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return response.BadRequest(echoCtx, "invalid id")
	}

	var req types.PrescriptionReq
	if err := echoCtx.Bind(&req); err != nil {
		return response.BadRequest(echoCtx, err.Error())
	}

	req.DoctorID = doc.ID

	resp, err := c.svc.Update(ctx, id, req)
	if err != nil {
		return response.InternalServerError(echoCtx, err.Error())
	}

	return response.Success(echoCtx, "prescription updated successfully", resp)
}

func (c *PrescriptionController) Get(echoCtx echo.Context) error {
	ctx := echoCtx.Request().Context()
	user, err := contextutil.GetUserFromContext(echoCtx)
	if err != nil {
		return response.Unauthorized(echoCtx, "Unauthorized")
	}

	doc, err := c.docRepo.GetDoctorByUserID(user.ID)
	if err != nil && user.Role != consts.RoleAdmin {
		return response.BadRequest(echoCtx, "doctor profile not found")
	}

	idParam := echoCtx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return response.BadRequest(echoCtx, "invalid id")
	}

	resp, err := c.svc.Get(ctx, id, doc.ID)
	if err != nil {
		return response.InternalServerError(echoCtx, err.Error())
	}

	return response.Success(echoCtx, "prescription fetched successfully", resp)
}

func (c *PrescriptionController) List(echoCtx echo.Context) error {
	ctx := echoCtx.Request().Context()
	user, err := contextutil.GetUserFromContext(echoCtx)
	if err != nil {
		return response.Unauthorized(echoCtx, "Unauthorized")
	}

	doc, err := c.docRepo.GetDoctorByUserID(user.ID)
	if err != nil && user.Role != consts.RoleAdmin {
		return response.BadRequest(echoCtx, "doctor profile not found")
	}

	var req types.PrescriptionListReq
	if err := echoCtx.Bind(&req); err != nil {
		return response.BadRequest(echoCtx, err.Error())
	}

	req.DoctorID = doc.ID
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 {
		req.Limit = 10
	}

	resp, err := c.svc.List(ctx, req)
	if err != nil {
		return response.InternalServerError(echoCtx, err.Error())
	}

	return response.Success(echoCtx, "prescriptions fetched successfully", resp)
}
