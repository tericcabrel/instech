//go:build integration

package http_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"tericcabrel/instech/internal/core"
	"tericcabrel/instech/internal/repository"
	"tericcabrel/instech/testutil"
	"testing"
	"time"

	"database/sql"

	"github.com/stretchr/testify/require"
)

func createTool(t *testing.T, router *core.HTTPRouter) map[string]any {
	t.Helper()

	bodyString := `{
		"name": "Golang",
		"slug": "golang",
		"category": "language",
		"subType": "backend",
		"prolang": "Go",
		"releaseYear": 2009,
		"devStatus": "active",
		"details": "Test Details",
		"usecases": ["api", "backend"],
		"tags": ["rest api", "server side", "cli"],
		"website": "https://go.dev",
		"github": "https://github.com/golang/go"
	}`
	body := bytes.NewBufferString(bodyString)
	req := httptest.NewRequest(http.MethodPost, "/tools", body)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.Initialize().ServeHTTP(rec, req)

	require.Equal(t, http.StatusCreated, rec.Code)

	responseBody, err := io.ReadAll(rec.Body)
	require.NoError(t, err)

	var response map[string]any
	err = json.Unmarshal(responseBody, &response)
	require.NoError(t, err)

	return response
}

func populateAlternatives(t *testing.T, db *sql.DB) {
	t.Helper()

	stmts := []string{
		"INSERT INTO tools (id, name, slug, category, sub_type, prolang, release_year, devstatus, details, use_cases, tags, website, github) VALUES (11, 'Golang', 'golang', 'language', 'backend', 'Go', 2009, 'active', 'Golang details', '[\"api\", \"backend\"]', '[\"rest api\", \"server side\", \"cli\"]', 'https://golang.org', 'https://github.com/golang/go')",
		"INSERT INTO tools (id, name, slug, category, sub_type, prolang, release_year, devstatus, details, use_cases, tags, website, github) VALUES (12, 'Node.js', 'nodejs', 'language', 'fullstack', 'JavaScript', 2009, 'active', 'Node.js details', '[\"api\", \"frontend\", \"fullstack\"]', '[\"web\", \"api\", \"frontend\"]', 'https://nodejs.org', 'https://github.com/nodejs/node')",
		"INSERT INTO tools (id, name, slug, category, sub_type, prolang, release_year, devstatus, details, use_cases, tags, website, github) VALUES (13, 'Python', 'python', 'language', 'backend', 'Python', 2009, 'active', 'Python details', '[\"api\", \"backend\"]', '[\"rest api\", \"server side\", \"cli\"]', 'https://python.org', 'https://github.com/python/python')",
		"INSERT INTO relationships (id, from_tool_id, to_tool_id, kind, metadata) VALUES (1, 11, 12, 'alternative_to', '{\"reason\": \"test relationship 1\"}')",
		"INSERT INTO relationships (id, from_tool_id, to_tool_id, kind, metadata) VALUES (2, 12, 11, 'alternative_to', '{\"reason\": \"test relationship 2\"}')",
		"INSERT INTO relationships (id, from_tool_id, to_tool_id, kind, metadata) VALUES (3, 11, 13, 'alternative_to', '{\"reason\": \"test relationship 3\"}')",
	}
	for _, stmt := range stmts {
		_, err := db.Exec(stmt)
		fmt.Printf("Error: %v", err)
		require.NoError(t, err)
	}
}

func cleanAlternatives(t *testing.T, db *sql.DB) {
	t.Helper()
	stmts := []string{
		"DELETE FROM tools WHERE id IN (11, 12, 13)",
		"DELETE FROM relationships WHERE id IN (1, 2, 3)",
	}
	for _, stmt := range stmts {
		_, err := db.Exec(stmt)
		require.NoError(t, err)
	}
}

