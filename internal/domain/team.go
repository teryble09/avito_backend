package domain

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

func (t *Team) ActiveMembers() []*User {
	active := make([]*User, 0)

	for _, member := range t.Members {
		if member.IsActive {
			active = append(active, member)
		}
	}

	return active
}

func (t *Team) GetMember(userID string) (*User, bool) {
	for _, member := range t.Members {
		if member.ID == userID {
			return member, true
		}
	}

	return nil, false
}
