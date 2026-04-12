package usecase

import (
	"context"
	"tericcabrel/instech/internal/domain"
	"tericcabrel/instech/internal/repository"
)

type AddToolInput struct {
	Name        string
	Slug        string
	Category    string
	SubType     string
	Prolang     string
	ReleaseYear int
	Devstatus   string
	Details     string
	UseCases    []string
	Tags        []string
	Website     string
	Github      string
}

func AddToolUsecase(toolRepository repository.ToolRepositoryInterface, input AddToolInput) (domain.Tool, error) {
	tool := domain.Tool{
		Name:        input.Name,
		Slug:        input.Slug,
		Category:    input.Category,
		SubType:     input.SubType,
		Prolang:     input.Prolang,
		ReleaseYear: input.ReleaseYear,
		Devstatus:   input.Devstatus,
		Details:     input.Details,
		UseCases:    input.UseCases,
		Tags:        input.Tags,
		Website:     input.Website,
		Github:      input.Github,
	}
	tool, err := toolRepository.CreateTool(context.Background(), tool)
	if err != nil {
		// TODO: Handle error
		return domain.Tool{}, err
	}

	return tool, nil
}
