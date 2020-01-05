package group

import (
	"github.com/google/uuid"
	"github.com/varrrro/pay-up/internal/gmicro/member"
)

// Group of people, each of which has a balance in the group.
type Group struct {
	ID      uuid.UUID       `json:"id" gorm:"type:uuid;primary_key"`
	Name    string          `json:"name"`
	Members []member.Member `json:"members" gorm:"foreignkey:GroupID"`
}
