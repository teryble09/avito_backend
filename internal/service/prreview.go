package service

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/teryble09/avito_backend/internal/domain"
)

type PRReviewerRepo interface {
	ReplaceReviewer(ctx context.Context, tx pgx.Tx, prID string, oldReviewerID string, newReviewerID string) error
	GetReviewers(ctx context.Context, prID string) ([]string, error)
}

type UserRepoForReviewer interface {
	GetUserByID(ctx context.Context, userID string) (*domain.User, error)
}

type TeamRepoForReviewer interface {
	GetTeamByName(ctx context.Context, teamName string) (*domain.Team, error)
}

type PullRequestRepoForReviewer interface {
	GetPRByID(ctx context.Context, prID string) (*domain.PullRequest, error)
}

type ReviewerService struct {
	db           *pgxpool.Pool
	reviewerRepo PRReviewerRepo
	userRepo     UserRepoForReviewer
	teamRepo     TeamRepoForReviewer
	prRepo       PullRequestRepoForReviewer
}

func NewReviewerService(
	db *pgxpool.Pool,
	reviewerRepo PRReviewerRepo,
	userRepo UserRepoForReviewer,
	teamRepo TeamRepoForReviewer,
	prRepo PullRequestRepoForReviewer,
) *ReviewerService {
	return &ReviewerService{
		db:           db,
		reviewerRepo: reviewerRepo,
		userRepo:     userRepo,
		teamRepo:     teamRepo,
		prRepo:       prRepo,
	}
}

func (s *ReviewerService) ReplaceReviewer(
	ctx context.Context,
	prID string,
	oldReviewerID string,
) (*domain.PullRequest, string, error) {
	user, err := s.userRepo.GetUserByID(ctx, oldReviewerID)
	if err != nil {
		return nil, "", fmt.Errorf("get author: %w", err)
	}

	team, err := s.teamRepo.GetTeamByName(ctx, user.TeamName)
	if err != nil {
		return nil, "", fmt.Errorf("get team: %w", err)
	}

	pr, err := s.prRepo.GetPRByID(ctx, prID)
	if err != nil {
		return nil, "", fmt.Errorf("get pull request: %w", err)
	}

	if pr.Status == domain.StatusMerged {
		return nil, "", domain.ErrPrMerged
	}

	newReviewerID, err := domain.SelectReplacementReviewer(team.Members, pr.AuthorID, oldReviewerID)
	if err != nil {
		return nil, "", fmt.Errorf("select replacement: %w", err)
	}

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, "", fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback(ctx) //nolint:errcheck //не страшно

	err = s.reviewerRepo.ReplaceReviewer(ctx, tx, prID, oldReviewerID, newReviewerID)
	if err != nil {
		return nil, "", fmt.Errorf("replace reviewer: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, "", fmt.Errorf("commit: %w", err)
	}

	reviewers, err := s.reviewerRepo.GetReviewers(ctx, prID)
	if err != nil {
		return nil, "", fmt.Errorf("get reviewers: %w", err)
	}

	pr.AssignedReviewers = reviewers

	return pr, newReviewerID, nil
}
