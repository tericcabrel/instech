package usecase

import (
	"context"
	"database/sql"
	"errors"
	"sort"
	"strconv"

	"tericcabrel/instech/internal/common"
	"tericcabrel/instech/internal/domain"
	"tericcabrel/instech/internal/repository"
)

const (
	MAX_GRAPH_DEPTH                 = 2
	MAX_GRAPH_NODES                 = 250
	MAX_GRAPH_LINKS                 = 500
	GRAPH_LAYOUT_MODE_CHRONOLOGICAL = "chronological"
	GRAPH_LAYOUT_MODE_FORCE         = "force"
)

type ToolGraphNode struct {
	ID          string `json:"id"`
	Slug        string `json:"slug"`
	Name        string `json:"name"`
	Category    string `json:"category"`
	SubType     string `json:"subType"`
	DevStatus   string `json:"devStatus"`
	Degree      int    `json:"degree"`
	IsFocus     bool   `json:"isFocus"`
	ReleaseYear int    `json:"releaseYear"`
}

type ToolGraphLink struct {
	Source string `json:"source"`
	Target string `json:"target"`
	Kind   string `json:"kind"`
	Reason string `json:"reason,omitempty"`
	ID     int    `json:"id"`
}

type ToolGraphMeta struct {
	LayoutMode   string   `json:"layoutMode"`
	KindsApplied []string `json:"kindsApplied"`
	Depth        int      `json:"depth"`
	TotalNodes   int      `json:"totalNodes"`
	TotalLinks   int      `json:"totalLinks"`
}

type ToolGraphResult struct {
	FocusNodeID string          `json:"focusNodeId"`
	Nodes       []ToolGraphNode `json:"nodes"`
	Links       []ToolGraphLink `json:"links"`
	Meta        ToolGraphMeta   `json:"meta"`
}

type GetToolGraphInput struct {
	LayoutMode string
	Kinds      []string
	Depth      int
}

type GetToolGraphUseCase struct {
	ToolRepository         repository.ToolRepositoryInterface
	RelationshipRepository repository.RelationshipRepositoryInterface
}

func IsLayoutModeValid(mode string) bool {
	return mode == GRAPH_LAYOUT_MODE_CHRONOLOGICAL || mode == GRAPH_LAYOUT_MODE_FORCE
}

func deduplicateKinds(kinds []string) []string {
	if len(kinds) == 0 {
		return append([]string{}, domain.RelationshipKinds...)
	}

	seen := make(map[string]bool)
	result := make([]string, 0, len(kinds))
	for _, kind := range kinds {
		if seen[kind] {
			continue
		}

		seen[kind] = true
		result = append(result, kind)
	}

	sort.Strings(result)

	return result
}

