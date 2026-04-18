package domain

import "fmt"

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
