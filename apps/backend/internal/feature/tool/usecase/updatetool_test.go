package usecase_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"tericcabrel/instech/internal/common"
	"tericcabrel/instech/internal/domain"
	"tericcabrel/instech/internal/feature/tool/usecase"
	"tericcabrel/instech/testutil"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestUpdateToolUseCase(t *testing.T) {
	t.Run("Update tool fails when the tool is not found", func(t *testing.T) {
		toolRepository := testutil.NewMockToolRepositoryInterface(t)
		toolRepository.EXPECT().
			GetToolBySlug(mock.Anything, mock.AnythingOfType("string")).
			Return(domain.Tool{}, sql.ErrNoRows)

		patchTool := usecase.PatchToolUseCase{ToolRepository: toolRepository}
		tool, err := patchTool.Execute("nodejs", usecase.PatchToolInput{})
		require.Error(t, err)
		if _, ok := err.(common.ErrResourceNotFound); !ok {
			t.Errorf("Expected ErrResourceNotFound, got %v", err)
		}
		require.Equal(t, domain.Tool{}, tool)
	})

	t.Run("Update tool fails when there is an error getting the tool", func(t *testing.T) {
		toolRepository := testutil.NewMockToolRepositoryInterface(t)
		toolRepository.EXPECT().
			GetToolBySlug(mock.Anything, mock.AnythingOfType("string")).
			Return(domain.Tool{}, errors.New("error getting the tool"))

		patchTool := usecase.PatchToolUseCase{ToolRepository: toolRepository}
		tool, err := patchTool.Execute("nodejs", usecase.PatchToolInput{})
		require.Error(t, err)
		require.Equal(t, domain.Tool{}, tool)
	})

	t.Run("Update tool fails when subtype is invalid", func(t *testing.T) {
		tool := testutil.CreateTestTool()
		tool.Id = 1

		toolRepository := testutil.NewMockToolRepositoryInterface(t)
		toolRepository.EXPECT().
			GetToolBySlug(mock.Anything, mock.AnythingOfType("string")).
			Return(tool, nil)

		patchTool := usecase.PatchToolUseCase{ToolRepository: toolRepository}
		input := usecase.PatchToolInput{SubType: common.PatchStringField{IsSet: true, Value: "invalid"}}

		returnedTool, err := patchTool.Execute("nodejs", input)
		require.Error(t, err)
		if _, ok := err.(domain.ErrInvalidToolSubType); !ok {
			t.Errorf("Expected ErrInvalidToolSubType, got %v", err)
		}
		require.Equal(t, domain.Tool{}, returnedTool)
	})

	t.Run("Update tool fails with invalid fields", func(t *testing.T) {
		tool := testutil.CreateTestTool()
		tool.Id = 1

		toolRepository := testutil.NewMockToolRepositoryInterface(t)
		toolRepository.EXPECT().
			GetToolBySlug(mock.Anything, mock.AnythingOfType("string")).
			Return(tool, nil)

		patchTool := usecase.PatchToolUseCase{ToolRepository: toolRepository}
		input := usecase.PatchToolInput{
			Name:        common.PatchStringField{IsSet: true, Value: ""},
			ReleaseYear: common.PatchIntField{IsSet: true, Value: 1922},
			Prolang:     common.PatchNullableStringField{IsSet: true, Value: new(string)},
			Details:     common.PatchNullableStringField{IsSet: true, Value: new(string)},
			Website:     common.PatchNullableStringField{IsSet: true, Value: new(string)},
			Github:      common.PatchNullableStringField{IsSet: true, Value: new(string)},
		}

		returnedTool, err := patchTool.Execute("nodejs", input)
		require.Error(t, err)
		if _, ok := err.(domain.ErrInvalidField); !ok {
			t.Errorf("Expected ErrInvalidField, got %v", err)
		}
		require.Equal(t, domain.Tool{}, returnedTool)
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

		patchTool := usecase.PatchToolUseCase{ToolRepository: toolRepository}
		input := usecase.PatchToolInput{Name: common.PatchStringField{IsSet: true, Value: "Node.js"}}

		returnedTool, err := patchTool.Execute("nodejs", input)

		require.Error(t, err)
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
			RunAndReturn(func(_ context.Context, patched domain.Tool) (domain.Tool, error) {
				return patched, nil
			})

		patchTool := usecase.PatchToolUseCase{ToolRepository: toolRepository}
		input := usecase.PatchToolInput{
			Name:        common.PatchStringField{IsSet: true, Value: "Express.js"},
			Category:    common.PatchStringField{IsSet: true, Value: "framework"},
			SubType:     common.PatchStringField{IsSet: true, Value: "backend"},
			Prolang:     common.PatchNullableStringField{IsSet: true, Value: new("JavaScript")},
			ReleaseYear: common.PatchIntField{IsSet: true, Value: 2012},
			Details:     common.PatchNullableStringField{IsSet: true, Value: nil},
			Github:      common.PatchNullableStringField{IsSet: false, Value: nil},
			UseCases:    common.PatchStringSliceField{IsSet: true, Value: []string{"backend"}},
			Website:     common.PatchNullableStringField{IsSet: true, Value: new("https://expressjs.com")},
			Tags:        common.PatchStringSliceField{IsSet: true, Value: []string{"web", "api", "framework"}},
			DevStatus:   common.PatchStringField{IsSet: true, Value: "deprecated"},
		}

		returnedTool, err := patchTool.Execute("golang", input)
		require.NoError(t, err)
		require.Equal(t, "Express.js", returnedTool.Name)
		require.Equal(t, "framework", returnedTool.Category)
		require.Equal(t, "backend", returnedTool.SubType)
		require.Equal(t, new("JavaScript"), returnedTool.Prolang)
		require.Equal(t, 2012, returnedTool.ReleaseYear)
		require.Nil(t, returnedTool.Details)
		require.Equal(t, tool.Github, returnedTool.Github)
		require.Equal(t, []string{"backend"}, returnedTool.UseCases)
		require.Equal(t, "https://expressjs.com", *returnedTool.Website)
		require.Equal(t, []string{"web", "api", "framework"}, returnedTool.Tags)
		require.Equal(t, "deprecated", returnedTool.DevStatus)
	})
}
