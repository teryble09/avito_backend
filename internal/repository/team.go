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

type TeamRepo struct {
	db *pgxpool.Pool
}

func NewTeamRepo(db *pgxpool.Pool) *TeamRepo {
	return &TeamRepo{
		db: db,
	}
}

func (r *TeamRepo) SaveNewTeam(ctx context.Context, tx pgx.Tx, team *domain.Team) error {
	teamEntity := entity.TeamFromDomain(team)

	query, args, err := squirrel.Insert("teams").
		PlaceholderFormat(squirrel.Dollar).
		Columns("team_name").
		Values(teamEntity.TeamName).
		ToSql()
	if err != nil {
		return fmt.Errorf("query build: %w", err)
	}

	_, err = tx.Exec(ctx, query, args...)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == pgerrcode.UniqueViolation {
				return domain.ErrTeamAlreadyExist
			}
		}

		return fmt.Errorf("query exec: %w", err)
	}

	return nil
}

func (r *TeamRepo) GetTeamByName(ctx context.Context, teamName string) (*domain.Team, error) {
	// Проверяем существование команды
	var exists bool

	err := r.db.QueryRow(ctx,
		`SELECT EXISTS(SELECT 1 FROM teams WHERE team_name = $1)`,
		teamName,
	).Scan(&exists)
	if err != nil {
		return nil, fmt.Errorf("check team exists: %w", err)
	}

	if !exists {
		return nil, domain.ErrTeamNotFound
	}

	query := `
		SELECT user_id, username, team_name, is_active 
		FROM users 
		WHERE team_name = $1
	`

	rows, err := r.db.Query(ctx, query, teamName)
	if err != nil {
		return nil, fmt.Errorf("select query: %w", err)
	}
	defer rows.Close()

	userEntities, err := pgx.CollectRows(rows, pgx.RowToStructByName[entity.User])
	if err != nil {
		return nil, fmt.Errorf("collect users: %w", err)
	}

	members := make([]*domain.User, 0, len(userEntities))
	for _, ue := range userEntities {
		members = append(members, ue.ToDomain())
	}

	return &domain.Team{
		Name:    teamName,
		Members: members,
	}, nil
}
