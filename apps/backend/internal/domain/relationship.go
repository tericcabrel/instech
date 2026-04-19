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
	Id         int                  `json:"id"`
	FromToolId int                  `json:"from_tool_id"`
	ToToolId   int                  `json:"to_tool_id"`
	Kind       string               `json:"kind"`
	Metadata   RelationshipMetadata `json:"metadata"`
	CreatedAt  time.Time            `json:"created_at"`
	UpdatedAt  time.Time            `json:"updated_at"`
}

type CreateRelationshipInput struct {
	FromToolId int
	ToToolId   int
	Kind       string
	Reason     string
}

type UpdateRelationshipInput struct {
	FromToolId int
	ToToolId   int
	Kind       string
	Metadata   RelationshipMetadata
}

var RELATIONSHIP_KINDS = []string{"built_on", "inspired_by", "alternative_to", "replaced_by", "used_with"}

func IsKindValid(kind string) bool {
	return slices.Contains(RELATIONSHIP_KINDS, kind)
}

func CreateRelationship(input CreateRelationshipInput) (Relationship, error) {
	var errors = make(map[string]string)

	if input.FromToolId <= 0 {
		errors["FromToolId"] = "The source tool ID is required"
	}
	if input.ToToolId <= 0 {
		errors["ToToolId"] = "The target tool ID is required"
	}

	if len(errors) > 0 {
		return Relationship{}, ErrInvalidField{Fields: errors}
	}

	if !IsKindValid(input.Kind) {
		return Relationship{}, ErrInvalidRelationshipKind{
			Kind:    input.Kind,
			Message: "The relationship kind is invalid. Valid kinds are: " + strings.Join(RELATIONSHIP_KINDS, ", "),
		}
	}

	relationship := Relationship{
		FromToolId: input.FromToolId,
		ToToolId:   input.ToToolId,
		Kind:       input.Kind,
		Metadata: RelationshipMetadata{
			Reason: input.Reason,
		},
	}

	return relationship, nil
}

func (relationship *Relationship) Update(input UpdateRelationshipInput) error {
	if !IsKindValid(input.Kind) {
		return ErrInvalidRelationshipKind{Kind: input.Kind, Message: "The relationship kind is invalid. Valid kinds are: " + strings.Join(RELATIONSHIP_KINDS, ", ")}
	}

	var errorsMap = make(map[string]string)

	if relationship.Kind != input.Kind {
		relationship.Kind = input.Kind
	}

	if relationship.Metadata != input.Metadata {
		relationship.Metadata = input.Metadata
	}

	if relationship.FromToolId != input.FromToolId {
		if input.FromToolId <= 0 {
			errorsMap["FromToolId"] = "The source tool ID is required"
		} else {
			relationship.FromToolId = input.FromToolId
		}
	}

	if relationship.ToToolId != input.ToToolId {
		if input.ToToolId <= 0 {
			errorsMap["ToToolId"] = "The target tool ID is required"
		} else {
			relationship.ToToolId = input.ToToolId
		}
	}

	if len(errorsMap) > 0 {
		return ErrInvalidField{Fields: errorsMap}
	}

	return nil
}
