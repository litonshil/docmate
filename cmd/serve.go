package cmd

import (
	"context"
	"docmate/client/conn"
	"docmate/internal/http/controllers"
	httpRoutes "docmate/internal/http/routes"
	httpServer "docmate/internal/http/server"
	"docmate/internal/repositories/db"
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

	_ = txsvc

	userController := controllers.NewUserController(ctx, usersvc)

	routes := httpRoutes.New(e, userController)
	routes.Init()
}

func RegisterServeCommand() {
	RegisterSubCommand(serveCmd)
}
