package handler

import (
	"context"
	"log/slog"

	api "github.com/teryble09/avito_backend/generated"
	"github.com/teryble09/avito_backend/internal/domain"
)

type TeamService interface {
	CreateTeam(ctx context.Context, team *domain.Team) (*domain.Team, error)
	GetTeam(ctx context.Context, teamName string) (*domain.Team, error)
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

func (oh *OgenHandler) TeamGetGet(ctx context.Context, params api.TeamGetGetParams) (api.TeamGetGetRes, error) {
	team, err := oh.teamService.GetTeam(ctx, params.TeamName)
	if err != nil {
		oh.logger.ErrorContext(ctx, "failed to get team",
			slog.String("team_name", params.TeamName),
			slog.String("error", err.Error()),
		)

		return ErrorToAPI(err), nil
	}

	apiTeam := TeamToAPI(team)

	return &apiTeam, nil
}
