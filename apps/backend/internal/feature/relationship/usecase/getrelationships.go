package usecase

import (
	"context"
	"strconv"
	"tericcabrel/instech/internal/common"
	"tericcabrel/instech/internal/domain"
	"tericcabrel/instech/internal/repository"
	"time"
)

type GetRelationshipsUseCase struct {
	RelationshipRepository repository.RelationshipRepositoryInterface
	ToolRepository         repository.ToolRepositoryInterface
}
type ClientRelationshipDataTool struct {
	Name string `json:"name"`
	Slug string `json:"slug"`
	ID   int    `json:"id"`
}

type ClientRelationshipPaginationMetadata struct {
	TotalCount int64 `json:"total_count"`
	NextCursor int64 `json:"next_cursor"`
	ItemsCount int64 `json:"items_count"`
}

type ClientRelationshipData struct {
	CreatedAt time.Time                   `json:"created_at"`
	UpdatedAt time.Time                   `json:"updated_at"`
	Kind      string                      `json:"kind"`
	Metadata  domain.RelationshipMetadata `json:"metadata"`
	FromTool  ClientRelationshipDataTool  `json:"from_tool"`
	ToTool    ClientRelationshipDataTool  `json:"to_tool"`
	ID        int                         `json:"id"`
}

type ClientRelationshipResult struct {
	Data []ClientRelationshipData             `json:"data"`
	Meta ClientRelationshipPaginationMetadata `json:"meta"`
}

type GetRelationshipsUseCaseParams struct {
	Kind   string `json:"kind"`
	Cursor int64  `json:"cursor"`
	ToolId int    `json:"tool_id"`
	Limit  int    `json:"limit"`
}

func (uc *GetRelationshipsUseCase) Execute(params GetRelationshipsUseCaseParams) (ClientRelationshipResult, error) {
	paginatedRelationships, err := uc.RelationshipRepository.GetRelationshipsAll(context.Background(), repository.GetRelationshipsAllParams{
		Cursor: params.Cursor,
		ToolId: params.ToolId,
		Kind:   params.Kind,
		Limit:  params.Limit,
	})
	if err != nil {
		return ClientRelationshipResult{}, err
	}

	var uniqueToolIds = domain.DedupeToolIdsFromRelationships(paginatedRelationships.Relationships)

	tools, err := uc.ToolRepository.GetToolByIDs(context.Background(), uniqueToolIds)
	if err != nil {
		return ClientRelationshipResult{}, err
	}

	var toolMap = map[int]domain.Tool{}
	for _, tool := range tools {
		toolMap[tool.Id] = tool
	}

	var relationships = make([]ClientRelationshipData, 0)
	for _, relationship := range paginatedRelationships.Relationships {
		fromTool, ok := toolMap[relationship.FromToolID]
		if !ok {
			return ClientRelationshipResult{}, common.ErrResourceNotFound{Id: strconv.Itoa(relationship.FromToolID), Message: "The from tool was not found"}
		}
		toTool, ok := toolMap[relationship.ToToolID]
		if !ok {
			return ClientRelationshipResult{}, common.ErrResourceNotFound{Id: strconv.Itoa(relationship.ToToolID), Message: "The to tool was not found"}
		}
		relationships = append(relationships, ClientRelationshipData{
			ID: relationship.ID,
			FromTool: ClientRelationshipDataTool{
				ID:   fromTool.Id,
				Name: fromTool.Name,
				Slug: fromTool.Slug,
			},
			ToTool: ClientRelationshipDataTool{
				ID:   toTool.Id,
				Name: toTool.Name,
				Slug: toTool.Slug,
			},
			Kind:      relationship.Kind,
			Metadata:  relationship.Metadata,
			CreatedAt: relationship.CreatedAt,
			UpdatedAt: relationship.UpdatedAt,
		})
	}

	result := ClientRelationshipResult{
		Data: relationships,
		Meta: ClientRelationshipPaginationMetadata{
			ItemsCount: paginatedRelationships.ItemsCount,
			TotalCount: paginatedRelationships.TotalCount,
			NextCursor: paginatedRelationships.NextCursor,
		},
	}

	return result, nil
}
