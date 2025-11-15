package service

import (
	"context"
	"fmt"

	"github.com/teryble09/avito_backend/internal/domain"
)

type UserRepo interface {
	SetIsActive(ctx context.Context, userID string, isActive bool) (*domain.User, error)
}

type UserService struct {
	userRepo UserRepo
}

func NewUserService(userRepo UserRepo) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

func (s *UserService) SetUserIsActive(ctx context.Context, userID string, isActive bool) (*domain.User, error) {
	user, err := s.userRepo.SetIsActive(ctx, userID, isActive)
	if err != nil {
		return nil, fmt.Errorf("set is_active: %w", err)
	}

	return user, nil
}
