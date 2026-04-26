//go:build integration

package http_test

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"tericcabrel/instech/internal/core"
	"tericcabrel/instech/internal/domain"
	"tericcabrel/instech/internal/repository"
	"tericcabrel/instech/testutil"
	"testing"

	"github.com/stretchr/testify/require"
)

func buildRouter(t *testing.T, db *sql.DB) *core.HTTPRouter {
	t.Helper()

	return &core.HTTPRouter{
		ToolRepository:         repository.NewToolRepository(db),
		RelationshipRepository: repository.NewRelationshipRepository(db),
	}
}

func createTools(t *testing.T, db *sql.DB) (nodejsId int, golangId int, pythonId int) {
	t.Helper()

	nodeTool, err := domain.CreateTool(domain.CreateToolInput{
		Name:        "Node.js",
		Slug:        "nodejs",
		Category:    "language",
		SubType:     "backend",
		Prolang:     "JavaScript",
		ReleaseYear: 2009,
		Devstatus:   "active",
	})
	require.NoError(t, err)

	golangTool, err := domain.CreateTool(domain.CreateToolInput{
		Name:        "Golang",
		Slug:        "golang",
		Category:    "language",
		SubType:     "backend",
		Prolang:     "Go",
		ReleaseYear: 2009,
		Devstatus:   "active",
	})
	require.NoError(t, err)

	pythonTool, err := domain.CreateTool(domain.CreateToolInput{
		Name:        "Python",
		Slug:        "python",
		Category:    "language",
		SubType:     "backend",
		Prolang:     "Python",
		ReleaseYear: 1995,
		Devstatus:   "active",
	})

	toolRepository := repository.NewToolRepository(db)
	createdNodeTool, err := toolRepository.CreateTool(context.Background(), nodeTool)
	require.NoError(t, err)
	createdGolangTool, err := toolRepository.CreateTool(context.Background(), golangTool)
	require.NoError(t, err)
	createdPythonTool, err := toolRepository.CreateTool(context.Background(), pythonTool)
	require.NoError(t, err)

	return createdNodeTool.Id, createdGolangTool.Id, createdPythonTool.Id
}

func createRelationship(t *testing.T, db *sql.DB, fromToolID int, toToolID int, kind string) int {
	t.Helper()

	relationship, err := domain.CreateRelationship(domain.CreateRelationshipInput{
		FromToolId: fromToolID,
		ToToolId:   toToolID,
		Kind:       kind,
		Reason:     "relationship test",
	})
	require.NoError(t, err)

	relationshipRepository := repository.NewRelationshipRepository(db)
	createdRelationship, err := relationshipRepository.CreateRelationship(context.Background(), relationship)
	require.NoError(t, err)

	return createdRelationship.Id
}

func deleteRelationships(t *testing.T, db *sql.DB, relationshipIds ...int) {
	t.Helper()

	relationshipRepository := repository.NewRelationshipRepository(db)
	for _, relationshipId := range relationshipIds {
		err := relationshipRepository.DeleteRelationship(context.Background(), relationshipId)
		require.NoError(t, err)
	}
}

func deleteTools(t *testing.T, db *sql.DB, toolIds ...int) {
	t.Helper()

	toolRepository := repository.NewToolRepository(db)
	tools, err := toolRepository.GetToolByIds(context.Background(), toolIds)
	require.NoError(t, err)
	for _, tool := range tools {
		err := toolRepository.DeleteTool(context.Background(), tool.Slug)
		require.NoError(t, err)
	}
}

