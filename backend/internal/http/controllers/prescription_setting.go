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

type PrescriptionSettingController struct {
	svc        model.PrescriptionSettingUseCase
	doctorRepo model.DoctorRepo
}

func NewPrescriptionSettingController(
	svc model.PrescriptionSettingUseCase,
	doctorRepo model.DoctorRepo,
) *PrescriptionSettingController {
	return &PrescriptionSettingController{
		svc:        svc,
		doctorRepo: doctorRepo,
	}
}

func (ctrl *PrescriptionSettingController) Upsert(c echo.Context) error {
	ctx := c.Request().Context()

	var req types.PrescriptionSettingReq
	if err := c.Bind(&req); err != nil {
		return response.BadRequest(c, err.Error())
	}

	doctorID, err := strconv.Atoi(c.Param("doctor_id"))
	if err != nil {
		return response.BadRequest(c, "invalid doctor id")
	}
	req.DoctorID = doctorID

	if err := req.Validate(); err != nil {
		return response.BadRequest(c, err.Error())
	}

	user, err := contextutil.GetUserFromContext(c)
	if err != nil {
		return response.Unauthorized(c, "Unauthorized")
	}

	if user.Role != consts.RoleAdmin {
		doctor, err := ctrl.doctorRepo.GetDoctorByUserID(user.ID)
		if err != nil || doctor.ID != req.DoctorID {
			return response.Unauthorized(c, "unauthorized to update settings for this doctor")
		}
	}

	resp, err := ctrl.svc.Upsert(ctx, req)
	if err != nil {
		return response.InternalServerError(c, err.Error())
	}

	return response.Success(c, "prescription settings saved successfully", resp)
}

func (ctrl *PrescriptionSettingController) GetByChamber(c echo.Context) error {
	ctx := c.Request().Context()

	doctorID, _ := strconv.Atoi(c.Param("doctor_id"))
	chamberID, _ := strconv.Atoi(c.QueryParam("chamber_id"))

	if doctorID == 0 || chamberID == 0 {
		return response.BadRequest(c, "doctor_id and chamber_id are required")
	}

	user, err := contextutil.GetUserFromContext(c)
	if err != nil {
		return response.Unauthorized(c, "Unauthorized")
	}

	if user.Role != consts.RoleAdmin {
		doctor, err := ctrl.doctorRepo.GetDoctorByUserID(user.ID)
		if err != nil || doctor.ID != doctorID {
			return response.Unauthorized(c, "unauthorized to view settings for this doctor")
		}
	}

	resp, err := ctrl.svc.GetByChamber(ctx, doctorID, chamberID)
	if err != nil {
		// If not found, we might want to return an empty response or a specific message
		return response.Success(c, "settings not found", nil)
	}

	return response.Success(c, "prescription settings fetched successfully", resp)
}
