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
