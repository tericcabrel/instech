//go:build integration

package repository_test

import (
	"context"
	"database/sql"
	"tericcabrel/instech/internal/domain"
	"tericcabrel/instech/internal/repository"
	"tericcabrel/instech/testutil"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func populateRelationships(t *testing.T, db *sql.DB) {
	t.Helper()
	stmts := []string{
		"INSERT INTO relationships (id, from_tool_id, to_tool_id, kind, metadata, created_at, updated_at) VALUES (65, 11, 12, 'built_on', '{\"reason\": \"reason number 1\"}', '2026-04-23 10:00:00', '2026-04-23 10:00:00')",
		"INSERT INTO relationships (id, from_tool_id, to_tool_id, kind, metadata, created_at, updated_at) VALUES (66, 12, 13, 'alternative_to', '{\"reason\": \"reason number 2\"}', '2026-04-23 10:01:00', '2026-04-23 10:01:00')",
		"INSERT INTO relationships (id, from_tool_id, to_tool_id, kind, metadata, created_at, updated_at) VALUES (67, 13, 11, 'inspired_by', '{\"reason\": \"reason number 3\"}', '2026-04-23 10:01:30', '2026-04-23 10:01:30')",
		"INSERT INTO relationships (id, from_tool_id, to_tool_id, kind, metadata, created_at, updated_at) VALUES (68, 14, 15, 'built_on', '{\"reason\": \"reason number 4\"}', '2026-04-23 10:02:00', '2026-04-23 10:02:00')",
		"INSERT INTO relationships (id, from_tool_id, to_tool_id, kind, metadata, created_at, updated_at) VALUES (69, 16, 20, 'replaced_by', '{\"reason\": \"reason number 5\"}', '2026-04-23 10:04:20', '2026-04-23 10:04:20')",
		"INSERT INTO relationships (id, from_tool_id, to_tool_id, kind, metadata, created_at, updated_at) VALUES (70, 14, 16, 'inspired_by', '{\"reason\": \"reason number 6\"}', '2026-04-23 10:05:05', '2026-04-23 10:05:05')",
		"INSERT INTO relationships (id, from_tool_id, to_tool_id, kind, metadata, created_at, updated_at) VALUES (71, 17, 18, 'built_on', '{\"reason\": \"reason number 7\"}', '2026-04-23 10:05:06', '2026-04-23 10:05:06')",
		"INSERT INTO relationships (id, from_tool_id, to_tool_id, kind, metadata, created_at, updated_at) VALUES (72, 19, 20, 'alternative_to', '{\"reason\": \"reason number 8\"}', '2026-04-23 10:06:10', '2026-04-23 10:06:10')",
		"INSERT INTO relationships (id, from_tool_id, to_tool_id, kind, metadata, created_at, updated_at) VALUES (73, 18, 19, 'used_with', '{\"reason\": \"reason number 9\"}', '2026-04-23 10:07:15', '2026-04-23 10:07:15')",
		"INSERT INTO relationships (id, from_tool_id, to_tool_id, kind, metadata, created_at, updated_at) VALUES (74, 17, 20, 'built_on', '{\"reason\": \"reason number 10\"}', '2026-04-23 10:08:50', '2026-04-23 10:08:50')",
		"INSERT INTO relationships (id, from_tool_id, to_tool_id, kind, metadata, created_at, updated_at) VALUES (75, 10, 9, 'inspired_by', '{\"reason\": \"reason number 11\"}', '2026-04-23 10:09:00', '2026-04-23 10:09:00')",
		"INSERT INTO relationships (id, from_tool_id, to_tool_id, kind, metadata, created_at, updated_at) VALUES (76, 11, 10, 'replaced_by', '{\"reason\": \"reason number 12\"}', '2026-04-23 10:10:25', '2026-04-23 10:10:25')",
		"INSERT INTO relationships (id, from_tool_id, to_tool_id, kind, metadata, created_at, updated_at) VALUES (77, 21, 22, 'built_on', '{\"reason\": \"reason number 13\"}', '2026-04-23 10:11:35', '2026-04-23 10:11:35')",
		"INSERT INTO relationships (id, from_tool_id, to_tool_id, kind, metadata, created_at, updated_at) VALUES (78, 23, 24, 'alternative_to', '{\"reason\": \"reason number 14\"}', '2026-04-23 10:12:45', '2026-04-23 10:12:45')",
		"INSERT INTO relationships (id, from_tool_id, to_tool_id, kind, metadata, created_at, updated_at) VALUES (79, 22, 25, 'inspired_by', '{\"reason\": \"reason number 15\"}', '2026-04-23 10:13:55', '2026-04-23 10:13:55')",
	}
	for _, q := range stmts {
		_, err := db.Exec(q)
		require.NoError(t, err)
	}
}

func cleanRelationships(t *testing.T, db *sql.DB) {
	t.Helper()
	_, err := db.Exec("DELETE FROM relationships WHERE id >= 65 AND id <= 79")
	require.NoError(t, err)
}

func TestRelationshipRepositoryIntegration(t *testing.T) {
	db := testutil.SetupTestDB(t)
	relationshipRepository := repository.NewRelationshipRepository(db)

	t.Run("CreateRelationship", func(t *testing.T) {
		relationship := testutil.CreateTestRelationship()
		relationship, err := relationshipRepository.CreateRelationship(context.Background(), relationship)
		require.NoError(t, err)
		require.NotZero(t, relationship.Id)

		createdRelationship, err := relationshipRepository.GetRelationshipById(context.Background(), relationship.Id)
		require.NoError(t, err)
		require.Equal(t, relationship, createdRelationship)

		db.Exec("DELETE FROM relationships WHERE id = ?", relationship.Id)
	})

	t.Run("UpdateRelationship", func(t *testing.T) {
		relationship := testutil.CreateTestRelationship()
		relationship, err := relationshipRepository.CreateRelationship(context.Background(), relationship)
		require.NoError(t, err)
		require.NotZero(t, relationship.Id)

		relationship.Kind = "alternative_to"
		relationship.Metadata = domain.RelationshipMetadata{Reason: "This is an updated test relationship"}
		updatedRelationship, err := relationshipRepository.UpdateRelationship(context.Background(), relationship)
		require.NoError(t, err)
		require.Equal(t, updatedRelationship.Kind, "alternative_to")
		require.Equal(t, updatedRelationship.Metadata.Reason, "This is an updated test relationship")

		db.Exec("DELETE FROM relationships WHERE id = ?", relationship.Id)
	})

	t.Run("DeleteRelationship", func(t *testing.T) {
		relationship := testutil.CreateTestRelationship()
		relationship, err := relationshipRepository.CreateRelationship(context.Background(), relationship)
		require.NoError(t, err)
		require.NotZero(t, relationship.Id)

		err = relationshipRepository.DeleteRelationship(context.Background(), relationship.Id)
		require.NoError(t, err)
	})

	t.Run("GetRelationshipById", func(t *testing.T) {
		relationship := testutil.CreateTestRelationship()
		relationship, err := relationshipRepository.CreateRelationship(context.Background(), relationship)
		require.NoError(t, err)
		require.NotZero(t, relationship.Id)

		createdRelationship, err := relationshipRepository.GetRelationshipById(context.Background(), relationship.Id)
		require.NoError(t, err)
		require.Equal(t, relationship, createdRelationship)

		db.Exec("DELETE FROM relationships WHERE id = ?", relationship.Id)
	})

	t.Run("GetRelationshipById fails when relationship is not found", func(t *testing.T) {
		_, err := relationshipRepository.GetRelationshipById(context.Background(), 1)
		require.Error(t, err)
		require.ErrorIs(t, err, sql.ErrNoRows)
	})

	t.Run("GetRelationshipsByToolID", func(t *testing.T) {
		firstRelationship := testutil.CreateTestRelationship() // FromToolId: 1, ToToolId: 2
		secondRelationship := testutil.CreateTestDynamicRelationship(89, domain.CreateRelationshipInput{
			FromToolId: 3,
			ToToolId:   2,
			Kind:       "inspired_by",
			Reason:     "parent to child",
		})
		thirdRelationship := testutil.CreateTestDynamicRelationship(90, domain.CreateRelationshipInput{
			FromToolId: 1,
			ToToolId:   3,
			Kind:       "used_with",
			Reason:     "child to parent",
		})

		firstRelationshipCreated, err := relationshipRepository.CreateRelationship(context.Background(), firstRelationship)
		require.NoError(t, err)
		require.NotZero(t, firstRelationshipCreated.Id)

		secondRelationshipCreated, err := relationshipRepository.CreateRelationship(context.Background(), secondRelationship)
		require.NoError(t, err)
		require.NotZero(t, secondRelationshipCreated.Id)

		thirdRelationshipCreated, err := relationshipRepository.CreateRelationship(context.Background(), thirdRelationship)
		require.NoError(t, err)
		require.NotZero(t, thirdRelationshipCreated.Id)

		relationships, err := relationshipRepository.GetRelationshipsByToolID(context.Background(), 1)
		require.NoError(t, err)
		require.Equal(t, len(relationships), 2)
		require.Equal(t, relationships[0], firstRelationshipCreated)
		require.Equal(t, relationships[1], thirdRelationshipCreated)

		db.Exec("DELETE FROM relationships WHERE id = ?", firstRelationshipCreated.Id)
		db.Exec("DELETE FROM relationships WHERE id = ?", secondRelationshipCreated.Id)
		db.Exec("DELETE FROM relationships WHERE id = ?", thirdRelationshipCreated.Id)
	})

	t.Run("GetRelationshipsByToolID succeeds when tool is not found", func(t *testing.T) {
		relationships, err := relationshipRepository.GetRelationshipsByToolID(context.Background(), 111)
		require.NoError(t, err)
		require.Equal(t, len(relationships), 0)
	})

	t.Run("GetToolAlternatives", func(t *testing.T) {
		relationship := testutil.CreateTestRelationship() // built_on 1 -> 2
		relationship, err := relationshipRepository.CreateRelationship(context.Background(), relationship)
		require.NoError(t, err)
		require.NotZero(t, relationship.Id)

		secondRelationship := testutil.CreateTestDynamicRelationship(89, domain.CreateRelationshipInput{
			FromToolId: 3,
			ToToolId:   2,
			Kind:       "alternative_to",
			Reason:     "parent to child",
		})
		secondRelationshipCreated, err := relationshipRepository.CreateRelationship(context.Background(), secondRelationship)
		require.NoError(t, err)
		require.NotZero(t, secondRelationshipCreated.Id)

		thirdRelationship := testutil.CreateTestDynamicRelationship(90, domain.CreateRelationshipInput{
			FromToolId: 2,
			ToToolId:   3,
			Kind:       "alternative_to",
			Reason:     "alternative to parent",
		})
		thirdRelationshipCreated, err := relationshipRepository.CreateRelationship(context.Background(), thirdRelationship)
		require.NoError(t, err)
		require.NotZero(t, thirdRelationshipCreated.Id)

		relationships, err := relationshipRepository.GetToolAlternatives(context.Background(), 2)
		require.NoError(t, err)
		require.Equal(t, len(relationships), 2)
		require.Equal(t, relationships[0], secondRelationshipCreated)
		require.Equal(t, relationships[1], thirdRelationshipCreated)

		db.Exec("DELETE FROM relationships WHERE id = ?", secondRelationshipCreated.Id)
		db.Exec("DELETE FROM relationships WHERE id = ?", thirdRelationshipCreated.Id)
	})

	t.Run("GetToolAlternatives succeeds when tool is not found", func(t *testing.T) {
		relationships, err := relationshipRepository.GetToolAlternatives(context.Background(), 111)
		require.NoError(t, err)
		require.Equal(t, len(relationships), 0)
	})

}

func TestGetRelationshipsAll(t *testing.T) {
	db := testutil.SetupTestDB(t)
	repo := repository.NewRelationshipRepository(db)
	ctx := context.Background()

	t.Run("no rows", func(t *testing.T) {
		res, err := repo.GetRelationshipsAll(ctx, repository.GetRelationshipsAllParams{})
		require.NoError(t, err)
		require.Empty(t, res.Relationships)
		require.Equal(t, int64(0), res.TotalCount)
		require.Equal(t, int64(0), res.ItemsCount)
		require.Equal(t, int64(-1), res.NextCursor)
	})

	t.Run("default limit returns full seed set ordered by created_at desc", func(t *testing.T) {
		populateRelationships(t, db)
		t.Cleanup(func() { cleanRelationships(t, db) })

		res, err := repo.GetRelationshipsAll(ctx, repository.GetRelationshipsAllParams{Cursor: 0, Limit: 0})
		require.NoError(t, err)
		require.Len(t, res.Relationships, 15)
		require.Equal(t, 79, res.Relationships[0].Id)
		require.Equal(t, 65, res.Relationships[14].Id)
		require.Equal(t, int64(15), res.TotalCount)
		require.Equal(t, int64(15), res.ItemsCount)
		require.Equal(t, int64(-1), res.NextCursor)
	})

	t.Run("negative cursor same as first page", func(t *testing.T) {
		populateRelationships(t, db)
		t.Cleanup(func() { cleanRelationships(t, db) })

		res, err := repo.GetRelationshipsAll(ctx, repository.GetRelationshipsAllParams{Cursor: -1, Limit: 0})
		require.NoError(t, err)
		require.Len(t, res.Relationships, 15)
		require.Equal(t, int64(15), res.TotalCount)
	})

	t.Run("limit 2 first page and cursor next page", func(t *testing.T) {
		populateRelationships(t, db)
		t.Cleanup(func() { cleanRelationships(t, db) })

		first, err := repo.GetRelationshipsAll(ctx, repository.GetRelationshipsAllParams{Limit: 2})
		require.NoError(t, err)
		require.Len(t, first.Relationships, 2)
		require.Equal(t, 79, first.Relationships[0].Id)
		require.Equal(t, 78, first.Relationships[1].Id)
		require.Equal(t, int64(15), first.TotalCount)
		require.Equal(t, int64(2), first.ItemsCount)
		// next cursor is the peeked (limit+1)th row’s time (id 77), not the last returned item
		peek20260423101135, err := time.ParseInLocation("2006-01-02 15:04:05", "2026-04-23 10:11:35", time.UTC)
		require.NoError(t, err)
		require.Equal(t, peek20260423101135.Unix(), first.NextCursor)

		second, err := repo.GetRelationshipsAll(ctx, repository.GetRelationshipsAllParams{
			Cursor: first.NextCursor,
			Limit:  2,
		})
		require.NoError(t, err)
		require.Len(t, second.Relationships, 2)
		require.Equal(t, 77, second.Relationships[0].Id)
		require.Equal(t, 76, second.Relationships[1].Id)
	})

	t.Run("filter by kind", func(t *testing.T) {
		populateRelationships(t, db)
		t.Cleanup(func() { cleanRelationships(t, db) })

		res, err := repo.GetRelationshipsAll(ctx, repository.GetRelationshipsAllParams{Kind: "built_on"})
		require.NoError(t, err)
		require.Len(t, res.Relationships, 5)
		require.Equal(t, 77, res.Relationships[0].Id)
		require.Equal(t, 65, res.Relationships[4].Id)
		require.Equal(t, int64(5), res.TotalCount)
	})

	t.Run("filter by tool id", func(t *testing.T) {
		populateRelationships(t, db)
		t.Cleanup(func() { cleanRelationships(t, db) })

		res, err := repo.GetRelationshipsAll(ctx, repository.GetRelationshipsAllParams{ToolId: 11})
		require.NoError(t, err)
		require.Len(t, res.Relationships, 3)
		require.Equal(t, 76, res.Relationships[0].Id)
		require.Equal(t, 67, res.Relationships[1].Id)
		require.Equal(t, 65, res.Relationships[2].Id)
		require.Equal(t, int64(3), res.TotalCount)
	})

	t.Run("total count is not affected by cursor while paginating", func(t *testing.T) {
		populateRelationships(t, db)
		t.Cleanup(func() { cleanRelationships(t, db) })

		p1, err := repo.GetRelationshipsAll(ctx, repository.GetRelationshipsAllParams{
			Kind:  "built_on",
			Limit: 1,
		})
		require.NoError(t, err)
		require.Len(t, p1.Relationships, 1)
		require.Equal(t, 77, p1.Relationships[0].Id)
		require.Equal(t, int64(5), p1.TotalCount)

		after, err := repo.GetRelationshipsAll(ctx, repository.GetRelationshipsAllParams{
			Kind:   "built_on",
			Cursor: p1.NextCursor,
			Limit:  1,
		})
		require.NoError(t, err)
		require.Equal(t, int64(5), after.TotalCount)
	})

	t.Run("limit capped to max", func(t *testing.T) {
		populateRelationships(t, db)
		t.Cleanup(func() { cleanRelationships(t, db) })

		res, err := repo.GetRelationshipsAll(ctx, repository.GetRelationshipsAllParams{Limit: 500})
		require.NoError(t, err)
		require.Len(t, res.Relationships, 15)
		require.Equal(t, int64(-1), res.NextCursor)
	})
}
