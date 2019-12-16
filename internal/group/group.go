package group

import "github.com/google/uuid"

// Group of people, each of which have a balance in the group.
type Group struct {
	ID      uuid.UUID
	Name    string
	Members map[string]float32
}

// New Group instance.
func New(name string) *Group {
	return &Group{
		ID:      uuid.New(),
		Name:    name,
		Members: make(map[string]float32),
	}
}

// AddMember to a group.
func (g *Group) AddMember(member string) error {
	if _, prs := g.Members[member]; prs {
		return &ExistingMemberError{g.ID.String(), member}
	}

	g.Members[member] = 0.0

	return nil
}

// DeleteMember from a group.
func (g *Group) DeleteMember(member string) error {
	if balance, prs := g.Members[member]; !prs {
		return &MemberNotFoundError{g.ID.String(), member}
	} else if balance != 0.0 {
		return &DeletingBalanceError{g.ID.String(), member, balance}
	}

	delete(g.Members, member)

	return nil
}
