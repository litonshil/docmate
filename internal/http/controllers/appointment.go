package controllers

import (
	"docmate/internal/consts"
	"docmate/internal/model"
	"docmate/response"
	"docmate/types"
	"docmate/utils/contextutil"
	"strconv"

	"github.com/labstack/echo/v4"
)

type AppointmentController struct {
	svc           model.AppointmentUseCase
	doctorRepo    model.DoctorRepo
	assistantRepo model.AssistantRepo
}

func NewAppointmentController(svc model.AppointmentUseCase, doctorRepo model.DoctorRepo, assistantRepo model.AssistantRepo) *AppointmentController {
	return &AppointmentController{
		svc:           svc,
		doctorRepo:    doctorRepo,
		assistantRepo: assistantRepo,
	}
}

func (ctrl *AppointmentController) Book(c echo.Context) error {
	ctx := c.Request().Context()
	user, err := contextutil.GetUserFromContext(c)
	if err != nil {
		return response.Unauthorized(c, "Unauthorized")
	}

	var doctorID int
	if user.Role == consts.RoleDoctor {
		doctor, err := ctrl.doctorRepo.GetDoctorByUserID(user.ID)
		if err != nil {
			return response.BadRequest(c, "Doctor profile not found")
		}
		doctorID = doctor.ID
	} else if user.Role == consts.RoleAssistant {
		assistant, err := ctrl.assistantRepo.GetAssistantByUserID(user.ID)
		if err != nil {
			return response.BadRequest(c, "Assistant profile not found")
		}
		if !assistant.IsActive {
			return response.Forbidden(c, "Assistant account is inactive")
		}
		doctorID = assistant.DoctorID
	} else {
		return response.Forbidden(c, "Forbidden")
	}

	var req types.AppointmentReq
	if err := c.Bind(&req); err != nil {
		return response.BadRequest(c, "Invalid request payload")
	}

	if user.Role == consts.RoleAssistant {
		// Verify assistant is assigned to the chamber
		assistant, _ := ctrl.assistantRepo.GetAssistantByUserID(user.ID)
		chambers, err := ctrl.assistantRepo.GetChambersByAssistantID(assistant.ID)
		if err != nil {
			return response.InternalServerError(c, "Failed to check chamber assignments")
		}
		assigned := false
		for _, ch := range chambers {
			if ch.ID == req.ChamberID {
				assigned = true

				break
			}
		}
		if !assigned {
			return response.Forbidden(c, "You are not assigned to manage this chamber")
		}
	}

	res, err := ctrl.svc.BookAppointment(ctx, req, doctorID)
	if err != nil {
		return response.InternalServerError(c, err.Error())
	}

	return response.Created(c, "Appointment booked successfully", res)
}

func (ctrl *AppointmentController) List(c echo.Context) error {
	ctx := c.Request().Context()
	user, err := contextutil.GetUserFromContext(c)
	if err != nil {
		return response.Unauthorized(c, "Unauthorized")
	}

	var doctorID int
	var chamberIDs []int
	if user.Role == consts.RoleDoctor {
		doctor, err := ctrl.doctorRepo.GetDoctorByUserID(user.ID)
		if err != nil {
			return response.BadRequest(c, "Doctor profile not found")
		}
		doctorID = doctor.ID
	} else if user.Role == consts.RoleAssistant {
		assistant, err := ctrl.assistantRepo.GetAssistantByUserID(user.ID)
		if err != nil {
			return response.BadRequest(c, "Assistant profile not found")
		}
		if !assistant.IsActive {
			return response.Forbidden(c, "Assistant account is inactive")
		}
		doctorID = assistant.DoctorID
		chambers, err := ctrl.assistantRepo.GetChambersByAssistantID(assistant.ID)
		if err != nil {
			return response.InternalServerError(c, "Failed to check chamber assignments")
		}
		for _, ch := range chambers {
			chamberIDs = append(chamberIDs, ch.ID)
		}
		// If assistant has no chambers assigned, return empty results
		if len(chamberIDs) == 0 {
			return response.Success(c, "Appointments fetched successfully", types.PaginatedResponse[types.AppointmentResp]{
				Pagination: types.Pagination{Total: 0, Page: 1, Limit: 10, LastPage: 1},
				Records:    []types.AppointmentResp{},
			})
		}
	} else if user.Role == consts.RoleAdmin {
		doctorID = 0
	} else {
		return response.Forbidden(c, "Forbidden")
	}

	dateFrom := c.QueryParam("date_from")
	dateTo := c.QueryParam("date_to")
	status := c.QueryParam("status")
	search := c.QueryParam("search")
	page, _ := strconv.Atoi(c.QueryParam("page"))
	limit, _ := strconv.Atoi(c.QueryParam("limit"))

	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	paginatedResp, err := ctrl.svc.ListAppointments(ctx, doctorID, chamberIDs, dateFrom, dateTo, status, search, page, limit)
	if err != nil {
		return response.InternalServerError(c, "Failed to list appointments")
	}

	return response.Success(c, "Appointments fetched successfully", paginatedResp)
}

