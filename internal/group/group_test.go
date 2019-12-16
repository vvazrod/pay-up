package group

import (
	"os"
	"testing"
)

var testGroup *Group

func TestMain(m *testing.M) {
	testGroup = &Group{
		Name: "Test Group",
		Members: map[string]float32{
			"Test1": 0.0,
			"Test2": 2.5,
			"Test3": 0.0,
		},
	}

	os.Exit(m.Run())
}

func TestNew(t *testing.T) {
	g := New("Test")

	if g.Name != "Test" {
		t.Error("Group created without the given name.")
	} else if g.Members == nil {
		t.Error("Group created with nil map.")
	}
}

func TestAddMember(t *testing.T) {
	m := "Test4"

	err := testGroup.AddMember(m)

	if err != nil {
		t.Errorf("Couldn't add new member. Error: " + err.Error())
	}

	if b := testGroup.Members[m]; b != 0 {
		t.Errorf("Balance of new member (%s) is not zero (%f).", m, b)
	}
}

func TestAddExistingMember(t *testing.T) {
	m := "Test1"
	startLen := len(testGroup.Members)

	testGroup.AddMember(m)
	err := testGroup.AddMember(m)

	if err == nil {
		t.Errorf("Adding duplicate member didn't return an error.")
	} else {
		t.Log(err.Error())
	}

	if len(testGroup.Members) > startLen {
		t.Errorf("Existing member was added to the group.")
	}
}

func TestDeleteMember(t *testing.T) {
	startLen := len(testGroup.Members)

	err := testGroup.DeleteMember("Test3")

	if err != nil {
		t.Errorf("Deleting a member returned an error: %s", err.Error())
	}

	if len(testGroup.Members) != startLen-1 {
		t.Errorf("Number of members isn't correct (%d).", len(testGroup.Members))
	}
}

func TestDeleteMemberNotFound(t *testing.T) {
	err := testGroup.DeleteMember("Test5")

	if err == nil {
		t.Errorf("Deleting a non-existing member didn't return an error.")
	} else {
		t.Log(err.Error())
	}
}

func TestDeleteMemberWithBalance(t *testing.T) {
	err := testGroup.DeleteMember("Test2")

	if err == nil {
		t.Errorf("Deleting a member with non-zero balance didn't return an error.")
	} else {
		t.Log(err.Error())
	}

}
