package gmicro_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/varrrro/pay-up/internal/gmicro/group"
	"github.com/varrrro/pay-up/internal/gmicro/member"
	"github.com/varrrro/pay-up/internal/tmicro/expense"
	"github.com/varrrro/pay-up/internal/tmicro/payment"
)

func TestCreateGroup(t *testing.T) {
	g := group.Group{ID: uuid.New(), Name: "test"}

	if err := gm.CreateGroup(&g); err != nil {
		t.Errorf("Couldn't create group. Error: %s", err.Error())
	}

	clearDB()
}

func TestFetchGroup(t *testing.T) {
	g := group.Group{ID: uuid.New(), Name: "test"}

	if err := gm.CreateGroup(&g); err != nil {
		t.Errorf("Couldn't create group. Error: %s", err.Error())
	}

	if g2, err := gm.FetchGroup(g.ID); err != nil {
		t.Errorf("Couldn't fetch group. Error: %s", err.Error())
	} else if g2.ID != g.ID {
		t.Error("Returned group isn't correct")
	} else if len(g2.Members) != 0 {
		t.Error("Groups' members array isn't empty")
	}

	clearDB()
}

func TestFetchGroupNotFound(t *testing.T) {
	if _, err := gm.FetchGroup(uuid.New()); err == nil {
		t.Error("Fetching non-existant group didn't return an error")
	}

	clearDB()
}

func TestUpdateGroup(t *testing.T) {
	g := group.Group{ID: uuid.New(), Name: "test"}

	if err := gm.CreateGroup(&g); err != nil {
		t.Errorf("Couldn't create group. Error: %s", err.Error())
	}

	g.Name = "updated"

	if err := gm.UpdateGroup(&g); err != nil {
		t.Errorf("Couldn't update group's name. Error: %s", err.Error())
	}

	if g2, err := gm.FetchGroup(g.ID); err != nil {
		t.Errorf("Couldn't fetch group. Error: %s", err.Error())
	} else if g2.Name != "updated" {
		t.Error("Group's name wasn't updated correctly")
	}

	clearDB()
}

func TestUpdateGroupNotFound(t *testing.T) {
	g := group.Group{ID: uuid.New(), Name: "updated"}

	if err := gm.UpdateGroup(&g); err == nil {
		t.Error("Updating non-existant group didn't return an error")
	}

	clearDB()
}

func TestRemoveGroup(t *testing.T) {
	g := group.Group{ID: uuid.New(), Name: "test"}

	if err := gm.CreateGroup(&g); err != nil {
		t.Errorf("Couldn't create group. Error: %s", err.Error())
	}

	if err := gm.RemoveGroup(g.ID); err != nil {
		t.Errorf("Couldn't delete group. Error: %s", err.Error())
	}

	if _, err := gm.FetchGroup(g.ID); err == nil {
		t.Error("Fetching group after deletion didn't return an error")
	}

	clearDB()
}

func TestRemoveGroupNotFound(t *testing.T) {
	if err := gm.RemoveGroup(uuid.New()); err == nil {
		t.Error("Deleting non-existant group didn't return an error")
	}

	clearDB()
}

func TestAddMember(t *testing.T) {
	g := group.Group{ID: uuid.New(), Name: "test"}

	if err := gm.CreateGroup(&g); err != nil {
		t.Errorf("Couldn't create group. Error: %s", err.Error())
	}

	m := member.Member{ID: uuid.New(), Name: "test"}

	if err := gm.AddMember(g.ID, &m); err != nil {
		t.Errorf("Couldn't create member. Error: %s", err.Error())
	}

	if g, err := gm.FetchGroup(g.ID); err != nil {
		t.Errorf("Couldn't fetch group. Error: %s", err.Error())
	} else if len(g.Members) != 1 {
		t.Error("Member wasn't added to group correctly")
	}

	clearDB()
}

func TestAddMemberNotFound(t *testing.T) {
	m := member.Member{ID: uuid.New(), Name: "test"}

	if err := gm.AddMember(uuid.New(), &m); err == nil {
		t.Error("Adding member to non-existant group didn't return an error")
	}

	clearDB()
}