func (ctrl *AppointmentController) UpdateStatus(c echo.Context) error {
	ctx := c.Request().Context()
	id, _ := strconv.Atoi(c.Param("id"))
	if id == 0 {
		return response.BadRequest(c, "Invalid appointment ID")
	}

	var req types.UpdateAppointmentStatusReq
	if err := c.Bind(&req); err != nil {
		return response.BadRequest(c, "Invalid request payload")
	}

	user, err := contextutil.GetUserFromContext(c)
	if err != nil {
		return response.Unauthorized(c, "Unauthorized")
	}

	if user.Role == consts.RoleAssistant {
		assistant, err := ctrl.assistantRepo.GetAssistantByUserID(user.ID)
		if err != nil {
			return response.BadRequest(c, "Assistant profile not found")
		}
		if !assistant.IsActive {
			return response.Forbidden(c, "Assistant account is inactive")
		}
		// Fetch appointment to check chamber assignment
		app, err := ctrl.svc.GetAppointment(ctx, id)
		if err != nil {
			return response.NotFound(c, "Appointment not found")
		}
		chambers, err := ctrl.assistantRepo.GetChambersByAssistantID(assistant.ID)
		if err != nil {
			return response.InternalServerError(c, "Failed to check chamber assignments")
		}
		assigned := false
		for _, ch := range chambers {
			if ch.ID == app.ChamberID {
				assigned = true

				break
			}
		}
		if !assigned {
			return response.Forbidden(c, "You are not assigned to manage this chamber")
		}
	}

	err = ctrl.svc.UpdateStatus(ctx, id, model.AppointmentStatus(req.Status))
	if err != nil {
		return response.InternalServerError(c, "Failed to update appointment status")
	}

	return response.Success(c, "Appointment status updated successfully", nil)
}

func (ctrl *AppointmentController) CollectFee(c echo.Context) error {
	ctx := c.Request().Context()
	id, _ := strconv.Atoi(c.Param("id"))
	if id == 0 {
		return response.BadRequest(c, "Invalid appointment ID")
	}

	var req types.CollectFeeReq
	if err := c.Bind(&req); err != nil {
		return response.BadRequest(c, "Invalid request body")
	}

	user, err := contextutil.GetUserFromContext(c)
	if err != nil {
		return response.Unauthorized(c, "Unauthorized")
	}

	if user.Role == consts.RoleAssistant {
		assistant, err := ctrl.assistantRepo.GetAssistantByUserID(user.ID)
		if err != nil {
			return response.BadRequest(c, "Assistant profile not found")
		}
		if !assistant.IsActive {
			return response.Forbidden(c, "Assistant account is inactive")
		}
		// Fetch appointment to check chamber assignment
		app, err := ctrl.svc.GetAppointment(ctx, id)
		if err != nil {
			return response.NotFound(c, "Appointment not found")
		}
		chambers, err := ctrl.assistantRepo.GetChambersByAssistantID(assistant.ID)
		if err != nil {
			return response.InternalServerError(c, "Failed to check chamber assignments")
		}
		assigned := false
		for _, ch := range chambers {
			if ch.ID == app.ChamberID {
				assigned = true

				break
			}
		}
		if !assigned {
			return response.Forbidden(c, "You are not assigned to manage this chamber")
		}
	}

	err = ctrl.svc.CollectFee(ctx, id, req.Amount)
	if err != nil {
		return response.InternalServerError(c, "Failed to collect fee")
	}

	return response.Success(c, "Fee collected successfully", nil)
}

func (ctrl *AppointmentController) GetDetails(c echo.Context) error {
	ctx := c.Request().Context()
	id, _ := strconv.Atoi(c.Param("id"))
	if id == 0 {
		return response.BadRequest(c, "Invalid appointment ID")
	}

	user, err := contextutil.GetUserFromContext(c)
	if err != nil {
		return response.Unauthorized(c, "Unauthorized")
	}

	appointment, err := ctrl.svc.GetAppointment(ctx, id)
	if err != nil {
		return response.NotFound(c, "Appointment not found")
	}

	if user.Role == consts.RoleAssistant {
		assistant, err := ctrl.assistantRepo.GetAssistantByUserID(user.ID)
		if err != nil {
			return response.BadRequest(c, "Assistant profile not found")
		}
		if !assistant.IsActive {
			return response.Forbidden(c, "Assistant account is inactive")
		}
		chambers, err := ctrl.assistantRepo.GetChambersByAssistantID(assistant.ID)
		if err != nil {
			return response.InternalServerError(c, "Failed to check chamber assignments")
		}
		assigned := false
		for _, ch := range chambers {
			if ch.ID == appointment.ChamberID {
				assigned = true

				break
			}
		}
		if !assigned {
			return response.Forbidden(c, "You are not assigned to manage this chamber")
		}
	}

	return response.Success(c, "Appointment details fetched successfully", appointment)
}
