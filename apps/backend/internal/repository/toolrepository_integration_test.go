//go:build integration

package repository_test

import (
	"context"
	"tericcabrel/instech/internal/domain"
	"tericcabrel/instech/internal/repository"
	"tericcabrel/instech/testutil"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestToolRepositoryIntegration(t *testing.T) {
	db := testutil.SetupTestDB(t)
	toolRepository := repository.NewToolRepository(db)

	t.Run("CreateTool", func(t *testing.T) {
		tool := testutil.CreateTestTool()
		tool, err := toolRepository.CreateTool(context.Background(), tool)

		require.NoError(t, err)
		require.NotZero(t, tool.Id)

		createdTool, err := toolRepository.GetToolByID(context.Background(), tool.Id)
		require.NoError(t, err)
		require.Equal(t, tool, createdTool)

		db.Exec("DELETE FROM tools WHERE id = ?", tool.Id)
	})

	t.Run("CreateTool fails when slug already exists", func(t *testing.T) {
		toolToCreate := testutil.CreateTestTool()
		firstTool, err := toolRepository.CreateTool(context.Background(), toolToCreate)
		require.NoError(t, err)
		require.NotZero(t, firstTool.Id)

		secondTool, err := toolRepository.CreateTool(context.Background(), toolToCreate)
		require.Error(t, err)
		require.Equal(t, err.Error(), "constraint failed: UNIQUE constraint failed: tools.slug (2067)")

		require.Equal(t, domain.Tool{}, secondTool)

		db.Exec("DELETE FROM tools WHERE id = ?", firstTool.Id)
	})

	t.Run("CreateTool persists omitted optional fields as NULL", func(t *testing.T) {
		tool, err := domain.CreateTool(domain.CreateToolInput{
			Name:        "React",
			Slug:        "react",
			Category:    "framework",
			SubType:     "frontend",
			DevStatus:   "active",
			ReleaseYear: 2013,
		})
		require.NoError(t, err)

		createdTool, err := toolRepository.CreateTool(context.Background(), tool)
		require.NoError(t, err)
		require.NotZero(t, createdTool.Id)
		require.Nil(t, createdTool.Prolang)
		require.Nil(t, createdTool.Details)
		require.Nil(t, createdTool.Website)
		require.Nil(t, createdTool.Github)

		storedTool, err := toolRepository.GetToolByID(context.Background(), createdTool.Id)
		require.NoError(t, err)
		require.Nil(t, storedTool.Prolang)
		require.Nil(t, storedTool.Details)
		require.Nil(t, storedTool.Website)
		require.Nil(t, storedTool.Github)

		db.Exec("DELETE FROM tools WHERE id = ?", createdTool.Id)
	})

	t.Run("UpdateTool", func(t *testing.T) {
		tool := testutil.CreateTestTool()
		tool, err := toolRepository.CreateTool(context.Background(), tool)
		require.NoError(t, err)
		require.NotZero(t, tool.Id)

		tool.Name = "Node.js"
		tool.Category = "language"
		tool.SubType = "backend"
		tool.Prolang = new("JavaScript")
		tool.ReleaseYear = 2009
		tool.DevStatus = "active"
		tool.Details = new("JavaScript runtime built on Chrome's V8")
		tool.UseCases = []string{"backend"}
		tool.Tags = []string{"JavaScript"}
		tool.Website = new("https://nodejs.org")
		tool.Github = new("https://github.com/nodejs/node")
		tool.Slug = "nodejs"

		updatedTool, err := toolRepository.UpdateTool(context.Background(), tool)

		require.NoError(t, err)
		require.Equal(t, tool, updatedTool)

		db.Exec("DELETE FROM tools WHERE id = ?", tool.Id)
	})

	t.Run("UpdateTool fails when tool is not found", func(t *testing.T) {
		tool := testutil.CreateTestTool()

		updatedTool, err := toolRepository.UpdateTool(context.Background(), tool)

		require.Error(t, err)
		require.Equal(t, err.Error(), "sql: no rows in result set")
		require.Equal(t, domain.Tool{}, updatedTool)

		db.Exec("DELETE FROM tools WHERE id = ?", tool.Id)
	})

	t.Run("UpdateTool fails when the slug already exists", func(t *testing.T) {
		firstTool := testutil.CreateTestTool()
		firstToolCreated, err := toolRepository.CreateTool(context.Background(), firstTool)
		require.NoError(t, err)
		require.NotZero(t, firstToolCreated.Id)

		secondTool := testutil.CreateTestDynamicTool(0, domain.CreateToolInput{
			Slug:        "nodejs",
			Name:        "Node.js",
			Category:    "language",
			SubType:     "backend",
			Prolang:     new("JavaScript"),
			ReleaseYear: 2009,
			DevStatus:   "active",
			Details:     new("JavaScript runtime built on Chrome's V8"),
			UseCases:    []string{"backend"},
			Tags:        []string{"JavaScript"},
			Website:     new("https://nodejs.org"),
			Github:      new("https://github.com/nodejs/node"),
		})

		secondToolCreated, err := toolRepository.CreateTool(context.Background(), secondTool)
		require.NoError(t, err)
		require.NotZero(t, secondToolCreated.Id)

		secondToolCreated.Slug = "golang"
		updatedTool, err := toolRepository.UpdateTool(context.Background(), secondToolCreated)
		require.Error(t, err)
		require.Equal(t, "constraint failed: UNIQUE constraint failed: tools.slug (2067)", err.Error())
		require.Equal(t, domain.Tool{}, updatedTool)

		db.Exec("DELETE FROM tools WHERE id IN (?, ?) ", firstToolCreated.Id, secondToolCreated.Id)
	})

	t.Run("DeleteTool", func(t *testing.T) {
		tool := testutil.CreateTestTool()
		tool, err := toolRepository.CreateTool(context.Background(), tool)
		require.NoError(t, err)
		require.NotZero(t, tool.Id)

		err = toolRepository.DeleteTool(context.Background(), tool.Slug)
		require.NoError(t, err)
	})

	t.Run("DeleteTool succeeds when tool is not found", func(t *testing.T) {
		err := toolRepository.DeleteTool(context.Background(), "nodejs")
		require.NoError(t, err)
	})

	t.Run("GetToolByIds", func(t *testing.T) {
		firstTool := testutil.CreateTestTool()
		firstToolCreated, err := toolRepository.CreateTool(context.Background(), firstTool)
		require.NoError(t, err)
		require.NotZero(t, firstToolCreated.Id)

		secondTool := testutil.CreateTestDynamicTool(0, domain.CreateToolInput{
			Slug:        "nodejs",
			Name:        "Node.js",
			Category:    "language",
			SubType:     "backend",
			Prolang:     new("JavaScript"),
			ReleaseYear: 2009,
			DevStatus:   "active",
		})
		secondToolCreated, err := toolRepository.CreateTool(context.Background(), secondTool)
		require.NoError(t, err)
		require.NotZero(t, secondToolCreated.Id)

		toolsIds := []int{firstToolCreated.Id, secondToolCreated.Id}
		tools, err := toolRepository.GetToolByIDs(context.Background(), toolsIds)
		require.NoError(t, err)
		require.Equal(t, []domain.Tool{firstToolCreated, secondToolCreated}, tools)
	})

	t.Run("GetToolByIds succeeds when tools are not found", func(t *testing.T) {
		toolsIds := []int{1, 2}
		tools, err := toolRepository.GetToolByIDs(context.Background(), toolsIds)
		require.NoError(t, err)
		require.Equal(t, []domain.Tool{}, tools)
	})
}
