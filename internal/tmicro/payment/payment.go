package payment

import (
	"time"

	"github.com/google/uuid"
)

// Payment made by one person to another.
type Payment struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key"`
	GroupID   uuid.UUID `json:"group_id" gorm:"type:uuid"`
	Date      time.Time `json:"date"`
	Amount    float32   `json:"amount"`
	Payer     uuid.UUID `json:"payer" gorm:"type:uuid"`
	Recipient uuid.UUID `json:"recipient" gorm:"type:uuid"`
}
