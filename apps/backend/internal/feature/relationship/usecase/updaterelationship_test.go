package usecase_test

import (
	"database/sql"
	"tericcabrel/instech/internal/common"
	"tericcabrel/instech/internal/domain"
	"tericcabrel/instech/internal/feature/relationship/usecase"
	"tericcabrel/instech/testutil"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestUpdateRelationshipUseCase(t *testing.T) {
	t.Run("Update relationship will fail if the relationship is not found", func(t *testing.T) {
		relationshipRepository := testutil.NewMockRelationshipRepositoryInterface(t)
		toolRepository := testutil.NewMockToolRepositoryInterface(t)

		relationshipRepository.EXPECT().
			GetRelationshipById(mock.Anything, 99).
			Return(domain.Relationship{}, sql.ErrNoRows)

		uc := usecase.UpdateRelationshipUseCase{
			RelationshipRepository: relationshipRepository,
			ToolRepository:         toolRepository,
		}
		rel, err := uc.Execute(99, usecase.UpdateRelationshipInput{
			FromToolId: 1,
			ToToolId:   2,
			Kind:       "built_on",
			Metadata:   domain.RelationshipMetadata{Reason: "x"},
		})

		require.Equal(t, domain.Relationship{}, rel)
		var notFound common.ErrResourceNotFound
		require.ErrorAs(t, err, &notFound)
		require.Equal(t, "99", notFound.Id)
		toolRepository.AssertNotCalled(t, "GetToolById", mock.Anything, mock.Anything)
		relationshipRepository.AssertNotCalled(t, "UpdateRelationship", mock.Anything, mock.Anything)
	})

	t.Run("Update relationship with invalid source tool", func(t *testing.T) {
		existing := testutil.CreateTestDynamicRelationship(1, domain.CreateRelationshipInput{
			FromToolId: 1,
			ToToolId:   2,
			Kind:       "built_on",
			Reason:     "original",
		})

		relationshipRepository := testutil.NewMockRelationshipRepositoryInterface(t)
		toolRepository := testutil.NewMockToolRepositoryInterface(t)

		relationshipRepository.EXPECT().
			GetRelationshipById(mock.Anything, 1).
			Return(existing, nil)
		toolRepository.EXPECT().
			GetToolById(mock.Anything, 99).
			Return(domain.Tool{}, sql.ErrNoRows)

		updateRelationship := usecase.UpdateRelationshipUseCase{
			RelationshipRepository: relationshipRepository,
			ToolRepository:         toolRepository,
		}
		relationship, err := updateRelationship.Execute(1, usecase.UpdateRelationshipInput{
			FromToolId: 99,
			ToToolId:   2,
			Kind:       "built_on",
			Metadata:   domain.RelationshipMetadata{Reason: "x"},
		})

		require.Equal(t, domain.Relationship{}, relationship)
		require.Error(t, err)
		if _, ok := err.(common.ErrResourceNotFound); !ok {
			t.Errorf("Expected ErrResourceNotFound, got %v", err)
		}
		if e, ok := err.(common.ErrResourceNotFound); ok {
			if e.Id != "99" {
				t.Errorf("Expected ID to be '99', got %s", e.Id)
			}
		}
		relationshipRepository.AssertNotCalled(t, "UpdateRelationship", mock.Anything, mock.Anything)
	})

	t.Run("Update relationship with invalid target tool", func(t *testing.T) {
		existing := testutil.CreateTestDynamicRelationship(1, domain.CreateRelationshipInput{
			FromToolId: 1,
			ToToolId:   2,
			Kind:       "built_on",
			Reason:     "original",
		})

		relationshipRepository := testutil.NewMockRelationshipRepositoryInterface(t)
		toolRepository := testutil.NewMockToolRepositoryInterface(t)

		relationshipRepository.EXPECT().
			GetRelationshipById(mock.Anything, 1).
			Return(existing, nil)
		toolRepository.EXPECT().
			GetToolById(mock.Anything, 99).
			Return(domain.Tool{}, sql.ErrNoRows)

		uc := usecase.UpdateRelationshipUseCase{
			RelationshipRepository: relationshipRepository,
			ToolRepository:         toolRepository,
		}
		rel, err := uc.Execute(1, usecase.UpdateRelationshipInput{
			FromToolId: 1,
			ToToolId:   99,
			Kind:       "built_on",
			Metadata:   domain.RelationshipMetadata{Reason: "x"},
		})

		require.Equal(t, domain.Relationship{}, rel)
		require.Error(t, err)
		if _, ok := err.(common.ErrResourceNotFound); !ok {
			t.Errorf("Expected ErrResourceNotFound, got %v", err)
		}
		if e, ok := err.(common.ErrResourceNotFound); ok {
			if e.Id != "99" {
				t.Errorf("Expected ID to be '99', got %s", e.Id)
			}
		}
		relationshipRepository.AssertNotCalled(t, "UpdateRelationship", mock.Anything, mock.Anything)
	})

	t.Run("Update relationship with invalid kind", func(t *testing.T) {
		existing := testutil.CreateTestDynamicRelationship(1, domain.CreateRelationshipInput{
			FromToolId: 1,
			ToToolId:   2,
			Kind:       "built_on",
			Reason:     "original",
		})

		relationshipRepository := testutil.NewMockRelationshipRepositoryInterface(t)
		toolRepository := testutil.NewMockToolRepositoryInterface(t)

		relationshipRepository.EXPECT().
			GetRelationshipById(mock.Anything, 1).
			Return(existing, nil)

		updateRelationship := usecase.UpdateRelationshipUseCase{
			RelationshipRepository: relationshipRepository,
			ToolRepository:         toolRepository,
		}
		_, err := updateRelationship.Execute(1, usecase.UpdateRelationshipInput{
			FromToolId: 1,
			ToToolId:   2,
			Kind:       "not_a_real_kind",
			Metadata:   domain.RelationshipMetadata{Reason: "x"},
		})

		require.Error(t, err)
		if _, ok := err.(domain.ErrInvalidRelationshipKind); !ok {
			t.Errorf("Expected ErrInvalidRelationshipKind, got %v", err)
		}
		if e, ok := err.(domain.ErrInvalidRelationshipKind); ok {
			if e.Kind != "not_a_real_kind" {
				t.Errorf("Expected kind to be 'not_a_real_kind', got %s", e.Kind)
			}
		}
		toolRepository.AssertNotCalled(t, "GetToolById", mock.Anything, mock.Anything)
		relationshipRepository.AssertNotCalled(t, "UpdateRelationship", mock.Anything, mock.Anything)
	})

	t.Run("Update relationship kind and metadata when source and target tool ids are the same", func(t *testing.T) {
		existing := testutil.CreateTestDynamicRelationship(1, domain.CreateRelationshipInput{
			FromToolId: 1,
			ToToolId:   2,
			Kind:       "built_on",
			Reason:     "original",
		})

		expected := existing
		expected.Kind = "inspired_by"
		expected.Metadata = domain.RelationshipMetadata{Reason: "updated reason"}

		relationshipRepository := testutil.NewMockRelationshipRepositoryInterface(t)
		toolRepository := testutil.NewMockToolRepositoryInterface(t)

		relationshipRepository.EXPECT().
			GetRelationshipById(mock.Anything, 1).
			Return(existing, nil)
		relationshipRepository.EXPECT().
			UpdateRelationship(mock.Anything, mock.AnythingOfType("domain.Relationship")).
			Return(expected, nil)

		updateRelationship := usecase.UpdateRelationshipUseCase{
			RelationshipRepository: relationshipRepository,
			ToolRepository:         toolRepository,
		}
		rel, err := updateRelationship.Execute(1, usecase.UpdateRelationshipInput{
			FromToolId: 1,
			ToToolId:   2,
			Kind:       "inspired_by",
			Metadata:   domain.RelationshipMetadata{Reason: "updated reason"},
		})

		require.NoError(t, err)
		require.Equal(t, expected, rel)
		toolRepository.AssertNotCalled(t, "GetToolById", mock.Anything, mock.Anything)
		relationshipRepository.AssertCalled(t, "UpdateRelationship", mock.Anything, expected)
	})

	t.Run("Update relationship when changing source and target tool ids", func(t *testing.T) {
		tool1 := testutil.CreateTestTool()
		tool3 := testutil.CreateTestDynamicTool(3, domain.CreateToolInput{
			Name:        "Other",
			Slug:        "other",
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

		existing := testutil.CreateTestDynamicRelationship(5, domain.CreateRelationshipInput{
			FromToolId: 1,
			ToToolId:   2,
			Kind:       "built_on",
			Reason:     "original",
		})

		expectedRelationship := existing
		expectedRelationship.FromToolId = tool3.Id
		expectedRelationship.ToToolId = tool1.Id
		expectedRelationship.Kind = "alternative_to"
		expectedRelationship.Metadata = domain.RelationshipMetadata{Reason: "swapped"}

		relationshipRepository := testutil.NewMockRelationshipRepositoryInterface(t)
		toolRepository := testutil.NewMockToolRepositoryInterface(t)

		relationshipRepository.EXPECT().
			GetRelationshipById(mock.Anything, 5).
			Return(existing, nil)
		toolRepository.EXPECT().
			GetToolById(mock.Anything, tool3.Id).
			Return(tool3, nil)
		toolRepository.EXPECT().
			GetToolById(mock.Anything, tool1.Id).
			Return(tool1, nil)
		relationshipRepository.EXPECT().
			UpdateRelationship(mock.Anything, expectedRelationship).
			Return(expectedRelationship, nil)

		updateRelationship := usecase.UpdateRelationshipUseCase{
			RelationshipRepository: relationshipRepository,
			ToolRepository:         toolRepository,
		}
		updatedRelationship, err := updateRelationship.Execute(5, usecase.UpdateRelationshipInput{
			FromToolId: tool3.Id,
			ToToolId:   tool1.Id,
			Kind:       "alternative_to",
			Metadata:   domain.RelationshipMetadata{Reason: "swapped"},
		})

		require.NoError(t, err)
		require.Equal(t, expectedRelationship, updatedRelationship)
		toolRepository.AssertCalled(t, "GetToolById", mock.Anything, tool3.Id)
		toolRepository.AssertCalled(t, "GetToolById", mock.Anything, tool1.Id)
		relationshipRepository.AssertCalled(t, "UpdateRelationship", mock.Anything, expectedRelationship)
	})
}
