package member

import "github.com/google/uuid"

// Member of a group.
type Member struct {
	ID      uuid.UUID
	Name    string
	Balance float32
}

// New Member instance.
func New(name string) *Member {
	return &Member{
		ID:      uuid.New(),
		Name:    name,
		Balance: 0.0,
	}
}
