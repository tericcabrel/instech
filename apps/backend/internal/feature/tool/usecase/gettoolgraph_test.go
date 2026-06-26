package usecase_test

import (
	"database/sql"
	"errors"
	"testing"
	"time"

	"tericcabrel/instech/internal/common"
	"tericcabrel/instech/internal/domain"
	"tericcabrel/instech/internal/feature/tool/usecase"
	"tericcabrel/instech/testutil"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGetToolGraphUseCase(t *testing.T) {
	t.Run("returns not found when focus tool does not exist", func(t *testing.T) {
		toolRepository := testutil.NewMockToolRepositoryInterface(t)
		relationshipRepository := testutil.NewMockRelationshipRepositoryInterface(t)
		toolRepository.EXPECT().
			GetToolBySlug(mock.Anything, "unknown").
			Return(domain.Tool{}, sql.ErrNoRows)

		getToolGraph := usecase.GetToolGraphUseCase{
			ToolRepository:         toolRepository,
			RelationshipRepository: relationshipRepository,
		}

		result, err := getToolGraph.Execute("unknown", usecase.GetToolGraphInput{
			Depth:      1,
			LayoutMode: usecase.GRAPH_LAYOUT_MODE_CHRONOLOGICAL,
		})

		require.Equal(t, usecase.ToolGraphResult{}, result)
		require.Error(t, err)
		var errResourceNotFound common.ErrResourceNotFound
		require.ErrorAs(t, err, &errResourceNotFound)
	})

	t.Run("builds graph with depth and kind filters", func(t *testing.T) {
		toolRepository := testutil.NewMockToolRepositoryInterface(t)
		relationshipRepository := testutil.NewMockRelationshipRepositoryInterface(t)

		golang := domain.Tool{Id: 1, Slug: "golang", Name: "Golang", Category: "language", SubType: "backend", DevStatus: "active", ReleaseYear: 2009}
		nodejs := domain.Tool{Id: 2, Slug: "nodejs", Name: "Node.js", Category: "language", SubType: "fullstack", DevStatus: "active", ReleaseYear: 2009}
		rust := domain.Tool{Id: 4, Slug: "rust", Name: "Rust", Category: "language", SubType: "backend", DevStatus: "active", ReleaseYear: 2010}

		toolRepository.EXPECT().
			GetToolBySlug(mock.Anything, "golang").
			Return(golang, nil)

		relationshipRepository.EXPECT().
			GetRelationshipsByToolID(mock.Anything, 1).
			Return([]domain.Relationship{
				{ID: 11, FromToolID: 1, ToToolID: 2, Kind: "alternative_to", Metadata: domain.RelationshipMetadata{Reason: "A"}},
				{ID: 12, FromToolID: 1, ToToolID: 3, Kind: "used_with", Metadata: domain.RelationshipMetadata{Reason: "B"}},
			}, nil)

		relationshipRepository.EXPECT().
			GetRelationshipsByToolID(mock.Anything, 2).
			Return([]domain.Relationship{
				{ID: 13, FromToolID: 2, ToToolID: 4, Kind: "alternative_to", Metadata: domain.RelationshipMetadata{Reason: "C"}},
			}, nil)

		toolRepository.EXPECT().
			GetToolByIDs(mock.Anything, mock.MatchedBy(func(ids []int) bool {
				seen := map[int]bool{}
				for _, id := range ids {
					seen[id] = true
				}

				return seen[1] && seen[2] && seen[4] && len(seen) == 3
			})).
			Return([]domain.Tool{golang, nodejs, rust}, nil)

		getToolGraph := usecase.GetToolGraphUseCase{
			ToolRepository:         toolRepository,
			RelationshipRepository: relationshipRepository,
		}

		result, err := getToolGraph.Execute("golang", usecase.GetToolGraphInput{
			Depth:      2,
			Kinds:      []string{"alternative_to"},
			LayoutMode: usecase.GRAPH_LAYOUT_MODE_CHRONOLOGICAL,
		})

		require.NoError(t, err)
		require.Equal(t, "golang", result.FocusNodeID)
		require.Len(t, result.Nodes, 3)
		require.Len(t, result.Links, 2)
		require.Equal(t, "alternative_to", result.Links[0].Kind)
		require.Equal(t, "golang", result.Links[0].Source)
		require.Equal(t, "nodejs", result.Links[0].Target)
		require.Equal(t, "alternative_to", result.Links[1].Kind)
		require.Equal(t, "nodejs", result.Links[1].Source)
		require.Equal(t, "rust", result.Links[1].Target)
		require.Equal(t, usecase.GRAPH_LAYOUT_MODE_CHRONOLOGICAL, result.Meta.LayoutMode)
		require.Equal(t, []string{"alternative_to"}, result.Meta.KindsApplied)
		require.Equal(t, 2, result.Meta.Depth)
		require.Equal(t, 3, result.Meta.TotalNodes)
		require.Equal(t, 2, result.Meta.TotalLinks)
		require.True(t, result.Nodes[0].IsFocus)
	})

	t.Run("returns error when linked tool is missing from batch lookup", func(t *testing.T) {
		toolRepository := testutil.NewMockToolRepositoryInterface(t)
		relationshipRepository := testutil.NewMockRelationshipRepositoryInterface(t)
		toolRepository.EXPECT().
			GetToolBySlug(mock.Anything, "golang").
			Return(domain.Tool{Id: 1, Slug: "golang", Name: "Golang"}, nil)
		relationshipRepository.EXPECT().
			GetRelationshipsByToolID(mock.Anything, 1).
			Return([]domain.Relationship{
				{ID: 77, FromToolID: 1, ToToolID: 9, Kind: "used_with"},
			}, nil)
		toolRepository.EXPECT().
			GetToolByIDs(mock.Anything, mock.Anything).
			Return([]domain.Tool{{Id: 1, Slug: "golang", Name: "Golang"}}, nil)

		getToolGraph := usecase.GetToolGraphUseCase{
			ToolRepository:         toolRepository,
			RelationshipRepository: relationshipRepository,
		}

		result, err := getToolGraph.Execute("golang", usecase.GetToolGraphInput{
			Depth:      1,
			LayoutMode: usecase.GRAPH_LAYOUT_MODE_FORCE,
		})

		require.Equal(t, usecase.ToolGraphResult{}, result)
		require.Error(t, err)
		var errResourceNotFound common.ErrResourceNotFound
		require.ErrorAs(t, err, &errResourceNotFound)
		require.Equal(t, "9", errResourceNotFound.Id)
	})

	t.Run("propagates repository errors", func(t *testing.T) {
		toolRepository := testutil.NewMockToolRepositoryInterface(t)
		relationshipRepository := testutil.NewMockRelationshipRepositoryInterface(t)
		expectedErr := errors.New("database unavailable")
		toolRepository.EXPECT().
			GetToolBySlug(mock.Anything, "golang").
			Return(domain.Tool{Id: 1, Slug: "golang", Name: "Golang", CreatedAt: time.Now(), UpdatedAt: time.Now()}, nil)
		relationshipRepository.EXPECT().
			GetRelationshipsByToolID(mock.Anything, 1).
			Return([]domain.Relationship{}, expectedErr)

		getToolGraph := usecase.GetToolGraphUseCase{
			ToolRepository:         toolRepository,
			RelationshipRepository: relationshipRepository,
		}

		result, err := getToolGraph.Execute("golang", usecase.GetToolGraphInput{
			Depth:      1,
			LayoutMode: usecase.GRAPH_LAYOUT_MODE_CHRONOLOGICAL,
		})

		require.Equal(t, usecase.ToolGraphResult{}, result)
		require.ErrorIs(t, err, expectedErr)
	})
}

func TestIsLayoutModeValid(t *testing.T) {
	require.True(t, usecase.IsLayoutModeValid(usecase.GRAPH_LAYOUT_MODE_CHRONOLOGICAL))
	require.True(t, usecase.IsLayoutModeValid(usecase.GRAPH_LAYOUT_MODE_FORCE))
	require.False(t, usecase.IsLayoutModeValid("grid"))
}

func TestNormalizeKindsPublicBehavior(t *testing.T) {
	result := usecase.GetToolGraphInput{
		Kinds: []string{"alternative_to", "alternative_to", "used_with"},
	}

	toolRepository := testutil.NewMockToolRepositoryInterface(t)
	relationshipRepository := testutil.NewMockRelationshipRepositoryInterface(t)
	toolRepository.EXPECT().
		GetToolBySlug(mock.Anything, "golang").
		Return(domain.Tool{Id: 1, Slug: "golang", Name: "Golang"}, nil)
	relationshipRepository.EXPECT().
		GetRelationshipsByToolID(mock.Anything, 1).
		Return([]domain.Relationship{}, nil)
	toolRepository.EXPECT().
		GetToolByIDs(mock.Anything, []int{1}).
		Return([]domain.Tool{{Id: 1, Slug: "golang", Name: "Golang"}}, nil)

	getToolGraph := usecase.GetToolGraphUseCase{
		ToolRepository:         toolRepository,
		RelationshipRepository: relationshipRepository,
	}

	graph, err := getToolGraph.Execute("golang", usecase.GetToolGraphInput{
		Depth:      1,
		Kinds:      result.Kinds,
		LayoutMode: usecase.GRAPH_LAYOUT_MODE_CHRONOLOGICAL,
	})
	require.NoError(t, err)
	require.Equal(t, []string{"alternative_to", "used_with"}, graph.Meta.KindsApplied)
}
