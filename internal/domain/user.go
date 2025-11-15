package domain

import "errors"

type User struct {
	ID       string
	Username string
	TeamName string
	IsActive bool
}

func NewUser(id, username, teamName string, isActive bool) *User {
	return &User{
		ID:       id,
		Username: username,
		TeamName: teamName,
		IsActive: isActive,
	}
}

var ErrUserNotFound = errors.New("user not found")

// CanReview проверяет, может ли пользователь быть ревьювером.
func (u *User) CanReview() bool {
	return u.IsActive
}
