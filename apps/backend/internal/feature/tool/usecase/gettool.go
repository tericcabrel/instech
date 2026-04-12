package usecase

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"tericcabrel/instech/internal/domain"
	"tericcabrel/instech/internal/repository"
)

type ErrToolNotFound struct {
	Slug string
}

func (e ErrToolNotFound) Error() string {
	return fmt.Sprintf("Tool not found: %s", e.Slug)
}

func GetToolBySlugUsecase(toolRepository repository.ToolRepositoryInterface, slug string) (domain.Tool, error) {
	tool, err := toolRepository.GetToolBySlug(context.Background(), slug)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Tool{}, ErrToolNotFound{Slug: slug}
		}
		return domain.Tool{}, err
	}

	return tool, nil
}
