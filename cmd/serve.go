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
	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use: "serve",
	Run: serve,
}

func serve(cmd *cobra.Command, args []string) {
	// base context
	baseContext := context.Background()

	dbClient := conn.Db()

	dbRepo := db.NewRepository(dbClient)

	txsvc := txservice.NewDBTransaction(dbRepo)
	usersvc := userservice.NewService(dbRepo)

	_ = txsvc

	// HttpServer
	var HttpServer = httpServer.New()

	userController := controllers.NewUserController(
		baseContext,
		usersvc,
	)

	var Routes = httpRoutes.New(
		HttpServer.Echo,
		userController,
	)

	// Spooling
	Routes.Init()
	HttpServer.Start()
}
