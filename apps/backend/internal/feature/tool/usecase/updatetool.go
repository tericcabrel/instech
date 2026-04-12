package usecase

import (
	"context"
	"database/sql"
	"errors"
	"tericcabrel/instech/internal/domain"
	"tericcabrel/instech/internal/repository"
)

type UpdateToolInput struct {
	Name        string
	Category    string
	SubType     string
	Prolang     string
	ReleaseYear int
	DevStatus   string
	Details     string
	UseCases    []string
	Tags        []string
	Website     string
	Github      string
}

func UpdateToolUsecase(toolRepository repository.ToolRepositoryInterface, slug string, input UpdateToolInput) (domain.Tool, error) {
	tool, err := toolRepository.GetToolBySlug(context.Background(), slug)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Tool{}, ErrToolNotFound{Slug: slug}
		}
		return domain.Tool{}, err
	}

	tool.Name = input.Name
	tool.Category = input.Category
	tool.SubType = input.SubType
	tool.Prolang = input.Prolang
	tool.ReleaseYear = input.ReleaseYear
	tool.Devstatus = input.DevStatus
	tool.Details = input.Details
	tool.UseCases = input.UseCases
	tool.Tags = input.Tags
	tool.Website = input.Website
	tool.Github = input.Github

	updatedTool, err := toolRepository.UpdateTool(context.Background(), tool)

	if err != nil {
		return domain.Tool{}, err
	}
	return updatedTool, nil
}
