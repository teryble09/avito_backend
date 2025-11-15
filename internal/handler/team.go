package handler

import (
	"context"
	"log/slog"

	api "github.com/teryble09/avito_backend/generated"
	"github.com/teryble09/avito_backend/internal/domain"
)

type TeamService interface {
	CreateTeam(ctx context.Context, team *domain.Team) (*domain.Team, error)
}

func (oh *OgenHandler) TeamAddPost(ctx context.Context, req *api.Team) (api.TeamAddPostRes, error) {
	domainTeam := TeamFromAPI(req)

	// Вызов сервиса
	createdTeam, err := oh.teamService.CreateTeam(ctx, domainTeam)
	if err != nil {
		oh.logger.ErrorContext(ctx, "failed to create team",
			slog.String("team_name", req.TeamName),
			slog.String("error", err.Error()),
		)

		return ErrorToAPI(err), nil
	}

	oh.logger.InfoContext(ctx, "create team",
		slog.String("team_name", req.TeamName),
	)

	apiTeam := TeamToAPI(createdTeam)

	return &api.TeamAddPostCreated{
		Team: api.NewOptTeam(apiTeam),
	}, nil
}
