package service

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/teryble09/avito_backend/internal/domain"
)

type TeamRepo interface {
	SaveNewTeam(ctx context.Context, tx pgx.Tx, team *domain.Team) error
}

type UserRepo interface {
	UpsertUsersBatch(ctx context.Context, tx pgx.Tx, team *domain.Team) error
}

type TeamService struct {
	db       *pgxpool.Pool
	teamRepo TeamRepo
	userRepo UserRepo
}

func NewTeamService(db *pgxpool.Pool, teamRepo TeamRepo, userRepo UserRepo) *TeamService {
	return &TeamService{
		db:       db,
		teamRepo: teamRepo,
		userRepo: userRepo,
	}
}

func (s *TeamService) CreateTeam(ctx context.Context, team *domain.Team) (*domain.Team, error) {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback(ctx) //nolint:errcheck //не страшно если roolback отдаст ошибку

	// сначала сохраняем команду
	err = s.teamRepo.SaveNewTeam(ctx, tx, team)
	if err != nil {
		return nil, fmt.Errorf("save team: %w", err)
	}

	// потом батч-вставка пользователей
	if err := s.userRepo.UpsertUsersBatch(ctx, tx, team); err != nil {
		return nil, fmt.Errorf("upsert users batch: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit transaction: %w", err)
	}

	return team, nil
}
