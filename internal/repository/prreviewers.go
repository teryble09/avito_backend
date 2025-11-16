package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/teryble09/avito_backend/internal/domain"
)

type PRReviewerRepo struct {
	db *pgxpool.Pool
}

func NewPRReviewerRepo(db *pgxpool.Pool) *PRReviewerRepo {
	return &PRReviewerRepo{db: db}
}

func (r *PRReviewerRepo) AssignReviewers(ctx context.Context, tx pgx.Tx, prID string, reviewerIDs []string) error {
	if len(reviewerIDs) == 0 {
		return nil
	}

	builder := squirrel.Insert("pr_reviewers").
		PlaceholderFormat(squirrel.Dollar).
		Columns("pull_request_id", "reviewer_id")

	for _, reviewerID := range reviewerIDs {
		builder = builder.Values(prID, reviewerID)
	}

	query, args, err := builder.ToSql()
	if err != nil {
		return fmt.Errorf("build query: %w", err)
	}

	_, err = tx.Exec(ctx, query, args...)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.ForeignKeyViolation {
			return domain.ErrUserNotFound // reviewer не существует
		}

		return fmt.Errorf("exec query: %w", err)
	}

	return nil
}
