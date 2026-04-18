package usecase

import (
	"context"
	"database/sql"
	"errors"
	"tericcabrel/instech/internal/common"
	"tericcabrel/instech/internal/domain"
	"tericcabrel/instech/internal/repository"
)

type GetToolBySlugUseCase struct {
	ToolRepository repository.ToolRepositoryInterface
}

func (uc *GetToolBySlugUseCase) Execute(slug string) (domain.Tool, error) {
	tool, err := uc.ToolRepository.GetToolBySlug(context.Background(), slug)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Tool{}, common.ErrResourceNotFound{Id: slug}
		}
		return domain.Tool{}, err
	}

	return tool, nil
}
