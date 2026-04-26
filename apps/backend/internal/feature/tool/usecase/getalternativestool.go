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
	Id          string                      `json:"id"`
	Name        string                      `json:"name"`
	Category    string                      `json:"category"`
	SubType     string                      `json:"sub_type"`
	Prolang     string                      `json:"prolang"`
	ReleaseYear int                         `json:"release_year"`
	DevStatus   string                      `json:"dev_status"`
	Details     string                      `json:"details"`
	UseCases    []string                    `json:"use_cases"`
	Tags        []string                    `json:"tags"`
	Website     string                      `json:"website"`
	Github      string                      `json:"github"`
	Metadata    domain.RelationshipMetadata `json:"metadata"`
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

	var uniqueToolIds []int = domain.DedupeToolIdsFromRelationships(relationships)

	if len(uniqueToolIds) == 0 {
		return []ToolAlternativesResult{}, nil
	}

	tools, err := uc.ToolRepository.GetToolByIds(context.Background(), uniqueToolIds)

	if err != nil {
		return []ToolAlternativesResult{}, err
	}

	var toolMap map[int]domain.Tool = map[int]domain.Tool{}
	for _, t := range tools {
		toolMap[t.Id] = t
	}

	var result []ToolAlternativesResult = make([]ToolAlternativesResult, 0)

	var processedToolIds map[int]bool = make(map[int]bool)

	for _, r := range relationships {
		var otherToolId int = r.FromToolId
		if tool.Id == r.FromToolId {
			otherToolId = r.ToToolId
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
			DevStatus:   otherTool.Devstatus,
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
