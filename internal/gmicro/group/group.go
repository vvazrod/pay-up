package group

import (
	"github.com/google/uuid"
	"github.com/varrrro/pay-up/internal/gmicro/member"
)

// Group of people, each of which has a balance in the group.
type Group struct {
	ID      string          `json:"id" gorm:"type:uuid;primary_key"`
	Name    string          `json:"name"`
	Members []member.Member `json:"members" gorm:"foreignkey:GroupID"`
}

// New Group instance.
func New(name string) *Group {
	return &Group{
		ID:      uuid.New().String(),
		Name:    name,
		Members: []member.Member{},
	}
}
