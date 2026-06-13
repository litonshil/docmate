package dashboard

import (
	"context"
	"docmate/internal/model"
	"docmate/types"
	"time"

	"gorm.io/gorm"
)

type dashboardRepo struct {
	db *gorm.DB
}

func NewDashboardRepo(db *gorm.DB) model.DashboardRepo {
	return &dashboardRepo{db: db}
}

func (r *dashboardRepo) GetTotalPatients(ctx context.Context, doctorID int) (int, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(&model.Patient{})
	if doctorID != 0 {
		query = query.Where("doctor_id = ?", doctorID)
	}
	err := query.Count(&count).Error

	return int(count), err
}

func (r *dashboardRepo) GetTodayVisits(ctx context.Context, doctorID int) (int, error) {
	var count int64
	today := time.Now().Truncate(24 * time.Hour)
	query := r.db.WithContext(ctx).Model(&model.Prescription{}).Where("created_at >= ?", today)
	if doctorID != 0 {
		query = query.Where("doctor_id = ?", doctorID)
	}
	err := query.Count(&count).Error

	return int(count), err
}

func (r *dashboardRepo) GetTotalPrescriptions(ctx context.Context, doctorID int) (int, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(&model.Prescription{})
	if doctorID != 0 {
		query = query.Where("doctor_id = ?", doctorID)
	}
	err := query.Count(&count).Error

	return int(count), err
}

func (r *dashboardRepo) GetActiveMedicines(ctx context.Context, doctorID int) (int, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.Medicine{}).Where("is_active = ?", true).Count(&count).Error

	return int(count), err
}

func (r *dashboardRepo) GetRecentPatients(ctx context.Context, doctorID int, limit int) ([]types.PatientSummary, error) {
	var patients []model.Patient
	query := r.db.WithContext(ctx).Model(&model.Patient{})
	if doctorID != 0 {
		query = query.Where("doctor_id = ?", doctorID)
	}
	err := query.Order("created_at desc").Limit(limit).Find(&patients).Error

	if err != nil {
		return nil, err
	}

	var summaries []types.PatientSummary
	for _, p := range patients {
		summaries = append(summaries, types.PatientSummary{
			ID:        p.ID,
			Name:      p.FullName,
			Gender:    p.Gender,
			Age:       p.Age,            // Simplification: Using raw age string from DB
			LastVisit: "Recently Added", // Simplification for now, would require joining with prescriptions
		})
	}

	return summaries, nil
}

func (r *dashboardRepo) GetTodaySchedule(ctx context.Context, doctorID int) ([]types.ScheduleSummary, error) {
	var appointments []model.Appointment
	today := time.Now().Truncate(24 * time.Hour)
	tomorrow := today.Add(24 * time.Hour)

	query := r.db.WithContext(ctx).Preload("Patient").
		Where("appointment_date >= ? AND appointment_date < ? AND status != ?",
			today, tomorrow, model.AppointmentStatusCancelled)
	if doctorID != 0 {
		query = query.Where("doctor_id = ?", doctorID)
	}
	err := query.Order("start_time asc").Find(&appointments).Error

	if err != nil {
		return nil, err
	}

	var summaries []types.ScheduleSummary
	for _, app := range appointments {
		summaries = append(summaries, types.ScheduleSummary{
			ID:          app.ID,
			PatientID:   app.PatientID,
			PatientName: app.Patient.FullName,
			Time:        app.StartTime,
			Type:        app.Reason,
			Status:      string(app.Status),
		})
	}

	return summaries, nil
}

func (r *dashboardRepo) GetTotalDoctors(ctx context.Context) (int, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.Doctor{}).Count(&count).Error

	return int(count), err
}
