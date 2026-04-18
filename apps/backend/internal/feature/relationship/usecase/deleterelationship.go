package usecase

import (
	"context"
	"database/sql"
	"errors"
	"strconv"
	"tericcabrel/instech/internal/common"
	"tericcabrel/instech/internal/repository"
)

func DeleteRelationshipUseCase(relationshipRepository repository.RelationshipRepositoryInterface, id int) error {
	_, err := relationshipRepository.GetRelationshipById(context.Background(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return common.ErrResourceNotFound{Id: strconv.Itoa(id), Message: "The relationship was not found"}
		}
		return err
	}

	return relationshipRepository.DeleteRelationship(context.Background(), id)
}