func TestToolRouter_CreateTool(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer db.Close()

	toolRepository := repository.NewToolRepository(db)
	relationshipRepository := repository.NewRelationshipRepository(db)

	t.Run("[POST] /tools", func(t *testing.T) {
		t.Run("return error 400 when the request body is an invalid JSON", func(t *testing.T) {
			router := core.HTTPRouter{
				ToolRepository:         toolRepository,
				RelationshipRepository: relationshipRepository,
			}
			h := router.Initialize()
			body := bytes.NewBufferString(`{
				"name": "Golang",
				"slug": "golang",
				"category": "invalid",
			}`)
			req := httptest.NewRequest(http.MethodPost, "/tools", body)
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			h.ServeHTTP(rec, req)
			require.Equal(t, http.StatusBadRequest, rec.Code)

			responseBody, err := io.ReadAll(rec.Body)
			require.NoError(t, err)

			var response map[string]any
			err = json.Unmarshal(responseBody, &response)
			require.NoError(t, err)

			require.Equal(t, "Bad Request", response["message"])
			require.Equal(t, "invalid character '}' looking for beginning of object key string", response["details"])
		})

		t.Run("return error 400 when the category is invalid", func(t *testing.T) {
			router := core.HTTPRouter{
				ToolRepository:         toolRepository,
				RelationshipRepository: relationshipRepository,
			}
			h := router.Initialize()
			body := bytes.NewBufferString(`{
				"name": "Golang",
				"slug": "golang",
				"category": "invalid"
			}`)
			req := httptest.NewRequest(http.MethodPost, "/tools", body)
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			h.ServeHTTP(rec, req)

			require.Equal(t, http.StatusBadRequest, rec.Code)

			responseBody, err := io.ReadAll(rec.Body)
			require.NoError(t, err)

			var response map[string]any
			err = json.Unmarshal(responseBody, &response)

			require.NoError(t, err)
			require.Equal(t, "Bad Request", response["message"])
			expectedDetails := map[string]interface{}{
				"Category": "invalid",
				"Message":  "The tool category is invalid. Valid categories are: language, framework, library",
			}
			require.Equal(t, expectedDetails, response["details"])
		})

		t.Run("return error 400 when the sub type is invalid", func(t *testing.T) {
			router := core.HTTPRouter{
				ToolRepository:         toolRepository,
				RelationshipRepository: relationshipRepository,
			}
			h := router.Initialize()
			body := bytes.NewBufferString(`{
			"name": "Golang",
			"slug": "golang",
			"category": "language",
			"subType": "invalid"
		}`)
			req := httptest.NewRequest(http.MethodPost, "/tools", body)
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			h.ServeHTTP(rec, req)
			require.Equal(t, http.StatusBadRequest, rec.Code)

			responseBody, err := io.ReadAll(rec.Body)
			require.NoError(t, err)

			var response map[string]any
			err = json.Unmarshal(responseBody, &response)
			require.NoError(t, err)
			require.Equal(t, "Bad Request", response["message"])
			expectedDetails := map[string]interface{}{
				"SubType": "invalid",
				"Message": "The tool sub type is invalid. Valid sub types are: backend, frontend, fullstack, mobile, desktop, game, other",
			}
			require.Equal(t, expectedDetails, response["details"])
		})

		t.Run("return error 400 when the dev status is invalid", func(t *testing.T) {
			router := core.HTTPRouter{
				ToolRepository:         toolRepository,
				RelationshipRepository: relationshipRepository,
			}
			h := router.Initialize()
			body := bytes.NewBufferString(`{
			"name": "Golang",
			"slug": "golang",
			"category": "language",
			"subType": "backend",
			"devStatus": "invalid"
		}`)
			req := httptest.NewRequest(http.MethodPost, "/tools", body)
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			h.ServeHTTP(rec, req)
			require.Equal(t, http.StatusBadRequest, rec.Code)

			responseBody, err := io.ReadAll(rec.Body)
			require.NoError(t, err)

			var response map[string]any
			err = json.Unmarshal(responseBody, &response)
			require.NoError(t, err)
			require.Equal(t, "Bad Request", response["message"])
			expectedDetails := map[string]interface{}{
				"Devstatus": "invalid",
				"Message":   "The tool dev status is invalid. Valid dev statuses are: active, deprecated",
			}
			require.Equal(t, expectedDetails, response["details"])
		})

		t.Run("return error 422 when the fields are invalid", func(t *testing.T) {
			router := core.HTTPRouter{
				ToolRepository:         toolRepository,
				RelationshipRepository: relationshipRepository,
			}
			h := router.Initialize()
			body := bytes.NewBufferString(`{
			"category": "language",
			"subType": "backend",
			"devStatus": "active",
			"website": "invalid",
			"github": "invalid"
		}`)
			req := httptest.NewRequest(http.MethodPost, "/tools", body)
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			h.ServeHTTP(rec, req)
			require.Equal(t, http.StatusUnprocessableEntity, rec.Code)

			responseBody, err := io.ReadAll(rec.Body)
			require.NoError(t, err)

			var response map[string]any
			err = json.Unmarshal(responseBody, &response)

			require.NoError(t, err)
			require.Equal(t, "Unprocessable Entity", response["message"])

			year := strconv.Itoa(time.Now().Year())
			expectedInvalidFields := map[string]interface{}{
				"Name":        "The tool name is required",
				"Prolang":     "The tool programming language is required",
				"Slug":        "The tool slug is required",
				"ReleaseYear": "The tool release year is invalid. Valid release years are between 1940 and " + year,
				"Website":     "The tool website is invalid. Valid websites must be a valid URL",
				"Github":      "The tool github is invalid. Valid github must be a valid URL",
			}
			expectedDetails := map[string]interface{}{
				"Fields": expectedInvalidFields,
			}
			require.Equal(t, expectedDetails, response["details"])
		})

		t.Run("Create tool with valid input", func(t *testing.T) {
			router := core.HTTPRouter{
				ToolRepository:         toolRepository,
				RelationshipRepository: relationshipRepository,
			}
			h := router.Initialize()
			body := bytes.NewBufferString(`{
				"name": "Golang",
				"slug": "golang",
				"category": "language",
				"subType": "backend",
				"prolang": "Go",
				"releaseYear": 2009,
				"devStatus": "active",
				"details": "Test Details",
				"usecases": ["api", "backend"],
				"tags": ["rest api", "server side", "cli"],
				"website": "https://go.dev",
				"github": "https://github.com/golang/go"
			}`)
			req := httptest.NewRequest(http.MethodPost, "/tools", body)
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			h.ServeHTTP(rec, req)

			require.Equal(t, http.StatusCreated, rec.Code)

			responseBody, err := io.ReadAll(rec.Body)
			require.NoError(t, err)

			var response map[string]any
			err = json.Unmarshal(responseBody, &response)
			require.NoError(t, err)

			require.Equal(t, float64(1), response["id"])
			require.Equal(t, "Golang", response["name"])
			require.Equal(t, "golang", response["slug"])
			require.Equal(t, "language", response["category"])
			require.Equal(t, "backend", response["sub_type"])
			require.Equal(t, "active", response["devstatus"])
			require.Equal(t, "Go", response["prolang"])
			require.Equal(t, float64(2009), response["release_year"])
			require.Equal(t, "Test Details", response["details"])
			require.Equal(t, []interface{}{"api", "backend"}, response["use_cases"])
			require.Equal(t, []interface{}{"rest api", "server side", "cli"}, response["tags"])
			require.Equal(t, "https://go.dev", response["website"])
			require.Equal(t, "https://github.com/golang/go", response["github"])
			require.Contains(t, response, "created_at")
			require.Contains(t, response, "updated_at")

			_, errDb := db.Exec("DELETE FROM tools WHERE slug = ?", response["slug"])
			require.NoError(t, errDb)
		})

		t.Run("return error 400 when the tool already exists", func(t *testing.T) {
			router := core.HTTPRouter{
				ToolRepository:         toolRepository,
				RelationshipRepository: relationshipRepository,
			}
			h := router.Initialize()
			bodyString := `{
				"name": "Node.js",
				"slug": "nodejs",
				"category": "language",
				"subType": "fullstack",
				"prolang": "JavaScript",
				"releaseYear": 2009,
				"devStatus": "active",
				"details": "Node.js details",
				"usecases": ["api", "frontend"],
				"tags": [],
				"website": "https://nodejs.org",
				"github": "https://github.com/nodejs/node"
			}`
			body := bytes.NewBufferString(bodyString)
			reqOne := httptest.NewRequest(http.MethodPost, "/tools", body)
			reqOne.Header.Set("Content-Type", "application/json")
			recOne := httptest.NewRecorder()
			h.ServeHTTP(recOne, reqOne)
			require.Equal(t, http.StatusCreated, recOne.Code)

			// Send the same request again
			bodyTwo := bytes.NewBufferString(bodyString)

			reqTwo := httptest.NewRequest(http.MethodPost, "/tools", bodyTwo)
			reqTwo.Header.Set("Content-Type", "application/json")
			recTwo := httptest.NewRecorder()
			h.ServeHTTP(recTwo, reqTwo)
			require.Equal(t, http.StatusBadRequest, recTwo.Code)

			responseBody, err := io.ReadAll(recTwo.Body)
			require.NoError(t, err)

			var response map[string]any
			err = json.Unmarshal(responseBody, &response)

			require.NoError(t, err)
			require.Equal(t, "Bad Request", response["message"])

			expectedDetails := map[string]interface{}{
				"Id":      "nodejs",
				"Message": "The tool already exists",
			}
			require.Equal(t, expectedDetails, response["details"])

			_, errDb := db.Exec("DELETE FROM tools WHERE slug = ?", "nodejs")
			require.NoError(t, errDb)
		})
	})

	t.Run("[PUT] /tools/{id}", func(t *testing.T) {
		router := &core.HTTPRouter{
			ToolRepository:         toolRepository,
			RelationshipRepository: relationshipRepository,
		}
		tool := createTool(t, router)

		t.Run("return error 400 when the request body is an invalid JSON", func(t *testing.T) {
			router := core.HTTPRouter{
				ToolRepository:         toolRepository,
				RelationshipRepository: relationshipRepository,
			}
			h := router.Initialize()
			body := bytes.NewBufferString(`{
				"name": "Golang",
				"slug": "golang",
				"category": "language",
				"subType": "backend",
				"devStatus": "activ
			}`)
			req := httptest.NewRequest(http.MethodPut, "/tools/golang", body)
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			h.ServeHTTP(rec, req)
			require.Equal(t, http.StatusBadRequest, rec.Code)

			responseBody, err := io.ReadAll(rec.Body)
			require.NoError(t, err)

			var response map[string]any
			err = json.Unmarshal(responseBody, &response)
			require.NoError(t, err)
			require.Equal(t, "Bad Request", response["message"])
			require.Equal(t, "invalid character '\\n' in string literal", response["details"])
		})

		t.Run("return error 400 when the category is invalid", func(t *testing.T) {
			router := core.HTTPRouter{
				ToolRepository:         toolRepository,
				RelationshipRepository: relationshipRepository,
			}
			h := router.Initialize()
			body := bytes.NewBufferString(`{
				"name": "Golang",
				"slug": "golang",
				"category": "invalid"
			}`)
			req := httptest.NewRequest(http.MethodPut, "/tools/golang", body)
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			h.ServeHTTP(rec, req)
			require.Equal(t, http.StatusBadRequest, rec.Code)

			responseBody, err := io.ReadAll(rec.Body)
			require.NoError(t, err)

			var response map[string]any
			err = json.Unmarshal(responseBody, &response)
			require.NoError(t, err)
			require.Equal(t, "Bad Request", response["message"])
			expectedDetails := map[string]interface{}{
				"Category": "invalid",
				"Message":  "The tool category is invalid. Valid categories are: language, framework, library",
			}
			require.Equal(t, expectedDetails, response["details"])
		})

		t.Run("return error 400 when the sub type is invalid", func(t *testing.T) {
			router := core.HTTPRouter{
				ToolRepository:         toolRepository,
				RelationshipRepository: relationshipRepository,
			}
			h := router.Initialize()
			body := bytes.NewBufferString(`{
				"name": "Golang",
				"slug": "golang",
				"category": "language",
				"subType": "invalid"
			}`)
			req := httptest.NewRequest(http.MethodPut, "/tools/golang", body)
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			h.ServeHTTP(rec, req)
			require.Equal(t, http.StatusBadRequest, rec.Code)

			responseBody, err := io.ReadAll(rec.Body)
			require.NoError(t, err)

			var response map[string]any
			err = json.Unmarshal(responseBody, &response)
			require.NoError(t, err)
			require.Equal(t, "Bad Request", response["message"])
			expectedDetails := map[string]interface{}{
				"SubType": "invalid",
				"Message": "The tool sub type is invalid. Valid sub types are: backend, frontend, fullstack, mobile, desktop, game, other",
			}
			require.Equal(t, expectedDetails, response["details"])
		})

		t.Run("return error 400 when the dev status is invalid", func(t *testing.T) {
			router := core.HTTPRouter{
				ToolRepository:         toolRepository,
				RelationshipRepository: relationshipRepository,
			}
			h := router.Initialize()
			body := bytes.NewBufferString(`{
				"name": "Golang",
				"slug": "golang",
				"category": "language",
				"subType": "backend",
				"devStatus": "invalid"
			}`)
			req := httptest.NewRequest(http.MethodPut, "/tools/golang", body)
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			h.ServeHTTP(rec, req)
			require.Equal(t, http.StatusBadRequest, rec.Code)
		})

		t.Run("return error 422 when the fields are invalid", func(t *testing.T) {
			router := core.HTTPRouter{
				ToolRepository:         toolRepository,
				RelationshipRepository: relationshipRepository,
			}
			h := router.Initialize()
			body := bytes.NewBufferString(`{
				"category": "language",
				"subType": "backend",
				"devStatus": "active",
				"website": "invalid",
				"github": "invalid"
			}`)
			req := httptest.NewRequest(http.MethodPut, "/tools/golang", body)
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			h.ServeHTTP(rec, req)
			require.Equal(t, http.StatusUnprocessableEntity, rec.Code)

			responseBody, err := io.ReadAll(rec.Body)
			require.NoError(t, err)

			var response map[string]any
			err = json.Unmarshal(responseBody, &response)
			require.NoError(t, err)
			require.Equal(t, "Unprocessable Entity", response["message"])

			year := strconv.Itoa(time.Now().Year())
			expectedInvalidFields := map[string]interface{}{
				"Name":        "The tool name is required",
				"Prolang":     "The tool programming language is required",
				"ReleaseYear": "The tool release year is invalid. Valid release years are between 1940 and " + year,
				"Website":     "The tool website is invalid. Valid websites must be a valid URL",
				"Github":      "The tool github is invalid. Valid github must be a valid URL",
			}
			expectedDetails := map[string]interface{}{
				"Fields": expectedInvalidFields,
			}
			require.Equal(t, expectedDetails, response["details"])
		})

		t.Run("Update tool with valid input", func(t *testing.T) {
			router := core.HTTPRouter{
				ToolRepository:         toolRepository,
				RelationshipRepository: relationshipRepository,
			}
			h := router.Initialize()

			bodyString := `{
				"name": "Node.js",
				"slug": "nodejs",
				"category": "language",
				"subType": "fullstack",
				"prolang": "JavaScript",
				"releaseYear": 2009,
				"devStatus": "active",
				"details": "Node.js details",
				"usecases": ["api", "frontend"],
				"tags": ["web"],
				"website": "https://nodejs.org",
				"github": "https://github.com/nodejs/node"
			}`
			createBody := bytes.NewBufferString(bodyString)
			reqOne := httptest.NewRequest(http.MethodPost, "/tools", createBody)
			reqOne.Header.Set("Content-Type", "application/json")
			recOne := httptest.NewRecorder()
			h.ServeHTTP(recOne, reqOne)

			respBody, err := io.ReadAll(recOne.Body)
			require.NoError(t, err)

			var resp map[string]any
			err = json.Unmarshal(respBody, &resp)
			fmt.Printf("Response: %+v", resp)
			require.Equal(t, http.StatusCreated, recOne.Code)

			updateBody := bytes.NewBufferString(`{
				"name": "Node JS",
				"slug": "golang",
				"category": "language",
				"subType": "backend",
				"prolang": "JS",
				"releaseYear": 2009,
				"devStatus": "deprecated",
				"details": "Node JS details",
				"usecases": ["api", "frontend", "fullstack"],
				"tags": ["web", "api", "frontend"],
				"website": "https://nodejs-updated.org",
				"github": "https://github.com/nodejs/node-updated"
			}`)
			req := httptest.NewRequest(http.MethodPut, "/tools/nodejs", updateBody)
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			h.ServeHTTP(rec, req)
			require.Equal(t, http.StatusOK, rec.Code)

			responseBody, err := io.ReadAll(rec.Body)
			require.NoError(t, err)

			var response map[string]any
			err = json.Unmarshal(responseBody, &response)
			require.NoError(t, err)

			require.Equal(t, float64(4), response["id"])
			require.Equal(t, "Node JS", response["name"])
			require.Equal(t, "nodejs", response["slug"]) // slug should not be updated
			require.Equal(t, "language", response["category"])
			require.Equal(t, "backend", response["sub_type"])
			require.Equal(t, "deprecated", response["devstatus"])
			require.Equal(t, "JS", response["prolang"])
			require.Equal(t, float64(2009), response["release_year"])
			require.Equal(t, "Node JS details", response["details"])
			require.Equal(t, []interface{}{"api", "frontend", "fullstack"}, response["use_cases"])
			require.Equal(t, []interface{}{"web", "api", "frontend"}, response["tags"])
			require.Equal(t, "https://nodejs-updated.org", response["website"])
			require.Equal(t, "https://github.com/nodejs/node-updated", response["github"])
			require.Contains(t, response, "created_at")
			require.Contains(t, response, "updated_at")
		})

		_, errDb := db.Exec("DELETE FROM tools WHERE slug = ?", tool["slug"])
		require.NoError(t, errDb)

		_, errDb = db.Exec("DELETE FROM tools WHERE slug = ?", "nodejs")
		require.NoError(t, errDb)
	})

	t.Run("[DELETE] /tools/{id}", func(t *testing.T) {
		router := &core.HTTPRouter{
			ToolRepository:         toolRepository,
			RelationshipRepository: relationshipRepository,
		}
		createTool(t, router)

		t.Run("Do not nothing when the tool is not found", func(t *testing.T) {
			router := &core.HTTPRouter{
				ToolRepository:         toolRepository,
				RelationshipRepository: relationshipRepository,
			}
			h := router.Initialize()
			req := httptest.NewRequest(http.MethodDelete, "/tools/not-found", nil)
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			h.ServeHTTP(rec, req)
			require.Equal(t, http.StatusNoContent, rec.Code)
		})

		t.Run("delete tool successfully", func(t *testing.T) {
			router := &core.HTTPRouter{
				ToolRepository:         toolRepository,
				RelationshipRepository: relationshipRepository,
			}
			h := router.Initialize()
			req := httptest.NewRequest(http.MethodDelete, "/tools/golang", nil)
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			h.ServeHTTP(rec, req)
			require.Equal(t, http.StatusNoContent, rec.Code)
		})
	})

	t.Run("[GET] /tools/{id}", func(t *testing.T) {
		router := &core.HTTPRouter{
			ToolRepository:         toolRepository,
			RelationshipRepository: relationshipRepository,
		}
		createTool(t, router)
		t.Run("return tool successfully", func(t *testing.T) {
			router := &core.HTTPRouter{
				ToolRepository:         toolRepository,
				RelationshipRepository: relationshipRepository,
			}
			h := router.Initialize()
			req := httptest.NewRequest(http.MethodGet, "/tools/golang", nil)
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			h.ServeHTTP(rec, req)
			require.Equal(t, http.StatusOK, rec.Code)

			responseBody, err := io.ReadAll(rec.Body)
			require.NoError(t, err)

			var response map[string]any
			err = json.Unmarshal(responseBody, &response)
			require.NoError(t, err)
			require.Equal(t, "Golang", response["name"])
			require.Equal(t, "golang", response["slug"])
			require.Equal(t, "language", response["category"])
			require.Equal(t, "backend", response["sub_type"])
			require.Equal(t, "Go", response["prolang"])
			require.Equal(t, float64(2009), response["release_year"])
			require.Equal(t, "active", response["devstatus"])
			require.Equal(t, "Test Details", response["details"])
			require.Equal(t, []interface{}{"api", "backend"}, response["use_cases"])
			require.Equal(t, []interface{}{"rest api", "server side", "cli"}, response["tags"])
			require.Equal(t, "https://go.dev", response["website"])
			require.Equal(t, "https://github.com/golang/go", response["github"])
			require.Contains(t, response, "created_at")
			require.Contains(t, response, "updated_at")

			_, errDb := db.Exec("DELETE FROM tools WHERE slug = ?", "golang")
			require.NoError(t, errDb)
		})

		t.Run("return error 404 when the tool is not found", func(t *testing.T) {
			router := &core.HTTPRouter{
				ToolRepository:         toolRepository,
				RelationshipRepository: relationshipRepository,
			}
			h := router.Initialize()
			req := httptest.NewRequest(http.MethodGet, "/tools/not-found", nil)
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			h.ServeHTTP(rec, req)
			require.Equal(t, http.StatusNotFound, rec.Code)
		})
	})

	t.Run("[GET] /tools/{id}/alternatives", func(t *testing.T) {
		populateAlternatives(t, db)

		t.Cleanup(func() {
			cleanAlternatives(t, db)
		})

		t.Run("return tool alternatives successfully", func(t *testing.T) {
			router := &core.HTTPRouter{
				ToolRepository:         toolRepository,
				RelationshipRepository: relationshipRepository,
			}
			h := router.Initialize()
			req := httptest.NewRequest(http.MethodGet, "/tools/golang/alternatives", nil)
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			h.ServeHTTP(rec, req)
			require.Equal(t, http.StatusOK, rec.Code)

			responseBody, err := io.ReadAll(rec.Body)
			require.NoError(t, err)

			var response []map[string]any
			err = json.Unmarshal(responseBody, &response)
			require.NoError(t, err)

			fmt.Printf("Response: %+v", response)

			require.Equal(t, 2, len(response))

			require.Equal(t, "nodejs", response[0]["id"])
			require.Equal(t, "Node.js", response[0]["name"])
			require.Equal(t, "language", response[0]["category"])
			require.Equal(t, "fullstack", response[0]["sub_type"])
			require.Equal(t, "JavaScript", response[0]["prolang"])
			require.Equal(t, float64(2009), response[0]["release_year"])
			require.Equal(t, "active", response[0]["dev_status"])
			require.Equal(t, "Node.js details", response[0]["details"])
			require.Equal(t, []interface{}{"api", "frontend", "fullstack"}, response[0]["use_cases"])
			require.Equal(t, []interface{}{"web", "api", "frontend"}, response[0]["tags"])
			require.Equal(t, "https://nodejs.org", response[0]["website"])
			require.Equal(t, "https://github.com/nodejs/node", response[0]["github"])
			require.Equal(t, map[string]interface{}{"reason": "test relationship 1"}, response[0]["metadata"])

			require.Equal(t, "python", response[1]["id"])
			require.Equal(t, "Python", response[1]["name"])
			require.Equal(t, "language", response[1]["category"])
			require.Equal(t, "backend", response[1]["sub_type"])
			require.Equal(t, "Python", response[1]["prolang"])
			require.Equal(t, float64(2009), response[1]["release_year"])
			require.Equal(t, "active", response[1]["dev_status"])
			require.Equal(t, "Python details", response[1]["details"])
			require.Equal(t, []interface{}{"api", "backend"}, response[1]["use_cases"])
			require.Equal(t, []interface{}{"rest api", "server side", "cli"}, response[1]["tags"])
			require.Equal(t, "https://python.org", response[1]["website"])
			require.Equal(t, "https://github.com/python/python", response[1]["github"])
			require.Equal(t, map[string]interface{}{"reason": "test relationship 3"}, response[1]["metadata"])
		})

		t.Run("return error 404 when the tool is not found", func(t *testing.T) {
			router := &core.HTTPRouter{
				ToolRepository:         toolRepository,
				RelationshipRepository: relationshipRepository,
			}
			h := router.Initialize()
			req := httptest.NewRequest(http.MethodGet, "/tools/not-found/alternatives", nil)
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			h.ServeHTTP(rec, req)
			require.Equal(t, http.StatusNotFound, rec.Code)
		})
	})
}
