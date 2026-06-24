package usecase

import (
	"context"
	"database/sql"
	"errors"
	"tericcabrel/instech/internal/common"
	"tericcabrel/instech/internal/domain"
	"tericcabrel/instech/internal/repository"
)

type AddToolUseCase struct {
	ToolRepository repository.ToolRepositoryInterface
}

type AddToolInput struct {
	Name        string
	Slug        string
	Category    string
	SubType     string
	Prolang     *string
	DevStatus   string
	Details     *string
	Website     *string
	Github      *string
	UseCases    []string
	Tags        []string
	ReleaseYear int
}

func (uc *AddToolUseCase) Execute(input AddToolInput) (domain.Tool, error) {
	tool, err := domain.CreateTool(domain.CreateToolInput{
		Name:        input.Name,
		Slug:        input.Slug,
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

	existingTool, err := uc.ToolRepository.GetToolBySlug(context.Background(), input.Slug)

	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return domain.Tool{}, err
		}
	}
	if existingTool.Id != 0 {
		return domain.Tool{}, common.ErrResourceAlreadyExists{Id: input.Slug, Message: "The tool already exists"}
	}

	return uc.ToolRepository.CreateTool(context.Background(), tool)
}
