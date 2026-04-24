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

func TestCreateRelationshipUseCase(t *testing.T) {
	t.Run("Create relationship with invalid source tool ID", func(t *testing.T) {
		relationshipRepository := testutil.NewMockRelationshipRepositoryInterface(t)
		toolRepository := testutil.NewMockToolRepositoryInterface(t)
		toolRepository.EXPECT().
			GetToolBySlug(mock.Anything, mock.AnythingOfType("string")).
			Return(domain.Tool{}, sql.ErrNoRows)
		createRelationship := usecase.CreateRelationshipUseCase{
			RelationshipRepository: relationshipRepository,
			ToolRepository:         toolRepository,
		}
		relationship, err := createRelationship.Execute(usecase.CreateRelationshipInput{
			FromToolId: "invalid",
			ToToolId:   "1",
			Kind:       "built_on",
			Metadata:   struct{ Reason string }{Reason: "This is a test relationship"},
		})
		if err == nil {
			t.Errorf("Expected error, got %v", err)
		}
		if _, ok := err.(common.ErrResourceNotFound); !ok {
			t.Errorf("Expected ErrResourceNotFound, got %v", err)
		}
		if e, ok := err.(common.ErrResourceNotFound); ok {
			if e.Id != "invalid" {
				t.Errorf("Expected ID to be 'invalid', got %s", e.Id)
			}
		}
		require.Equal(t, domain.Relationship{}, relationship)
		relationshipRepository.AssertNotCalled(t, "CreateRelationship", mock.Anything, mock.AnythingOfType("domain.Relationship"))
		toolRepository.AssertNumberOfCalls(t, "GetToolBySlug", 1)
	})

	t.Run("Create relationship with invalid target tool ID", func(t *testing.T) {
		tool := testutil.CreateTestTool()

		relationshipRepository := testutil.NewMockRelationshipRepositoryInterface(t)
		toolRepository := testutil.NewMockToolRepositoryInterface(t)

		toolRepository.EXPECT().
			GetToolBySlug(mock.Anything, tool.Slug).
			Return(tool, nil)
		toolRepository.EXPECT().
			GetToolBySlug(mock.Anything, "invalid").
			Return(domain.Tool{}, sql.ErrNoRows)

		createRelationship := usecase.CreateRelationshipUseCase{
			RelationshipRepository: relationshipRepository,
			ToolRepository:         toolRepository,
		}
		relationship, err := createRelationship.Execute(usecase.CreateRelationshipInput{
			FromToolId: tool.Slug,
			ToToolId:   "invalid",
			Kind:       "built_on",
			Metadata:   struct{ Reason string }{Reason: "This is a test relationship"},
		})
		if err == nil {
			t.Errorf("Expected error, got %v", err)
		}
		if _, ok := err.(common.ErrResourceNotFound); !ok {
			t.Errorf("Expected ErrResourceNotFound, got %v", err)
		}
		if e, ok := err.(common.ErrResourceNotFound); ok {
			if e.Id != "invalid" {
				t.Errorf("Expected ID to be 'invalid', got %s", e.Id)
			}
		}
		require.Equal(t, domain.Relationship{}, relationship)
		relationshipRepository.AssertNotCalled(t, "CreateRelationship", mock.Anything, mock.AnythingOfType("domain.Relationship"))
		toolRepository.AssertNumberOfCalls(t, "GetToolBySlug", 2)
	})

	t.Run("Create relationship with invalid kind", func(t *testing.T) {
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

		toolRepository.EXPECT().
			GetToolBySlug(mock.Anything, tool1.Slug).
			Return(tool1, nil)
		toolRepository.EXPECT().
			GetToolBySlug(mock.Anything, tool2.Slug).
			Return(tool2, nil)

		createRelationship := usecase.CreateRelationshipUseCase{
			RelationshipRepository: relationshipRepository,
			ToolRepository:         toolRepository,
		}
		relationship, err := createRelationship.Execute(usecase.CreateRelationshipInput{
			FromToolId: tool1.Slug,
			ToToolId:   tool2.Slug,
			Kind:       "invalid",
			Metadata:   struct{ Reason string }{Reason: "This is a test relationship"},
		})

		require.Error(t, err)
		if _, ok := err.(domain.ErrInvalidRelationshipKind); !ok {
			t.Errorf("Expected ErrInvalidRelationshipKind, got %v", err)
		}
		if e, ok := err.(domain.ErrInvalidRelationshipKind); ok {
			if e.Kind != "invalid" {
				t.Errorf("Expected kind to be 'invalid', got %s", e.Kind)
			}
		}
		require.Equal(t, domain.Relationship{}, relationship)
		relationshipRepository.AssertNotCalled(t, "CreateRelationship", mock.Anything, mock.AnythingOfType("domain.Relationship"))
		toolRepository.AssertNumberOfCalls(t, "GetToolBySlug", 2)
	})

	t.Run("Create relationship with valid input", func(t *testing.T) {
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

		expectedRelationship := testutil.CreateTestDynamicRelationship(1, domain.CreateRelationshipInput{
			FromToolId: tool1.Id,
			ToToolId:   tool2.Id,
			Kind:       "built_on",
			Reason:     "This is a test relationship",
		})

		relationshipRepository := testutil.NewMockRelationshipRepositoryInterface(t)
		toolRepository := testutil.NewMockToolRepositoryInterface(t)

		toolRepository.EXPECT().
			GetToolBySlug(mock.Anything, tool1.Slug).
			Return(tool1, nil)
		toolRepository.EXPECT().
			GetToolBySlug(mock.Anything, tool2.Slug).
			Return(tool2, nil)
		relationshipRepository.EXPECT().
			CreateRelationship(mock.Anything, mock.AnythingOfType("domain.Relationship")).
			Return(expectedRelationship, nil)

		createRelationship := usecase.CreateRelationshipUseCase{
			RelationshipRepository: relationshipRepository,
			ToolRepository:         toolRepository,
		}
		relationship, err := createRelationship.Execute(usecase.CreateRelationshipInput{
			FromToolId: tool1.Slug,
			ToToolId:   tool2.Slug,
			Kind:       "built_on",
			Metadata:   struct{ Reason string }{Reason: "This is a test relationship"},
		})

		require.NoError(t, err)
		require.Equal(t, relationship, expectedRelationship)
		relationshipRepository.AssertCalled(t, "CreateRelationship", mock.Anything, domain.Relationship{
			FromToolId: tool1.Id,
			ToToolId:   tool2.Id,
			Kind:       "built_on",
			Metadata: domain.RelationshipMetadata{
				Reason: "This is a test relationship",
			},
		})
		toolRepository.AssertNumberOfCalls(t, "GetToolBySlug", 2)
	})
}
