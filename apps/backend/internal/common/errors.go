package common

import "fmt"

type ErrResourceNotFound struct {
	Id      string
	Message string
}

func (e ErrResourceNotFound) Error() string {
	var msg = "Resource not found"
	if e.Message != "" {
		msg = e.Message
	}

	return fmt.Sprintf("%s: %s", msg, e.Id)
}

type ErrResourceAlreadyExists struct {
	Id      string
	Message string
}

func (e ErrResourceAlreadyExists) Error() string {
	var msg = "Resource already exists"
	if e.Message != "" {
		msg = e.Message
	}

	return fmt.Sprintf("%s: %s", msg, e.Id)
}
