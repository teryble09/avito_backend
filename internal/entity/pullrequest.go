package entity

import "github.com/teryble09/avito_backend/internal/domain"

type PullRequest struct {
	PullRequestID   string `db:"pull_request_id"`
	PullRequestName string `db:"pull_request_name"`
	AuthorID        string `db:"author_id"`
	Status          string `db:"status"`
}

func (e *PullRequest) ToDomain() *domain.PullRequest {
	return &domain.PullRequest{
		PullRequestID:   e.PullRequestID,
		PullRequestName: e.PullRequestName,
		AuthorID:        e.AuthorID,
		Status:          domain.PullRequestStatus(e.Status),
	}
}

func PullRequestFromDomain(pr *domain.PullRequest) *PullRequest {
	return &PullRequest{
		PullRequestID:   pr.PullRequestID,
		PullRequestName: pr.PullRequestName,
		AuthorID:        pr.AuthorID,
		Status:          string(pr.Status),
	}
}
