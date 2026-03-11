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
	err := r.db.WithContext(ctx).Model(&model.Patient{}).Where("doctor_id = ?", doctorID).Count(&count).Error

	return int(count), err
}

func (r *dashboardRepo) GetTodayVisits(ctx context.Context, doctorID int) (int, error) {
	var count int64
	today := time.Now().Truncate(24 * time.Hour)
	err := r.db.WithContext(ctx).Model(&model.Prescription{}).
		Where("doctor_id = ? AND created_at >= ?", doctorID, today).
		Count(&count).Error

	return int(count), err
}

func (r *dashboardRepo) GetTotalPrescriptions(ctx context.Context, doctorID int) (int, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.Prescription{}).Where("doctor_id = ?", doctorID).Count(&count).Error

	return int(count), err
}

func (r *dashboardRepo) GetActiveMedicines(ctx context.Context, doctorID int) (int, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.Medicine{}).Where("is_active = ?", true).Count(&count).Error

	return int(count), err
}

func (r *dashboardRepo) GetRecentPatients(ctx context.Context, doctorID int, limit int) ([]types.PatientSummary, error) {
	var patients []model.Patient
	err := r.db.WithContext(ctx).
		Where("doctor_id = ?", doctorID).
		Order("created_at desc").
		Limit(limit).
		Find(&patients).Error

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
	var prescriptions []model.Prescription
	today := time.Now().Truncate(24 * time.Hour)
	tomorrow := today.Add(24 * time.Hour)

	// Fetch prescriptions where today is the FollowUpDate OR created today
	err := r.db.WithContext(ctx).
		Where("doctor_id = ? AND ((follow_up_date >= ? AND follow_up_date < ?) OR (created_at >= ? AND created_at < ?))",
			doctorID, today, tomorrow, today, tomorrow).
		Order("created_at asc").
		Find(&prescriptions).Error

	if err != nil {
		return nil, err
	}

	var summaries []types.ScheduleSummary
	for _, p := range prescriptions {
		typ := "Checkup"
		if p.FollowUpDate != nil && p.FollowUpDate.After(today) && p.FollowUpDate.Before(tomorrow) {
			typ = "Follow-up"
		}

		summaries = append(summaries, types.ScheduleSummary{
			PrescriptionID: p.ID,
			PatientID:      p.PatientID,
			PatientName:    p.PatientName,
			Time:           p.CreatedAt.Format("03:04 PM"),
			Type:           typ,
		})
	}

	return summaries, nil
}
