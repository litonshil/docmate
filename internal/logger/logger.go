package logger

import (
	"log/slog"
	"os"
)

var (
	defaultLogger *slog.Logger
)

func init() {
	opts := &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true,
	}

	handler := slog.NewJSONHandler(os.Stdout, opts)
	defaultLogger = slog.New(handler)
	slog.SetDefault(defaultLogger)
}

// Debug logs at debug level
func Debug(msg string, args ...any) {
	defaultLogger.Debug(msg, args...)
}

// Info logs at info level
func Info(msg string, args ...any) {
	defaultLogger.Info(msg, args...)
}

// Warn logs at warn level
func Warn(msg string, args ...any) {
	defaultLogger.Warn(msg, args...)
}

// Error logs at error level
func Error(msg string, args ...any) {
	defaultLogger.Error(msg, args...)
}

// With returns a new Logger with the args appended to its context
func With(args ...any) *slog.Logger {
	return defaultLogger.With(args...)
}

// WithGroup returns a new Logger that starts a group
func WithGroup(name string) *slog.Logger {
	return defaultLogger.WithGroup(name)
}
