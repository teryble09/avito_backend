package handler

import (
	"errors"

	api "github.com/teryble09/avito_backend/generated"
	"github.com/teryble09/avito_backend/internal/domain"
)

// для маппинга ошибок.
var ErrInternal = errors.New("internal error")

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

// особенность ogen, если 1 одна ошибка в спеке, то
// ErrorResponse будет имплементировать ответ эндпоинта, а если 2,
// то уже будут обертки для ErrorReponse, и приходится писать другой маппер
// что есть внизу, этого можно избежать используя convenient error
// то есть чуть переделать спеку, но решил ее не трогать

func ErrorToAPI(err error) (*api.ErrorResponse, error) {
	switch {
	case errors.Is(err, domain.ErrTeamAlreadyExist):
		return &api.ErrorResponse{
			Error: api.ErrorResponseError{
				Code:    api.ErrorResponseErrorCodeTEAMEXISTS,
				Message: "team already exists",
			},
		}, nil

	case errors.Is(err, domain.ErrTeamNotFound):
		return &api.ErrorResponse{
			Error: api.ErrorResponseError{
				Code:    api.ErrorResponseErrorCodeNOTFOUND,
				Message: "team not found",
			},
		}, nil

	case errors.Is(err, domain.ErrUserNotFound):
		return &api.ErrorResponse{
			Error: api.ErrorResponseError{
				Code:    api.ErrorResponseErrorCodeNOTFOUND,
				Message: "user not found",
			},
		}, nil

	case errors.Is(err, domain.ErrPrAlreadyExists):
		return &api.ErrorResponse{
			Error: api.ErrorResponseError{
				Code:    api.ErrorResponseErrorCodePREXISTS,
				Message: "PR id already exists",
			},
		}, nil

	default:
		return nil, ErrInternal
	}
}

// маппер для эндпоинта с 2-умя ошибками

func PullRequestCreateErrorToAPI(err error) (api.PullRequestCreatePostRes, error) {
	switch {
	case errors.Is(err, domain.ErrPrAlreadyExists):
		return &api.PullRequestCreatePostConflict{
			Error: api.ErrorResponseError{
				Code:    api.ErrorResponseErrorCodePREXISTS,
				Message: "pr already exist",
			},
		}, nil

	case errors.Is(err, domain.ErrUserNotFound):
		return &api.PullRequestCreatePostNotFound{
			Error: api.ErrorResponseError{
				Code:    api.ErrorResponseErrorCodeNOTFOUND,
				Message: "user not found",
			},
		}, nil

	default:
		return nil, ErrInternal
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

func PullRequestToAPI(pr *domain.PullRequest) api.PullRequest {
	return api.PullRequest{
		PullRequestID:     pr.PullRequestID,
		PullRequestName:   pr.PullRequestName,
		AuthorID:          pr.AuthorID,
		Status:            api.PullRequestStatus(pr.Status),
		AssignedReviewers: pr.AssignedReviewers,
	}
}

func PullRequestShortToAPI(pr *domain.PullRequestShort) api.PullRequestShort {
	return api.PullRequestShort{
		PullRequestID:   pr.PullRequestID,
		PullRequestName: pr.PullRequestName,
		AuthorID:        pr.AuthorID,
		Status:          api.PullRequestShortStatus(pr.Status),
	}
}

func PullRequestCreateFromAPI(pr *api.PullRequestCreatePostReq) *domain.PullRequest {
	return &domain.PullRequest{
		PullRequestID:   pr.PullRequestID,
		PullRequestName: pr.PullRequestName,
		AuthorID:        pr.AuthorID,
	}
}
