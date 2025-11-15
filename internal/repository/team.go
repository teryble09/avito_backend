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
			switch pgErr.Code {
			case pgerrcode.UniqueViolation:
				return domain.ErrTeamAlreadyExist
			}
		}

		return fmt.Errorf("query exec: %w", err)
	}

	return nil
}
