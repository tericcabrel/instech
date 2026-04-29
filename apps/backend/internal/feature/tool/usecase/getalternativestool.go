package usecase

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"tericcabrel/instech/internal/common"
	"tericcabrel/instech/internal/domain"
	"tericcabrel/instech/internal/repository"
)

type GetToolAlternativesUseCase struct {
	ToolRepository         repository.ToolRepositoryInterface
	RelationshipRepository repository.RelationshipRepositoryInterface
}

type ToolAlternativesResult struct {
	DevStatus   string                      `json:"dev_status"`
	Name        string                      `json:"name"`
	Category    string                      `json:"category"`
	SubType     string                      `json:"sub_type"`
	Prolang     string                      `json:"prolang"`
	Id          string                      `json:"id"`
	Details     string                      `json:"details"`
	Website     string                      `json:"website"`
	Github      string                      `json:"github"`
	Metadata    domain.RelationshipMetadata `json:"metadata"`
	UseCases    []string                    `json:"use_cases"`
	Tags        []string                    `json:"tags"`
	ReleaseYear int                         `json:"release_year"`
}

func (uc *GetToolAlternativesUseCase) Execute(toolSlug string) ([]ToolAlternativesResult, error) {
	tool, err := uc.ToolRepository.GetToolBySlug(context.Background(), toolSlug)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []ToolAlternativesResult{}, common.ErrResourceNotFound{Id: toolSlug, Message: "The tool was not found"}
		}

		return []ToolAlternativesResult{}, err
	}

	relationships, err := uc.RelationshipRepository.GetToolAlternatives(context.Background(), tool.Id)
	if err != nil {
		return []ToolAlternativesResult{}, err
	}

	fmt.Printf("Relationships: %+v", relationships)

	var uniqueToolIds = domain.DedupeToolIdsFromRelationships(relationships)

	if len(uniqueToolIds) == 0 {
		return []ToolAlternativesResult{}, nil
	}

	tools, err := uc.ToolRepository.GetToolByIDs(context.Background(), uniqueToolIds)

	if err != nil {
		return []ToolAlternativesResult{}, err
	}

	var toolMap = map[int]domain.Tool{}
	for _, t := range tools {
		toolMap[t.Id] = t
	}

	var result = make([]ToolAlternativesResult, 0)

	var processedToolIds = make(map[int]bool)

	for _, r := range relationships {
		var otherToolId = r.FromToolID
		if tool.Id == r.FromToolID {
			otherToolId = r.ToToolID
		}

		otherTool, ok := toolMap[otherToolId]
		if !ok {
			return []ToolAlternativesResult{}, common.ErrResourceNotFound{Id: strconv.Itoa(otherToolId), Message: "The tool was not found"}
		}

		if processedToolIds[otherToolId] {
			continue
		}

		result = append(result, ToolAlternativesResult{
			Id:          otherTool.Slug,
			Name:        otherTool.Name,
			Category:    otherTool.Category,
			SubType:     otherTool.SubType,
			Prolang:     otherTool.Prolang,
			ReleaseYear: otherTool.ReleaseYear,
			DevStatus:   otherTool.DevStatus,
			Details:     otherTool.Details,
			UseCases:    otherTool.UseCases,
			Tags:        otherTool.Tags,
			Website:     otherTool.Website,
			Github:      otherTool.Github,
			Metadata:    r.Metadata,
		})

		processedToolIds[otherToolId] = true
	}

	return result, nil
}
