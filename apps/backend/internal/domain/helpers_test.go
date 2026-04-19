package domain_test

import (
	"slices"
	"tericcabrel/instech/internal/domain"
	"testing"
)

func TestHelpers(t *testing.T) {
	t.Run("DedupeToolIdsFromRelationships with multiple relationships", func(t *testing.T) {
		relationships := []domain.Relationship{
			{FromToolId: 1, ToToolId: 2},
			{FromToolId: 2, ToToolId: 3},
			{FromToolId: 1, ToToolId: 3},
		}

		toolIds := domain.DedupeToolIdsFromRelationships(relationships)
		if len(toolIds) != 3 {
			t.Errorf("Expected 3 tool IDs, got %d", len(toolIds))
		}
		if !slices.Contains(toolIds, 1) {
			t.Errorf("Expected tool ID 1 to be present, got %v", toolIds)
		}
		if !slices.Contains(toolIds, 2) {
			t.Errorf("Expected tool ID 2 to be present, got %v", toolIds)
		}
		if !slices.Contains(toolIds, 3) {
			t.Errorf("Expected tool ID 3 to be present, got %v", toolIds)
		}
	})

	t.Run("DedupeToolIdsFromRelationships with empty relationships", func(t *testing.T) {
		relationships := []domain.Relationship{}
		toolIds := domain.DedupeToolIdsFromRelationships(relationships)
		if len(toolIds) != 0 {
			t.Errorf("Expected 0 tool IDs, got %d", len(toolIds))
		}
	})
}
