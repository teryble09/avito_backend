package handler

import (
	"errors"

	api "github.com/teryble09/avito_backend/generated"
	"github.com/teryble09/avito_backend/internal/domain"
)

func TeamFromAPI(req *api.Team) *domain.Team {
	members := make([]*domain.User, 0, len(req.Members))
	for _, m := range req.Members {
		members = append(members, &domain.User{
			ID:       m.UserID,
			Username: m.Username,
			TeamName: req.TeamName,
			IsActive: m.IsActive,
		})
	}

	return &domain.Team{
		Name:    req.TeamName,
		Members: members,
	}
}

func TeamToAPI(team *domain.Team) api.Team {
	members := make([]api.TeamMember, 0, len(team.Members))
	for _, m := range team.Members {
		members = append(members, api.TeamMember{
			UserID:   m.ID,
			Username: m.Username,
			IsActive: m.IsActive,
		})
	}

	return api.Team{
		TeamName: team.Name,
		Members:  members,
	}
}

func ErrorToAPI(err error) *api.ErrorResponse {
	switch {
	case errors.Is(err, domain.ErrTeamAlreadyExist):
		return &api.ErrorResponse{
			Error: api.ErrorResponseError{
				Code:    api.ErrorResponseErrorCodeTEAMEXISTS,
				Message: "team already exists",
			},
		}

	case errors.Is(err, domain.ErrTeamNotFound):
		return &api.ErrorResponse{
			Error: api.ErrorResponseError{
				Code:    api.ErrorResponseErrorCodeNOTFOUND,
				Message: "team not found",
			},
		}

	case errors.Is(err, domain.ErrUserNotFound):
		return &api.ErrorResponse{
			Error: api.ErrorResponseError{
				Code:    api.ErrorResponseErrorCodeNOTFOUND,
				Message: "user not found",
			},
		}

	default:
		return &api.ErrorResponse{
			Error: api.ErrorResponseError{
				Code:    api.ErrorResponseErrorCodeNOTFOUND,
				Message: "unknown error",
			},
		}
	}
}

func UserToAPI(u *domain.User) api.User {
	return api.User{
		UserID:   u.ID,
		Username: u.Username,
		TeamName: u.TeamName,
		IsActive: u.IsActive,
	}
}

func UserFromAPI(u *api.User) *domain.User {
	return &domain.User{
		ID:       u.UserID,
		Username: u.Username,
		TeamName: u.TeamName,
		IsActive: u.IsActive,
	}
}
