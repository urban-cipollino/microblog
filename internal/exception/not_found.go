package exception

import (
	"strings"
)

var errObjectNotFound = "#OBJECT# not found"

type NotFoundError struct {
	Object string
}

func (e *NotFoundError) Error() string {
	return strings.Replace(errObjectNotFound, "#OBJECT#", e.Object, 1)
}
