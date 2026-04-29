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
	Kind       string
	Metadata   domain.RelationshipMetadata
	FromToolID int
	ToToolID   int
}

func (uc *UpdateRelationshipUseCase) Execute(Id int, input UpdateRelationshipInput) (domain.Relationship, error) {
	var err error

	relationship, err := uc.RelationshipRepository.GetRelationshipByID(context.Background(), Id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Relationship{}, common.ErrResourceNotFound{Id: strconv.Itoa(Id), Message: "The relationship was not found"}
		}

		return domain.Relationship{}, err
	}

	if relationship.FromToolID != input.FromToolID {
		_, err = uc.ToolRepository.GetToolByID(context.Background(), input.FromToolID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return domain.Relationship{}, common.ErrResourceNotFound{Id: strconv.Itoa(input.FromToolID), Message: "The source tool was not found"}
			}

			return domain.Relationship{}, err
		}
	}

	if relationship.ToToolID != input.ToToolID {
		_, err = uc.ToolRepository.GetToolByID(context.Background(), input.ToToolID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return domain.Relationship{}, common.ErrResourceNotFound{Id: strconv.Itoa(input.ToToolID), Message: "The target tool was not found"}
			}

			return domain.Relationship{}, err
		}
	}

	err = relationship.Update(domain.UpdateRelationshipInput{
		FromToolID: input.FromToolID,
		ToToolID:   input.ToToolID,
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