func TestRelationshipRouter_Integration(t *testing.T) {
	db := testutil.SetupTestDB(t)
	nodejsId, golangId, pythonId := createTools(t, db)

	t.Cleanup(func() {
		deleteTools(t, db, nodejsId, golangId, pythonId)
	})
	t.Run("[POST] /relationships", func(t *testing.T) {
		router := buildRouter(t, db)

		t.Run("return error 400 when request body is invalid JSON", func(t *testing.T) {
			body := bytes.NewBufferString(`{
				"fromToolId": "golang",
				"toToolId": "nodejs",
				"kind": "alternative_to"
			`)
			req := httptest.NewRequest(http.MethodPost, "/relationships", body)
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			router.Initialize().ServeHTTP(rec, req)

			require.Equal(t, http.StatusBadRequest, rec.Code)

			responseBody, err := io.ReadAll(rec.Body)

			require.NoError(t, err)

			var response map[string]any
			err = json.Unmarshal(responseBody, &response)

			require.NoError(t, err)
			require.Equal(t, "Bad Request", response["message"])
			require.Equal(t, "unexpected EOF", response["details"])
		})

		t.Run("return error 404 when source tool does not exist", func(t *testing.T) {
			body := bytes.NewBufferString(`{
				"fromToolId": "not-found",
				"toToolId": "nodejs",
				"kind": "alternative_to",
				"metadata": { "reason": "same purpose" }
			}`)
			req := httptest.NewRequest(http.MethodPost, "/relationships", body)
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			router.Initialize().ServeHTTP(rec, req)

			require.Equal(t, http.StatusNotFound, rec.Code)

			responseBody, err := io.ReadAll(rec.Body)

			require.NoError(t, err)

			var response map[string]any
			err = json.Unmarshal(responseBody, &response)

			require.NoError(t, err)
			require.Equal(t, "Resource Not Found", response["message"])
			require.Equal(t, map[string]any{"id": "not-found"}, response["details"])
		})

		t.Run("return error 400 when relationship kind is invalid", func(t *testing.T) {
			body := bytes.NewBufferString(`{
				"fromToolId": "golang",
				"toToolId": "nodejs",
				"kind": "invalid_kind",
				"metadata": { "reason": "same purpose" }
			}`)
			req := httptest.NewRequest(http.MethodPost, "/relationships", body)
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			router.Initialize().ServeHTTP(rec, req)

			require.Equal(t, http.StatusBadRequest, rec.Code)

			responseBody, err := io.ReadAll(rec.Body)

			require.NoError(t, err)

			var response map[string]any
			err = json.Unmarshal(responseBody, &response)

			require.NoError(t, err)
			require.Equal(t, "Bad Request", response["message"])
			require.Equal(t, map[string]any{
				"kind":    "invalid_kind",
				"message": "The relationship kind is invalid. Valid kinds are: built_on, inspired_by, alternative_to, replaced_by, used_with",
			}, response["details"])
		})

		t.Run("return error 404 when required fields are missing", func(t *testing.T) {
			body := bytes.NewBufferString(`{
				"kind": "alternative_to",
				"metadata": { "reason": "same purpose" }
			}`)
			req := httptest.NewRequest(http.MethodPost, "/relationships", body)
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			router.Initialize().ServeHTTP(rec, req)

			require.Equal(t, http.StatusNotFound, rec.Code)

			responseBody, err := io.ReadAll(rec.Body)

			require.NoError(t, err)

			var response map[string]any
			err = json.Unmarshal(responseBody, &response)

			require.NoError(t, err)
			require.Equal(t, "Resource Not Found", response["message"])
			require.Equal(t, map[string]any{"id": ""}, response["details"])
		})

		t.Run("create relationship with valid input", func(t *testing.T) {
			body := bytes.NewBufferString(`{
				"fromToolId": "golang",
				"toToolId": "nodejs",
				"kind": "used_with",
				"metadata": { "reason": "same purpose" }
			}`)
			req := httptest.NewRequest(http.MethodPost, "/relationships", body)
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			router.Initialize().ServeHTTP(rec, req)

			require.Equal(t, http.StatusCreated, rec.Code)

			responseBody, err := io.ReadAll(rec.Body)

			require.NoError(t, err)

			var response map[string]any
			err = json.Unmarshal(responseBody, &response)

			require.NoError(t, err)

			require.Equal(t, float64(1), response["id"])
			require.Equal(t, float64(golangId), response["from_tool_id"])
			require.Equal(t, float64(nodejsId), response["to_tool_id"])
			require.Equal(t, "used_with", response["kind"])
			require.Equal(t, map[string]any{"reason": "same purpose"}, response["metadata"])

			deleteRelationships(t, db, int(response["id"].(float64)))
		})
	})

	t.Run("[PUT] /relationships/{id}", func(t *testing.T) {
		router := buildRouter(t, db)
		createdRelationshipId := createRelationship(t, db, nodejsId, golangId, "alternative_to")

		t.Cleanup(func() {
			deleteRelationships(t, db, createdRelationshipId)
		})

		t.Run("return error 400 when relationship ID is invalid", func(t *testing.T) {
			body := bytes.NewBufferString(`{
				"fromToolId": 1,
				"toToolId": 2,
				"kind": "used_with",
				"metadata": { "reason": "often combined" }
			}`)
			req := httptest.NewRequest(http.MethodPut, "/relationships/abc", body)
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			router.Initialize().ServeHTTP(rec, req)

			require.Equal(t, http.StatusBadRequest, rec.Code)

			responseBody, err := io.ReadAll(rec.Body)

			require.NoError(t, err)

			var response map[string]any
			err = json.Unmarshal(responseBody, &response)

			require.NoError(t, err)
			require.Equal(t, "Bad Request", response["message"])
			require.Equal(t, map[string]any{
				"message": "Invalid relationship ID",
			}, response["details"])
		})

		t.Run("return error 400 when request body is invalid JSON", func(t *testing.T) {
			body := bytes.NewBufferString(`{
				"fromToolId": 1,
				"toToolId": 2,
				"kind": "used_with",
			}`)
			req := httptest.NewRequest(http.MethodPut, "/relationships/1", body)
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			router.Initialize().ServeHTTP(rec, req)

			require.Equal(t, http.StatusBadRequest, rec.Code)

			responseBody, err := io.ReadAll(rec.Body)

			require.NoError(t, err)

			var response map[string]any
			err = json.Unmarshal(responseBody, &response)

			require.NoError(t, err)
			require.Equal(t, "Bad Request", response["message"])
			require.Equal(t, "invalid character '}' looking for beginning of object key string", response["details"])
		})

		t.Run("return error 404 when relationship does not exist", func(t *testing.T) {
			body := bytes.NewBufferString(`{
				"fromToolId": 1,
				"toToolId": 2,
				"kind": "used_with",
				"metadata": { "reason": "often combined" }
			}`)
			req := httptest.NewRequest(http.MethodPut, "/relationships/999", body)
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			router.Initialize().ServeHTTP(rec, req)

			require.Equal(t, http.StatusNotFound, rec.Code)

			responseBody, err := io.ReadAll(rec.Body)

			require.NoError(t, err)

			var response map[string]any
			err = json.Unmarshal(responseBody, &response)

			require.NoError(t, err)
			require.Equal(t, "Resource Not Found", response["message"])
			require.Equal(t, map[string]any{"id": "999"}, response["details"])
		})

		t.Run("return error 404 when source tool does not exist", func(t *testing.T) {
			body := bytes.NewBufferString(`{
				"fromToolId": 999,
				"toToolId": 2,
				"kind": "used_with",
				"metadata": { "reason": "often combined" }
			}`)
			req := httptest.NewRequest(http.MethodPut, "/relationships/"+strconv.Itoa(createdRelationshipId), body)
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			router.Initialize().ServeHTTP(rec, req)

			require.Equal(t, http.StatusNotFound, rec.Code)

			responseBody, err := io.ReadAll(rec.Body)

			require.NoError(t, err)

			var response map[string]any
			err = json.Unmarshal(responseBody, &response)

			require.NoError(t, err)
			require.Equal(t, "Resource Not Found", response["message"])
			require.Equal(t, map[string]any{"id": "999"}, response["details"])
		})

		t.Run("return error 404 when target tool does not exist", func(t *testing.T) {
			body := bytes.NewBufferString(`{
				"fromToolId": 1,
				"toToolId": 999,
				"kind": "used_with",
				"metadata": { "reason": "often combined" }
			}`)
			req := httptest.NewRequest(http.MethodPut, "/relationships/"+strconv.Itoa(createdRelationshipId), body)
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			router.Initialize().ServeHTTP(rec, req)

			require.Equal(t, http.StatusNotFound, rec.Code)

			responseBody, err := io.ReadAll(rec.Body)

			require.NoError(t, err)

			var response map[string]any
			err = json.Unmarshal(responseBody, &response)

			require.NoError(t, err)
			require.Equal(t, "Resource Not Found", response["message"])
			require.Equal(t, map[string]any{"id": "999"}, response["details"])
		})

		t.Run("return error 400 when kind is invalid", func(t *testing.T) {
			body := bytes.NewBufferString(`{
				"fromToolId": 1,
				"toToolId": 2,
				"kind": "invalid_kind",
				"metadata": { "reason": "often combined" }
			}`)
			req := httptest.NewRequest(http.MethodPut, "/relationships/"+strconv.Itoa(createdRelationshipId), body)
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			router.Initialize().ServeHTTP(rec, req)

			require.Equal(t, http.StatusBadRequest, rec.Code)

			responseBody, err := io.ReadAll(rec.Body)

			require.NoError(t, err)

			var response map[string]any
			err = json.Unmarshal(responseBody, &response)

			require.NoError(t, err)
			require.Equal(t, "Bad Request", response["message"])
			require.Equal(t, map[string]any{
				"kind":    "invalid_kind",
				"message": "The relationship kind is invalid. Valid kinds are: built_on, inspired_by, alternative_to, replaced_by, used_with",
			}, response["details"])
		})

		t.Run("update relationship with valid input", func(t *testing.T) {
			body := `{
				"fromToolId": ` + strconv.Itoa(pythonId) + `,
				"toToolId": ` + strconv.Itoa(golangId) + `,
				"kind": "used_with",
				"metadata": { "reason": "often combined" }
			}`

			req := httptest.NewRequest(http.MethodPut, "/relationships/"+strconv.Itoa(createdRelationshipId), bytes.NewBufferString(body))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			router.Initialize().ServeHTTP(rec, req)

			require.Equal(t, http.StatusOK, rec.Code)

			responseBody, err := io.ReadAll(rec.Body)

			require.NoError(t, err)

			var response map[string]any
			err = json.Unmarshal(responseBody, &response)

			require.NoError(t, err)

			require.Equal(t, float64(createdRelationshipId), response["id"])
			require.Equal(t, float64(pythonId), response["from_tool_id"])
			require.Equal(t, float64(golangId), response["to_tool_id"])
			require.Equal(t, "used_with", response["kind"])
			require.Equal(t, map[string]any{"reason": "often combined"}, response["metadata"])
		})
	})

	t.Run("[DELETE] /relationships/{id}", func(t *testing.T) {
		router := buildRouter(t, db)
		createdRelationshipId := createRelationship(t, db, nodejsId, golangId, "alternative_to")

		t.Run("return error 400 when relationship ID is invalid", func(t *testing.T) {
			req := httptest.NewRequest(http.MethodDelete, "/relationships/invalid", nil)
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			router.Initialize().ServeHTTP(rec, req)

			require.Equal(t, http.StatusBadRequest, rec.Code)

			responseBody, err := io.ReadAll(rec.Body)

			require.NoError(t, err)

			var response map[string]any
			err = json.Unmarshal(responseBody, &response)

			require.NoError(t, err)
			require.Equal(t, "Bad Request", response["message"])
			require.Equal(t, map[string]any{
				"message": "Invalid relationship Id",
			}, response["details"])
		})

		t.Run("return error 404 when relationship does not exist", func(t *testing.T) {
			req := httptest.NewRequest(http.MethodDelete, "/relationships/999", nil)
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			router.Initialize().ServeHTTP(rec, req)

			require.Equal(t, http.StatusNotFound, rec.Code)

			responseBody, err := io.ReadAll(rec.Body)

			require.NoError(t, err)

			var response map[string]any
			err = json.Unmarshal(responseBody, &response)

			require.NoError(t, err)
			require.Equal(t, "Resource Not Found", response["message"])
			require.Equal(t, map[string]any{"id": "999"}, response["details"])
		})

		t.Run("delete relationship successfully", func(t *testing.T) {
			req := httptest.NewRequest(http.MethodDelete, "/relationships/"+strconv.Itoa(createdRelationshipId), nil)
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			router.Initialize().ServeHTTP(rec, req)

			require.Equal(t, http.StatusNoContent, rec.Code)
		})
	})

	t.Run("[GET] /relationships/query", func(t *testing.T) {
		router := buildRouter(t, db)

		relId1 := createRelationship(t, db, nodejsId, golangId, "alternative_to")
		relId2 := createRelationship(t, db, golangId, pythonId, "used_with")
		relId3 := createRelationship(t, db, pythonId, nodejsId, "alternative_to")

		t.Cleanup(func() {
			deleteRelationships(t, db, relId1, relId2, relId3)
		})

		t.Run("return error 400 when tool_id query param is invalid", func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/relationships/query?tool_id=abc", nil)
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			router.Initialize().ServeHTTP(rec, req)

			require.Equal(t, http.StatusBadRequest, rec.Code)

			responseBody, err := io.ReadAll(rec.Body)
			require.NoError(t, err)

			var response map[string]any
			err = json.Unmarshal(responseBody, &response)

			require.NoError(t, err)
			require.Equal(t, "Bad Request", response["message"])
			require.Equal(t, map[string]any{
				"message": "Invalid tool Id",
			}, response["details"])
		})

		t.Run("return error 400 when kind query param is invalid", func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/relationships/query?kind=invalid_kind", nil)
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			router.Initialize().ServeHTTP(rec, req)

			require.Equal(t, http.StatusBadRequest, rec.Code)

			responseBody, err := io.ReadAll(rec.Body)
			require.NoError(t, err)

			var response map[string]any
			err = json.Unmarshal(responseBody, &response)

			require.NoError(t, err)
			require.Equal(t, "Bad Request", response["message"])
			require.Equal(t, map[string]any{
				"kind":    "invalid_kind",
				"message": "Invalid kind",
			}, response["details"])
		})

		t.Run("return error 400 when cursor query param is invalid", func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/relationships/query?cursor=abc", nil)
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			router.Initialize().ServeHTTP(rec, req)

			require.Equal(t, http.StatusBadRequest, rec.Code)

			responseBody, err := io.ReadAll(rec.Body)

			require.NoError(t, err)

			var response map[string]any
			err = json.Unmarshal(responseBody, &response)

			require.NoError(t, err)

			require.Equal(t, "Bad Request", response["message"])
			require.Equal(t, map[string]any{
				"message": "Invalid cursor",
			}, response["details"])
		})

		t.Run("return error 400 when limit query param is invalid", func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/relationships/query?limit=abc", nil)
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			router.Initialize().ServeHTTP(rec, req)

			require.Equal(t, http.StatusBadRequest, rec.Code)

			responseBody, err := io.ReadAll(rec.Body)

			require.NoError(t, err)

			var response map[string]any
			err = json.Unmarshal(responseBody, &response)

			require.NoError(t, err)
			require.Equal(t, "Bad Request", response["message"])
			require.Equal(t, map[string]any{
				"message": "Invalid limit",
			}, response["details"])
		})

		t.Run("return empty result when tool_id does not match any tool", func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/relationships/query?tool_id=999", nil)
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			router.Initialize().ServeHTTP(rec, req)

			require.Equal(t, http.StatusOK, rec.Code)
		})

		t.Run("return relationships successfully with valid kind and limit query params", func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/relationships/query?kind=alternative_to&limit=10", nil)
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			router.Initialize().ServeHTTP(rec, req)

			require.Equal(t, http.StatusOK, rec.Code)

			responseBody, err := io.ReadAll(rec.Body)
			require.NoError(t, err)

			var response map[string]any
			err = json.Unmarshal(responseBody, &response)
			require.NoError(t, err)

			data, ok := response["data"].([]any)
			require.True(t, ok)
			require.Len(t, data, 2)

			meta, ok := response["meta"].(map[string]any)
			require.True(t, ok)
			require.Equal(t, float64(2), meta["items_count"])
			require.Equal(t, float64(2), meta["total_count"])
			require.Equal(t, float64(-1), meta["next_cursor"])
		})

		t.Run("return relationships successfully with valid tool_id and kind query params", func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/relationships/query?tool_id=1&kind=alternative_to", nil)
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			router.Initialize().ServeHTTP(rec, req)

			require.Equal(t, http.StatusOK, rec.Code)

			responseBody, err := io.ReadAll(rec.Body)
			require.NoError(t, err)

			var response map[string]any
			err = json.Unmarshal(responseBody, &response)
			require.NoError(t, err)

			data, ok := response["data"].([]any)
			require.True(t, ok)
			require.Len(t, data, 2)

			meta, ok := response["meta"].(map[string]any)
			require.True(t, ok)
			require.Equal(t, float64(2), meta["items_count"])
			require.Equal(t, float64(2), meta["total_count"])
			require.Equal(t, float64(-1), meta["next_cursor"])
		})

		t.Run("return relationships successfully with a valid next cursor in the response", func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/relationships/query?limit=2", nil)
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			router.Initialize().ServeHTTP(rec, req)

			require.Equal(t, http.StatusOK, rec.Code)

			responseBody, err := io.ReadAll(rec.Body)
			require.NoError(t, err)

			var response map[string]any
			err = json.Unmarshal(responseBody, &response)
			require.NoError(t, err)

			data, ok := response["data"].([]any)
			require.True(t, ok)
			require.Len(t, data, 2)

			meta, ok := response["meta"].(map[string]any)
			require.True(t, ok)
			require.Equal(t, float64(2), meta["items_count"])
			require.Equal(t, float64(3), meta["total_count"])
			require.NotEqual(t, float64(-1), meta["next_cursor"])
		})
	})
}
