package usecase_test

import (
	"context"
	"database/sql"
	"errors"
	"tericcabrel/instech/internal/common"
	"tericcabrel/instech/internal/domain"
	"tericcabrel/instech/internal/feature/tool/usecase"
	"tericcabrel/instech/testutil"

	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestUpdateToolUseCase(t *testing.T) {
	t.Run("Update tool fails when the tool is not found", func(t *testing.T) {
		toolRepository := testutil.NewMockToolRepositoryInterface(t)
		toolRepository.EXPECT().
			GetToolBySlug(mock.Anything, mock.AnythingOfType("string")).
			Return(domain.Tool{}, sql.ErrNoRows)

		updateTool := usecase.UpdateToolUseCase{
			ToolRepository: toolRepository,
		}

		input := usecase.UpdateToolInput{
			Name:        "Node.js",
			Category:    "language",
			SubType:     "backend",
			Prolang:     "JavaScript",
			ReleaseYear: 2009,
			DevStatus:   "active",
			Details:     "JavaScript runtime built on Chrome's V8",
		}

		tool, err := updateTool.Execute("nodejs", input)
		require.Error(t, err)
		if _, ok := err.(common.ErrResourceNotFound); !ok {
			t.Errorf("Expected ErrResourceNotFound, got %v", err)
		}
		require.Equal(t, domain.Tool{}, tool)

		toolRepository.AssertCalled(t, "GetToolBySlug", mock.Anything, mock.AnythingOfType("string"))
		toolRepository.AssertNotCalled(t, "UpdateTool", mock.Anything, mock.AnythingOfType("domain.Tool"))
	})

	t.Run("Update tool fails when there is an error getting the tool", func(t *testing.T) {
		toolRepository := testutil.NewMockToolRepositoryInterface(t)
		toolRepository.EXPECT().
			GetToolBySlug(mock.Anything, mock.AnythingOfType("string")).
			Return(domain.Tool{}, errors.New("error getting the tool"))

		updateTool := usecase.UpdateToolUseCase{
			ToolRepository: toolRepository,
		}

		input := usecase.UpdateToolInput{
			Name:        "Node.js",
			Category:    "language",
			SubType:     "backend",
			Prolang:     "JavaScript",
			ReleaseYear: 2009,
			DevStatus:   "active",
			Details:     "JavaScript runtime built on Chrome's V8",
		}

		tool, err := updateTool.Execute("nodejs", input)
		require.Error(t, err)
		if _, ok := err.(common.ErrResourceNotFound); ok {
			t.Errorf("Expected error not to be ErrResourceNotFound, got %v", err)
		}
		require.Equal(t, domain.Tool{}, tool)

		toolRepository.AssertCalled(t, "GetToolBySlug", mock.Anything, mock.AnythingOfType("string"))
		toolRepository.AssertNotCalled(t, "UpdateTool", mock.Anything, mock.AnythingOfType("domain.Tool"))
	})

	t.Run("Update tool fails when the subtype is invalid", func(t *testing.T) {
		tool := testutil.CreateTestTool()
		tool.Id = 1

		toolRepository := testutil.NewMockToolRepositoryInterface(t)
		toolRepository.EXPECT().
			GetToolBySlug(mock.Anything, mock.AnythingOfType("string")).
			Return(tool, nil)

		updateTool := usecase.UpdateToolUseCase{
			ToolRepository: toolRepository,
		}

		input := usecase.UpdateToolInput{
			Category: "language",
			SubType:  "invalid",
		}

		returnedTool, err := updateTool.Execute("nodejs", input)
		require.Error(t, err)
		if _, ok := err.(domain.ErrInvalidToolSubType); !ok {
			t.Errorf("Expected ErrInvalidToolSubType, got %v", err)
		}
		require.Equal(t, domain.Tool{}, returnedTool)

		toolRepository.AssertCalled(t, "GetToolBySlug", mock.Anything, mock.AnythingOfType("string"))
		toolRepository.AssertNotCalled(t, "UpdateTool", mock.Anything, mock.AnythingOfType("domain.Tool"))
	})

	t.Run("Update tool fails with invalid fields", func(t *testing.T) {
		tool := testutil.CreateTestTool()
		tool.Id = 1

		toolRepository := testutil.NewMockToolRepositoryInterface(t)
		toolRepository.EXPECT().
			GetToolBySlug(mock.Anything, mock.AnythingOfType("string")).
			Return(tool, nil)

		updateTool := usecase.UpdateToolUseCase{
			ToolRepository: toolRepository,
		}

		input := usecase.UpdateToolInput{
			Name:        "",
			Category:    "language",
			SubType:     "backend",
			Prolang:     "",
			ReleaseYear: 1922,
			DevStatus:   "active",
			Details:     "",
			UseCases:    []string{""},
			Tags:        []string{""},
			Website:     "invalid",
			Github:      "invalid",
		}

		returnedTool, err := updateTool.Execute("nodejs", input)
		require.Error(t, err)
		if _, ok := err.(domain.ErrInvalidField); !ok {
			t.Errorf("Expected ErrInvalidField, got %v", err)
		}
		if e, ok := err.(domain.ErrInvalidField); ok {
			if _, exist := e.Fields["Name"]; !exist {
				t.Errorf("Expected the field \"Name\" to be present")
			}
			if _, exist := e.Fields["ReleaseYear"]; !exist {
				t.Errorf("Expected the field \"ReleaseYear\" to be present")
			}
			if _, exist := e.Fields["Prolang"]; !exist {
				t.Errorf("Expected the field \"Prolang\" to be present")
			}
			if _, exist := e.Fields["Website"]; !exist {
				t.Errorf("Expected the field \"Website\" to be present")
			}
			if _, exist := e.Fields["Github"]; !exist {
				t.Errorf("Expected the field \"Github\" to be present")
			}
		}
		require.Equal(t, domain.Tool{}, returnedTool)

		toolRepository.AssertCalled(t, "GetToolBySlug", mock.Anything, mock.AnythingOfType("string"))
		toolRepository.AssertNotCalled(t, "UpdateTool", mock.Anything, mock.AnythingOfType("domain.Tool"))
	})

	t.Run("Update tool fails when there is an error updating the tool", func(t *testing.T) {
		tool := testutil.CreateTestTool()
		tool.Id = 1

		toolRepository := testutil.NewMockToolRepositoryInterface(t)
		toolRepository.EXPECT().
			GetToolBySlug(mock.Anything, mock.AnythingOfType("string")).
			Return(tool, nil)
		toolRepository.EXPECT().
			UpdateTool(mock.Anything, mock.AnythingOfType("domain.Tool")).
			Return(domain.Tool{}, errors.New("error updating the tool"))

		updateTool := usecase.UpdateToolUseCase{
			ToolRepository: toolRepository,
		}

		input := usecase.UpdateToolInput{
			Name:        "Node.js",
			Category:    "language",
			SubType:     "backend",
			Prolang:     "JavaScript",
			ReleaseYear: 2009,
			DevStatus:   "active",
			Details:     "JavaScript runtime built on Chrome's V8",
		}

		returnedTool, err := updateTool.Execute("nodejs", input)
		require.Error(t, err)
		if _, ok := err.(common.ErrResourceNotFound); ok {
			t.Errorf("Expected error not to be ErrResourceNotFound, got %v", err)
		}
		require.Equal(t, domain.Tool{}, returnedTool)

		toolRepository.AssertCalled(t, "GetToolBySlug", mock.Anything, mock.AnythingOfType("string"))
		toolRepository.AssertCalled(t, "UpdateTool", mock.Anything, mock.AnythingOfType("domain.Tool"))
	})

	t.Run("Update tool succeeds", func(t *testing.T) {
		tool := testutil.CreateTestTool()
		tool.Id = 1

		toolRepository := testutil.NewMockToolRepositoryInterface(t)
		toolRepository.EXPECT().
			GetToolBySlug(mock.Anything, mock.AnythingOfType("string")).
			Return(tool, nil)
		toolRepository.EXPECT().
			UpdateTool(mock.Anything, mock.AnythingOfType("domain.Tool")).
			RunAndReturn(func(_ context.Context, tool domain.Tool) (domain.Tool, error) {
				tool.Name = "Node.js"
				tool.Category = "language"
				tool.SubType = "backend"
				tool.Prolang = "JavaScript"
				tool.ReleaseYear = 2009
				tool.Devstatus = "active"
				tool.Details = "JavaScript runtime built on Chrome's V8"
				tool.UseCases = []string{"backend"}
				tool.Tags = []string{"JavaScript"}
				tool.Website = "https://nodejs.org"
				tool.Github = "https://github.com/nodejs/node"
				return tool, nil
			})

		updateTool := usecase.UpdateToolUseCase{
			ToolRepository: toolRepository,
		}

		input := usecase.UpdateToolInput{
			Name:        "Node.js",
			Category:    "language",
			SubType:     "backend",
			Prolang:     "JavaScript",
			ReleaseYear: 2009,
			DevStatus:   "active",
			Details:     "JavaScript runtime built on Chrome's V8",
			UseCases:    []string{"backend"},
			Tags:        []string{"JavaScript"},
			Website:     "https://nodejs.org",
			Github:      "https://github.com/nodejs/node",
		}

		returnedTool, err := updateTool.Execute("golang", input)
		require.NoError(t, err)
		require.Equal(t, input.Category, returnedTool.Category)
		require.Equal(t, input.SubType, returnedTool.SubType)
		require.Equal(t, input.Prolang, returnedTool.Prolang)
		require.Equal(t, input.ReleaseYear, returnedTool.ReleaseYear)
		require.Equal(t, input.DevStatus, returnedTool.Devstatus)
		require.Equal(t, input.Details, returnedTool.Details)
		require.Equal(t, input.UseCases, returnedTool.UseCases)
		require.Equal(t, input.Tags, returnedTool.Tags)
		require.Equal(t, input.Website, returnedTool.Website)
		require.Equal(t, input.Github, returnedTool.Github)

		toolRepository.AssertCalled(t, "GetToolBySlug", mock.Anything, mock.AnythingOfType("string"))
		toolRepository.AssertCalled(t, "UpdateTool", mock.Anything, mock.AnythingOfType("domain.Tool"))
	})
}
