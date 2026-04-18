package usecase

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"tericcabrel/instech/internal/common"
	"tericcabrel/instech/internal/domain"
	"tericcabrel/instech/internal/repository"
)

type CreateRelationshipInput struct {
	FromToolId string
	ToToolId   string
	Kind       string
	Metadata   struct {
		Reason string
	}
}

func CreateRelationshipUseCase(relationshipRepository repository.RelationshipRepositoryInterface, toolRepository repository.ToolRepositoryInterface, input CreateRelationshipInput) (domain.Relationship, error) {
	fromTool, err := toolRepository.GetToolBySlug(context.Background(), input.FromToolId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Relationship{}, common.ErrResourceNotFound{Id: input.FromToolId, Message: "The source tool was not found"}
		}
		return domain.Relationship{}, err
	}

	toTool, err := toolRepository.GetToolBySlug(context.Background(), input.ToToolId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Relationship{}, common.ErrResourceNotFound{Id: input.ToToolId, Message: "The target tool was not found"}
		}
		return domain.Relationship{}, err
	}

	relationship := domain.Relationship{
		FromToolId: fromTool.Id,
		ToToolId:   toTool.Id,
		Kind:       input.Kind,
		Metadata: domain.RelationshipMetadata{
			Reason: input.Metadata.Reason,
		},
	}

	if !relationship.IsKindValid() {
		return domain.Relationship{}, common.ErrInvalidRelationshipKind{Kind: input.Kind, Message: "The relationship kind is invalid. Valid kinds are: " + strings.Join(domain.RELATIONSHIP_KINDS, ", ")}
	}

	createdRelationship, err := relationshipRepository.CreateRelationship(context.Background(), relationship)
	if err != nil {
		return domain.Relationship{}, err
	}

	return createdRelationship, nil
}
