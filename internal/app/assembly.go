package app

import (
	"log/slog"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/teryble09/avito_backend/tree/dev/internal/config"
)

// собираем слои приложения и отдаем handler
func assemblyLayers(cfg *config.Config, db *pgxpool.Pool, logger *slog.Logger) (http.Handler, error) {
	return nil, nil
}
