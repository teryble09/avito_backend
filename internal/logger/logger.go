package logger

import (
	"log/slog"
	"os"

	"github.com/teryble09/avito_backend/tree/dev/internal/config"
)

func Setup(cfg *config.Config) *slog.Logger {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	return logger
}
