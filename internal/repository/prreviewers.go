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

func (r *PRReviewerRepo) GetReviewers(ctx context.Context, prID string) ([]string, error) {
	query, args, err := squirrel.Select("reviewer_id").
		PlaceholderFormat(squirrel.Dollar).
		From("pr_reviewers").
		Where(squirrel.Eq{"pull_request_id": prID}).
		OrderBy("assigned_at").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build query: %w", err)
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}
	defer rows.Close()

	reviewers := make([]string, 0)

	for rows.Next() {
		var reviewerID string
		if err := rows.Scan(&reviewerID); err != nil {
			return nil, fmt.Errorf("scan: %w", err)
		}

		reviewers = append(reviewers, reviewerID)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return reviewers, nil
}

func (r *PRReviewerRepo) ReplaceReviewer(ctx context.Context, tx pgx.Tx, prID, oldReviewerID, newReviewerID string) error {
	deleteQuery, deleteArgs, err := squirrel.Delete("pr_reviewers").
		PlaceholderFormat(squirrel.Dollar).
		Where(squirrel.Eq{
			"pull_request_id": prID,
			"reviewer_id":     oldReviewerID,
		}).
		ToSql()
	if err != nil {
		return fmt.Errorf("build delete query: %w", err)
	}

	ct, err := tx.Exec(ctx, deleteQuery, deleteArgs...)
	if err != nil {
		return fmt.Errorf("delete old reviewer: %w", err)
	}

	if ct.RowsAffected() == 0 {
		return domain.ErrReviewerNotAssigned
	}

	insertQuery, insertArgs, err := squirrel.Insert("pr_reviewers").
		PlaceholderFormat(squirrel.Dollar).
		Columns("pull_request_id", "reviewer_id").
		Values(prID, newReviewerID).
		ToSql()
	if err != nil {
		return fmt.Errorf("build insert query: %w", err)
	}

	_, err = tx.Exec(ctx, insertQuery, insertArgs...)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return domain.ErrPrAlreadyExists
		}

		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.ForeignKeyViolation {
			return domain.ErrUserNotFound
		}

		return fmt.Errorf("insert new reviewer: %w", err)
	}

	return nil
}
