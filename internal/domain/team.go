package domain

import "errors"

type Team struct {
	Name    string
	Members []*User
}

func NewTeam(name string, members []*User) *Team {
	return &Team{
		Name:    name,
		Members: members,
	}
}

var (
	ErrTeamAlreadyExist = errors.New("team already exists")
	ErrTeamNotFound     = errors.New("team not found")
)

// ActiveMembers возвращает активных участников команды.
func (t *Team) ActiveMembers() []*User {
	active := make([]*User, 0)

	for _, member := range t.Members {
		if member.IsActive {
			active = append(active, member)
		}
	}

	return active
}

// GetMember возвращает участника по ID.
func (t *Team) GetMember(userID string) (*User, bool) {
	for _, member := range t.Members {
		if member.ID == userID {
			return member, true
		}
	}

	return nil, false
}
