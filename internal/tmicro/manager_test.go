package tmicro_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/varrrro/pay-up/internal/tmicro/expense"
	"github.com/varrrro/pay-up/internal/tmicro/payment"
)

func TestCreateExpense(t *testing.T) {
	e := expense.Expense{
		ID:          uuid.New(),
		GroupID:     uuid.New(),
		Date:        time.Now(),
		Amount:      25.3,
		Description: "test",
		Payer:       uuid.New(),
		Recipients:  uuid.New().String() + ";" + uuid.New().String(),
	}

	if err := tm.CreateExpense(&e); err != nil {
		t.Errorf("Couldn't create expense. Error: %s", err.Error())
	}

	clearDB()
}

func TestCreateExpenseParseError(t *testing.T) {
	e := expense.Expense{
		ID:          uuid.New(),
		GroupID:     uuid.New(),
		Date:        time.Now(),
		Amount:      25.3,
		Description: "test",
		Payer:       uuid.New(),
		Recipients:  "fail",
	}

	if err := tm.CreateExpense(&e); err == nil {
		t.Error("Creating expense with wrong recipient ID didn't return an error.")
	}

	clearDB()
}

func TestRemoveLastExpense(t *testing.T) {
	e := expense.Expense{
		ID:          uuid.New(),
		GroupID:     uuid.New(),
		Date:        time.Now(),
		Amount:      25.3,
		Description: "test",
		Payer:       uuid.New(),
		Recipients:  uuid.New().String() + ";" + uuid.New().String(),
	}

	if err := tm.CreateExpense(&e); err != nil {
		t.Errorf("Couldn't create expense. Error: %s", err.Error())
	}

	if e2, err := tm.RemoveLastExpense(e.GroupID); err != nil {
		t.Errorf("Couldn't remove last expense. Error: %s", err.Error())
	} else if e2.ID != e.ID {
		t.Error("Returned expense doesn't match original.")
	}

	clearDB()
}

func TestRemoveLastExpenseNotFound(t *testing.T) {
	if _, err := tm.RemoveLastExpense(uuid.New()); err == nil {
		t.Error("Removing expense from non-existant group didn't return an error.")
	}

	clearDB()
}

func TestCreatePayment(t *testing.T) {
	p := payment.Payment{
		ID:        uuid.New(),
		GroupID:   uuid.New(),
		Date:      time.Now(),
		Amount:    27.3,
		Payer:     uuid.New(),
		Recipient: uuid.New(),
	}

	if err := tm.CreatePayment(&p); err != nil {
		t.Errorf("Creating a payment returned an error. Error: %s", err.Error())
	}

	clearDB()
}

func TestRemoveLastPayment(t *testing.T) {
	p := payment.Payment{
		ID:        uuid.New(),
		GroupID:   uuid.New(),
		Date:      time.Now(),
		Amount:    27.3,
		Payer:     uuid.New(),
		Recipient: uuid.New(),
	}

	if err := tm.CreatePayment(&p); err != nil {
		t.Errorf("Couldn't create payment. Error: %s", err.Error())
	}

	if p2, err := tm.RemoveLastPayment(p.GroupID); err != nil {
		t.Errorf("Couldn't remove last payment. Error: %s", err.Error())
	} else if p2.ID != p.ID {
		t.Error("Returned payment doesn't match original.")
	}

	clearDB()
}

func TestRemoveLastPaymentNotFound(t *testing.T) {
	if _, err := tm.RemoveLastPayment(uuid.New()); err == nil {
		t.Error("Removing payment from non-existant group didn't return an error.")
	}

	clearDB()
}
