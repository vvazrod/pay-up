package group

import "fmt"

// An ExistingMemberError is used when trying to add an existing member to a group.
type ExistingMemberError struct {
	GroupID int
	Member  string
}

func (e *ExistingMemberError) Error() string {
	return fmt.Sprintf("Group %d: Member %q is already present", e.GroupID, e.Member)
}

// A DeletingBalanceError is used when trying to delete a member with non-zero balance in the group.
type DeletingBalanceError struct {
	GroupID int
	Member  string
	Balance float32
}

func (e *DeletingBalanceError) Error() string {
	return fmt.Sprintf("Group %d: Member %q has non-zero balance (%f)", e.GroupID, e.Member, e.Balance)
}

// A MemberNotFoundError is used when trying to access a member not present in the group.
type MemberNotFoundError struct {
	GroupID int
	Member  string
}

func (e *MemberNotFoundError) Error() string {
	return fmt.Sprintf("Group %d: Member %q not found", e.GroupID, e.Member)
}
