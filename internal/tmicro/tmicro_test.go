package tmicro_test

import (
	"os"
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/varrrro/pay-up/internal/tmicro"
	"github.com/varrrro/pay-up/internal/tmicro/expense"
	"github.com/varrrro/pay-up/internal/tmicro/payment"
)

var db *gorm.DB
var tm *tmicro.TransactionsManager

func TestMain(m *testing.M) {
	// Open connection to test DB
	db, _ = gorm.Open("sqlite3", ":memory:")
	defer db.Close()

	// Create tables
	db.CreateTable(&expense.Expense{}, &payment.Payment{})

	// Create manager with test DB connection
	tm = tmicro.NewManager(db)

	// Run tests
	os.Exit(m.Run())
}

func clearDB() {
	db.Delete(&expense.Expense{})
	db.Delete(&payment.Payment{})
}
