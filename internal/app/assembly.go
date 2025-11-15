package app

import (
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/teryble09/avito_backend/internal/handler"
	"github.com/teryble09/avito_backend/internal/repository"
	"github.com/teryble09/avito_backend/internal/service"
)

// собираем слои приложения и отдаем handler.
func assemblyLayers(db *pgxpool.Pool, logger *slog.Logger) *handler.OgenHandler {
	teamRepo := repository.NewTeamRepo(db)
	userRepo := repository.NewUserRepo(db)

	teamService := service.NewTeamService(db, teamRepo, userRepo)
	userService := service.NewUserService(userRepo)

	handler := handler.NewOgenHandler(logger, teamService, userService)

	return handler
}
