package member

import "github.com/google/uuid"

// Member of a group.
type Member struct {
	ID      uuid.UUID `json:"id" gorm:"type:uuid;primary_key"`
	Name    string    `json:"name"`
	Balance float32   `json:"balance"`
	GroupID uuid.UUID `json:"group_id" gorm:"type:uuid"`
}
