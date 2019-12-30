package tmicro

import (
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/varrrro/pay-up/internal/tmicro/expense"
	"github.com/varrrro/pay-up/internal/tmicro/payment"
)

// TransactionsManager that works as single source of truth.
type TransactionsManager struct {
	DB *gorm.DB
}

// NewManager with the given database connection.
func NewManager(db *gorm.DB) *TransactionsManager {
	return &TransactionsManager{DB: db}
}

// AddExpense to the given group.
func (tm *TransactionsManager) AddExpense(description string, amount float32, groupid, payer uuid.UUID, recipients *[]uuid.UUID) error {
	e := expense.New(description, amount, groupid, payer, recipients)

	tm.DB.Create(e)

	return nil
}

// DeleteLastExpense in the given group.
func (tm *TransactionsManager) DeleteLastExpense(groupid uuid.UUID) error {
	var e expense.Expense

	tm.DB.Where("group_id = ?", groupid).Order("date DESC").First(&e)

	if e.GroupID != groupid.String() {
		return &NotFoundError{"No expense found", groupid.String()}
	}

	tm.DB.Delete(&e)

	return nil
}

// AddPayment to the given group.
func (tm *TransactionsManager) AddPayment(amount float32, groupid, payer, recipient uuid.UUID) error {
	p := payment.New(amount, groupid, payer, recipient)

	tm.DB.Create(p)

	return nil
}

// DeleteLastPayment in the given group.
func (tm *TransactionsManager) DeleteLastPayment(groupid uuid.UUID) error {
	var p payment.Payment

	tm.DB.Where("group_id = ?", groupid).Order("date DESC").First(&p)

	if p.GroupID != groupid.String() {
		return &NotFoundError{"No payment found", groupid.String()}
	}

	tm.DB.Delete(&p)

	return nil
}
