package group

import "github.com/google/uuid"

import "errors"

var (
	errExistingMember = errors.New("member already present in the group")
	errMemberNotFound = errors.New("couldn't find the member")
	errDeleteBalance  = errors.New("tried to delete a member with non-zero balance")
)

// Group of people, each of which has a balance in the group.
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
		return errExistingMember
	}

	g.Members[member] = 0.0

	return nil
}

// DeleteMember from a group.
func (g *Group) DeleteMember(member string) error {
	if balance, prs := g.Members[member]; !prs {
		return errMemberNotFound
	} else if balance != 0.0 {
		return errDeleteBalance
	}

	delete(g.Members, member)

	return nil
}
