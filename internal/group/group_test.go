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
	}

	if len(g.Members) > 1 {
		t.Errorf("Existing member was added to the group.")
	}
}

func TestAddMembers(t *testing.T) {
	g := Group{ID: 123, Name: "Test Group"}

	err := g.AddMembers([]string{"Test1", "Test2"})

	if err != nil {
		t.Errorf("Couldn't add new members. Error: " + err.Error())
	}

	for k, v := range g.Members {
		if v != 0.0 {
			t.Errorf("Balance of new member (%s) is not zero (%f).", k, v)
		}
	}
}

func TestAddExistingMembers(t *testing.T) {
	g := Group{ID: 123, Name: "Test Group"}

	err := g.AddMembers([]string{"Test1", "Test2", "Test2"})

	if err == nil {
		t.Errorf("Adding duplicate member didn't return an error")
	}

	if len(g.Members) != 2 {
		t.Errorf("Existing member was added to the group.")
	}
}
