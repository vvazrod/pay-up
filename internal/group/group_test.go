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
		t.Errorf("Adding duplicate member didn't return an error.")
	}

	if len(g.Members) != 2 {
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
	}
}

func TestDeleteMembers(t *testing.T) {
	members := map[string]float32{
		"Test1": 0.0,
		"Test2": 2.5,
		"Test3": 0.0,
	}
	g := Group{ID: 123, Name: "Test Group", Members: members}

	err := g.DeleteMembers([]string{"Test1", "Test3"})

	if err != nil {
		t.Errorf("Deleting members returned an error: %s", err.Error())
	}

	if len(g.Members) != 1 {
		t.Errorf("Number of members isn't correct (%d).", len(g.Members))
	}
}

func TestDeleteMembersNotFound(t *testing.T) {
	members := map[string]float32{
		"Test1": 0.0,
		"Test2": 2.5,
		"Test3": 0.0,
	}
	g := Group{ID: 123, Name: "Test Group", Members: members}

	err := g.DeleteMembers([]string{"Test1", "Test4"})

	if err == nil {
		t.Errorf("Deleting a non-existing member didn't return an error.")
	}

	if len(g.Members) != 2 {
		t.Errorf("Number of members isn't correct (%d).", len(g.Members))
	}
}

func TestDeleteMembersWithBalance(t *testing.T) {
	members := map[string]float32{
		"Test1": 0.0,
		"Test2": 2.5,
		"Test3": 0.0,
	}
	g := Group{ID: 123, Name: "Test Group", Members: members}

	err := g.DeleteMembers([]string{"Test1", "Test2"})

	if err == nil {
		t.Errorf("Deleting a non-existing member didn't return an error.")
	}

	if len(g.Members) != 2 {
		t.Errorf("Number of members isn't correct (%d).", len(g.Members))
	}
}

func TestDeleteMembersCombined(t *testing.T) {
	members := map[string]float32{
		"Test1": 0.0,
		"Test2": 2.5,
		"Test3": 0.0,
	}
	g := Group{ID: 123, Name: "Test Group", Members: members}

	err := g.DeleteMembers([]string{"Test1", "Test2", "Test4"})

	if err == nil {
		t.Errorf("Deleting non-existing and non-zero balance members didn't return an error.")
	}

	if len(g.Members) != 2 {
		t.Errorf("Number of members isn't correct (%d).", len(g.Members))
	}
}
