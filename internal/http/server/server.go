package server

import (
	"context"
	"docmate/config"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
)

type Server struct {
	Echo *echo.Echo
}

func (s *Server) Start() {
	e := s.Echo

	go func() {
		e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", config.App().Port)))
	}()

	GracefulShutdown(e)
}

func New() *Server {
	return &Server{Echo: echo.New()}
}

// GracefulShutdown server will gracefully shut down within 5 sec.
func GracefulShutdown(e *echo.Echo) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	<-ch
	slog.Info("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_ = e.Shutdown(ctx)
	slog.Info("server shutdowns gracefully")
}
