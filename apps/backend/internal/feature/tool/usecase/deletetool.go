package usecase

import (
	"context"
	"tericcabrel/instech/internal/repository"
)

func DeleteToolUsecase(toolRepository repository.ToolRepositoryInterface, slug string) error {
	err := toolRepository.DeleteTool(context.Background(), slug)
	if err != nil {
		return err
	}
	return nil
}
