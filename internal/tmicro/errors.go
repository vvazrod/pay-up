package tmicro

import "fmt"

// NotFoundError used when an item isn't found.
type NotFoundError struct {
	msg string
	id  string
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("%s [Group ID]: %s", e.msg, e.id)
}
