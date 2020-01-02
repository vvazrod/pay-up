package expense_test

import (
	"github.com/google/uuid"
	"github.com/varrrro/pay-up/internal/tmicro/expense"
	"testing"
)

func TestNew(t *testing.T) {
	testid := uuid.New()
	groupid := uuid.New()
	payerid := uuid.New()
	recipientids := []uuid.UUID{testid, payerid}

	e := expense.New("Test", 25.9, groupid, payerid, &recipientids)

	if e.GroupID != groupid.String() {
		t.Error("Group ID wasn't parsed correctly.")
	} else if e.Payer != payerid.String() {
		t.Error("Payer ID wasn't parsed correctly.")
	} else if e.Recipients != testid.String()+";"+payerid.String() {
		t.Error("Recipients string wasn't created correctly.")
	}
}
