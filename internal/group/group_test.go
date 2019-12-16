package group

import (
	"testing"
)

func TestAddMember(t *testing.T) {
	g := Group{ID: 123, Name: "Test Group"}
	m := "Test Member"

	err := g.AddMember(m)

	if err != nil {
		t.Errorf("Couldn't add new member. Error: " + err.Error())
	}

	if b := g.Members[m]; b != 0 {
		t.Errorf("Balance of new member (%s) is not zero (%f).", m, b)
	}
}

func TestAddExistingMember(t *testing.T) {
	g := Group{ID: 123, Name: "Test Group"}
	m := "Test Member"

	g.AddMember(m)
	err := g.AddMember(m)

	if err == nil {
		t.Errorf("Adding duplicate member didn't return an error.")
	} else {
		t.Log(err.Error())
	}

	if len(g.Members) > 1 {
		t.Errorf("Existing member was added to the group.")
	}
}

func TestDeleteMember(t *testing.T) {
	members := map[string]float32{
		"Test1": 0.0,
		"Test2": 2.5,
	}
	g := Group{ID: 123, Name: "Test Group", Members: members}

	err := g.DeleteMember("Test1")

	if err != nil {
		t.Errorf("Deleting a member returned an error: %s", err.Error())
	}

	if len(g.Members) != 1 {
		t.Errorf("Number of members isn't correct (%d).", len(g.Members))
	}
}

func TestDeleteMemberNotFound(t *testing.T) {
	members := map[string]float32{
		"Test1": 0.0,
		"Test2": 2.5,
	}
	g := Group{ID: 123, Name: "Test Group", Members: members}

	err := g.DeleteMember("Test3")

	if err == nil {
		t.Errorf("Deleting a non-existing member didn't return an error.")
	} else {
		t.Log(err.Error())
	}
}

func TestDeleteMemberWithBalance(t *testing.T) {
	members := map[string]float32{
		"Test1": 0.0,
		"Test2": 2.5,
	}
	g := Group{ID: 123, Name: "Test Group", Members: members}

	err := g.DeleteMember("Test2")

	if err == nil {
		t.Errorf("Deleting a member with non-zero balance didn't return an error.")
	} else {
		t.Log(err.Error())
	}
}
