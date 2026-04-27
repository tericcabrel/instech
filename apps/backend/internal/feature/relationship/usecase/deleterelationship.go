package usecase

import (
	"context"
	"database/sql"
	"errors"
	"strconv"
	"tericcabrel/instech/internal/common"
	"tericcabrel/instech/internal/repository"
)

type DeleteRelationshipUseCase struct {
	RelationshipRepository repository.RelationshipRepositoryInterface
}

func (uc *DeleteRelationshipUseCase) Execute(id int) error {
	_, err := uc.RelationshipRepository.GetRelationshipByID(context.Background(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return common.ErrResourceNotFound{Id: strconv.Itoa(id), Message: "The relationship was not found"}
		}
		return err
	}

	return uc.RelationshipRepository.DeleteRelationship(context.Background(), id)
}
