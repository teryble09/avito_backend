package handler

import (
	"log/slog"

	api "github.com/teryble09/avito_backend/generated"
)

type OgenHandler struct {
	api.UnimplementedHandler

	logger *slog.Logger

	teamService TeamService
}

func NewOgenHandler(logger *slog.Logger, teamService TeamService) *OgenHandler {
	return &OgenHandler{
		logger:      logger,
		teamService: teamService,
	}
}
