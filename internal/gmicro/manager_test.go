package gmicro_test

import (
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/varrrro/pay-up/internal/gmicro"
	"github.com/varrrro/pay-up/internal/gmicro/group"
	"github.com/varrrro/pay-up/internal/gmicro/member"
)

var manager *gmicro.GroupsManager

func TestMain(m *testing.M) {
	// Open connection to test DB
	db, _ := gorm.Open("sqlite3", ":memory:")
	defer db.Close()

	// Create tables
	db.CreateTable(&group.Group{}, &member.Member{})

	// Create manager with test DB connection
	manager = gmicro.NewManager(db)

	// Build handlers with manager
	httpHandlers := gmicro.NewHTTPHandlers(manager)

	// Build router with handlers
	router = gmicro.NewRouter(httpHandlers)

	os.Exit(m.Run())
}

func TestCreateGroup(t *testing.T) {
	g, err := manager.CreateGroup("test")

	if err != nil {
		t.Errorf("Couldn't create group. Error: %s", err.Error())
	}

	if g.Name != "test" {
		t.Error("Returned group doesn't have the wanted name.")
	}
}

func TestFetchGroup(t *testing.T) {
	g, _ := manager.CreateGroup("test")

	_, err := manager.FetchGroup(uuid.MustParse(g.ID))

	if err != nil {
		t.Errorf("Couldn't fetch group. Error: %s", err.Error())
	}
}

func TestFetchGroupNotFound(t *testing.T) {
	testID := uuid.New()

	_, err := manager.FetchGroup(testID)

	if err == nil {
		t.Error("Fetching non-existant group didn't return an error.")
	}
}

func TestUpdateGroup(t *testing.T) {
	g, _ := manager.CreateGroup("test")

	err := manager.UpdateGroup(uuid.MustParse(g.ID), "updated")

	if err != nil {
		t.Errorf("Couldn't update group's name. Error: %s", err.Error())
	}

	g, _ = manager.FetchGroup(uuid.MustParse(g.ID))

	if g.Name != "updated" {
		t.Error("Group's name wasn't updated correctly.")
	}
}

func TestUpdateGroupNotFound(t *testing.T) {
	testID := uuid.New()

	err := manager.UpdateGroup(testID, "updated")

	if err == nil {
		t.Error("Updating non-existant group didn't return an error.")
	}
}

func TestRemoveGroup(t *testing.T) {
	g, _ := manager.CreateGroup("test")

	err := manager.RemoveGroup(uuid.MustParse(g.ID))

	if err != nil {
		t.Errorf("Couldn't delete group. Error: %s", err.Error())
	}

	_, err = manager.FetchGroup(uuid.MustParse(g.ID))

	if err == nil {
		t.Error("Fetching group after deletion didn't return an error.")
	}
}

func TestRemoveGroupNotFound(t *testing.T) {
	testID := uuid.New()

	err := manager.RemoveGroup(testID)

	if err == nil {
		t.Error("Deleting non-existant group didn't return an error.")
	}
}

func TestAddMember(t *testing.T) {
	g, _ := manager.CreateGroup("test")

	m, err := manager.AddMember(uuid.MustParse(g.ID), "test")

	if err != nil {
		t.Errorf("Couldn't create member. Error: %s", err.Error())
	}

	if m.Name != "test" {
		t.Error("Returned member doesn't have the wanted name.")
	}

	g, _ = manager.FetchGroup(uuid.MustParse(g.ID))

	if len(g.Members) != 1 {
		t.Error("Member wasn't added to group correctly.")
	}
}

func TestAddMemberNotFound(t *testing.T) {
	testID := uuid.New()

	_, err := manager.AddMember(testID, "test")

	if err == nil {
		t.Error("Adding member to non-existant group didn't return an error.")
	}
}

func TestAddMemberRepeated(t *testing.T) {
	g, _ := manager.CreateGroup("test")

	manager.AddMember(uuid.MustParse(g.ID), "test")
	_, err := manager.AddMember(uuid.MustParse(g.ID), "test")

	g, _ = manager.FetchGroup(uuid.MustParse(g.ID))

	if err == nil {
		t.Error("Adding duplicate member didn't return an error.")
	}
}

