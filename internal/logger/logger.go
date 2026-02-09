package logger

import (
	"log/slog"
	"os"
	"strings"
)

func New(env string) *slog.Logger {
	level := parseLevel(os.Getenv("LOG_LEVEL"))
	if level == nil {
		if env == "dev" {
			l := slog.LevelDebug
			level = &l
		} else {
			l := slog.LevelInfo
			level = &l
		}
	}

	h := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: *level})
	return slog.New(h)
}

func parseLevel(v string) *slog.Level {
	switch strings.ToLower(strings.TrimSpace(v)) {
	case "debug":
		l := slog.LevelDebug
		return &l
	case "info":
		l := slog.LevelInfo
		return &l
	case "warn", "warning":
		l := slog.LevelWarn
		return &l
	case "error":
		l := slog.LevelError
		return &l
	default:
		return nil
	}
}
