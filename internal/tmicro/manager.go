package tmicro

import (
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/varrrro/pay-up/internal/tmicro/expense"
	"github.com/varrrro/pay-up/internal/tmicro/payment"
	"strings"
)

// Manager interface for the transactions microservice.
type Manager interface {
	CreateExpense(e *expense.Expense) error
	RemoveLastExpense(groupid uuid.UUID) (*expense.Expense, error)
	CreatePayment(p *payment.Payment) error
	RemoveLastPayment(groupid uuid.UUID) (*payment.Payment, error)
}

// TransactionsManager that works as single source of truth.
type TransactionsManager struct {
	DB *gorm.DB
}

// NewManager with the given database connection.
func NewManager(db *gorm.DB) *TransactionsManager {
	return &TransactionsManager{DB: db}
}

// CreateExpense in the given group.
func (tm *TransactionsManager) CreateExpense(e *expense.Expense) error {
	if _, err := uuid.Parse(e.ID); err != nil {
		return &UUIDParseError{"Couldn't parse expense ID", e.ID, err}
	} else if _, err := uuid.Parse(e.GroupID); err != nil {
		return &UUIDParseError{"Couldn't parse expense group ID", e.GroupID, err}
	} else if _, err := uuid.Parse(e.Payer); err != nil {
		return &UUIDParseError{"Couldn't parse expense payer ID", e.Payer, err}
	}

	rec := strings.Split(e.Recipients, ";")
	for _, r := range rec {
		if _, err := uuid.Parse(r); err != nil {
			return &UUIDParseError{"Couldn't parse expense recipient ID", r, err}
		}
	}

	tm.DB.Create(e)

	return nil
}

// RemoveLastExpense from the given group.
func (tm *TransactionsManager) RemoveLastExpense(groupid uuid.UUID) (*expense.Expense, error) {
	var e expense.Expense

	tm.DB.Where("group_id = ?", groupid).Order("date DESC").First(&e)

	if e.GroupID != groupid.String() {
		return nil, &NotFoundError{"No expense found", groupid.String()}
	}

	tm.DB.Delete(&e)

	return &e, nil
}

// CreatePayment in the given group.
func (tm *TransactionsManager) CreatePayment(p *payment.Payment) error {
	if _, err := uuid.Parse(p.ID); err != nil {
		return &UUIDParseError{"Couldn't parse payment ID", p.ID, err}
	} else if _, err := uuid.Parse(p.GroupID); err != nil {
		return &UUIDParseError{"Couldn't parse payment group ID", p.GroupID, err}
	} else if _, err := uuid.Parse(p.Payer); err != nil {
		return &UUIDParseError{"Couldn't parse payment payer ID", p.Payer, err}
	} else if _, err := uuid.Parse(p.Recipient); err != nil {
		return &UUIDParseError{"Couldn't parse payment recipient ID", p.Recipient, err}
	}

	tm.DB.Create(p)

	return nil
}

// RemoveLastPayment from the given group.
func (tm *TransactionsManager) RemoveLastPayment(groupid uuid.UUID) (*payment.Payment, error) {
	var p payment.Payment

	tm.DB.Where("group_id = ?", groupid).Order("date DESC").First(&p)

	if p.GroupID != groupid.String() {
		return nil, &NotFoundError{"No payment found", groupid.String()}
	}

	tm.DB.Delete(&p)

	return &p, nil
}
