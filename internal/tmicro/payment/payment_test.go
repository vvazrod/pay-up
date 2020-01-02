package payment_test

import (
	"github.com/google/uuid"
	"github.com/varrrro/pay-up/internal/tmicro/payment"
	"testing"
)

func TestNew(t *testing.T) {
	groupid := uuid.New()
	payerid := uuid.New()
	recipientid := uuid.New()

	p := payment.New(23.4, groupid, payerid, recipientid)

	if p.GroupID != groupid.String() {
		t.Error("Group ID wasn't parsed correctly.")
	} else if p.Payer != payerid.String() {
		t.Error("Payer ID wasn't parsed correctly.")
	} else if p.Recipient != recipientid.String() {
		t.Error("Recipient ID wasn't parsed correctly.")
	}
}
