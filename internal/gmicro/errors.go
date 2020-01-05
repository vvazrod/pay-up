package gmicro

import "fmt"

import "github.com/google/uuid"

// NotFoundError used when an item isn't found.
type NotFoundError struct {
	msg string
	id  uuid.UUID
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("%s [ID]: %v", e.msg, e.id)
}

// AlreadyPresentError used when trying to insert a member name already in use.
type AlreadyPresentError struct {
	msg     string
	groupid uuid.UUID
	name    string
}

func (e *AlreadyPresentError) Error() string {
	return fmt.Sprintf("%s [GroupID]: %v [Name]: %s", e.msg, e.groupid, e.name)
}

// BalanceError used when trying to delete a member with non-zero balance.
type BalanceError struct {
	msg      string
	groupid  uuid.UUID
	memberid uuid.UUID
	balance  float32
}

func (e *BalanceError) Error() string {
	return fmt.Sprintf("%s [GroupID]: %v [MemberID]: %v [Balance]: %f", e.msg, e.groupid, e.memberid, e.balance)
}
