package group_test

import (
	"testing"

	"github.com/varrrro/pay-up/internal/gmicro/group"
)

func TestNew(t *testing.T) {
	name := "Test"

	g := group.New(name)

	if g.Name != name {
		t.Error("Group created without the given name.")
	} else if g.Members == nil {
		t.Error("Group created with nil map.")
	}
}
