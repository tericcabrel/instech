package usecase

import (
	"context"
	"database/sql"
	"errors"
	"tericcabrel/instech/internal/common"
	"tericcabrel/instech/internal/domain"
	"tericcabrel/instech/internal/repository"
)

type UpdateToolUseCase struct {
	ToolRepository repository.ToolRepositoryInterface
}
type UpdateToolInput struct {
	Name        string
	Category    string
	SubType     string
	Prolang     string
	DevStatus   string
	Details     string
	Website     string
	Github      string
	UseCases    []string
	Tags        []string
	ReleaseYear int
}

func (uc *UpdateToolUseCase) Execute(slug string, input UpdateToolInput) (domain.Tool, error) {
	tool, err := uc.ToolRepository.GetToolBySlug(context.Background(), slug)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Tool{}, common.ErrResourceNotFound{Id: slug}
		}

		return domain.Tool{}, err
	}

	err = tool.Update(domain.UpdateToolInput{
		Name:        input.Name,
		Slug:        slug,
		Category:    input.Category,
		SubType:     input.SubType,
		Prolang:     input.Prolang,
		ReleaseYear: input.ReleaseYear,
		DevStatus:   input.DevStatus,
		Details:     input.Details,
		UseCases:    input.UseCases,
		Tags:        input.Tags,
		Website:     input.Website,
		Github:      input.Github,
	})

	if err != nil {
		return domain.Tool{}, err
	}

	updatedTool, err := uc.ToolRepository.UpdateTool(context.Background(), tool)

	if err != nil {
		return domain.Tool{}, err
	}

	return updatedTool, nil
}
