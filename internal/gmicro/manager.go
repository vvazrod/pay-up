package gmicro

import (
	"math"
	"strings"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/varrrro/pay-up/internal/gmicro/group"
	"github.com/varrrro/pay-up/internal/gmicro/member"
	"github.com/varrrro/pay-up/internal/tmicro/expense"
	"github.com/varrrro/pay-up/internal/tmicro/payment"
)

// Manager interface for the groups microservice.
type Manager interface {
	CreateGroup(g *group.Group) error
	FetchGroup(id uuid.UUID) (group.Group, error)
	UpdateGroup(g *group.Group) error
	RemoveGroup(id uuid.UUID) error
	AddMember(gid uuid.UUID, m *member.Member) error
	FetchMember(gid uuid.UUID, mid uuid.UUID) (member.Member, error)
	UpdateMember(gid uuid.UUID, m *member.Member) error
	RemoveMember(gid uuid.UUID, mid uuid.UUID) error
	AddExpense(e *expense.Expense) error
	RemoveExpense(e *expense.Expense) error
	AddPayment(p *payment.Payment) error
	RemovePayment(p *payment.Payment) error
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
func (gm *GroupsManager) CreateGroup(g *group.Group) error {
	gm.DB.Create(g)

	return nil
}

// FetchGroup with the given ID.
func (gm *GroupsManager) FetchGroup(id uuid.UUID) (group.Group, error) {
	var g group.Group

	gm.DB.Preload("Members").First(&g, "id = ?", id)

	if g.ID != id {
		return g, &NotFoundError{"No group found", id}
	}

	return g, nil
}

// UpdateGroup with a new name.
func (gm *GroupsManager) UpdateGroup(g *group.Group) error {
	var prevg group.Group

	gm.DB.First(&prevg, "id = ?", g.ID)

	if prevg.ID != g.ID {
		return &NotFoundError{"No group found", g.ID}
	}

	gm.DB.Model(&prevg).Update("name", g.Name)

	return nil
}

// RemoveGroup with the given ID.
func (gm *GroupsManager) RemoveGroup(id uuid.UUID) error {
	var g group.Group

	gm.DB.First(&g, "id = ?", id)

	if g.ID != id {
		return &NotFoundError{"No group found", id}
	}

	gm.DB.Delete(&g)

	return nil
}

// AddMember to the given group.
func (gm *GroupsManager) AddMember(gid uuid.UUID, m *member.Member) error {
	var g group.Group

	gm.DB.Preload("Members").First(&g, "id = ?", gid)

	if g.ID != gid {
		return &NotFoundError{"No group found", gid}
	}

	for _, prevm := range g.Members {
		if prevm.Name == m.Name {
			return &AlreadyPresentError{"Member already present in the group", gid, m.Name}
		}
	}

	gm.DB.Model(&g).Association("Members").Append(m)

	return nil
}

// FetchMember with the given ID and group ID.
func (gm *GroupsManager) FetchMember(gid, mid uuid.UUID) (member.Member, error) {
	var m member.Member

	gm.DB.First(&m, "id = ? AND group_id = ?", mid, gid)

	if m.ID != mid {
		return m, &NotFoundError{"No member found", mid}
	}

	return m, nil
}

// TODO: Check if new name is already in use
// UpdateMember with a new name.
func (gm *GroupsManager) UpdateMember(gid uuid.UUID, m *member.Member) error {
	var prevm member.Member

	gm.DB.First(&prevm, "id = ? AND group_id = ?", m.ID, gid)

	if prevm.ID != m.ID {
		return &NotFoundError{"No member found", m.ID}
	}

	gm.DB.Model(&prevm).Update("name", m.Name)

	return nil
}

// RemoveMember with the given ID and group ID.
func (gm *GroupsManager) RemoveMember(gid, mid uuid.UUID) error {
	var m member.Member

	gm.DB.First(&m, "id = ? AND group_id = ?", mid, gid)

	if m.ID != mid {
		return &NotFoundError{"No member found", mid}
	}

	if m.Balance != 0.0 {
		return &BalanceError{"Can't delete member with balance", gid, mid, m.Balance}
	}

	gm.DB.Delete(&m)

	return nil
}

// AddExpense to a group, updating the balance of the members involved.
func (gm *GroupsManager) AddExpense(e *expense.Expense) error {
	tx := gm.DB.Begin()

	// Update payer's balance
	if err := updateBalance(tx, e.GroupID, e.Payer, e.Amount); err != nil {
		tx.Rollback()
		return err
	}

	// Update recipients' balances
	rec := strings.Split(e.Recipients, ";")
	recAmount := e.Amount / float32(len(rec))
	recAmount = float32(math.Floor(float64(recAmount*100))) / 100
	for _, r := range rec {
		rid, err := uuid.Parse(r)
		if err != nil {
			tx.Rollback()
			return err
		}

		if err := updateBalance(tx, e.GroupID, rid, -recAmount); err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()
	return nil
}

// RemoveExpense from a group, updating the balance of the members involved.
func (gm *GroupsManager) RemoveExpense(e *expense.Expense) error {
	tx := gm.DB.Begin()

	// Update payer's balance
	if err := updateBalance(tx, e.GroupID, e.Payer, -e.Amount); err != nil {
		tx.Rollback()
		return err
	}

	// Update recipients' balances
	rec := strings.Split(e.Recipients, ";")
	recAmount := e.Amount / float32(len(rec))
	recAmount = float32(math.Floor(float64(recAmount*100))) / 100
	for _, r := range rec {
		rid, err := uuid.Parse(r)
		if err != nil {
			tx.Rollback()
			return err
		}

		if err := updateBalance(tx, e.GroupID, rid, recAmount); err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()
	return nil
}

// AddPayment to a group, updating the balance of the members involved.
func (gm *GroupsManager) AddPayment(p *payment.Payment) error {
	tx := gm.DB.Begin()

	// Update payer's balance
	if err := updateBalance(tx, p.GroupID, p.Payer, p.Amount); err != nil {
		tx.Rollback()
		return err
	}

	// Update recipient's balance
	if err := updateBalance(tx, p.GroupID, p.Recipient, -p.Amount); err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

// RemovePayment from a group, updating the balance of the members involved.
func (gm *GroupsManager) RemovePayment(p *payment.Payment) error {
	tx := gm.DB.Begin()

	// Update payer's balance
	if err := updateBalance(tx, p.GroupID, p.Payer, -p.Amount); err != nil {
		tx.Rollback()
		return err
	}

	// Update recipient's balance
	if err := updateBalance(tx, p.GroupID, p.Recipient, p.Amount); err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func updateBalance(tx *gorm.DB, gid, mid uuid.UUID, amount float32) error {
	var m member.Member

	tx.First(&m, "id = ? AND group_id = ?", mid, gid)

	if m.ID != mid {
		return &NotFoundError{"No member found", mid}
	}

	tx.Model(&m).Update("balance", m.Balance+amount)

	return nil
}
