package entity

import (
	"github.com/teryble09/avito_backend/internal/domain"
)

type User struct {
	UserID   string `db:"user_id"`
	Username string `db:"username"`
	TeamName string `db:"team_name"`
	IsActive bool   `db:"is_active"`
}

// ToDomain конвертирует entity в domain model.
func (e *User) ToDomain() *domain.User {
	return domain.NewUser(e.UserID, e.Username, e.TeamName, e.IsActive)
}

// UserFromDomain создает entity из domain model.
func UserFromDomain(u *domain.User) *User {
	return &User{
		UserID:   u.ID,
		Username: u.Username,
		TeamName: u.TeamName,
		IsActive: u.IsActive,
	}
}
