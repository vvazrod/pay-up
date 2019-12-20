package group_test

import (
	"os"
	"testing"

	"github.com/varrrro/pay-up/internal/gmicro/group"
	"github.com/varrrro/pay-up/internal/gmicro/member"
)

var testGroup *group.Group

func TestMain(m *testing.M) {
	testGroup = &group.Group{Name: "Test Group"}

	os.Exit(m.Run())
}

func teardown(t *testing.T) {
	testGroup.Members = nil
}

func TestNew(t *testing.T) {
	name := "Test"

	g := group.New(name)

	if g.Name != name {
		t.Error("Group created without the given name.")
	} else if g.Members == nil {
		t.Error("Group created with nil map.")
	}
}

func TestAddMember(t *testing.T) {
	defer teardown(t)

	m := member.Member{Name: "Test"}

	err := testGroup.AddMember(&m)

	if err != nil {
		t.Errorf("Couldn't add new member. Error: %s", err.Error())
	}

	if len(testGroup.Members) != 1 {
		t.Error("The member wasn't added correctly.")
	}
}

func TestAddExistingMember(t *testing.T) {
	defer teardown(t)

	m := member.Member{Name: "Test"}

	testGroup.AddMember(&m)
	err := testGroup.AddMember(&m)

	if err == nil {
		t.Error("Adding duplicate member didn't return an error.")
	}

	if len(testGroup.Members) != 1 {
		t.Error("Existing member was added to the group.")
	}
}
