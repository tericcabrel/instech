package usecase

import (
	"context"
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
	Prolang     string
	ReleaseYear int
	Devstatus   string
	Details     string
	UseCases    []string
	Tags        []string
	Website     string
	Github      string
}

func (uc *AddToolUseCase) Execute(input AddToolInput) (domain.Tool, error) {
	tool, err := domain.CreateTool(domain.CreateToolInput{
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
	})
	if err != nil {
		return domain.Tool{}, err
	}

	return uc.ToolRepository.CreateTool(context.Background(), tool)
}
