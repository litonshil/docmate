package logger

import (
	"docmate/config"
	"log/slog"
	"os"

	"gopkg.in/natefinch/lumberjack.v2"

	slogmulti "github.com/samber/slog-multi"
)

// logger is the default logger used by the application.
var logger *slog.Logger

// Set sets the logger configuration based on the environment.
func Set(config config.AppConfig) {
	logger = slog.New(
		slog.NewJSONHandler(os.Stderr, nil),
	)

	if config.ENV == "production" {
		logRotate := &lumberjack.Logger{
			Filename:   "log/app.log",
			MaxSize:    100, // megabytes
			MaxBackups: 3,
			MaxAge:     28, // days
			Compress:   true,
		}

		logger = slog.New(
			slogmulti.Fanout(
				slog.NewJSONHandler(logRotate, nil),
				slog.NewTextHandler(os.Stderr, nil),
			),
		)
	}

	slog.SetDefault(logger)
}
