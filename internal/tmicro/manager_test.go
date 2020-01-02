package tmicro_test

import (
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/varrrro/pay-up/internal/tmicro"
	"github.com/varrrro/pay-up/internal/tmicro/expense"
	"github.com/varrrro/pay-up/internal/tmicro/payment"
)

var manager *tmicro.TransactionsManager

func TestMain(m *testing.M) {
	// Open connection to test DB
	db, _ := gorm.Open("sqlite3", ":memory:")
	defer db.Close()

	// Create tables
	db.CreateTable(&expense.Expense{}, &payment.Payment{})

	// Create manager with test DB connection
	manager = tmicro.NewManager(db)

	os.Exit(m.Run())
}

func TestCreateExpense(t *testing.T) {
	e := expense.New("test", 25.2, uuid.New(), uuid.New(), &[]uuid.UUID{uuid.New()})

	err := manager.CreateExpense(e)

	if err != nil {
		t.Errorf("Couldn't create expense. Error: %s", err.Error())
	}
}

func TestCreateExpenseParseErrorID(t *testing.T) {
	e := &expense.Expense{
		ID:          "fail",
		GroupID:     uuid.New().String(),
		Description: "test",
		Amount:      27.3,
		Date:        time.Now(),
		Payer:       uuid.New().String(),
		Recipients:  uuid.New().String(),
	}

	err := manager.CreateExpense(e)

	if err == nil {
		t.Error("Creating expense with wrong ID didn't return an error.")
	}
}

func TestCreateExpenseParseErrorGroup(t *testing.T) {
	e := &expense.Expense{
		ID:          uuid.New().String(),
		GroupID:     "fail",
		Description: "test",
		Amount:      27.3,
		Date:        time.Now(),
		Payer:       uuid.New().String(),
		Recipients:  uuid.New().String(),
	}

	err := manager.CreateExpense(e)

	if err == nil {
		t.Error("Creating expense with wrong group ID didn't return an error.")
	}
}

func TestCreateExpenseParseErrorPayer(t *testing.T) {
	e := &expense.Expense{
		ID:          uuid.New().String(),
		GroupID:     uuid.New().String(),
		Description: "test",
		Amount:      27.3,
		Date:        time.Now(),
		Payer:       "fail",
		Recipients:  uuid.New().String(),
	}

	err := manager.CreateExpense(e)

	if err == nil {
		t.Error("Creating expense with wrong payer ID didn't return an error.")
	}
}

func TestCreateExpenseParseErrorRecipient(t *testing.T) {
	e := &expense.Expense{
		ID:          uuid.New().String(),
		GroupID:     uuid.New().String(),
		Description: "test",
		Amount:      27.3,
		Date:        time.Now(),
		Payer:       uuid.New().String(),
		Recipients:  "fail",
	}

	err := manager.CreateExpense(e)

	if err == nil {
		t.Error("Creating expense with wrong recipient ID didn't return an error.")
	}
}

func TestRemoveLastExpense(t *testing.T) {
	groupid := uuid.New()
	e := expense.New("test", 25.2, groupid, uuid.New(), &[]uuid.UUID{uuid.New()})
	err := manager.CreateExpense(e)
	if err != nil {
		t.Errorf("Couldn't create expense. Error: %s", err.Error())
	}

	e2, err := manager.RemoveLastExpense(groupid)
	if err != nil {
		t.Errorf("Couldn't remove last expense. Error: %s", err.Error())
	} else if e2.ID != e.ID {
		t.Error("Returned expense doesn't match original.")
	}
}

func TestRemoveLastExpenseNotFound(t *testing.T) {
	_, err := manager.RemoveLastExpense(uuid.New())
	if err == nil {
		t.Error("Removing expense from non-existant group didn't return an error.")
	}
}

func TestCreatePayment(t *testing.T) {
	p := payment.New(27.3, uuid.New(), uuid.New(), uuid.New())

	err := manager.CreatePayment(p)

	if err != nil {
		t.Errorf("Creating a payment returned an error. Error: %s", err.Error())
	}
}

func TestCreatePaymentParseErrorID(t *testing.T) {
	p := &payment.Payment{
		ID:        "fail",
		GroupID:   uuid.New().String(),
		Amount:    27.3,
		Date:      time.Now(),
		Payer:     uuid.New().String(),
		Recipient: uuid.New().String(),
	}

	err := manager.CreatePayment(p)

	if err == nil {
		t.Error("Creating payment with wrong ID didn't return an error.")
	}
}

func TestCreatePaymentParseErrorGroup(t *testing.T) {
	p := &payment.Payment{
		ID:        uuid.New().String(),
		GroupID:   "fail",
		Amount:    27.3,
		Date:      time.Now(),
		Payer:     uuid.New().String(),
		Recipient: uuid.New().String(),
	}

	err := manager.CreatePayment(p)

	if err == nil {
		t.Error("Creating payment with wrong group ID didn't return an error.")
	}
}

func TestCreatePaymentParseErrorPayer(t *testing.T) {
	p := &payment.Payment{
		ID:        uuid.New().String(),
		GroupID:   uuid.New().String(),
		Amount:    27.3,
		Date:      time.Now(),
		Payer:     "fail",
		Recipient: uuid.New().String(),
	}

	err := manager.CreatePayment(p)

	if err == nil {
		t.Error("Creating payment with wrong payer ID didn't return an error.")
	}
}

func TestCreatePaymentParseErrorRecipient(t *testing.T) {
	p := &payment.Payment{
		ID:        uuid.New().String(),
		GroupID:   uuid.New().String(),
		Amount:    27.3,
		Date:      time.Now(),
		Payer:     uuid.New().String(),
		Recipient: "fail",
	}

	err := manager.CreatePayment(p)

	if err == nil {
		t.Error("Creating payment with wrong recipient ID didn't return an error.")
	}
}

func TestRemoveLastPayment(t *testing.T) {
	groupid := uuid.New()
	p := payment.New(25.2, groupid, uuid.New(), uuid.New())
	err := manager.CreatePayment(p)
	if err != nil {
		t.Errorf("Couldn't create payment. Error: %s", err.Error())
	}

	p2, err := manager.RemoveLastPayment(groupid)
	if err != nil {
		t.Errorf("Couldn't remove last payment. Error: %s", err.Error())
	} else if p2.ID != p.ID {
		t.Error("Returned payment doesn't match original.")
	}
}

func TestRemoveLastPaymentNotFound(t *testing.T) {
	_, err := manager.RemoveLastPayment(uuid.New())
	if err == nil {
		t.Error("Removing payment from non-existant group didn't return an error.")
	}
}
