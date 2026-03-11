package dashboard

import (
	"context"
	"docmate/internal/model"
	"docmate/types"

	"golang.org/x/sync/errgroup"
)

type dashboardService struct {
	repo model.DashboardRepo
}

func NewDashboardService(repo model.DashboardRepo) model.DashboardUseCase {
	return &dashboardService{repo: repo}
}

func (s *dashboardService) GetSummary(ctx context.Context, doctorID int) (types.DashboardSummaryResp, error) {
	var resp types.DashboardSummaryResp
	g, ctx := errgroup.WithContext(ctx)

	// Fetch Total Patients
	g.Go(func() error {
		count, err := s.repo.GetTotalPatients(ctx, doctorID)
		if err == nil {
			resp.TotalPatients = count
		}

		return err
	})

	// Fetch Today's Visits
	g.Go(func() error {
		count, err := s.repo.GetTodayVisits(ctx, doctorID)
		if err == nil {
			resp.TodayVisits = count
		}

		return err
	})

	// Fetch Total Prescriptions
	g.Go(func() error {
		count, err := s.repo.GetTotalPrescriptions(ctx, doctorID)
		if err == nil {
			resp.TotalPrescriptions = count
		}

		return err
	})

	// Fetch Active Medicines
	g.Go(func() error {
		count, err := s.repo.GetActiveMedicines(ctx, doctorID)
		if err == nil {
			resp.ActiveMedicines = count
		}

		return err
	})

	// Fetch Recent Patients
	g.Go(func() error {
		patients, err := s.repo.GetRecentPatients(ctx, doctorID, 5) // Fetch top 5
		if err == nil {
			resp.RecentPatients = patients
		}

		return err
	})

	// Fetch Today's Schedule
	g.Go(func() error {
		schedule, err := s.repo.GetTodaySchedule(ctx, doctorID)
		if err == nil {
			resp.TodaySchedule = schedule
		}

		return err
	})

	// Wait for all queries to finish concurrently
	if err := g.Wait(); err != nil {
		return types.DashboardSummaryResp{}, err
	}

	// Ensure slices are initialized to [] instead of null in JSON response if empty
	if resp.RecentPatients == nil {
		resp.RecentPatients = []types.PatientSummary{}
	}
	if resp.TodaySchedule == nil {
		resp.TodaySchedule = []types.ScheduleSummary{}
	}

	return resp, nil
}
