package usecase

import (
	"context"
	"database/sql"
	"errors"
	"strconv"
	"tericcabrel/instech/internal/common"
	"tericcabrel/instech/internal/domain"
	"tericcabrel/instech/internal/repository"
)

type UpdateRelationshipUseCase struct {
	RelationshipRepository repository.RelationshipRepositoryInterface
	ToolRepository         repository.ToolRepositoryInterface
}

type UpdateRelationshipInput struct {
	FromToolId int
	ToToolId   int
	Kind       string
	Metadata   domain.RelationshipMetadata
}

func (uc *UpdateRelationshipUseCase) Execute(Id int, input UpdateRelationshipInput) (domain.Relationship, error) {
	var err error

	relationship, err := uc.RelationshipRepository.GetRelationshipById(context.Background(), Id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Relationship{}, common.ErrResourceNotFound{Id: strconv.Itoa(Id), Message: "The relationship was not found"}
		}
		return domain.Relationship{}, err
	}

	if relationship.FromToolId != input.FromToolId {
		_, err = uc.ToolRepository.GetToolById(context.Background(), input.FromToolId)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return domain.Relationship{}, common.ErrResourceNotFound{Id: strconv.Itoa(input.FromToolId), Message: "The source tool was not found"}
			}
			return domain.Relationship{}, err
		}
	}

	if relationship.ToToolId != input.ToToolId {
		_, err = uc.ToolRepository.GetToolById(context.Background(), input.ToToolId)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return domain.Relationship{}, common.ErrResourceNotFound{Id: strconv.Itoa(input.ToToolId), Message: "The target tool was not found"}
			}
			return domain.Relationship{}, err
		}
	}

	err = relationship.Update(domain.UpdateRelationshipInput{
		FromToolId: input.FromToolId,
		ToToolId:   input.ToToolId,
		Kind:       input.Kind,
		Metadata:   input.Metadata,
	})
	if err != nil {
		return domain.Relationship{}, err
	}

	updatedRelationship, err := uc.RelationshipRepository.UpdateRelationship(context.Background(), relationship)
	if err != nil {
		return domain.Relationship{}, err
	}

	return updatedRelationship, nil
}
