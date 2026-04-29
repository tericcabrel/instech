package usecase_test

import (
	"context"
	"tericcabrel/instech/internal/common"
	"tericcabrel/instech/internal/domain"
	"tericcabrel/instech/internal/feature/tool/usecase"
	"tericcabrel/instech/testutil"

	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestAddToolUseCase(t *testing.T) {
	t.Run("Add tool with valid input", func(t *testing.T) {
		toolRepository := testutil.NewMockToolRepositoryInterface(t)
		toolRepository.EXPECT().
			GetToolBySlug(mock.Anything, mock.AnythingOfType("string")).
			Return(domain.Tool{}, nil)
		toolRepository.EXPECT().
			CreateTool(mock.Anything, mock.AnythingOfType("domain.Tool")).
			RunAndReturn(func(_ context.Context, tool domain.Tool) (domain.Tool, error) {
				out := tool
				out.Id = 99

				return out, nil
			})

		addTool := usecase.AddToolUseCase{
			ToolRepository: toolRepository,
		}

		input := usecase.AddToolInput{
			Name:        "Node.js",
			Slug:        "nodejs",
			Category:    "language",
			SubType:     "backend",
			Prolang:     "JavaScript",
			ReleaseYear: 2009,
			DevStatus:   "active",
			Details:     "JavaScript runtime built on Chrome's V8",
			UseCases:    []string{"Backend", "Frontend", "Fullstack"},
			Tags:        []string{"JavaScript", "Node.js", "Backend", "Frontend", "Fullstack"},
			Website:     "https://nodejs.org",
			Github:      "https://github.com/nodejs/node",
		}

		tool, err := addTool.Execute(input)

		require.NoError(t, err)
		require.Equal(t, 99, tool.Id)
	})

	t.Run("Add tool with invalid category will fail", func(t *testing.T) {
		toolRepository := testutil.NewMockToolRepositoryInterface(t)
		addTool := usecase.AddToolUseCase{
			ToolRepository: toolRepository,
		}
		input := usecase.AddToolInput{}

		tool, err := addTool.Execute(input)
		require.Error(t, err)

		if _, ok := err.(domain.ErrInvalidToolCategory); !ok {
			t.Errorf("Expected ErrInvalidToolCategory, got %v", err)
		}
		require.Equal(t, domain.Tool{}, tool)

		toolRepository.AssertNotCalled(t, "CreateTool", mock.Anything, mock.AnythingOfType("domain.Tool"))
	})

	t.Run("Add tool with invalid fields will fail", func(t *testing.T) {
		toolRepository := testutil.NewMockToolRepositoryInterface(t)
		addTool := usecase.AddToolUseCase{
			ToolRepository: toolRepository,
		}
		input := usecase.AddToolInput{
			Name:        "",
			Slug:        "",
			Category:    "language",
			SubType:     "backend",
			Prolang:     "",
			ReleaseYear: 1922,
			DevStatus:   "active",
			Details:     "8",
			UseCases:    []string{""},
			Tags:        []string{""},
			Website:     "invalid",
			Github:      "invalid",
		}

		tool, err := addTool.Execute(input)
		require.Error(t, err)

		if _, ok := err.(domain.ErrInvalidField); !ok {
			t.Errorf("Expected ErrInvalidField, got %v", err)
		}
		if e, ok := err.(domain.ErrInvalidField); ok {
			fields := []string{"Name", "Slug", "ReleaseYear", "Prolang", "Website", "Github"}
			for _, field := range fields {
				if _, exist := e.Fields[field]; !exist {
					t.Errorf("Expected the field \"%s\" to be present", field)
				}
			}
		}
		require.Equal(t, domain.Tool{}, tool)

		toolRepository.AssertNotCalled(t, "CreateTool", mock.Anything, mock.AnythingOfType("domain.Tool"))
	})

	t.Run("Fail to add tool when it already exists", func(t *testing.T) {
		toolRepository := testutil.NewMockToolRepositoryInterface(t)
		toolRepository.EXPECT().
			GetToolBySlug(mock.Anything, mock.AnythingOfType("string")).
			Return(domain.Tool{Id: 1}, nil)

		addTool := usecase.AddToolUseCase{
			ToolRepository: toolRepository,
		}

		input := usecase.AddToolInput{
			Name:        "Node.js",
			Slug:        "nodejs",
			Category:    "language",
			SubType:     "backend",
			Prolang:     "JavaScript",
			ReleaseYear: 2009,
			DevStatus:   "active",
			Details:     "JavaScript runtime built on Chrome's V8",
			UseCases:    []string{"Backend", "Frontend", "Fullstack"},
			Tags:        []string{"JavaScript", "Node.js", "Backend", "Frontend", "Fullstack"},
			Website:     "https://nodejs.org",
			Github:      "https://github.com/nodejs/node",
		}

		tool, err := addTool.Execute(input)

		require.Error(t, err)
		if _, ok := err.(common.ErrResourceAlreadyExists); !ok {
			t.Errorf("Expected ErrResourceAlreadyExists, got %v", err)
		}
		require.Equal(t, domain.Tool{}, tool)
		toolRepository.AssertNotCalled(t, "CreateTool", mock.Anything, mock.AnythingOfType("domain.Tool"))
	})
}
