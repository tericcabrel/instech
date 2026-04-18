package usecase

import (
	"context"
	"slices"
	"tericcabrel/instech/internal/domain"
	"tericcabrel/instech/internal/repository"
	"time"
)

type ClientRelationshipDataTool struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

type ClientRelationshipPaginationMetadata struct {
	TotalCount int64 `json:"total_count"`
	NextCursor int64 `json:"next_cursor"`
	ItemsCount int64 `json:"items_count"`
}

type ClientRelationshipData struct {
	Id        int                         `json:"id"`
	FromTool  ClientRelationshipDataTool  `json:"from_tool"`
	ToTool    ClientRelationshipDataTool  `json:"to_tool"`
	Kind      string                      `json:"kind"`
	Metadata  domain.RelationshipMetadata `json:"metadata"`
	CreatedAt time.Time                   `json:"created_at"`
	UpdatedAt time.Time                   `json:"updated_at"`
}

type ClientRelationshipResult struct {
	Data []ClientRelationshipData             `json:"data"`
	Meta ClientRelationshipPaginationMetadata `json:"meta"`
}

type GetRelationshipsUseCaseParams struct {
	Cursor int64  `json:"cursor"`
	ToolId int    `json:"tool_id"`
	Kind   string `json:"kind"`
	Limit  int    `json:"limit"`
}

func GetRelationshipsUsecase(relationshipRepository repository.RelationshipRepositoryInterface, toolRepository repository.ToolRepositoryInterface, params GetRelationshipsUseCaseParams) (ClientRelationshipResult, error) {
	paginatedRelationships, err := relationshipRepository.GetRelationshipsAll(context.Background(), repository.GetRelationshipsAllParams{
		Cursor: params.Cursor,
		ToolId: params.ToolId,
		Kind:   params.Kind,
		Limit:  params.Limit,
	})
	if err != nil {
		return ClientRelationshipResult{}, err
	}

	var uniqueToolIds []int = make([]int, 0)
	for _, relationship := range paginatedRelationships.Relationships {
		if !slices.Contains(uniqueToolIds, relationship.FromToolId) {
			uniqueToolIds = append(uniqueToolIds, relationship.FromToolId)
		}
		if !slices.Contains(uniqueToolIds, relationship.ToToolId) {
			uniqueToolIds = append(uniqueToolIds, relationship.ToToolId)
		}
	}

	tools, err := toolRepository.GetToolByIds(context.Background(), uniqueToolIds)
	if err != nil {
		return ClientRelationshipResult{}, err
	}

	var toolMap map[int]domain.Tool = map[int]domain.Tool{}
	for _, tool := range tools {
		toolMap[tool.Id] = tool
	}

	var relationships []ClientRelationshipData = make([]ClientRelationshipData, 0)
	for _, relationship := range paginatedRelationships.Relationships {
		relationships = append(relationships, ClientRelationshipData{
			Id: relationship.Id,
			FromTool: ClientRelationshipDataTool{
				Id:   relationship.FromToolId,
				Name: toolMap[relationship.FromToolId].Name,
				Slug: toolMap[relationship.FromToolId].Slug,
			},
			ToTool: ClientRelationshipDataTool{
				Id:   relationship.ToToolId,
				Name: toolMap[relationship.ToToolId].Name,
				Slug: toolMap[relationship.ToToolId].Slug,
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
