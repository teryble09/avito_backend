package logger

import (
	"log/slog"
	"os"
	"strings"

	"github.com/teryble09/avito_backend/internal/config"
)

func Setup(cfg *config.Config) (*slog.Logger, error) {
	var level slog.Level

	switch strings.ToLower(cfg.LogLevel) {
	case "debug":
		level = slog.LevelDebug
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo // по умолчанию
	}

	h := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: level})

	return slog.New(h), nil
}
