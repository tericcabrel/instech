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

type ErrInvalidRelationshipKind struct {
	Kind    string
	Message string
}

func (e ErrInvalidRelationshipKind) Error() string {
	return fmt.Sprintf("Invalid relationship kind: %s", e.Kind)
}
