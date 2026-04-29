package usecase

import (
	"context"
	"database/sql"
	"errors"

	"tericcabrel/instech/internal/common"
	"tericcabrel/instech/internal/domain"
	"tericcabrel/instech/internal/repository"
)

type CreateRelationshipUseCase struct {
	RelationshipRepository repository.RelationshipRepositoryInterface
	ToolRepository         repository.ToolRepositoryInterface
}

type CreateRelationshipInput struct {
	FromToolID string
	ToToolID   string
	Kind       string
	Metadata   struct {
		Reason string
	}
}

func (uc *CreateRelationshipUseCase) Execute(input CreateRelationshipInput) (domain.Relationship, error) {
	fromTool, err := uc.ToolRepository.GetToolBySlug(context.Background(), input.FromToolID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Relationship{}, common.ErrResourceNotFound{Id: input.FromToolID, Message: "The source tool was not found"}
		}

		return domain.Relationship{}, err
	}

	toTool, err := uc.ToolRepository.GetToolBySlug(context.Background(), input.ToToolID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Relationship{}, common.ErrResourceNotFound{Id: input.ToToolID, Message: "The target tool was not found"}
		}

		return domain.Relationship{}, err
	}

	relationship, err := domain.CreateRelationship(domain.CreateRelationshipInput{
		FromToolID: fromTool.Id,
		ToToolID:   toTool.Id,
		Kind:       input.Kind,
		Reason:     input.Metadata.Reason,
	})

	if err != nil {
		return domain.Relationship{}, err
	}

	createdRelationship, err := uc.RelationshipRepository.CreateRelationship(context.Background(), relationship)
	if err != nil {
		return domain.Relationship{}, err
	}

	return createdRelationship, nil
}
