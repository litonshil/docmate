package cmd

import (
	"context"
	"docmate/client/conn"
	"docmate/internal/http/controllers"
	"docmate/internal/http/middlewares"
	httpRoutes "docmate/internal/http/routes"
	httpServer "docmate/internal/http/server"
	dashboardrepo "docmate/internal/repositories/dashboard"
	"docmate/internal/repositories/db"
	chamberservice "docmate/internal/services/chamber"
	dashboardservice "docmate/internal/services/dashboard"
	doctorservice "docmate/internal/services/doctor"
	medicineservice "docmate/internal/services/medicine"
	patientservice "docmate/internal/services/patient"
	prescriptionservice "docmate/internal/services/prescription"
	prescriptionsettingservice "docmate/internal/services/prescription_setting"
	txservice "docmate/internal/services/transaction"
	userservice "docmate/internal/services/user"
	"os"
	"os/signal"
	"syscall"

	"github.com/labstack/echo/v4"
	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the Docmate HTTP server",
	Run: func(cmd *cobra.Command, args []string) {
		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
		defer stop()

		HttpServer := httpServer.New()
		middlewares.Init(HttpServer.Echo)
		MountRoutes(ctx, HttpServer.Echo)
		HttpServer.Start()
	},
}

func MountRoutes(ctx context.Context, e *echo.Echo) {
	dbClient := conn.Db()
	dbRepo := db.NewRepository(dbClient)

	txsvc := txservice.NewDBTransaction(dbRepo)
	usersvc := userservice.NewService(dbRepo, dbRepo)
	doctorsvc := doctorservice.NewService(dbRepo)
	patientsvc := patientservice.NewService(dbRepo)
	chambersvc := chamberservice.NewService(dbRepo)
	medicinesvc := medicineservice.NewService(dbRepo)
	prescriptionsvc := prescriptionservice.NewService(dbRepo)
	prescriptionsettingsvc := prescriptionsettingservice.NewService(dbRepo)

	dashrepo := dashboardrepo.NewDashboardRepo(dbClient)
	dashboardsvc := dashboardservice.NewDashboardService(dashrepo)

	_ = txsvc

	userController := controllers.NewUserController(ctx, usersvc)
	doctorController := controllers.NewDoctorController(ctx, doctorsvc)
	patientController := controllers.NewPatientController(ctx, patientsvc, dbRepo)
	chamberController := controllers.NewChamberController(ctx, chambersvc, dbRepo)
	medicineController := controllers.NewMedicineController(ctx, medicinesvc)
	prescriptionController := controllers.NewPrescriptionController(ctx, prescriptionsvc, dbRepo)
	prescriptionSettingController := controllers.NewPrescriptionSettingController(prescriptionsettingsvc, dbRepo)
	dashboardController := controllers.NewDashboardController(ctx, dashboardsvc, dbRepo)

	routes := httpRoutes.New(e, userController, doctorController, patientController, chamberController, medicineController, prescriptionController, prescriptionSettingController, dashboardController)
	routes.Init()
}

func RegisterServeCommand() {
	RegisterSubCommand(serveCmd)
}
