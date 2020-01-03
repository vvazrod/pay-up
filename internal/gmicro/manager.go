package gmicro

import (
	"math"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/varrrro/pay-up/internal/gmicro/group"
	"github.com/varrrro/pay-up/internal/gmicro/member"
)

// Manager interface for the groups microservice.
type Manager interface {
	CreateGroup(name string) (group.Group, error)
	FetchGroup(id uuid.UUID) (group.Group, error)
	UpdateGroup(id uuid.UUID, name string) error
	RemoveGroup(id uuid.UUID) error
	AddMember(groupid uuid.UUID, name string) (member.Member, error)
	FetchMember(groupid uuid.UUID, memberid uuid.UUID) (member.Member, error)
	UpdateMember(groupid uuid.UUID, memberid uuid.UUID, name string) error
	RemoveMember(groupid uuid.UUID, memberid uuid.UUID) error
	AddExpense(amount float32, groupid, payerid uuid.UUID, recipients *[]uuid.UUID) error
	RemoveExpense(amount float32, groupid, payerid uuid.UUID, recipients *[]uuid.UUID) error
	AddPayment(amount float32, groupid, payerid, recipientid uuid.UUID) error
	RemovePayment(amount float32, groupid, payerid, recipientid uuid.UUID) error
}

// GroupsManager that works as single source of truth.
type GroupsManager struct {
	DB *gorm.DB
}

// NewManager with the given database connection.
func NewManager(db *gorm.DB) *GroupsManager {
	return &GroupsManager{DB: db}
}

// CreateGroup with the given name.
func (gm *GroupsManager) CreateGroup(name string) (group.Group, error) {
	g := group.New(name)

	gm.DB.Create(g)

	return *g, nil
}

// FetchGroup with the given ID.
func (gm *GroupsManager) FetchGroup(id uuid.UUID) (group.Group, error) {
	var g group.Group

	gm.DB.Preload("Members").First(&g, "id = ?", id.String())

	if g.ID != id.String() {
		return g, &NotFoundError{"No group found", id.String()}
	}

	return g, nil
}

// UpdateGroup with a new name.
func (gm *GroupsManager) UpdateGroup(id uuid.UUID, name string) error {
	var g group.Group

	gm.DB.First(&g, "id = ?", id)

	if g.ID != id.String() {
		return &NotFoundError{"No group found", id.String()}
	}

	g.Name = name
	gm.DB.Save(&g)

	return nil
}

// RemoveGroup with the given ID.
func (gm *GroupsManager) RemoveGroup(id uuid.UUID) error {
	var g group.Group

	gm.DB.First(&g, "id = ?", id)

	if g.ID != id.String() {
		return &NotFoundError{"No group found", id.String()}
	}

	gm.DB.Delete(&g)

	return nil
}

// AddMember to the given group.
func (gm *GroupsManager) AddMember(groupid uuid.UUID, name string) (member.Member, error) {
	var g group.Group

	gm.DB.Preload("Members").First(&g, "id = ?", groupid)

	if g.ID != groupid.String() {
		return member.Member{}, &NotFoundError{"No group found", groupid.String()}
	}

	for _, m := range g.Members {
		if m.Name == name {
			return member.Member{}, &AlreadyPresentError{"Member already present in the group", groupid.String(), name}
		}
	}

	m := member.New(name)
	gm.DB.Model(&g).Association("Members").Append(m)

	return *m, nil
}

// FetchMember with the given ID and group ID.
func (gm *GroupsManager) FetchMember(groupid, memberid uuid.UUID) (member.Member, error) {
	var m member.Member

	gm.DB.First(&m, "id = ? AND group_id = ?", memberid, groupid)

	if m.ID != memberid.String() {
		return m, &NotFoundError{"No member found", memberid.String()}
	}

	return m, nil
}

// UpdateMember with a new name.
func (gm *GroupsManager) UpdateMember(groupid, memberid uuid.UUID, name string) error {
	var m member.Member

	gm.DB.First(&m, "id = ? AND group_id = ?", memberid, groupid)

	if m.ID != memberid.String() {
		return &NotFoundError{"No member found", memberid.String()}
	}

	m.Name = name
	gm.DB.Save(&m)

	return nil
}

// RemoveMember with the given ID and group ID.
func (gm *GroupsManager) RemoveMember(groupid, memberid uuid.UUID) error {
	var m member.Member

	gm.DB.First(&m, "id = ? AND group_id = ?", memberid, groupid)

	if m.ID != memberid.String() {
		return &NotFoundError{"No member found", memberid.String()}
	}

	if m.Balance != 0.0 {
		return &BalanceError{"Can't delete member with balance", groupid.String(), memberid.String(), m.Balance}
	}

	gm.DB.Delete(&m)

	return nil
}

// AddExpense to a group, updating the balance of the members involved.
func (gm *GroupsManager) AddExpense(amount float32, groupid, payerid uuid.UUID, recipients *[]uuid.UUID) error {
	tx := gm.DB.Begin()

	// Update payer's balance
	if err := updateBalance(tx, groupid, payerid, amount); err != nil {
		tx.Rollback()
		return err
	}

	// Update recipients' balances
	recAmount := amount / float32(len(*recipients))
	recAmount = float32(math.Floor(float64(recAmount*100))) / 100
	for _, r := range *recipients {
		if err := updateBalance(tx, groupid, r, -recAmount); err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()
	return nil
}

// RemoveExpense from a group, updating the balance of the members involved.
func (gm *GroupsManager) RemoveExpense(amount float32, groupid, payerid uuid.UUID, recipients *[]uuid.UUID) error {
	tx := gm.DB.Begin()

	// Update payer's balance
	if err := updateBalance(tx, groupid, payerid, -amount); err != nil {
		tx.Rollback()
		return err
	}

	// Update recipients' balances
	recAmount := amount / float32(len(*recipients))
	recAmount = float32(math.Floor(float64(recAmount*100))) / 100
	for _, r := range *recipients {
		if err := updateBalance(tx, groupid, r, recAmount); err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()
	return nil
}

// AddPayment to a group, updating the balance of the members involved.
func (gm *GroupsManager) AddPayment(amount float32, groupid, payerid, recipientid uuid.UUID) error {
	tx := gm.DB.Begin()

	// Update payer's balance
	if err := updateBalance(tx, groupid, payerid, amount); err != nil {
		tx.Rollback()
		return err
	}

	// Update recipient's balance
	if err := updateBalance(tx, groupid, recipientid, -amount); err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

// RemovePayment from a group, updating the balance of the members involved.
func (gm *GroupsManager) RemovePayment(amount float32, groupid, payerid, recipientid uuid.UUID) error {
	tx := gm.DB.Begin()

	// Update payer's balance
	if err := updateBalance(tx, groupid, payerid, -amount); err != nil {
		tx.Rollback()
		return err
	}

	// Update recipient's balance
	if err := updateBalance(tx, groupid, recipientid, amount); err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func updateBalance(tx *gorm.DB, groupid, memberid uuid.UUID, amount float32) error {
	var m member.Member

	tx.First(&m, "id = ? AND group_id = ?", memberid, groupid)

	if m.ID != memberid.String() {
		return &NotFoundError{"No member found", memberid.String()}
	}

	tx.Model(&m).Update("balance", m.Balance+amount)

	return nil
}
