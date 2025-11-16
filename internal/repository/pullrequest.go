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

func (r *PullRequestRepo) GetPRsByReviewer(ctx context.Context, reviewerID string) ([]*domain.PullRequestShort, error) {
	query, args, err := squirrel.Select(
		"pr.pull_request_id",
		"pr.pull_request_name",
		"pr.author_id",
		"pr.status",
	).
		PlaceholderFormat(squirrel.Dollar).
		From("pull_requests pr").
		InnerJoin("pr_reviewers rev ON rev.pull_request_id = pr.pull_request_id").
		Where(squirrel.Eq{
			"rev.reviewer_id": reviewerID,
			"pr.status":       "OPEN",
		}).
		OrderBy("pr.created_at DESC").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build query: %w", err)
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}
	defer rows.Close()

	prEntities, err := pgx.CollectRows(rows, pgx.RowToStructByName[entity.PullRequest])
	if err != nil {
		return nil, fmt.Errorf("collect rows: %w", err)
	}

	prs := make([]*domain.PullRequestShort, 0, len(prEntities))
	for _, prEntity := range prEntities {
		prs = append(prs, prEntity.ToDomainShort())
	}

	return prs, nil
}

func (r *PullRequestRepo) MergePR(ctx context.Context, prID string) (*domain.PullRequest, error) {
	updateQuery, updateArgs, err := squirrel.Update("pull_requests").
		PlaceholderFormat(squirrel.Dollar).
		Set("status", "MERGED").
		Set("merged_at", squirrel.Expr("COALESCE(merged_at, NOW())")).
		Where(squirrel.Eq{"pull_request_id": prID}).
		Suffix("RETURNING pull_request_id, pull_request_name, author_id, status, merged_at").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build update query: %w", err)
	}

	rows, err := r.db.Query(ctx, updateQuery, updateArgs...)
	if err != nil {
		return nil, fmt.Errorf("exec update: %w", err)
	}
	defer rows.Close()

	prEntity, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[entity.PullRequest])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrPrNotFound
		}

		return nil, fmt.Errorf("scan pr: %w", err)
	}

	return prEntity.ToDomain(), nil
}
