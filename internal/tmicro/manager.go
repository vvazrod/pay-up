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
	RemoveLastExpense(gid uuid.UUID) (*expense.Expense, error)
	CreatePayment(p *payment.Payment) error
	RemoveLastPayment(gid uuid.UUID) (*payment.Payment, error)
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
func (tm *TransactionsManager) RemoveLastExpense(gid uuid.UUID) (*expense.Expense, error) {
	var e expense.Expense

	tm.DB.Where("group_id = ?", gid).Order("date DESC").First(&e)

	if e.GroupID != gid {
		return nil, &NotFoundError{"No expense found", gid}
	}

	tm.DB.Delete(&e)

	return &e, nil
}

// CreatePayment in the given group.
func (tm *TransactionsManager) CreatePayment(p *payment.Payment) error {
	tm.DB.Create(p)

	return nil
}

// RemoveLastPayment from the given group.
func (tm *TransactionsManager) RemoveLastPayment(gid uuid.UUID) (*payment.Payment, error) {
	var p payment.Payment

	tm.DB.Where("group_id = ?", gid).Order("date DESC").First(&p)

	if p.GroupID != gid {
		return nil, &NotFoundError{"No payment found", gid}
	}

	tm.DB.Delete(&p)

	return &p, nil
}
