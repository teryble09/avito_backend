package handler

import (
	"context"
	"log/slog"

	api "github.com/teryble09/avito_backend/generated"
	"github.com/teryble09/avito_backend/internal/domain"
)

type PullRequestService interface {
	CreatePullRequest(ctx context.Context, prCreateReq *domain.PullRequest) (*domain.PullRequest, error)
	GetReviewerPRs(ctx context.Context, userID string) ([]*domain.PullRequestShort, error)
	MergePullRequest(ctx context.Context, prID string) (*domain.PullRequest, error)
}

func (oh *OgenHandler) PullRequestCreatePost(
	ctx context.Context,
	req *api.PullRequestCreatePostReq,
) (api.PullRequestCreatePostRes, error) {
	prCreateReq := PullRequestCreateFromAPI(req)

	pr, err := oh.prService.CreatePullRequest(ctx, prCreateReq)
	if err != nil {
		oh.logger.ErrorContext(ctx, "failed to create pull request",
			slog.String("pr_id", req.PullRequestID),
			slog.String("pr_name", req.PullRequestName),
			slog.String("author_id", req.AuthorID),
			slog.String("error", err.Error()),
		)

		return PullRequestCreateErrorToAPI(err)
	}

	oh.logger.InfoContext(ctx, "create pr",
		slog.String("pr_id", req.PullRequestID),
		slog.String("pr_name", req.PullRequestName),
		slog.String("author_id", req.AuthorID),
	)

	prAPI := PullRequestToAPI(pr)

	return &api.PullRequestCreatePostCreated{
		Pr: api.NewOptPullRequest(prAPI),
	}, nil
}

func (oh *OgenHandler) PullRequestMergePost(
	ctx context.Context,
	req *api.PullRequestMergePostReq,
) (api.PullRequestMergePostRes, error) {
	pr, err := oh.prService.MergePullRequest(ctx, req.PullRequestID)
	if err != nil {
		oh.logger.ErrorContext(ctx, "failed to merge pull request",
			slog.String("pr_id", req.PullRequestID),
			slog.String("error", err.Error()),
		)

		return ErrorToAPI(err)
	}

	oh.logger.InfoContext(ctx, "create pr",
		slog.String("pr_id", req.PullRequestID),
	)

	return &api.PullRequestMergePostOK{
		Pr: api.NewOptPullRequest(
			PullRequestToAPI(pr),
		),
	}, nil
}
