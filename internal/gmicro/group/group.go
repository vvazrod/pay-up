package group

import (
	"errors"

	"github.com/google/uuid"
	"github.com/varrrro/pay-up/internal/member"
)

var (
	errExistingMember = errors.New("member already present in the group")
	errMemberNotFound = errors.New("couldn't find the member")
	errDeleteBalance  = errors.New("tried to delete a member with non-zero balance")
)

// Group of people, each of which has a balance in the group.
type Group struct {
	ID      uuid.UUID       `json:"id"`
	Name    string          `json:"name"`
	Members []member.Member `json:"members"`
}

// New Group instance.
func New(name string) *Group {
	return &Group{
		ID:      uuid.New(),
		Name:    name,
		Members: []member.Member{},
	}
}

// GetMember with the given name.
func (g *Group) GetMember(name string) (member.Member, error) {
	for _, m := range g.Members {
		if m.Name == name {
			return m, nil
		}
	}

	return member.Member{}, errMemberNotFound
}

// AddMember to a group.
func (g *Group) AddMember(name string) error {
	for _, m := range g.Members {
		if m.Name == name {
			return errExistingMember
		}
	}

	g.Members = append(g.Members, *member.New(name))

	return nil
}

// DeleteMember from a group.
func (g *Group) DeleteMember(name string) error {
	for i, m := range g.Members {
		if m.Name == name {
			if m.Balance != 0.0 {
				return errDeleteBalance
			}

			g.Members[i] = g.Members[len(g.Members)-1]
			g.Members = g.Members[:len(g.Members)-1]

			return nil
		}
	}

	return errMemberNotFound
}
