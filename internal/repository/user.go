package repository

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/teryble09/avito_backend/internal/domain"
	"github.com/teryble09/avito_backend/internal/entity"
)

type UserRepo struct {
	db *pgxpool.Pool
}

func NewUserRepo(db *pgxpool.Pool) *UserRepo {
	return &UserRepo{
		db: db,
	}
}

func (r *UserRepo) UpsertUsersBatch(ctx context.Context, tx pgx.Tx, team *domain.Team) error {
	if len(team.Members) == 0 {
		return nil
	}

	builder := squirrel.Insert("users").
		PlaceholderFormat(squirrel.Dollar).
		Columns("user_id", "username", "team_name", "is_active")

	teamEntity := entity.TeamFromDomain(team)

	// Добавляем все строки в один запрос
	for _, user := range team.Members {
		userEntity := entity.UserFromDomain(user)
		builder = builder.Values(userEntity.UserID, userEntity.Username, teamEntity.TeamName, userEntity.IsActive)
	}

	builder = builder.Suffix(`
		ON CONFLICT (user_id) DO UPDATE SET
			username = EXCLUDED.username,
			team_name = EXCLUDED.team_name,
			is_active = EXCLUDED.is_active
	`)

	query, args, err := builder.ToSql()
	if err != nil {
		return fmt.Errorf("query build: %w", err)
	}

	_, err = tx.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("query exec: %w", err)
	}

	return nil
}
