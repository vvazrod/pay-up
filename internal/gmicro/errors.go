package gmicro

import "fmt"

// NotFoundError used when an item isn't found.
type NotFoundError struct {
	msg string
	id  string
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("%s [ID]: %s", e.msg, e.id)
}

// AlreadyPresentError used when trying to insert a member name already in use.
type AlreadyPresentError struct {
	msg     string
	groupid string
	name    string
}

func (e *AlreadyPresentError) Error() string {
	return fmt.Sprintf("%s [GroupID]: %s [Name]: %s", e.msg, e.groupid, e.name)
}

// BalanceError used when trying to delete a member with non-zero balance.
type BalanceError struct {
	msg      string
	groupid  string
	memberid string
	balance  float32
}

func (e *BalanceError) Error() string {
	return fmt.Sprintf("%s [GroupID]: %s [MemberID]: %s [Balance]: %f", e.msg, e.groupid, e.memberid, e.balance)
}
