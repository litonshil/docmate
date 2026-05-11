package appointment

import (
	"context"
	"docmate/internal/model"
	"docmate/types"
	"errors"
	"fmt"
	"log/slog"
	"time"
)

type appointmentService struct {
	repo model.AppointmentRepo
}

func NewAppointmentService(repo model.AppointmentRepo) model.AppointmentUseCase {
	return &appointmentService{
		repo: repo,
	}
}

func (s *appointmentService) BookAppointment(ctx context.Context, req types.AppointmentReq, doctorID int) (types.AppointmentResp, error) {
	if err := req.Validate(); err != nil {
		return types.AppointmentResp{}, err
	}

	date, err := time.Parse("2006-01-02", req.AppointmentDate)
	if err != nil {
		return types.AppointmentResp{}, fmt.Errorf("invalid date format: %w", err)
	}

	// Time validation: Ensure appointment is in the future
	// StartTime format: "10:30 AM"
	timeStr := req.StartTime
	t, err := time.Parse("03:04 PM", timeStr)
	if err != nil {
		return types.AppointmentResp{}, fmt.Errorf("invalid time format: %w", err)
	}

	appointmentTime := time.Date(date.Year(), date.Month(), date.Day(), t.Hour(), t.Minute(), 0, 0, time.Local)
	if appointmentTime.Before(time.Now().Add(-1 * time.Minute)) { // Allow 1 min grace for network lag
		return types.AppointmentResp{}, errors.New("appointment time must be in the future")
	}

	appointment := &model.Appointment{
		DoctorID:        doctorID,
		PatientID:       req.PatientID,
		ChamberID:       req.ChamberID,
		AppointmentDate: date,
		StartTime:       req.StartTime,
		EndTime:         req.EndTime,
		Reason:          req.Reason,
		Notes:           req.Notes,
		VisitingFee:     req.VisitingFee,
		Status:          model.AppointmentStatusPending,
	}

	// Handle Quick Add Patient
	if appointment.PatientID == 0 && req.QuickPatient != nil {
		quickPatient := model.Patient{
			DoctorID: doctorID,
			FullName: req.QuickPatient.FullName,
			Phone:    req.QuickPatient.Phone,
			Gender:   req.QuickPatient.Gender,
			Age:      req.QuickPatient.Age,
		}

		patientRepo := s.repo.(model.PatientRepo)
		p, err := patientRepo.CreatePatient(quickPatient)
		if err != nil {
			slog.Error("failed to quick add patient", "error", err.Error())

			return types.AppointmentResp{}, err
		}
		appointment.PatientID = p.ID
	}

	if appointment.PatientID == 0 {
		return types.AppointmentResp{}, errors.New("patient ID is required")
	}

	err = s.repo.CreateAppointment(appointment)
	if err != nil {
		slog.Error("failed to create appointment", "error", err.Error())

		return types.AppointmentResp{}, err
	}

	fullApp, err := s.repo.GetAppointmentByID(appointment.ID)
	if err != nil {
		return types.AppointmentResp{}, err
	}

	return mapToAppointmentResponse(*fullApp), nil
}

func (s *appointmentService) UpdateStatus(ctx context.Context, id int, status model.AppointmentStatus) error {
	appointment, err := s.repo.GetAppointmentByID(id)
	if err != nil {
		return err
	}

	appointment.Status = status
	err = s.repo.UpdateAppointment(appointment)
	if err != nil {
		slog.Error("failed to update appointment status", "id", id, "status", status, "error", err.Error())
	}

	return err
}

func (s *appointmentService) CollectFee(ctx context.Context, id int, amount float64) error {
	appointment, err := s.repo.GetAppointmentByID(id)
	if err != nil {
		return err
	}

	updates := map[string]interface{}{
		"is_fee_collected": true,
	}
	if amount > 0 {
		updates["visiting_fee"] = amount
	} else if appointment.VisitingFee > 0 {
		updates["visiting_fee"] = appointment.VisitingFee
	}

	err = s.repo.UpdateAppointmentFields(id, updates)
	if err != nil {
		slog.Error("failed to collect fee", "id", id, "error", err.Error())
	}

	return err
}

func (s *appointmentService) GetAppointment(ctx context.Context, id int) (types.AppointmentResp, error) {
	app, err := s.repo.GetAppointmentByID(id)
	if err != nil {
		return types.AppointmentResp{}, err
	}

	return mapToAppointmentResponse(*app), nil
}

func (s *appointmentService) ListAppointments(ctx context.Context, doctorID int, dateFromStr, dateToStr string, status string, search string, page, limit int) (types.PaginatedResponse[types.AppointmentResp], error) {
	var dateFrom, dateTo *time.Time

	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	if dateFromStr != "" {
		t, err := time.Parse("2006-01-02", dateFromStr)
		if err == nil {
			dateFrom = &t
		}
	}

	if dateToStr != "" {
		t, err := time.Parse("2006-01-02", dateToStr)
		if err == nil {
			dateTo = &t
		}
	}

	appointments, total, err := s.repo.ListAppointments(doctorID, dateFrom, dateTo, status, search, page, limit)
	if err != nil {
		slog.Error("failed to list appointments", "doctor_id", doctorID, "error", err.Error())

		return types.PaginatedResponse[types.AppointmentResp]{}, err
	}

	var records []types.AppointmentResp
	for _, app := range appointments {
		records = append(records, mapToAppointmentResponse(app))
	}

	lastPage := (int(total) + limit - 1) / limit
	if lastPage == 0 {
		lastPage = 1
	}

	return types.PaginatedResponse[types.AppointmentResp]{
		Pagination: types.Pagination{
			Total:    total,
			Page:     page,
			Limit:    limit,
			LastPage: lastPage,
		},
		Records: records,
	}, nil
}

func mapToAppointmentResponse(app model.Appointment) types.AppointmentResp {
	resp := types.AppointmentResp{
		ID:              app.ID,
		DoctorID:        app.DoctorID,
		PatientID:       app.PatientID,
		ChamberID:       app.ChamberID,
		AppointmentDate: app.AppointmentDate,
		StartTime:       app.StartTime,
		EndTime:         app.EndTime,
		Status:          string(app.Status),
		Reason:          app.Reason,
		Notes:           app.Notes,
		VisitingFee:     app.VisitingFee,
		IsFeeCollected:  app.IsFeeCollected,
		CreatedAt:       app.CreatedAt,
		UpdatedAt:       app.UpdatedAt,
	}

	if app.Patient.ID != 0 {
		resp.Patient = &types.PatientResp{
			ID:       app.Patient.ID,
			FullName: app.Patient.FullName,
			Phone:    app.Patient.Phone,
		}
	}

	if app.Chamber.ID != 0 {
		resp.Chamber = &types.ChamberResp{
			ID:   app.Chamber.ID,
			Name: app.Chamber.Name,
		}
	}

	return resp
}