func TestAddMemberRepeated(t *testing.T) {
	g := group.Group{ID: uuid.New(), Name: "test"}

	if err := gm.CreateGroup(&g); err != nil {
		t.Errorf("Couldn't create group. Error: %s", err.Error())
	}

	m := member.Member{ID: uuid.New(), Name: "test"}

	if err := gm.AddMember(g.ID, &m); err != nil {
		t.Errorf("Couldn't create member. Error: %s", err.Error())
	}

	if err := gm.AddMember(g.ID, &m); err == nil {
		t.Error("Adding duplicate member didn't return an error")
	}

	if g, err := gm.FetchGroup(g.ID); err != nil {
		t.Errorf("Couldn't fetch group. Error: %s", err.Error())
	} else if len(g.Members) != 1 {
		t.Error("Member was added to the group nonetheless")
	}

	clearDB()
}

func TestFetchMember(t *testing.T) {
	g := group.Group{ID: uuid.New(), Name: "test"}

	if err := gm.CreateGroup(&g); err != nil {
		t.Errorf("Couldn't create group. Error: %s", err.Error())
	}

	m := member.Member{ID: uuid.New(), Name: "test"}

	if err := gm.AddMember(g.ID, &m); err != nil {
		t.Errorf("Couldn't create member. Error: %s", err.Error())
	}

	if m2, err := gm.FetchMember(g.ID, m.ID); err != nil {
		t.Errorf("Couldn't fetch member. Error: %s", err.Error())
	} else if m2.ID != m.ID || m2.GroupID != g.ID {
		t.Error("Returned member isn't correct")
	}

	clearDB()
}

func TestFetchMemberNotFound(t *testing.T) {
	g := group.Group{ID: uuid.New(), Name: "test"}

	if err := gm.CreateGroup(&g); err != nil {
		t.Errorf("Couldn't create group. Error: %s", err.Error())
	}

	if _, err := gm.FetchMember(g.ID, uuid.New()); err == nil {
		t.Error("Fetching non-existant member didn't return an error")
	}

	clearDB()
}

func TestUpdateMember(t *testing.T) {
	g := group.Group{ID: uuid.New(), Name: "test"}

	if err := gm.CreateGroup(&g); err != nil {
		t.Errorf("Couldn't create group. Error: %s", err.Error())
	}

	m := member.Member{ID: uuid.New(), Name: "test"}

	if err := gm.AddMember(g.ID, &m); err != nil {
		t.Errorf("Couldn't create member. Error: %s", err.Error())
	}

	m.Name = "updated"

	if err := gm.UpdateMember(g.ID, &m); err != nil {
		t.Errorf("Couldn't update member's name. Error: %s", err.Error())
	}

	if m2, err := gm.FetchMember(g.ID, m.ID); err != nil {
		t.Errorf("Couldn't fetch member. Error: %s", err.Error())
	} else if m2.Name != "updated" {
		t.Error("Member's name wasn't updated correctly")
	}

	clearDB()
}

func TestUpdateMemberNotFound(t *testing.T) {
	g := group.Group{ID: uuid.New(), Name: "test"}

	if err := gm.CreateGroup(&g); err != nil {
		t.Errorf("Couldn't create group. Error: %s", err.Error())
	}

	m := member.Member{ID: uuid.New(), Name: "updated"}

	if err := gm.UpdateMember(g.ID, &m); err == nil {
		t.Error("Updating non-existant member didn't return an error")
	}

	clearDB()
}

func TestUpdateMemberAlreadyPresent(t *testing.T) {
	g := group.Group{ID: uuid.New(), Name: "test"}

	if err := gm.CreateGroup(&g); err != nil {
		t.Errorf("Couldn't create group. Error: %s", err.Error())
	}

	m1 := member.Member{ID: uuid.New(), Name: "inuse"}

	if err := gm.AddMember(g.ID, &m1); err != nil {
		t.Errorf("Couldn't create member. Error: %s", err.Error())
	}

	m2 := member.Member{ID: uuid.New(), Name: "test"}

	if err := gm.AddMember(g.ID, &m2); err != nil {
		t.Errorf("Couldn't create member. Error: %s", err.Error())
	}

	m2.Name = "inuse"

	if err := gm.UpdateMember(g.ID, &m2); err == nil {
		t.Error("Updating member with name already in use didn't return an error")
	}

	if m3, err := gm.FetchMember(g.ID, m2.ID); err != nil {
		t.Errorf("Couldn't fetch member. Error: %s", err.Error())
	} else if m3.Name == "inuse" {
		t.Error("Member's name was updated nonetheless")
	}

	clearDB()
}

