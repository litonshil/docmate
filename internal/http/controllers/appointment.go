package controllers

import (
	"docmate/internal/model"
	"docmate/response"
	"docmate/types"
	"docmate/utils/contextutil"
	"strconv"

	"github.com/labstack/echo/v4"
)

type AppointmentController struct {
	svc        model.AppointmentUseCase
	doctorRepo model.DoctorRepo
}

func NewAppointmentController(svc model.AppointmentUseCase, doctorRepo model.DoctorRepo) *AppointmentController {
	return &AppointmentController{
		svc:        svc,
		doctorRepo: doctorRepo,
	}
}

func (ctrl *AppointmentController) Book(c echo.Context) error {
	ctx := c.Request().Context()
	user, err := contextutil.GetUserFromContext(c)
	if err != nil {
		return response.Unauthorized(c, "Unauthorized")
	}

	doctor, err := ctrl.doctorRepo.GetDoctorByUserID(user.ID)
	if err != nil {
		return response.BadRequest(c, "Doctor profile not found")
	}

	var req types.AppointmentReq
	if err := c.Bind(&req); err != nil {
		return response.BadRequest(c, "Invalid request payload")
	}

	res, err := ctrl.svc.BookAppointment(ctx, req, doctor.ID)
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

	doctor, err := ctrl.doctorRepo.GetDoctorByUserID(user.ID)
	if err != nil {
		return response.BadRequest(c, "Doctor profile not found")
	}

	dateFrom := c.QueryParam("date_from")
	dateTo := c.QueryParam("date_to")
	status := c.QueryParam("status")
	page, _ := strconv.Atoi(c.QueryParam("page"))
	limit, _ := strconv.Atoi(c.QueryParam("limit"))

	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	paginatedResp, err := ctrl.svc.ListAppointments(ctx, doctor.ID, dateFrom, dateTo, status, page, limit)
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

	err := ctrl.svc.UpdateStatus(ctx, id, model.AppointmentStatus(req.Status))
	if err != nil {
		return response.InternalServerError(c, "Failed to update appointment status")
	}

	return response.Success(c, "Appointment status updated successfully", nil)
}

func (ctrl *AppointmentController) GetDetails(c echo.Context) error {
	ctx := c.Request().Context()
	id, _ := strconv.Atoi(c.Param("id"))
	if id == 0 {
		return response.BadRequest(c, "Invalid appointment ID")
	}

	appointment, err := ctrl.svc.GetAppointment(ctx, id)
	if err != nil {
		return response.NotFound(c, "Appointment not found")
	}

	return response.Success(c, "Appointment details fetched successfully", appointment)
}
