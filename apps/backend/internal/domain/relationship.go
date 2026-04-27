package domain

import (
	"slices"
	"strings"
	"time"
)

type RelationshipMetadata struct {
	Reason string `json:"reason"`
}

type Relationship struct {
	CreatedAt  time.Time            `json:"created_at"`
	UpdatedAt  time.Time            `json:"updated_at"`
	Kind       string               `json:"kind"`
	Metadata   RelationshipMetadata `json:"metadata"`
	ID         int                  `json:"id"`
	FromToolID int                  `json:"from_tool_id"`
	ToToolID   int                  `json:"to_tool_id"`
}

type CreateRelationshipInput struct {
	Kind       string
	Reason     string
	FromToolID int
	ToToolID   int
}

type UpdateRelationshipInput struct {
	Kind       string
	Metadata   RelationshipMetadata
	FromToolID int
	ToToolID   int
}

var RelationshipKinds = []string{"built_on", "inspired_by", "alternative_to", "replaced_by", "used_with"}

func IsKindValid(kind string) bool {
	return slices.Contains(RelationshipKinds, kind)
}

func CreateRelationship(input CreateRelationshipInput) (Relationship, error) {
	var errors = make(map[string]string)

	if input.FromToolID <= 0 {
		errors["FromToolId"] = "The source tool ID is required"
	}
	if input.ToToolID <= 0 {
		errors["ToToolId"] = "The target tool ID is required"
	}

	if len(errors) > 0 {
		return Relationship{}, ErrInvalidField{Fields: errors}
	}

	if !IsKindValid(input.Kind) {
		return Relationship{}, ErrInvalidRelationshipKind{
			Kind:    input.Kind,
			Message: "The relationship kind is invalid. Valid kinds are: " + strings.Join(RelationshipKinds, ", "),
		}
	}

	relationship := Relationship{
		FromToolID: input.FromToolID,
		ToToolID:   input.ToToolID,
		Kind:       input.Kind,
		Metadata: RelationshipMetadata{
			Reason: input.Reason,
		},
	}

	return relationship, nil
}

func (relationship *Relationship) Update(input UpdateRelationshipInput) error {
	if !IsKindValid(input.Kind) {
		return ErrInvalidRelationshipKind{Kind: input.Kind, Message: "The relationship kind is invalid. Valid kinds are: " + strings.Join(RelationshipKinds, ", ")}
	}

	var errorsMap = make(map[string]string)

	if relationship.Kind != input.Kind {
		relationship.Kind = input.Kind
	}

	if relationship.Metadata != input.Metadata {
		relationship.Metadata = input.Metadata
	}

	if relationship.FromToolID != input.FromToolID {
		if input.FromToolID <= 0 {
			errorsMap["FromToolId"] = "The source tool ID is required"
		} else {
			relationship.FromToolID = input.FromToolID
		}
	}

	if relationship.ToToolID != input.ToToolID {
		if input.ToToolID <= 0 {
			errorsMap["ToToolId"] = "The target tool ID is required"
		} else {
			relationship.ToToolID = input.ToToolID
		}
	}

	if len(errorsMap) > 0 {
		return ErrInvalidField{Fields: errorsMap}
	}

	return nil
}
