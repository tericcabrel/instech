package usecase_test

import (
	"errors"
	"tericcabrel/instech/internal/common"
	"tericcabrel/instech/internal/domain"
	"tericcabrel/instech/internal/feature/relationship/usecase"
	"tericcabrel/instech/internal/repository"
	"tericcabrel/instech/testutil"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGetRelationshipsUseCase(t *testing.T) {
	params := usecase.GetRelationshipsUseCaseParams{
		Cursor: 10,
		ToolId: 3,
		Kind:   "built_on",
		Limit:  25,
	}
	repoParams := repository.GetRelationshipsAllParams{
		Cursor: params.Cursor,
		ToolId: params.ToolId,
		Kind:   params.Kind,
		Limit:  params.Limit,
	}

	t.Run("Get relationships propagates error from GetRelationshipsAll", func(t *testing.T) {
		relationshipRepository := testutil.NewMockRelationshipRepositoryInterface(t)
		toolRepository := testutil.NewMockToolRepositoryInterface(t)

		dbErr := errors.New("database unavailable")
		relationshipRepository.EXPECT().
			GetRelationshipsAll(mock.Anything, repoParams).
			Return(repository.PaginatedRelationshipsResult{}, dbErr)

		getRelationships := usecase.GetRelationshipsUseCase{
			RelationshipRepository: relationshipRepository,
			ToolRepository:         toolRepository,
		}
		out, err := getRelationships.Execute(params)

		require.Equal(t, usecase.ClientRelationshipResult{}, out)
		require.ErrorIs(t, err, dbErr)
		toolRepository.AssertNotCalled(t, "GetToolByIds", mock.Anything, mock.Anything)
	})

	t.Run("Get relationships propagates error from GetToolByIds", func(t *testing.T) {
		rel := testutil.CreateTestDynamicRelationship(1, domain.CreateRelationshipInput{
			FromToolID: 1,
			ToToolID:   2,
			Kind:       "built_on",
			Reason:     "r",
		})

		relationshipRepository := testutil.NewMockRelationshipRepositoryInterface(t)
		toolRepository := testutil.NewMockToolRepositoryInterface(t)

		toolsErr := errors.New("tools query failed")
		relationshipRepository.EXPECT().
			GetRelationshipsAll(mock.Anything, repoParams).
			Return(repository.PaginatedRelationshipsResult{
				Relationships: []domain.Relationship{rel},
				TotalCount:    1,
				ItemsCount:    1,
				NextCursor:    0,
			}, nil)
		toolRepository.EXPECT().
			GetToolByIDs(mock.Anything, mock.MatchedBy(func(ids []int) bool {
				return len(ids) == 2
			})).
			Return(nil, toolsErr)

		getRelationships := usecase.GetRelationshipsUseCase{
			RelationshipRepository: relationshipRepository,
			ToolRepository:         toolRepository,
		}
		out, err := getRelationships.Execute(params)

		require.Equal(t, usecase.ClientRelationshipResult{}, out)
		require.ErrorIs(t, err, toolsErr)
	})

	t.Run("Get relationships returns error when missing from tool in tool map", func(t *testing.T) {
		rel := testutil.CreateTestDynamicRelationship(1, domain.CreateRelationshipInput{
			FromToolID: 1,
			ToToolID:   2,
			Kind:       "built_on",
			Reason:     "r",
		})
		tool2 := testutil.CreateTestDynamicTool(2, domain.CreateToolInput{
			Name:        "Only To",
			Slug:        "only-to",
			Category:    "language",
			SubType:     "backend",
			Devstatus:   "active",
			Details:     "d",
			UseCases:    []string{"u"},
			Tags:        []string{"t"},
			Website:     "https://example.com",
			Github:      "https://github.com/example",
			ReleaseYear: 2020,
			Prolang:     "Go",
		})

		relationshipRepository := testutil.NewMockRelationshipRepositoryInterface(t)
		toolRepository := testutil.NewMockToolRepositoryInterface(t)

		relationshipRepository.EXPECT().
			GetRelationshipsAll(mock.Anything, repoParams).
			Return(repository.PaginatedRelationshipsResult{
				Relationships: []domain.Relationship{rel},
				TotalCount:    1,
				ItemsCount:    1,
				NextCursor:    0,
			}, nil)
		toolRepository.EXPECT().
			GetToolByIDs(mock.Anything, mock.Anything).
			Return([]domain.Tool{tool2}, nil)

		getRelationships := usecase.GetRelationshipsUseCase{
			RelationshipRepository: relationshipRepository,
			ToolRepository:         toolRepository,
		}
		out, err := getRelationships.Execute(params)

		require.Equal(t, usecase.ClientRelationshipResult{}, out)
		var notFound common.ErrResourceNotFound
		require.ErrorAs(t, err, &notFound)
		require.Equal(t, "1", notFound.Id)
	})

	t.Run("Get relationships returns empty result", func(t *testing.T) {
		relationshipRepository := testutil.NewMockRelationshipRepositoryInterface(t)
		toolRepository := testutil.NewMockToolRepositoryInterface(t)

		relationshipRepository.EXPECT().
			GetRelationshipsAll(mock.Anything, repoParams).
			Return(repository.PaginatedRelationshipsResult{
				Relationships: nil,
				TotalCount:    0,
				ItemsCount:    0,
				NextCursor:    0,
			}, nil)
		toolRepository.EXPECT().
			GetToolByIDs(mock.Anything, mock.MatchedBy(func(ids []int) bool {
				return len(ids) == 0
			})).
			Return([]domain.Tool{}, nil)

		uc := usecase.GetRelationshipsUseCase{
			RelationshipRepository: relationshipRepository,
			ToolRepository:         toolRepository,
		}
		out, err := uc.Execute(params)

		require.NoError(t, err)
		require.Empty(t, out.Data)
		require.Equal(t, int64(0), out.Meta.TotalCount)
		require.Equal(t, int64(0), out.Meta.ItemsCount)
		require.Equal(t, int64(0), out.Meta.NextCursor)
	})

	t.Run("Get relationships with many relationships and tools", func(t *testing.T) {
		created := time.Date(2024, 1, 2, 3, 4, 5, 0, time.UTC)
		updated := time.Date(2024, 2, 3, 4, 5, 6, 0, time.UTC)

		rel := testutil.CreateTestDynamicRelationship(7, domain.CreateRelationshipInput{
			FromToolID: 1,
			ToToolID:   2,
			Kind:       "inspired_by",
			Reason:     "because",
		})
		rel.CreatedAt = created
		rel.UpdatedAt = updated

		tool1 := testutil.CreateTestTool()
		tool2 := testutil.CreateTestDynamicTool(2, domain.CreateToolInput{
			Name:        "Node.js",
			Slug:        "nodejs",
			Category:    "language",
			SubType:     "backend",
			Devstatus:   "active",
			Details:     "Test Details",
			UseCases:    []string{"Test Use Case"},
			Tags:        []string{"Test Tag"},
			Website:     "https://nodejs.org",
			Github:      "https://github.com/nodejs/node",
			ReleaseYear: 2009,
			Prolang:     "JavaScript",
		})

		relationshipRepository := testutil.NewMockRelationshipRepositoryInterface(t)
		toolRepository := testutil.NewMockToolRepositoryInterface(t)

		relationshipRepository.EXPECT().
			GetRelationshipsAll(mock.Anything, repoParams).
			Return(repository.PaginatedRelationshipsResult{
				Relationships: []domain.Relationship{rel},
				TotalCount:    100,
				ItemsCount:    1,
				NextCursor:    42,
			}, nil)
		toolRepository.EXPECT().
			GetToolByIDs(mock.Anything, mock.MatchedBy(func(ids []int) bool {
				return len(ids) == 2
			})).
			Return([]domain.Tool{tool1, tool2}, nil)

		getRelationships := usecase.GetRelationshipsUseCase{
			RelationshipRepository: relationshipRepository,
			ToolRepository:         toolRepository,
		}
		out, err := getRelationships.Execute(params)

		require.NoError(t, err)
		require.Len(t, out.Data, 1)
		require.Equal(t, 7, out.Data[0].ID)
		require.Equal(t, "inspired_by", out.Data[0].Kind)
		require.Equal(t, domain.RelationshipMetadata{Reason: "because"}, out.Data[0].Metadata)
		require.Equal(t, created, out.Data[0].CreatedAt)
		require.Equal(t, updated, out.Data[0].UpdatedAt)
		require.Equal(t, tool1.Id, out.Data[0].FromTool.ID)
		require.Equal(t, tool1.Name, out.Data[0].FromTool.Name)
		require.Equal(t, tool1.Slug, out.Data[0].FromTool.Slug)
		require.Equal(t, tool2.Id, out.Data[0].ToTool.ID)
		require.Equal(t, tool2.Name, out.Data[0].ToTool.Name)
		require.Equal(t, tool2.Slug, out.Data[0].ToTool.Slug)
		require.Equal(t, int64(100), out.Meta.TotalCount)
		require.Equal(t, int64(1), out.Meta.ItemsCount)
		require.Equal(t, int64(42), out.Meta.NextCursor)
	})
}
