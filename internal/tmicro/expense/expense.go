package expense

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

// Expense paid by a person on behalf of others.
type Expense struct {
	ID          string    `json:"id" gorm:"type:uuid;primary_key"`
	GroupID     string    `json:"group_id" gorm:"type:uuid"`
	Description string    `json:"description"`
	Amount      float32   `json:"amount"`
	Date        time.Time `json:"date"`
	Payer       string    `json:"payer" gorm:"type:uuid"`
	Recipients  string    `json:"recipients"`
}

// New Expense instance.
func New(description string, amount float32, groupid, payer uuid.UUID, recipients *[]uuid.UUID) *Expense {
	var recStrings []string
	for _, r := range *recipients {
		recStrings = append(recStrings, r.String())
	}

	return &Expense{
		ID:          uuid.New().String(),
		GroupID:     groupid.String(),
		Description: description,
		Amount:      amount,
		Date:        time.Now(),
		Payer:       payer.String(),
		Recipients:  strings.Join(recStrings, ";"),
	}
}
