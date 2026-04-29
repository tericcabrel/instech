package usecase

import (
	"context"
	"tericcabrel/instech/internal/repository"
)

type DeleteToolUseCase struct {
	ToolRepository repository.ToolRepositoryInterface
}

func (uc *DeleteToolUseCase) Execute(slug string) error {
	err := uc.ToolRepository.DeleteTool(context.Background(), slug)
	if err != nil {
		return err
	}

	return nil
}