func (uc *GetToolGraphUseCase) Execute(toolSlug string, input GetToolGraphInput) (ToolGraphResult, error) {
	tool, err := uc.ToolRepository.GetToolBySlug(context.Background(), toolSlug)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ToolGraphResult{}, common.ErrResourceNotFound{Id: toolSlug, Message: "The tool was not found"}
		}

		return ToolGraphResult{}, err
	}

	appliedKinds := deduplicateKinds(input.Kinds)
	allowedKinds := make(map[string]bool, len(appliedKinds))
	for _, kind := range appliedKinds {
		allowedKinds[kind] = true
	}

	visited := map[int]bool{
		tool.Id: true,
	}
	frontier := []int{tool.Id}
	relationshipByID := make(map[int]domain.Relationship)

	for range input.Depth {
		if len(frontier) == 0 {
			break
		}

		nextFrontierSeen := make(map[int]bool)
		nextFrontier := make([]int, 0)

		for _, toolID := range frontier {
			relationships, relErr := uc.RelationshipRepository.GetRelationshipsByToolID(context.Background(), toolID)
			if relErr != nil {
				return ToolGraphResult{}, relErr
			}

			for _, relationship := range relationships {
				if !allowedKinds[relationship.Kind] {
					continue
				}

				relationshipByID[relationship.ID] = relationship
				if len(relationshipByID) >= MAX_GRAPH_LINKS {
					break
				}

				candidateIDs := []int{relationship.FromToolID, relationship.ToToolID}
				for _, candidateID := range candidateIDs {
					if visited[candidateID] || nextFrontierSeen[candidateID] {
						continue
					}

					nextFrontierSeen[candidateID] = true
					nextFrontier = append(nextFrontier, candidateID)
				}
			}
		}

		for _, nextID := range nextFrontier {
			visited[nextID] = true
			if len(visited) >= MAX_GRAPH_NODES {
				break
			}
		}

		if len(visited) >= MAX_GRAPH_NODES || len(relationshipByID) >= MAX_GRAPH_LINKS {
			break
		}

		frontier = nextFrontier
	}

	toolIDs := make([]int, 0, len(visited))
	for id := range visited {
		toolIDs = append(toolIDs, id)
	}

	tools, err := uc.ToolRepository.GetToolByIDs(context.Background(), toolIDs)
	if err != nil {
		return ToolGraphResult{}, err
	}

	toolMap := make(map[int]domain.Tool, len(tools))
	for _, item := range tools {
		toolMap[item.Id] = item
	}

	links := make([]ToolGraphLink, 0, len(relationshipByID))
	degreeBySlug := make(map[string]int)

	for _, relationship := range relationshipByID {
		sourceTool, sourceOk := toolMap[relationship.FromToolID]
		targetTool, targetOk := toolMap[relationship.ToToolID]
		if !sourceOk {
			return ToolGraphResult{}, common.ErrResourceNotFound{
				Id:      strconv.Itoa(relationship.FromToolID),
				Message: "The source tool was not found",
			}
		}
		if !targetOk {
			return ToolGraphResult{}, common.ErrResourceNotFound{
				Id:      strconv.Itoa(relationship.ToToolID),
				Message: "The target tool was not found",
			}
		}

		link := ToolGraphLink{
			ID:     relationship.ID,
			Source: sourceTool.Slug,
			Target: targetTool.Slug,
			Kind:   relationship.Kind,
			Reason: relationship.Metadata.Reason,
		}

		links = append(links, link)
		degreeBySlug[sourceTool.Slug]++
		degreeBySlug[targetTool.Slug]++
	}

	sort.Slice(links, func(i, j int) bool {
		left := links[i]
		right := links[j]

		if left.Kind != right.Kind {
			return left.Kind < right.Kind
		}
		if left.Source != right.Source {
			return left.Source < right.Source
		}
		if left.Target != right.Target {
			return left.Target < right.Target
		}

		return left.ID < right.ID
	})

	nodes := make([]ToolGraphNode, 0, len(toolMap))
	for _, item := range toolMap {
		nodes = append(nodes, ToolGraphNode{
			ID:          item.Slug,
			Slug:        item.Slug,
			Name:        item.Name,
			Category:    item.Category,
			SubType:     item.SubType,
			DevStatus:   item.DevStatus,
			ReleaseYear: item.ReleaseYear,
			Degree:      degreeBySlug[item.Slug],
			IsFocus:     item.Slug == tool.Slug,
		})
	}

	sort.Slice(nodes, func(i, j int) bool {
		left := nodes[i]
		right := nodes[j]

		if left.IsFocus != right.IsFocus {
			return left.IsFocus
		}
		if left.Degree != right.Degree {
			return left.Degree > right.Degree
		}

		return left.Name < right.Name
	})

	return ToolGraphResult{
		FocusNodeID: tool.Slug,
		Nodes:       nodes,
		Links:       links,
		Meta: ToolGraphMeta{
			Depth:        input.Depth,
			TotalNodes:   len(nodes),
			TotalLinks:   len(links),
			KindsApplied: appliedKinds,
			LayoutMode:   input.LayoutMode,
		},
	}, nil
}
