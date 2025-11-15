package handler

import (
	"log/slog"

	api "github.com/teryble09/avito_backend/generated"
)

type OgenHandler struct {
	api.UnimplementedHandler

	logger *slog.Logger

	teamService TeamService
	userService UserService
}

func NewOgenHandler(
	logger *slog.Logger,
	teamService TeamService,
	userService UserService,
) *OgenHandler {
	return &OgenHandler{
		logger:      logger,
		teamService: teamService,
		userService: userService,
	}
}
