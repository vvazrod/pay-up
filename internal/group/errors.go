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
