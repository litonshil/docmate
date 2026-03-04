package cmd

import (
	"context"
	"docmate/client/conn"
	"docmate/internal/http/controllers"
	httpRoutes "docmate/internal/http/routes"
	httpServer "docmate/internal/http/server"
	"docmate/internal/repositories/db"
	chamberservice "docmate/internal/services/chamber"
	doctorservice "docmate/internal/services/doctor"
	patientservice "docmate/internal/services/patient"
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
		MountRoutes(ctx, HttpServer.Echo)
		HttpServer.Start()
	},
}

func MountRoutes(ctx context.Context, e *echo.Echo) {
	dbClient := conn.Db()
	dbRepo := db.NewRepository(dbClient)

	txsvc := txservice.NewDBTransaction(dbRepo)
	usersvc := userservice.NewService(dbRepo)
	doctorsvc := doctorservice.NewService(dbRepo)
	patientsvc := patientservice.NewService(dbRepo)
	chambersvc := chamberservice.NewService(dbRepo)

	_ = txsvc

	userController := controllers.NewUserController(ctx, usersvc)
	doctorController := controllers.NewDoctorController(ctx, doctorsvc)
	patientController := controllers.NewPatientController(ctx, patientsvc, dbRepo)
	chamberController := controllers.NewChamberController(ctx, chambersvc, dbRepo)

	routes := httpRoutes.New(e, userController, doctorController, patientController, chamberController)
	routes.Init()
}

func RegisterServeCommand() {
	RegisterSubCommand(serveCmd)
}
