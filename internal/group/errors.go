package group

import "fmt"

// An ExistingMembersError is used when adding existing members to a group.
type ExistingMembersError struct {
	GroupID int
	Members []string
}

func (e *ExistingMembersError) Error() string {
	return fmt.Sprintf("The following members were already present in the group (%d): %v", e.GroupID, e.Members)
}

// A DeletingBalanceError is used when trying to delete a member with non zero balance from a group.
type DeletingBalanceError struct {
	GroupID int
	Members []string
}

func (e *DeletingBalanceError) Error() string {
	return fmt.Sprintf("The following members have a non zero balance in the group (%d) and can't be deleted: %v", e.GroupID, e.Members)
}

// A MembersNotFoundError is used when trying to access members that are not present in a group.
type MembersNotFoundError struct {
	GroupID int
	Members []string
}

func (e *MembersNotFoundError) Error() string {
	return fmt.Sprintf("The following members are not present in the group (%d): %v", e.GroupID, e.Members)
}

// A DeleteMembersError is used to combine DeletingBalanceError and MembersNotFoundError.
type DeleteMembersError struct {
	GroupID         int
	NotFoundMembers []string
	BalanceMembers  []string
}

func (e *DeleteMembersError) Error() string {
	return fmt.Sprintf("Couldn't delete some members from the group (%d), reason:\n"+
		"The following members are not present: %v\n"+
		"The following members have a non zero balance: %v", e.GroupID, e.NotFoundMembers, e.BalanceMembers)
}
