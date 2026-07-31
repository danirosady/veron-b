package logger

import (
	"log/slog"
	"os"
	"strings"
)

var log *slog.Logger

func InitLogger(env string) {
	opts := &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}

	if strings.ToLower(env) == "development" {
		opts.Level = slog.LevelDebug
	}

	handler := slog.NewJSONHandler(os.Stdout, opts)
	log = slog.New(handler)
	slog.SetDefault(log)
}

func Info(msg string, args ...any) {
	log.Info(msg, args...)
}

func Warn(msg string, args ...any) {
	log.Warn(msg, args...)
}

func Error(msg string, args ...any) {
	log.Error(msg, args...)
}

func Debug(msg string, args ...any) {
	log.Debug(msg, args...)
}

func With(args ...any) *slog.Logger {
	return log.With(args...)
}
