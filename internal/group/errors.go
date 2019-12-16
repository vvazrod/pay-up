package group

import "fmt"

// An ExistingMemberError is used when trying to add an existing member to a group.
type ExistingMemberError struct {
	GroupID string
	Member  string
}

func (e *ExistingMemberError) Error() string {
	return fmt.Sprintf("Group %s: Member %q is already present", e.GroupID, e.Member)
}

// A DeletingBalanceError is used when trying to delete a member with non-zero balance in the group.
type DeletingBalanceError struct {
	GroupID string
	Member  string
	Balance float32
}

func (e *DeletingBalanceError) Error() string {
	return fmt.Sprintf("Group %s: Member %q has non-zero balance (%f)", e.GroupID, e.Member, e.Balance)
}

// A MemberNotFoundError is used when trying to access a member not present in the group.
type MemberNotFoundError struct {
	GroupID string
	Member  string
}

func (e *MemberNotFoundError) Error() string {
	return fmt.Sprintf("Group %s: Member %q not found", e.GroupID, e.Member)
}
