package handler

import (
	"context"
	"log/slog"

	api "github.com/teryble09/avito_backend/generated"
	"github.com/teryble09/avito_backend/internal/domain"
)

type UserService interface {
	SetUserIsActive(ctx context.Context, userID string, isActive bool) (*domain.User, error)
}

func (oh *OgenHandler) UsersSetIsActivePost(
	ctx context.Context,
	req *api.UsersSetIsActivePostReq,
) (api.UsersSetIsActivePostRes, error) {
	user, err := oh.userService.SetUserIsActive(ctx, req.UserID, req.IsActive)
	if err != nil {
		oh.logger.ErrorContext(ctx, "failed to set user is_active",
			slog.String("user_id", req.UserID),
			slog.Bool("is_active", req.IsActive),
			slog.String("error", err.Error()),
		)

		return ErrorToAPI(err)
	}

	oh.logger.InfoContext(ctx, "set is_active",
		slog.String("UserID", req.UserID),
	)

	userApi := UserToAPI(user)

	return &api.UsersSetIsActivePostOK{
		User: api.NewOptUser(userApi),
	}, nil
}

func (oh *OgenHandler) UsersGetReviewGet(ctx context.Context, params api.UsersGetReviewGetParams) (*api.UsersGetReviewGetOK, error) {
	prs, err := oh.prService.GetReviewerPRs(ctx, params.UserID)
	if err != nil {
		oh.logger.ErrorContext(ctx, "user get review",
			slog.String("user_id", params.UserID),
			slog.String("error", err.Error()),
		)

		return nil, ErrInternal
	}

	prsApi := make([]api.PullRequestShort, 0)

	for i := range len(prs) {
		prsApi = append(prsApi, PullRequestShortToAPI(prs[i]))
	}

	oh.logger.InfoContext(ctx, "got reviews",
		slog.String("user_id", params.UserID),
	)

	return &api.UsersGetReviewGetOK{
		UserID:       params.UserID,
		PullRequests: prsApi,
	}, nil
}
