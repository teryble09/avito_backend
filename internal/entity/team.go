package entity

import (
	"github.com/teryble09/avito_backend/internal/domain"
)

type Team struct {
	TeamName string `db:"team_name"`
}

// FromDomain создает entity из domain model.
func TeamFromDomain(t *domain.Team) *Team {
	return &Team{
		TeamName: t.Name,
	}
}
