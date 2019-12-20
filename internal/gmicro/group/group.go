package group

import (
	"errors"

	"github.com/google/uuid"
	"github.com/varrrro/pay-up/internal/gmicro/member"
)

var (
	errExistingMember = errors.New("member already present in the group")
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

// AddMember to a group.
func (g *Group) AddMember(member *member.Member) error {
	for _, m := range g.Members {
		if m.Name == member.Name {
			return errExistingMember
		}
	}

	g.Members = append(g.Members, *member)

	return nil
}
