package member

import "github.com/google/uuid"

// Member of a group.
type Member struct {
	ID      string  `json:"id" gorm:"type:uuid;primary_key"`
	Name    string  `json:"name"`
	Balance float32 `json:"balance"`
	GroupID string  `json:"group_id" gorm:"type:uuid"`
}

// New Member instance.
func New(name string) *Member {
	return &Member{
		ID:      uuid.New().String(),
		Name:    name,
		Balance: 0.0,
	}
}
