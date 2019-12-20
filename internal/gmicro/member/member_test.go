package member_test

import (
	"testing"

	"github.com/varrrro/pay-up/internal/member"
)

func TestNew(t *testing.T) {
	m := member.New("Test")

	if m.Name != "Test" {
		t.Error("Member created without the given name.")
	} else if m.Balance != 0.0 {
		t.Error("Member created without zero balance.")
	}
}
