package gmicro

import (
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/varrrro/pay-up/internal/gmicro/group"
	"github.com/varrrro/pay-up/internal/gmicro/member"
)

// GroupsManager that works as single source of truth.
type GroupsManager struct {
	DB *gorm.DB
}

// NewManager with the given database connection.
func NewManager(db *gorm.DB) *GroupsManager {
	return &GroupsManager{DB: db}
}

// CreateGroup with the given name.
func (gm *GroupsManager) CreateGroup(name string) error {
	g := group.New(name)

	gm.DB.Create(g)

	return nil
}

// FetchGroup with the given ID.
//
// Returns an empty Group if none are found.
func (gm *GroupsManager) FetchGroup(id uuid.UUID) group.Group {
	var g group.Group

	gm.DB.First(&g, "id = ?", id)

	return g
}

// UpdateGroup with a new name.
func (gm *GroupsManager) UpdateGroup(id uuid.UUID, name string) error {
	var g group.Group

	gm.DB.First(&g, "id = ?", id)
	gm.DB.Model(&g).Update("name", name)

	return nil
}

// DeleteGroup with the given ID.
func (gm *GroupsManager) DeleteGroup(id uuid.UUID) error {
	var g group.Group

	gm.DB.First(&g, "id = ?", id)
	gm.DB.Delete(&g)

	return nil
}

// AddMember to the given group.
func (gm *GroupsManager) AddMember(groupid uuid.UUID, name string) error {
	var g group.Group

	gm.DB.First(&g, "id = ?", groupid)
	g.AddMember(name)
	gm.DB.Save(&g)

	return nil
}

// FetchMember with the given ID and group ID.
func (gm *GroupsManager) FetchMember(groupid, memberid uuid.UUID) member.Member {
	var m member.Member

	gm.DB.First(&m, "id = ? AND groupid = ?", memberid, groupid)

	return m
}

// UpdateMember with a new name.
func (gm *GroupsManager) UpdateMember(id uuid.UUID, name string) error {
	var m member.Member

	gm.DB.First(&m, "id = ?", id)
	gm.DB.Model(&m).Update("name", name)

	return nil
}

// DeleteMember with the given ID and group ID.
func (gm *GroupsManager) DeleteMember(groupid, memberid uuid.UUID) error {
	var m member.Member

	gm.DB.First(&m, "id = ? AND groupid = ?", memberid, groupid)
	gm.DB.Delete(&m)

	return nil
}
