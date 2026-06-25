package usecase

import (
	"context"
	"database/sql"
	"errors"
	"tericcabrel/instech/internal/common"
	"tericcabrel/instech/internal/domain"
	"tericcabrel/instech/internal/repository"
)

type PatchToolUseCase struct {
	ToolRepository repository.ToolRepositoryInterface
}

type PatchToolInput struct {
	Name        common.PatchStringField         `json:"name"`
	Category    common.PatchStringField         `json:"category"`
	SubType     common.PatchStringField         `json:"subType"`
	Prolang     common.PatchNullableStringField `json:"prolang"`
	DevStatus   common.PatchStringField         `json:"devStatus"`
	Details     common.PatchNullableStringField `json:"details"`
	Website     common.PatchNullableStringField `json:"website"`
	Github      common.PatchNullableStringField `json:"github"`
	UseCases    common.PatchStringSliceField    `json:"useCases"`
	Tags        common.PatchStringSliceField    `json:"tags"`
	ReleaseYear common.PatchIntField            `json:"releaseYear"`
}

func (uc *PatchToolUseCase) Execute(slug string, input PatchToolInput) (domain.Tool, error) {
	tool, err := uc.ToolRepository.GetToolBySlug(context.Background(), slug)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Tool{}, common.ErrResourceNotFound{Id: slug}
		}

		return domain.Tool{}, err
	}

	updateInput := domain.UpdateToolInput{
		Name:        tool.Name,
		Slug:        slug,
		Category:    tool.Category,
		SubType:     tool.SubType,
		Prolang:     tool.Prolang,
		ReleaseYear: tool.ReleaseYear,
		DevStatus:   tool.DevStatus,
		Details:     tool.Details,
		UseCases:    tool.UseCases,
		Tags:        tool.Tags,
		Website:     tool.Website,
		Github:      tool.Github,
	}

	if input.Name.IsSet {
		updateInput.Name = input.Name.Value
	}
	if input.Category.IsSet {
		updateInput.Category = input.Category.Value
	}
	if input.SubType.IsSet {
		updateInput.SubType = input.SubType.Value
	}
	if input.Prolang.IsSet {
		updateInput.Prolang = input.Prolang.Value
	}
	if input.DevStatus.IsSet {
		updateInput.DevStatus = input.DevStatus.Value
	}
	if input.Details.IsSet {
		updateInput.Details = input.Details.Value
	}
	if input.Website.IsSet {
		updateInput.Website = input.Website.Value
	}
	if input.Github.IsSet {
		updateInput.Github = input.Github.Value
	}
	if input.UseCases.IsSet {
		updateInput.UseCases = input.UseCases.Value
	}
	if input.Tags.IsSet {
		updateInput.Tags = input.Tags.Value
	}
	if input.ReleaseYear.IsSet {
		updateInput.ReleaseYear = input.ReleaseYear.Value
	}

	err = tool.Update(updateInput)

	if err != nil {
		return domain.Tool{}, err
	}

	updatedTool, err := uc.ToolRepository.UpdateTool(context.Background(), tool)

	if err != nil {
		return domain.Tool{}, err
	}

	return updatedTool, nil
}
