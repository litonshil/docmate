package model

import (
	"context"
	"docmate/types"
)

type DashboardUseCase interface {
	GetSummary(ctx context.Context, doctorID int) (types.DashboardSummaryResp, error)
}

type DashboardRepo interface {
	GetTotalPatients(ctx context.Context, doctorID int) (int, error)
	GetTodayVisits(ctx context.Context, doctorID int) (int, error)
	GetTotalPrescriptions(ctx context.Context, doctorID int) (int, error)
	GetActiveMedicines(ctx context.Context, doctorID int) (int, error)
	GetRecentPatients(ctx context.Context, doctorID int, limit int) ([]types.PatientSummary, error)
	GetTodaySchedule(ctx context.Context, doctorID int) ([]types.ScheduleSummary, error)
}
