package logger

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/teryble09/avito_backend/tree/dev/internal/config"
)

func Setup(cfg *config.Config) (*slog.Logger, error) {
	level := slog.LevelInfo

	err := level.UnmarshalText([]byte(strings.ToUpper(cfg.LogLevel)))
	if err == nil {
		return nil, fmt.Errorf("incorrect log level format: %s", cfg.LogLevel)
	}

	h := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: level})

	return slog.New(h), nil
}
