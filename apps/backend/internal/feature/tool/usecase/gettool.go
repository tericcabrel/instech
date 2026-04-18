package usecase

import (
	"context"
	"database/sql"
	"errors"
	"tericcabrel/instech/internal/common"
	"tericcabrel/instech/internal/domain"
	"tericcabrel/instech/internal/repository"
)

func GetToolBySlugUsecase(toolRepository repository.ToolRepositoryInterface, slug string) (domain.Tool, error) {
	tool, err := toolRepository.GetToolBySlug(context.Background(), slug)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Tool{}, common.ErrResourceNotFound{Id: slug}
		}
		return domain.Tool{}, err
	}

	return tool, nil
}
