package member

import "github.com/google/uuid"

// Member of a group.
type Member struct {
	ID      uuid.UUID `json:"id"`
	Name    string    `json:"name"`
	Balance float32   `json:"balance"`
	GroupID uuid.UUID `json:"group_id"`
}

// New Member instance.
func New(name string) *Member {
	return &Member{
		ID:      uuid.New(),
		Name:    name,
		Balance: 0.0,
	}
}