func TestFetchMember(t *testing.T) {
	g, _ := manager.CreateGroup("test")
	m, _ := manager.AddMember(uuid.MustParse(g.ID), "test")

	_, err := manager.FetchMember(uuid.MustParse(g.ID), uuid.MustParse(m.ID))

	if err != nil {
		t.Errorf("Couldn't fetch member. Error: %s", err.Error())
	}
}

func TestFetchMemberNotFound(t *testing.T) {
	g, _ := manager.CreateGroup("test")
	testID := uuid.New()

	_, err := manager.FetchMember(uuid.MustParse(g.ID), testID)

	if err == nil {
		t.Error("Fetching non-existant member didn't return an error.")
	}
}

func TestUpdateMember(t *testing.T) {
	g, _ := manager.CreateGroup("test")
	m, _ := manager.AddMember(uuid.MustParse(g.ID), "test")

	err := manager.UpdateMember(uuid.MustParse(g.ID), uuid.MustParse(m.ID), "updated")

	if err != nil {
		t.Errorf("Couldn't update member's name. Error: %s", err.Error())
	}

	m, _ = manager.FetchMember(uuid.MustParse(g.ID), uuid.MustParse(m.ID))

	if m.Name != "updated" {
		t.Error("Member's name wasn't updated correctly.")
	}
}

func TestUpdateMemberNotFound(t *testing.T) {
	g, _ := manager.CreateGroup("test")
	testID := uuid.New()

	err := manager.UpdateMember(uuid.MustParse(g.ID), testID, "updated")

	if err == nil {
		t.Error("Updating non-existant member didn't return an error.")
	}
}

func TestRemoveMember(t *testing.T) {
	g, _ := manager.CreateGroup("test")
	m, _ := manager.AddMember(uuid.MustParse(g.ID), "test")

	err := manager.RemoveMember(uuid.MustParse(g.ID), uuid.MustParse(m.ID))

	if err != nil {
		t.Errorf("Couldn't delete member. Error: %s", err.Error())
	}

	_, err = manager.FetchMember(uuid.MustParse(g.ID), uuid.MustParse(m.ID))

	if err == nil {
		t.Error("Fetching member after deletion didn't return an error.")
	}

	g, _ = manager.FetchGroup(uuid.MustParse(g.ID))

	if len(g.Members) != 0 {
		t.Error("Member wasn't removed from group correctly.")
	}
}

func TestRemoveMemberNotFound(t *testing.T) {
	g, _ := manager.CreateGroup("test")
	testID := uuid.New()

	err := manager.RemoveMember(uuid.MustParse(g.ID), testID)

	if err == nil {
		t.Error("Deleting non-existant member didn't return an error.")
	}
}

func TestAddExpense(t *testing.T) {
	g, _ := manager.CreateGroup("test")
	m1, _ := manager.AddMember(uuid.MustParse(g.ID), "test1")
	m2, _ := manager.AddMember(uuid.MustParse(g.ID), "test2")
	m3, _ := manager.AddMember(uuid.MustParse(g.ID), "test3")

	if err := manager.AddExpense(
		23.3,
		uuid.MustParse(g.ID),
		uuid.MustParse(m1.ID),
		&[]uuid.UUID{uuid.MustParse(m2.ID), uuid.MustParse(m3.ID)},
	); err != nil {
		t.Errorf("Couldn't update balances with new expense. Error: %s", err.Error())
	}

	m1, _ = manager.FetchMember(uuid.MustParse(g.ID), uuid.MustParse(m1.ID))
	if m1.Balance != 23.3 {
		t.Errorf("Balance wasn't updated correctly. [Expected]: %f [Actual]: %f", 23.3, m1.Balance)
	}

	m2, _ = manager.FetchMember(uuid.MustParse(g.ID), uuid.MustParse(m2.ID))
	if m2.Balance != -11.65 {
		t.Errorf("Balance wasn't updated correctly. [Expected]: %f [Actual]: %f", -11.65, m2.Balance)
	}

	m3, _ = manager.FetchMember(uuid.MustParse(g.ID), uuid.MustParse(m3.ID))
	if m2.Balance != -11.65 {
		t.Errorf("Balance wasn't updated correctly. [Expected]: %f [Actual]: %f", -11.65, m3.Balance)
	}
}

