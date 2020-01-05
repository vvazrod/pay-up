package expense

import (
	"time"

	"github.com/google/uuid"
)

// Expense paid by a person on behalf of others.
type Expense struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primary_key"`
	GroupID     uuid.UUID `json:"group_id" gorm:"type:uuid"`
	Date        time.Time `json:"date"`
	Amount      float32   `json:"amount"`
	Description string    `json:"description"`
	Payer       uuid.UUID `json:"payer" gorm:"type:uuid"`
	Recipients  string    `json:"recipients"`
}
