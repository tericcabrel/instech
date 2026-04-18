package usecase

import (
	"context"
	"strconv"
	"strings"
	"tericcabrel/instech/internal/common"
	"tericcabrel/instech/internal/domain"
	"tericcabrel/instech/internal/repository"
)

type UpdateRelationshipInput struct {
	FromToolId int
	ToToolId   int
	Kind       string
	Metadata   domain.RelationshipMetadata
}

func UpdateRelationshipUseCase(relationshipRepository repository.RelationshipRepositoryInterface, toolRepository repository.ToolRepositoryInterface, Id int, input UpdateRelationshipInput) (domain.Relationship, error) {
	relationship, err := relationshipRepository.GetRelationshipById(context.Background(), Id)
	if err != nil {
		return domain.Relationship{}, common.ErrResourceNotFound{Id: strconv.Itoa(input.FromToolId), Message: "The relationship was not found"}
	}

	if relationship.FromToolId != input.FromToolId {
		fromTool, err := toolRepository.GetToolById(context.Background(), input.FromToolId)
		if err != nil {
			return domain.Relationship{}, common.ErrResourceNotFound{Id: strconv.Itoa(input.FromToolId), Message: "The source tool was not found"}
		}
		relationship.FromToolId = fromTool.Id
	}

	if relationship.ToToolId != input.ToToolId {
		toTool, err := toolRepository.GetToolById(context.Background(), input.ToToolId)
		if err != nil {
			return domain.Relationship{}, common.ErrResourceNotFound{Id: strconv.Itoa(input.ToToolId), Message: "The target tool was not found"}
		}
		relationship.ToToolId = toTool.Id
	}

	if relationship.Kind != input.Kind {
		relationship.Kind = input.Kind
	}

	if !relationship.IsKindValid() {
		return domain.Relationship{}, common.ErrInvalidRelationshipKind{Kind: input.Kind, Message: "The relationship kind is invalid. Valid kinds are: " + strings.Join(domain.RELATIONSHIP_KINDS, ", ")}
	}

	if relationship.Metadata != input.Metadata {
		relationship.Metadata = input.Metadata
	}

	updatedRelationship, err := relationshipRepository.UpdateRelationship(context.Background(), relationship)
	if err != nil {
		return domain.Relationship{}, err
	}

	return updatedRelationship, nil
}
