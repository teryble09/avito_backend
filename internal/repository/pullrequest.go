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
	"github.com/teryble09/avito_backend/internal/entity"
)

type PullRequestRepo struct {
	db *pgxpool.Pool
}

func NewPullRequestRepo(db *pgxpool.Pool) *PullRequestRepo {
	return &PullRequestRepo{db: db}
}

func (r *PullRequestRepo) CreatePR(ctx context.Context, tx pgx.Tx, pr *domain.PullRequest) error {
	prEntity := entity.PullRequestFromDomain(pr)

	query, args, err := squirrel.Insert("pull_requests").
		PlaceholderFormat(squirrel.Dollar).
		Columns("pull_request_id", "pull_request_name", "author_id", "status").
		Values(prEntity.PullRequestID, prEntity.PullRequestName, prEntity.AuthorID, prEntity.Status).
		ToSql()
	if err != nil {
		return fmt.Errorf("build pr insert query: %w", err)
	}

	_, err = tx.Exec(ctx, query, args...)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return domain.ErrPrAlreadyExists
		}

		return fmt.Errorf("insert pr: %w", err)
	}

	return nil
}
