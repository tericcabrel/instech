package domain

import (
	"fmt"
	"strings"
)

type ErrInvalidRelationshipKind struct {
	Kind    string
	Message string
}

func (e ErrInvalidRelationshipKind) Error() string {
	return fmt.Sprintf("Invalid relationship kind: %s", e.Kind)
}

type ErrInvalidToolCategory struct {
	Category string
	Message  string
}

func (e ErrInvalidToolCategory) Error() string {
	return fmt.Sprintf("Invalid tool category: %s", e.Category)
}

type ErrInvalidToolSubType struct {
	SubType string
	Message string
}

func (e ErrInvalidToolSubType) Error() string {
	return fmt.Sprintf("Invalid tool sub type: %s", e.SubType)
}

type ErrInvalidToolDevstatus struct {
	Devstatus string
	Message   string
}

func (e ErrInvalidToolDevstatus) Error() string {
	return fmt.Sprintf("Invalid tool dev status: %s", e.Devstatus)
}

type ErrInvalidField struct {
	Fields map[string]string
}

func (e ErrInvalidField) Error() string {
	fields := []string{}
	for field, message := range e.Fields {
		fields = append(fields, fmt.Sprintf("%s: %s", field, message))
	}
	return fmt.Sprintf("Invalid fields: %s", strings.Join(fields, ", "))
}
