package usecase_test

import (
	"database/sql"
	"errors"
	"tericcabrel/instech/internal/common"
	"tericcabrel/instech/internal/domain"
	"tericcabrel/instech/internal/feature/tool/usecase"
	"tericcabrel/instech/internal/repository"
	"tericcabrel/instech/testutil"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGetToolAlternativesUseCase(t *testing.T) {
	t.Run("Get tool alternatives fails when the tool is not found", func(t *testing.T) {
		toolRepository := repository.NewMockToolRepositoryInterface(t)
		toolRepository.EXPECT().
			GetToolBySlug(mock.Anything, mock.AnythingOfType("string")).
			Return(domain.Tool{}, sql.ErrNoRows)

		getToolAlternatives := usecase.GetToolAlternativesUseCase{
			ToolRepository: toolRepository,
		}

		toolAlternatives, err := getToolAlternatives.Execute("nodejs")
		require.Error(t, err)
		if _, ok := err.(common.ErrResourceNotFound); !ok {
			t.Errorf("Expected ErrResourceNotFound, got %v", err)
		}
		require.Equal(t, []usecase.ToolAlternativesResult{}, toolAlternatives)

		toolRepository.AssertCalled(t, "GetToolBySlug", mock.Anything, mock.AnythingOfType("string"))
		toolRepository.AssertNotCalled(t, "GetToolAlternatives", mock.Anything, mock.AnythingOfType("int"))
	})

	t.Run("Get tool alternatives fails when there is an error getting the tool alternatives", func(t *testing.T) {
		toolRepository := repository.NewMockToolRepositoryInterface(t)
		toolRepository.EXPECT().
			GetToolBySlug(mock.Anything, mock.AnythingOfType("string")).
			Return(domain.Tool{}, errors.New("error getting the tool alternatives"))

		getToolAlternatives := usecase.GetToolAlternativesUseCase{
			ToolRepository: toolRepository,
		}
		toolAlternatives, err := getToolAlternatives.Execute("nodejs")
		require.Error(t, err)
		require.Equal(t, []usecase.ToolAlternativesResult{}, toolAlternatives)

		toolRepository.AssertCalled(t, "GetToolBySlug", mock.Anything, mock.AnythingOfType("string"))
		toolRepository.AssertNotCalled(t, "GetToolAlternatives", mock.Anything, mock.AnythingOfType("int"))
	})

	t.Run("Get tool alternatives is empty when the tool has no alternatives", func(t *testing.T) {
		tool := testutil.CreateTestTool()
		tool.Id = 1

		toolRepository := repository.NewMockToolRepositoryInterface(t)
		toolRepository.EXPECT().
			GetToolBySlug(mock.Anything, mock.AnythingOfType("string")).
			Return(tool, nil)

		relationshipRepository := repository.NewMockRelationshipRepositoryInterface(t)
		relationshipRepository.EXPECT().
			GetToolAlternatives(mock.Anything, mock.AnythingOfType("int")).
			Return([]domain.Relationship{}, nil)

		getToolAlternatives := usecase.GetToolAlternativesUseCase{
			ToolRepository:         toolRepository,
			RelationshipRepository: relationshipRepository,
		}

		toolAlternatives, err := getToolAlternatives.Execute("nodejs")
		require.NoError(t, err)
		require.Equal(t, []usecase.ToolAlternativesResult{}, toolAlternatives)

		toolRepository.AssertCalled(t, "GetToolBySlug", mock.Anything, mock.AnythingOfType("string"))
		toolRepository.AssertNotCalled(t, "GetToolAlternatives", mock.Anything, mock.AnythingOfType("int"))
	})

	t.Run("Get tool alternatives succeeds when the tool has alternatives", func(t *testing.T) {
		tool1 := testutil.CreateTestDynamicTool(1, domain.CreateToolInput{
			Name:        "Test Tool 1",
			Slug:        "test-tool-1",
			Category:    "language",
			SubType:     "backend",
			Devstatus:   "active",
			Details:     "Test Details for tool 1",
			UseCases:    []string{"use-case-java"},
			Tags:        []string{"tag-java"},
			Website:     "https://test-tool-1.com",
			Github:      "https://github.com/test-tool-1",
			ReleaseYear: 1995,
			Prolang:     "Java",
		})

		tool2 := testutil.CreateTestDynamicTool(2, domain.CreateToolInput{
			Name:        "Test Tool 2",
			Slug:        "test-tool-2",
			Category:    "language",
			SubType:     "fullstack",
			Devstatus:   "active",
			Details:     "Test Details for tool 2",
			UseCases:    []string{"use-case-javascript"},
			Tags:        []string{"tag-javascript"},
			Website:     "https://test-tool-2.com",
			Github:      "https://github.com/test-tool-2",
			ReleaseYear: 1995,
			Prolang:     "JavaScript",
		})

		tool3 := testutil.CreateTestDynamicTool(3, domain.CreateToolInput{
			Name:        "Test Tool 3",
			Slug:        "test-tool-3",
			Category:    "language",
			SubType:     "backend",
			Devstatus:   "active",
			Details:     "Test Details for tool 3",
			UseCases:    []string{"use-case-python"},
			Tags:        []string{"tag-python"},
			Website:     "https://test-tool-3.com",
			Github:      "https://github.com/test-tool-3",
			ReleaseYear: 1995,
			Prolang:     "Python",
		})

		relationship1 := testutil.CreateTestDynamicRelationship(1, domain.CreateRelationshipInput{
			FromToolId: tool1.Id,
			ToToolId:   tool2.Id,
			Kind:       "alternative_to",
			Reason:     "This is a test relationship for tool 1 and tool 2",
		})

		relationship2 := testutil.CreateTestDynamicRelationship(2, domain.CreateRelationshipInput{
			FromToolId: tool3.Id,
			ToToolId:   tool1.Id,
			Kind:       "alternative_to",
			Reason:     "This is a test relationship for tool 3 and tool 1",
		})

		relationshipRepository := repository.NewMockRelationshipRepositoryInterface(t)
		relationshipRepository.EXPECT().
			GetToolAlternatives(mock.Anything, mock.AnythingOfType("int")).
			Return([]domain.Relationship{relationship1, relationship2}, nil)

		toolRepository := repository.NewMockToolRepositoryInterface(t)
		toolRepository.EXPECT().
			GetToolBySlug(mock.Anything, mock.AnythingOfType("string")).
			Return(tool1, nil)
		toolRepository.EXPECT().
			GetToolByIds(mock.Anything, mock.AnythingOfType("[]int")).
			Return([]domain.Tool{tool1, tool2, tool3}, nil)

		getToolAlternatives := usecase.GetToolAlternativesUseCase{
			ToolRepository:         toolRepository,
			RelationshipRepository: relationshipRepository,
		}

		toolAlternatives, err := getToolAlternatives.Execute("test-tool-1")
		require.NoError(t, err)
		require.Equal(t, []usecase.ToolAlternativesResult{
			{Id: "test-tool-2", Name: "Test Tool 2", Category: "language", SubType: "fullstack", DevStatus: "active", Details: "Test Details for tool 2", UseCases: []string{"use-case-javascript"}, Tags: []string{"tag-javascript"}, Website: "https://test-tool-2.com", Github: "https://github.com/test-tool-2", ReleaseYear: 1995, Prolang: "JavaScript", Metadata: domain.RelationshipMetadata{Reason: "This is a test relationship for tool 1 and tool 2"}},
			{Id: "test-tool-3", Name: "Test Tool 3", Category: "language", SubType: "backend", DevStatus: "active", Details: "Test Details for tool 3", UseCases: []string{"use-case-python"}, Tags: []string{"tag-python"}, Website: "https://test-tool-3.com", Github: "https://github.com/test-tool-3", ReleaseYear: 1995, Prolang: "Python", Metadata: domain.RelationshipMetadata{Reason: "This is a test relationship for tool 3 and tool 1"}},
		}, toolAlternatives)

		toolRepository.AssertCalled(t, "GetToolBySlug", mock.Anything, mock.AnythingOfType("string"))
		toolRepository.AssertCalled(t, "GetToolByIds", mock.Anything, mock.AnythingOfType("[]int"))
		relationshipRepository.AssertCalled(t, "GetToolAlternatives", mock.Anything, mock.AnythingOfType("int"))
	})
}
