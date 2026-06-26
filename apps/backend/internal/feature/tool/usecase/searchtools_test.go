package usecase_test

import (
	"context"
	"errors"
	"tericcabrel/instech/internal/feature/tool/usecase"
	"tericcabrel/instech/internal/repository"
	"tericcabrel/instech/testutil"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestSearchToolsUseCase(t *testing.T) {
	t.Run("empty keyword returns empty slice", func(t *testing.T) {
		toolRepository := testutil.NewMockToolRepositoryInterface(t)
		searchTools := usecase.SearchToolsUseCase{
			ToolRepository: toolRepository,
		}

		testKeywords := []string{"", "   ", "\t\n"}
		for _, keyword := range testKeywords {
			results, err := searchTools.Execute(keyword)
			require.NoError(t, err)
			require.Equal(t, []usecase.SearchToolsResult{}, results)
		}

		toolRepository.AssertNotCalled(t, "SearchTools", mock.Anything, mock.Anything)
	})

	t.Run("keyword is trimmed before repository call", func(t *testing.T) {
		toolRepository := testutil.NewMockToolRepositoryInterface(t)
		toolRepository.EXPECT().
			SearchTools(mock.Anything, mock.AnythingOfType("string")).
			Run(func(_ context.Context, keyword string) {
				require.Equal(t, "node", keyword)
			}).
			Return([]repository.ToolSearchResult{}, nil)

		searchTools := usecase.SearchToolsUseCase{
			ToolRepository: toolRepository,
		}

		results, err := searchTools.Execute("  node  ")
		require.NoError(t, err)
		require.Equal(t, []usecase.SearchToolsResult{}, results)
	})

	t.Run("repository error is propagated", func(t *testing.T) {
		toolRepository := testutil.NewMockToolRepositoryInterface(t)
		dbErr := errors.New("database unavailable")
		toolRepository.EXPECT().
			SearchTools(mock.Anything, "go").
			Return([]repository.ToolSearchResult{}, dbErr)

		searchTools := usecase.SearchToolsUseCase{
			ToolRepository: toolRepository,
		}

		results, err := searchTools.Execute("go")
		require.ErrorIs(t, err, dbErr)
		require.Equal(t, []usecase.SearchToolsResult{}, results)
	})

	t.Run("maps repository dto to use case dto", func(t *testing.T) {
		toolRepository := testutil.NewMockToolRepositoryInterface(t)
		toolRepository.EXPECT().
			SearchTools(mock.Anything, "node").
			Return([]repository.ToolSearchResult{
				{
					Id:       1,
					Slug:     "nodejs",
					Name:     "Node.js",
					Category: "language",
					SubType:  "fullstack",
				},
				{
					Id:       2,
					Slug:     "express",
					Name:     "Express.js",
					Category: "framework",
					SubType:  "backend",
				},
			}, nil)

		searchTools := usecase.SearchToolsUseCase{
			ToolRepository: toolRepository,
		}

		results, err := searchTools.Execute("node")
		require.NoError(t, err)
		require.Equal(t, []usecase.SearchToolsResult{
			{
				Id:       1,
				Slug:     "nodejs",
				Name:     "Node.js",
				Category: "language",
				SubType:  "fullstack",
			},
			{
				Id:       2,
				Slug:     "express",
				Name:     "Express.js",
				Category: "framework",
				SubType:  "backend",
			},
		}, results)
	})
}
