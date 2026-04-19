package domain_test

import (
	"tericcabrel/instech/internal/domain"
	"tericcabrel/instech/testutil"
	"testing"
)

func TestRelationship(t *testing.T) {
	t.Run("Create relationship with invalid kind", func(t *testing.T) {
		_, err := domain.CreateRelationship(domain.CreateRelationshipInput{
			FromToolId: 1,
			ToToolId:   1,
			Kind:       "invalid",
			Reason:     "",
		})
		if err == nil {
			t.Errorf("Expected error, got %v", err)
		}
		if _, ok := err.(domain.ErrInvalidRelationshipKind); !ok {
			t.Errorf("Expected ErrInvalidRelationshipKind, got %v", err)
		}
		if e, ok := err.(domain.ErrInvalidRelationshipKind); ok {
			if e.Kind != "invalid" {
				t.Errorf("Expected kind to be 'invalid', got %s", e.Kind)
			}
		}
	})

	t.Run("Create relationship with invalid fields", func(t *testing.T) {
		_, err := domain.CreateRelationship(domain.CreateRelationshipInput{
			FromToolId: 0,
			ToToolId:   -3,
			Kind:       "built_on",
			Reason:     "",
		})
		if err == nil {
			t.Errorf("Expected error, got %v", err)
		}
		if _, ok := err.(domain.ErrInvalidField); !ok {
			t.Errorf("Expected ErrInvalidField, got %v", err)
		}
		if e, ok := err.(domain.ErrInvalidField); ok {
			if _, exist := e.Fields["FromToolId"]; !exist {
				t.Errorf("Expected the field \"FromToolId\" to be present")
			}
			if _, exist := e.Fields["ToToolId"]; !exist {
				t.Errorf("Expected the field \"ToToolId\" to be present")
			}
		}
	})

	t.Run("Create relationship with valid input", func(t *testing.T) {
		relationship, err := domain.CreateRelationship(domain.CreateRelationshipInput{
			FromToolId: 1,
			ToToolId:   1,
			Kind:       "built_on",
			Reason:     "This is a test relationship",
		})
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if relationship.FromToolId != 1 {
			t.Errorf("Expected source tool ID to be 1, got %d", relationship.FromToolId)
		}
		if relationship.ToToolId != 1 {
			t.Errorf("Expected target tool ID to be 1, got %d", relationship.ToToolId)
		}
		if relationship.Kind != "built_on" {
			t.Errorf("Expected kind to be 'built_on', got %s", relationship.Kind)
		}
		if relationship.Metadata.Reason != "This is a test relationship" {
			t.Errorf("Expected metadata reason to be 'This is a test relationship', got %s", relationship.Metadata.Reason)
		}
	})

	t.Run("Update relationship with invalid kind", func(t *testing.T) {
		relationship := testutil.CreateTestRelationship()
		err := relationship.Update(domain.UpdateRelationshipInput{
			FromToolId: 1,
			ToToolId:   1,
			Kind:       "invalid",
			Metadata: domain.RelationshipMetadata{
				Reason: "This is a test relationship",
			},
		})
		if err == nil {
			t.Errorf("Expected error, got %v", err)
		}
		if _, ok := err.(domain.ErrInvalidRelationshipKind); !ok {
			t.Errorf("Expected ErrInvalidRelationshipKind, got %v", err)
		}
		if e, ok := err.(domain.ErrInvalidRelationshipKind); ok {
			if e.Kind != "invalid" {
				t.Errorf("Expected kind to be 'invalid', got %s", e.Kind)
			}
		}
	})

	t.Run("Update relationship with invalid fields", func(t *testing.T) {
		relationship := testutil.CreateTestRelationship()
		err := relationship.Update(domain.UpdateRelationshipInput{
			FromToolId: -3,
			ToToolId:   0,
			Kind:       "built_on",
			Metadata: domain.RelationshipMetadata{
				Reason: "This is a test relationship",
			},
		})
		if err == nil {
			t.Errorf("Expected error, got %v", err)
		}
		if _, ok := err.(domain.ErrInvalidField); !ok {
			t.Errorf("Expected ErrInvalidField, got %v", err)
		}
		if e, ok := err.(domain.ErrInvalidField); ok {
			if _, exist := e.Fields["FromToolId"]; !exist {
				t.Errorf("Expected the field \"FromToolId\" to be present")
			}
			if _, exist := e.Fields["ToToolId"]; !exist {
				t.Errorf("Expected the field \"ToToolId\" to be present")
			}
		}
	})

	t.Run("Update relationship with valid input", func(t *testing.T) {
		relationship := testutil.CreateTestRelationship()
		err := relationship.Update(domain.UpdateRelationshipInput{
			FromToolId: 2,
			ToToolId:   3,
			Kind:       "inspired_by",
			Metadata: domain.RelationshipMetadata{
				Reason: "reason updated",
			},
		})
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if relationship.FromToolId != 2 {
			t.Errorf("Expected source tool ID to be 2, got %d", relationship.FromToolId)
		}
		if relationship.ToToolId != 3 {
			t.Errorf("Expected target tool ID to be 3, got %d", relationship.ToToolId)
		}
		if relationship.Kind != "inspired_by" {
			t.Errorf("Expected kind to be 'inspired_by', got %s", relationship.Kind)
		}
		if relationship.Metadata.Reason != "reason updated" {
			t.Errorf("Expected metadata reason to be 'reason updated', got %s", relationship.Metadata.Reason)
		}
	})
}
