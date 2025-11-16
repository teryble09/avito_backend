package handler

import (
	"context"

	api "github.com/teryble09/avito_backend/generated"
	"github.com/teryble09/avito_backend/internal/domain"
)

type ReviewerService interface {
	ReplaceReviewer(ctx context.Context, prID string, oldReviewerID string) (pr *domain.PullRequest, newReviewerID string, err error)
}

func (oh *OgenHandler) PullRequestReassignPost(
	ctx context.Context,
	req *api.PullRequestReassignPostReq,
) (api.PullRequestReassignPostRes, error) {
	pr, newReviewer, err := oh.reviewerService.ReplaceReviewer(ctx, req.PullRequestID, req.OldUserID)
	if err != nil {
		return PullRequestReassignErrorToAPI(err)
	}

	prAPI := PullRequestToAPI(pr)

	return &api.PullRequestReassignPostOK{
		Pr:         prAPI,
		ReplacedBy: newReviewer,
	}, nil
}
