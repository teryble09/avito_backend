package repository

import (
	"context"
	"errors"
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

func (r *UserRepo) SetIsActive(ctx context.Context, userID string, isActive bool) (*domain.User, error) {
	query, args, err := squirrel.Update("users").
		PlaceholderFormat(squirrel.Dollar).
		Set("is_active", isActive).
		Where(squirrel.Eq{"user_id": userID}).
		Suffix("RETURNING user_id, username, team_name, is_active").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build query: %w", err)
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("exec query: %w", err)
	}
	defer rows.Close()

	userEntity, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[entity.User])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}

		return nil, fmt.Errorf("scan user: %w", err)
	}

	userDomain := userEntity.ToDomain()

	return userDomain, nil
}

func (r *UserRepo) GetUserByID(ctx context.Context, userID string) (*domain.User, error) {
	query := `
		SELECT user_id, username, team_name, is_active
		FROM users
		WHERE user_id = $1
	`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}
	defer rows.Close()

	userEntity, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[entity.User])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}

		return nil, fmt.Errorf("scan user: %w", err)
	}

	return userEntity.ToDomain(), nil
}

func (r *UserRepo) GetActiveTeamMembers(ctx context.Context, teamName string) ([]*domain.User, error) {
	builder := squirrel.Select("user_id", "username", "team_name", "is_active").
		PlaceholderFormat(squirrel.Dollar).
		From("users").
		Where(squirrel.Eq{
			"team_name": teamName,
			"is_active": true,
		})

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("build query: %w", err)
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("run query: %w", err)
	}
	defer rows.Close()

	userEntities, err := pgx.CollectRows(rows, pgx.RowToStructByName[entity.User])
	if err != nil {
		return nil, fmt.Errorf("collect rows: %w", err)
	}

	users := make([]*domain.User, 0, len(userEntities))
	for _, ue := range userEntities {
		users = append(users, ue.ToDomain())
	}

	return users, nil
}
