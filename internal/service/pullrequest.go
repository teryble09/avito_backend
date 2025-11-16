package service

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/teryble09/avito_backend/internal/domain"
)

type PullRequestRepo interface {
	CreatePR(ctx context.Context, tx pgx.Tx, pr *domain.PullRequest) error
	GetPRsByReviewer(ctx context.Context, reviewerID string) ([]*domain.PullRequestShort, error)
}

type PRReviewerRepo interface {
	AssignReviewers(ctx context.Context, tx pgx.Tx, prID string, reviewerIDs []string) error
}

type UserRepoForPR interface {
	GetUserByID(ctx context.Context, userID string) (*domain.User, error)
	GetActiveTeamMembers(ctx context.Context, teamName string) ([]*domain.User, error)
}

type PullRequestService struct {
	db           *pgxpool.Pool
	prRepo       PullRequestRepo
	reviewerRepo PRReviewerRepo
	userRepo     UserRepoForPR
}

func NewPullRequestService(
	db *pgxpool.Pool,
	prRepo PullRequestRepo,
	reviewerRepo PRReviewerRepo,
	userRepo UserRepoForPR,
) *PullRequestService {
	return &PullRequestService{
		db:           db,
		prRepo:       prRepo,
		reviewerRepo: reviewerRepo,
		userRepo:     userRepo,
	}
}

func (s *PullRequestService) CreatePullRequest(
	ctx context.Context,
	pr *domain.PullRequest,
) (*domain.PullRequest, error) {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback(ctx) //nolint:errcheck //не страшно если rollback не сработает

	author, err := s.userRepo.GetUserByID(ctx, pr.AuthorID)
	if err != nil {
		return nil, fmt.Errorf("get author: %w", err)
	}

	teamMembers, err := s.userRepo.GetActiveTeamMembers(ctx, author.TeamName)
	if err != nil {
		return nil, fmt.Errorf("get team members: %w", err)
	}

	selectedReviewers := domain.SelectRandomReviewersExcludingAuthor(teamMembers, 2, author)

	selectedReviewersIds := make([]string, len(selectedReviewers))
	for i := range len(selectedReviewers) {
		selectedReviewersIds[i] = selectedReviewers[i].Username
	}

	pr = &domain.PullRequest{
		PullRequestID:     pr.PullRequestID,
		PullRequestName:   pr.PullRequestName,
		AuthorID:          pr.AuthorID,
		Status:            domain.StatusOpen,
		AssignedReviewers: selectedReviewersIds,
	}

	err = s.prRepo.CreatePR(ctx, tx, pr)
	if err != nil {
		return nil, fmt.Errorf("create pr: %w", err)
	}

	if len(selectedReviewers) > 0 {
		err = s.reviewerRepo.AssignReviewers(ctx, tx, pr.PullRequestID, selectedReviewersIds)
		if err != nil {
			return nil, fmt.Errorf("assign reviewers: %w", err)
		}
	} else {
		// Тут получается, что человек работает solo,
		// если ему доверяют, то надо merdge
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}

	return pr, nil
}

func (s *PullRequestService) GetReviewerPRs(ctx context.Context, userID string) ([]*domain.PullRequestShort, error) {
	prs, err := s.prRepo.GetPRsByReviewer(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get reviewer PRs: %w", err)
	}

	return prs, nil
}
