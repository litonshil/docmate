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
	baseCtx       context.Context
	chamberSvc    model.ChamberUseCase
	doctorRepo    model.DoctorRepo
	assistantRepo model.AssistantRepo
}

func NewChamberController(
	baseCtx context.Context,
	chamberSvc model.ChamberUseCase,
	doctorRepo model.DoctorRepo,
	assistantRepo model.AssistantRepo,
) *ChamberController {
	return &ChamberController{
		baseCtx:       baseCtx,
		chamberSvc:    chamberSvc,
		doctorRepo:    doctorRepo,
		assistantRepo: assistantRepo,
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
	user, err := contextutil.GetUserFromContext(c)
	if err != nil {
		return response.Unauthorized(c, "Unauthorized")
	}

	var req types.ChamberListReq
	if err := c.Bind(&req); err != nil {
		return response.BadRequest(c, err.Error())
	}

	// For doctors, force listing only their own chambers.
	if user.Role == consts.RoleDoctor {
		doctor, err := controller.doctorRepo.GetDoctorByUserID(user.ID)
		if err != nil {
			return response.BadRequest(c, "Doctor profile not found for user")
		}
		req.DoctorID = doctor.ID
	} else if user.Role == consts.RoleAssistant {
		assistant, err := controller.assistantRepo.GetAssistantByUserID(user.ID)
		if err != nil {
			return response.BadRequest(c, "Assistant profile not found for user")
		}
		chambers, err := controller.assistantRepo.GetChambersByAssistantID(assistant.ID)
		if err != nil {
			return response.InternalServerError(c, "Failed to get assigned chambers")
		}
		records := make([]types.ChamberResp, 0)
		for _, ch := range chambers {
			records = append(records, types.ChamberResp{
				ID:        ch.ID,
				DoctorID:  ch.DoctorID,
				Name:      ch.Name,
				Address:   ch.Address,
				Phone:     ch.Phone,
				IsActive:  ch.IsActive,
				CreatedAt: ch.CreatedAt,
				UpdatedAt: ch.UpdatedAt,
			})
		}

		return response.Success(c, "chambers fetched successfully", types.PaginatedChamberResp{
			Pagination: types.Pagination{
				Total:    int64(len(records)),
				Page:     1,
				Limit:    1000,
				LastPage: 1,
			},
			Records: records,
		})
	} else if user.Role == consts.RoleAdmin {
		// For admins, doctor ID is optional. If specified, validate existence.
		if req.DoctorID > 0 {
			if _, err := controller.doctorRepo.GetDoctorByID(req.DoctorID); err != nil {
				return response.BadRequest(c, "Doctor not found")
			}
		}
	} else {
		return response.Forbidden(c, "Forbidden")
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

func (controller *ChamberController) Delete(c echo.Context) error {
	ctx := c.Request().Context()

	var filter types.ChamberFilter
	if err := c.Bind(&filter); err != nil {
		return response.BadRequest(c, "invalid parameters")
	}

	if err := filter.Validate(); err != nil {
		return response.BadRequest(c, err.Error())
	}

	// 1. Get existing chamber to verify existence and doctor association
	existing, err := controller.chamberSvc.Get(ctx, filter)
	if err != nil {
		return response.NotFound(c, "chamber not found")
	}

	// Verify that the chamber belongs to the doctor in the URL
	if existing.DoctorID != filter.DoctorID {
		return response.BadRequest(c, "chamber does not belong to this doctor")
	}

	// 2. Authorize (Owner or Admin)
	if err := controller.authorizeDoctor(c, existing.DoctorID); err != nil {
		return err
	}

	// 3. Delete
	if err := controller.chamberSvc.Delete(ctx, filter.ID); err != nil {
		return response.InternalServerError(c, err.Error())
	}

	return response.Success(c, "chamber deleted successfully", nil)
}
