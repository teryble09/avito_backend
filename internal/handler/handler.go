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
	prService   PullRequestService
}

func NewOgenHandler(
	logger *slog.Logger,
	teamService TeamService,
	userService UserService,
	prService PullRequestService,
) *OgenHandler {
	return &OgenHandler{
		logger:      logger,
		teamService: teamService,
		userService: userService,
		prService:   prService,
	}
}
