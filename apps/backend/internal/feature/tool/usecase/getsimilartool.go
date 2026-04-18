package usecase

import (
	"context"
	"database/sql"
	"errors"
	"slices"
	"tericcabrel/instech/internal/common"
	"tericcabrel/instech/internal/domain"
	"tericcabrel/instech/internal/repository"
)

type SimilarToolResult struct {
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

func GetSimilarToolUsecase(toolRepository repository.ToolRepositoryInterface, relationshipRepository repository.RelationshipRepositoryInterface, toolSlug string) ([]SimilarToolResult, error) {
	tool, err := toolRepository.GetToolBySlug(context.Background(), toolSlug)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []SimilarToolResult{}, common.ErrResourceNotFound{Id: toolSlug, Message: err.Error()}
		}
		return []SimilarToolResult{}, err
	}

	relationships, err := relationshipRepository.GetToolSimilar(context.Background(), tool.Id)
	if err != nil {
		return []SimilarToolResult{}, err
	}

	var uniqueToolIds []int = make([]int, 0)
	for _, relationship := range relationships {
		if !slices.Contains(uniqueToolIds, relationship.FromToolId) {
			uniqueToolIds = append(uniqueToolIds, relationship.FromToolId)
		}
		if !slices.Contains(uniqueToolIds, relationship.ToToolId) {
			uniqueToolIds = append(uniqueToolIds, relationship.ToToolId)
		}
	}

	if len(uniqueToolIds) == 0 {
		return []SimilarToolResult{}, nil
	}

	tools, err := toolRepository.GetToolByIds(context.Background(), uniqueToolIds)
	if err != nil {
		return []SimilarToolResult{}, err
	}

	var toolMap map[int]domain.Tool = map[int]domain.Tool{}
	for _, t := range tools {
		toolMap[t.Id] = t
	}

	var result []SimilarToolResult = make([]SimilarToolResult, 0)

	for _, r := range relationships {
		var otherToolId int = r.FromToolId
		if tool.Id == r.FromToolId {
			otherToolId = r.ToToolId
		}

		result = append(result, SimilarToolResult{
			Id:          toolMap[otherToolId].Slug,
			Name:        toolMap[otherToolId].Name,
			Category:    toolMap[otherToolId].Category,
			SubType:     toolMap[otherToolId].SubType,
			Prolang:     toolMap[otherToolId].Prolang,
			ReleaseYear: toolMap[otherToolId].ReleaseYear,
			DevStatus:   toolMap[otherToolId].Devstatus,
			Details:     toolMap[otherToolId].Details,
			UseCases:    toolMap[otherToolId].UseCases,
			Tags:        toolMap[otherToolId].Tags,
			Website:     toolMap[otherToolId].Website,
			Github:      toolMap[otherToolId].Github,
			Metadata:    r.Metadata,
		})
	}

	return result, nil
}
