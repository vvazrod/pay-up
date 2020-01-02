package payment

import (
	"time"

	"github.com/google/uuid"
)

// Payment made by one person to another.
type Payment struct {
	ID        string    `json:"id" gorm:"type:uuid;primary_key"`
	GroupID   string    `json:"group_id" gorm:"type:uuid"`
	Date      time.Time `json:"date"`
	Amount    float32   `json:"amount"`
	Payer     string    `json:"payer" gorm:"type:uuid"`
	Recipient string    `json:"recipient" gorm:"type:uuid"`
}

// New Payment instance.
func New(amount float32, groupid, payer, recipient uuid.UUID) *Payment {
	return &Payment{
		ID:        uuid.New().String(),
		GroupID:   groupid.String(),
		Amount:    amount,
		Date:      time.Now(),
		Payer:     payer.String(),
		Recipient: recipient.String(),
	}
}
