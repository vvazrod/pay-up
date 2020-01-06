package tmicro

import (
	"fmt"

	"github.com/google/uuid"
)

// NotFoundError used when an item isn't found.
type NotFoundError struct {
	msg string
	id  uuid.UUID
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("%s [Group ID]: %v", e.msg, e.id)
}

// UUIDParseError used when a string can't be parsed as a UUID.
type UUIDParseError struct {
	msg string
	val string
	err error
}

func (e *UUIDParseError) Error() string {
	return fmt.Sprintf("%s [Value]: %s [Error]: %s", e.msg, e.val, e.err.Error())
}