func TestRemoveExpense(t *testing.T) {
	g, _ := manager.CreateGroup("test")
	m1, _ := manager.AddMember(uuid.MustParse(g.ID), "test1")
	m2, _ := manager.AddMember(uuid.MustParse(g.ID), "test2")
	m3, _ := manager.AddMember(uuid.MustParse(g.ID), "test3")

	if err := manager.RemoveExpense(
		23.3,
		uuid.MustParse(g.ID),
		uuid.MustParse(m1.ID),
		&[]uuid.UUID{uuid.MustParse(m2.ID), uuid.MustParse(m3.ID)},
	); err != nil {
		t.Errorf("Couldn't update balances with new expense. Error: %s", err.Error())
	}

	m1, _ = manager.FetchMember(uuid.MustParse(g.ID), uuid.MustParse(m1.ID))
	if m1.Balance != -23.3 {
		t.Errorf("Balance wasn't updated correctly. [Expected]: %f [Actual]: %f", -23.3, m1.Balance)
	}

	m2, _ = manager.FetchMember(uuid.MustParse(g.ID), uuid.MustParse(m2.ID))
	if m2.Balance != 11.65 {
		t.Errorf("Balance wasn't updated correctly. [Expected]: %f [Actual]: %f", 11.65, m2.Balance)
	}

	m3, _ = manager.FetchMember(uuid.MustParse(g.ID), uuid.MustParse(m3.ID))
	if m2.Balance != 11.65 {
		t.Errorf("Balance wasn't updated correctly. [Expected]: %f [Actual]: %f", 11.65, m3.Balance)
	}
}

func TestAddPayment(t *testing.T) {
	g, _ := manager.CreateGroup("test")
	m1, _ := manager.AddMember(uuid.MustParse(g.ID), "test1")
	m2, _ := manager.AddMember(uuid.MustParse(g.ID), "test2")

	if err := manager.AddPayment(
		25.3,
		uuid.MustParse(g.ID),
		uuid.MustParse(m1.ID),
		uuid.MustParse(m2.ID),
	); err != nil {
		t.Errorf("Couldn't update balances with new expense. Error: %s", err.Error())
	}

	m1, _ = manager.FetchMember(uuid.MustParse(g.ID), uuid.MustParse(m1.ID))
	if m1.Balance != 25.3 {
		t.Errorf("Balance wasn't updated correctly. [Expected]: %f [Actual]: %f", 25.3, m1.Balance)
	}

	m2, _ = manager.FetchMember(uuid.MustParse(g.ID), uuid.MustParse(m2.ID))
	if m2.Balance != -25.3 {
		t.Errorf("Balance wasn't updated correctly. [Expected]: %f [Actual]: %f", -25.3, m2.Balance)
	}
}

func TestRemovePayment(t *testing.T) {
	g, _ := manager.CreateGroup("test")
	m1, _ := manager.AddMember(uuid.MustParse(g.ID), "test1")
	m2, _ := manager.AddMember(uuid.MustParse(g.ID), "test2")

	if err := manager.RemovePayment(
		25.3,
		uuid.MustParse(g.ID),
		uuid.MustParse(m1.ID),
		uuid.MustParse(m2.ID),
	); err != nil {
		t.Errorf("Couldn't update balances with new expense. Error: %s", err.Error())
	}

	m1, _ = manager.FetchMember(uuid.MustParse(g.ID), uuid.MustParse(m1.ID))
	if m1.Balance != -25.3 {
		t.Errorf("Balance wasn't updated correctly. [Expected]: %f [Actual]: %f", -25.3, m1.Balance)
	}

	m2, _ = manager.FetchMember(uuid.MustParse(g.ID), uuid.MustParse(m2.ID))
	if m2.Balance != 25.3 {
		t.Errorf("Balance wasn't updated correctly. [Expected]: %f [Actual]: %f", 25.3, m2.Balance)
	}
}