func TestRemoveMember(t *testing.T) {
	g := group.Group{ID: uuid.New(), Name: "test"}

	if err := gm.CreateGroup(&g); err != nil {
		t.Errorf("Couldn't create group. Error: %s", err.Error())
	}

	m := member.Member{ID: uuid.New(), Name: "test"}

	if err := gm.AddMember(g.ID, &m); err != nil {
		t.Errorf("Couldn't create member. Error: %s", err.Error())
	}

	if err := gm.RemoveMember(g.ID, m.ID); err != nil {
		t.Errorf("Couldn't delete member. Error: %s", err.Error())
	}

	if _, err := gm.FetchMember(g.ID, m.ID); err == nil {
		t.Error("Fetching member after deletion didn't return an error")
	}

	if g2, err := gm.FetchGroup(g.ID); err != nil {
		t.Errorf("Couldn't fetch group. Error: %s", err.Error())
	} else if len(g2.Members) != 0 {
		t.Error("Member wasn't removed from group correctly")
	}

	clearDB()
}

func TestRemoveMemberNotFound(t *testing.T) {
	g := group.Group{ID: uuid.New(), Name: "test"}

	if err := gm.CreateGroup(&g); err != nil {
		t.Errorf("Couldn't create group. Error: %s", err.Error())
	}

	if err := gm.RemoveMember(g.ID, uuid.New()); err == nil {
		t.Error("Deleting non-existant member didn't return an error")
	}

	clearDB()
}

func TestAddExpense(t *testing.T) {
	g := group.Group{ID: uuid.New(), Name: "test"}

	if err := gm.CreateGroup(&g); err != nil {
		t.Errorf("Couldn't create group. Error: %s", err.Error())
	}

	m1 := member.Member{ID: uuid.New(), Name: "test1"}

	if err := gm.AddMember(g.ID, &m1); err != nil {
		t.Errorf("Couldn't create member. Error: %s", err.Error())
	}

	m2 := member.Member{ID: uuid.New(), Name: "test2"}

	if err := gm.AddMember(g.ID, &m2); err != nil {
		t.Errorf("Couldn't create member. Error: %s", err.Error())
	}

	m3 := member.Member{ID: uuid.New(), Name: "test3"}

	if err := gm.AddMember(g.ID, &m3); err != nil {
		t.Errorf("Couldn't create member. Error: %s", err.Error())
	}

	e := expense.Expense{
		GroupID:    g.ID,
		Amount:     23.3,
		Payer:      m1.ID,
		Recipients: m2.ID.String() + ";" + m3.ID.String(),
	}

	if err := gm.AddExpense(&e); err != nil {
		t.Errorf("Couldn't update balances with new expense. Error: %s", err.Error())
	}

	if m1, err := gm.FetchMember(g.ID, m1.ID); err != nil {
		t.Errorf("Couldn't fetch member. Error: %s", err.Error())
	} else if m1.Balance != 23.3 {
		t.Errorf("Balance wasn't updated correctly. [Expected]: %f [Actual]: %f", 23.3, m1.Balance)
	}

	if m2, err := gm.FetchMember(g.ID, m2.ID); err != nil {
		t.Errorf("Couldn't fetch member. Error: %s", err.Error())
	} else if m2.Balance != -11.65 {
		t.Errorf("Balance wasn't updated correctly. [Expected]: %f [Actual]: %f", -11.65, m2.Balance)
	}

	if m3, err := gm.FetchMember(g.ID, m3.ID); err != nil {
		t.Errorf("Couldn't fetch member. Error: %s", err.Error())
	} else if m3.Balance != -11.65 {
		t.Errorf("Balance wasn't updated correctly. [Expected]: %f [Actual]: %f", -11.65, m3.Balance)
	}

	clearDB()
}

