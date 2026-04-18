package common

import "fmt"

type ErrResourceNotFound struct {
	Id      string
	Message string
}

func (e ErrResourceNotFound) Error() string {
	var msg string = "Resource not found"
	if e.Message != "" {
		msg = e.Message
	}
	return fmt.Sprintf("%s: %s", msg, e.Id)
}
