package group_test

import (
	"os"
	"testing"

	"github.com/varrrro/pay-up/internal/group"
	"github.com/varrrro/pay-up/internal/member"
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

func TestGetMember(t *testing.T) {
	defer teardown(t)

	m := "Test"

	testGroup.AddMember(m)
	_, err := testGroup.GetMember(m)

	if err != nil {
		t.Errorf("Couldn't retrieve member. Error: %s", err.Error())
	}
}

func TestGetNonExistantMember(t *testing.T) {
	defer teardown(t)

	m := "Test"

	_, err := testGroup.GetMember(m)

	if err == nil {
		t.Error("Retrieving a non-existant member didn't return an error.")
	}
}

func TestAddMember(t *testing.T) {
	defer teardown(t)

	m := "Test"

	err := testGroup.AddMember(m)

	if err != nil {
		t.Errorf("Couldn't add new member. Error: %s", err.Error())
	}

	if len(testGroup.Members) != 1 {
		t.Error("The member wasn't added correctly.")
	}
}

func TestAddExistingMember(t *testing.T) {
	defer teardown(t)

	m := "Test"

	testGroup.AddMember(m)
	err := testGroup.AddMember(m)

	if err == nil {
		t.Error("Adding duplicate member didn't return an error.")
	}

	if len(testGroup.Members) != 1 {
		t.Error("Existing member was added to the group.")
	}
}

func TestDeleteMember(t *testing.T) {
	defer teardown(t)

	m := "Test"

	testGroup.AddMember(m)
	err := testGroup.DeleteMember(m)

	if err != nil {
		t.Errorf("Deleting a member returned an error: %s", err.Error())
	}

	if len(testGroup.Members) != 0 {
		t.Errorf("Number of members isn't correct (%d).", len(testGroup.Members))
	}
}

func TestDeleteMemberNotFound(t *testing.T) {
	err := testGroup.DeleteMember("Test")

	if err == nil {
		t.Error("Deleting a non-existing member didn't return an error.")
	}
}

func TestDeleteMemberWithBalance(t *testing.T) {
	defer teardown(t)

	m := "Test"
	testGroup.Members = append(testGroup.Members, member.Member{
		Name:    m,
		Balance: 2.5,
	})

	err := testGroup.DeleteMember(m)

	if err == nil {
		t.Error("Deleting a member with non-zero balance didn't return an error.")
	}
}