func TestRemoveExpense(t *testing.T) {
	g := group.Group{ID: uuid.New(), Name: "test"}

	if err := gm.CreateGroup(&g); err != nil {
		t.Errorf("Couldn't create group. Error: %s", err.Error())
	}

	m1 := member.Member{ID: uuid.New(), Name: "test1"}

	if err := gm.AddMember(g.ID, &m1); err != nil {
		t.Errorf("Couldn't create member. Error: %s", err.Error())
	}

	m2 := member.Member{ID: uuid.New(), Name: "test2"}

	if err := gm.AddMember(g.ID, &m2); err != nil {
		t.Errorf("Couldn't create member. Error: %s", err.Error())
	}

	m3 := member.Member{ID: uuid.New(), Name: "test3"}

	if err := gm.AddMember(g.ID, &m3); err != nil {
		t.Errorf("Couldn't create member. Error: %s", err.Error())
	}

	e := expense.Expense{
		GroupID:    g.ID,
		Amount:     23.3,
		Payer:      m1.ID,
		Recipients: m2.ID.String() + ";" + m3.ID.String(),
	}

	if err := gm.RemoveExpense(&e); err != nil {
		t.Errorf("Couldn't update balances with new expense. Error: %s", err.Error())
	}

	if m1, err := gm.FetchMember(g.ID, m1.ID); err != nil {
		t.Errorf("Couldn't fetch member. Error: %s", err.Error())
	} else if m1.Balance != -23.3 {
		t.Errorf("Balance wasn't updated correctly. [Expected]: %f [Actual]: %f", -23.3, m1.Balance)
	}

	if m2, err := gm.FetchMember(g.ID, m2.ID); err != nil {
		t.Errorf("Couldn't fetch member. Error: %s", err.Error())
	} else if m2.Balance != 11.65 {
		t.Errorf("Balance wasn't updated correctly. [Expected]: %f [Actual]: %f", 11.65, m2.Balance)
	}

	if m3, err := gm.FetchMember(g.ID, m3.ID); err != nil {
		t.Errorf("Couldn't fetch member. Error: %s", err.Error())
	} else if m3.Balance != 11.65 {
		t.Errorf("Balance wasn't updated correctly. [Expected]: %f [Actual]: %f", 11.65, m3.Balance)
	}

	clearDB()
}

func TestAddPayment(t *testing.T) {
	g := group.Group{ID: uuid.New(), Name: "test"}

	if err := gm.CreateGroup(&g); err != nil {
		t.Errorf("Couldn't create group. Error: %s", err.Error())
	}

	m1 := member.Member{ID: uuid.New(), Name: "test1"}

	if err := gm.AddMember(g.ID, &m1); err != nil {
		t.Errorf("Couldn't create member. Error: %s", err.Error())
	}

	m2 := member.Member{ID: uuid.New(), Name: "test2"}

	if err := gm.AddMember(g.ID, &m2); err != nil {
		t.Errorf("Couldn't create member. Error: %s", err.Error())
	}

	p := payment.Payment{
		GroupID:   g.ID,
		Amount:    25.3,
		Payer:     m1.ID,
		Recipient: m2.ID,
	}

	if err := gm.AddPayment(&p); err != nil {
		t.Errorf("Couldn't update balances with new expense. Error: %s", err.Error())
	}

	if m1, err := gm.FetchMember(g.ID, m1.ID); err != nil {
		t.Errorf("Couldn't fetch member. Error: %s", err.Error())
	} else if m1.Balance != 25.3 {
		t.Errorf("Balance wasn't updated correctly. [Expected]: %f [Actual]: %f", 25.3, m1.Balance)
	}

	if m2, err := gm.FetchMember(g.ID, m2.ID); err != nil {
		t.Errorf("Couldn't fetch member. Error: %s", err.Error())
	} else if m2.Balance != -25.3 {
		t.Errorf("Balance wasn't updated correctly. [Expected]: %f [Actual]: %f", -25.3, m2.Balance)
	}

	clearDB()
}

func TestRemovePayment(t *testing.T) {
	g := group.Group{ID: uuid.New(), Name: "test"}

	if err := gm.CreateGroup(&g); err != nil {
		t.Errorf("Couldn't create group. Error: %s", err.Error())
	}

	m1 := member.Member{ID: uuid.New(), Name: "test1"}

	if err := gm.AddMember(g.ID, &m1); err != nil {
		t.Errorf("Couldn't create member. Error: %s", err.Error())
	}

	m2 := member.Member{ID: uuid.New(), Name: "test2"}

	if err := gm.AddMember(g.ID, &m2); err != nil {
		t.Errorf("Couldn't create member. Error: %s", err.Error())
	}

	p := payment.Payment{
		GroupID:   g.ID,
		Amount:    25.3,
		Payer:     m1.ID,
		Recipient: m2.ID,
	}

	if err := gm.RemovePayment(&p); err != nil {
		t.Errorf("Couldn't update balances with new expense. Error: %s", err.Error())
	}

	if m1, err := gm.FetchMember(g.ID, m1.ID); err != nil {
		t.Errorf("Couldn't fetch member. Error: %s", err.Error())
	} else if m1.Balance != -25.3 {
		t.Errorf("Balance wasn't updated correctly. [Expected]: %f [Actual]: %f", -25.3, m1.Balance)
	}

	if m2, err := gm.FetchMember(g.ID, m2.ID); err != nil {
		t.Errorf("Couldn't fetch member. Error: %s", err.Error())
	} else if m2.Balance != 25.3 {
		t.Errorf("Balance wasn't updated correctly. [Expected]: %f [Actual]: %f", 25.3, m2.Balance)
	}

	clearDB()
}
