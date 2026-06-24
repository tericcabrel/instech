//go:build integration

package http_test

import (
	"bytes"
	"context"
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
	"time"

	"database/sql"

	"github.com/stretchr/testify/require"
)

func buildRouter(t *testing.T, db *sql.DB) *core.HTTPRouter {
	t.Helper()

	return &core.HTTPRouter{
		ToolRepository:         repository.NewToolRepository(db),
		RelationshipRepository: repository.NewRelationshipRepository(db),
	}
}

func createTool(t *testing.T, db *sql.DB) int {
	t.Helper()

	tool, err := domain.CreateTool(domain.CreateToolInput{
		Name:        "Golang",
		Slug:        "golang",
		Category:    "language",
		SubType:     "backend",
		Prolang:     new("Go"),
		ReleaseYear: 2009,
		DevStatus:   "active",
		Details:     new("Test Details"),
		UseCases:    []string{"api", "backend"},
		Tags:        []string{"rest api", "server side", "cli"},
		Website:     new("https://go.dev"),
		Github:      new("https://github.com/golang/go"),
	})
	require.NoError(t, err)

	toolRepository := repository.NewToolRepository(db)
	createdTool, err := toolRepository.CreateTool(context.Background(), tool)
	require.NoError(t, err)

	return createdTool.Id
}

func deleteTools(t *testing.T, db *sql.DB, toolIds ...int) {
	t.Helper()

	toolRepository := repository.NewToolRepository(db)

	tools, err := toolRepository.GetToolByIDs(context.Background(), toolIds)
	require.NoError(t, err)

	for _, tool := range tools {
		err := toolRepository.DeleteTool(context.Background(), tool.Slug)
		require.NoError(t, err)
	}
}

func populateAlternatives(t *testing.T, db *sql.DB) []string {
	t.Helper()

	toolRepository := repository.NewToolRepository(db)
	relationshipRepository := repository.NewRelationshipRepository(db)

	golangTool, err := domain.CreateTool(domain.CreateToolInput{
		Name:        "Golang",
		Slug:        "golang",
		Category:    "language",
		SubType:     "backend",
		DevStatus:   "active",
		Prolang:     new("Go"),
		ReleaseYear: 2009,
		Details:     new("Golang details"),
		UseCases:    []string{"api", "backend"},
		Tags:        []string{"rest api", "server side", "cli"},
		Website:     new("https://golang.org"),
		Github:      new("https://github.com/golang/go"),
	})
	require.NoError(t, err)
	createdGolangTool, err := toolRepository.CreateTool(context.Background(), golangTool)
	require.NoError(t, err)

	nodejsTool, err := domain.CreateTool(domain.CreateToolInput{
		Name:        "Node.js",
		Slug:        "nodejs",
		Category:    "language",
		SubType:     "fullstack",
		DevStatus:   "active",
		Prolang:     new("JavaScript"),
		Details:     new("Node.js details"),
		UseCases:    []string{"api", "frontend", "fullstack"},
		Tags:        []string{"web", "api", "frontend"},
		Website:     new("https://nodejs.org"),
		Github:      new("https://github.com/nodejs/node"),
		ReleaseYear: 2009,
	})
	require.NoError(t, err)
	createdNodejsTool, err := toolRepository.CreateTool(context.Background(), nodejsTool)
	require.NoError(t, err)

	pythonTool, err := domain.CreateTool(domain.CreateToolInput{
		Name:        "Python",
		Slug:        "python",
		Category:    "language",
		SubType:     "backend",
		DevStatus:   "active",
		Prolang:     new("Python"),
		Details:     new("Python details"),
		UseCases:    []string{"api", "backend"},
		Tags:        []string{"rest api", "server side", "cli"},
		Website:     new("https://python.org"),
		Github:      new("https://github.com/python/python"),
		ReleaseYear: 1995,
	})
	require.NoError(t, err)
	createdPythonTool, err := toolRepository.CreateTool(context.Background(), pythonTool)
	require.NoError(t, err)

	rel1, err := domain.CreateRelationship(domain.CreateRelationshipInput{
		FromToolID: createdGolangTool.Id,
		ToToolID:   createdNodejsTool.Id,
		Kind:       "alternative_to",
		Reason:     "test relationship 1",
	})
	require.NoError(t, err)
	_, err = relationshipRepository.CreateRelationship(context.Background(), rel1)
	require.NoError(t, err)

	rel2, err := domain.CreateRelationship(domain.CreateRelationshipInput{
		FromToolID: createdNodejsTool.Id,
		ToToolID:   createdGolangTool.Id,
		Kind:       "alternative_to",
		Reason:     "test relationship 2",
	})
	require.NoError(t, err)
	_, err = relationshipRepository.CreateRelationship(context.Background(), rel2)
	require.NoError(t, err)

	rel3, err := domain.CreateRelationship(domain.CreateRelationshipInput{
		FromToolID: createdGolangTool.Id,
		ToToolID:   createdPythonTool.Id,
		Kind:       "alternative_to",
		Reason:     "test relationship 3",
	})
	require.NoError(t, err)
	_, err = relationshipRepository.CreateRelationship(context.Background(), rel3)
	require.NoError(t, err)

	return []string{createdGolangTool.Slug, createdNodejsTool.Slug, createdPythonTool.Slug}
}

func cleanAlternatives(t *testing.T, db *sql.DB, toolSlugs []string) {
	t.Helper()

	toolRepository := repository.NewToolRepository(db)
	for _, slug := range toolSlugs {
		err := toolRepository.DeleteTool(context.Background(), slug)
		require.NoError(t, err)
	}

	relationshipRepository := repository.NewRelationshipRepository(db)
	relationships, err := relationshipRepository.GetRelationshipsAll(context.Background(), repository.GetRelationshipsAllParams{
		Limit: 100,
	})
	require.NoError(t, err)
	for _, relationship := range relationships.Relationships {
		err := relationshipRepository.DeleteRelationship(context.Background(), relationship.ID)
		require.NoError(t, err)
	}
}

func TestToolRouter_CreateTool(t *testing.T) {
	t.Run("[POST] /tools", func(t *testing.T) {
		db := testutil.SetupTestDB(t)
		t.Cleanup(func() {
			db.Close()
		})
		router := buildRouter(t, db)
		t.Run("return error 400 when the request body is an invalid JSON", func(t *testing.T) {
			body := bytes.NewBufferString(`{
				"name": "Golang",
				"slug": "golang",
				"category": "invalid",
			}`)
			req := httptest.NewRequest(http.MethodPost, "/tools", body)
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

		t.Run("return error 400 when the category is invalid", func(t *testing.T) {
			body := bytes.NewBufferString(`{
				"name": "Golang",
				"slug": "golang",
				"category": "invalid"
			}`)
			req := httptest.NewRequest(http.MethodPost, "/tools", body)
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
			expectedDetails := map[string]interface{}{
				"Category": "invalid",
				"Message":  "The category is invalid. Valid categories are: language, framework, library",
			}
			require.Equal(t, expectedDetails, response["details"])
		})

		t.Run("return error 400 when the sub type is invalid", func(t *testing.T) {
			body := bytes.NewBufferString(`{
				"name": "Golang",
				"slug": "golang",
				"category": "language",
				"subType": "invalid"
			}`)
			req := httptest.NewRequest(http.MethodPost, "/tools", body)
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
			expectedDetails := map[string]interface{}{
				"SubType": "invalid",
				"Message": "The sub type is invalid. Valid sub types are: backend, frontend, fullstack, mobile, desktop, game, other",
			}
			require.Equal(t, expectedDetails, response["details"])
		})

		t.Run("return error 400 when the dev status is invalid", func(t *testing.T) {
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
			router.Initialize().ServeHTTP(rec, req)

			require.Equal(t, http.StatusBadRequest, rec.Code)

			responseBody, err := io.ReadAll(rec.Body)

			require.NoError(t, err)

			var response map[string]any
			err = json.Unmarshal(responseBody, &response)

			require.NoError(t, err)
			require.Equal(t, "Bad Request", response["message"])
			expectedDetails := map[string]interface{}{
				"DevStatus": "invalid",
				"Message":   "The dev status is invalid. Valid dev statuses are: active, deprecated",
			}
			require.Equal(t, expectedDetails, response["details"])
		})

		t.Run("return error 422 when the fields are invalid", func(t *testing.T) {
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
			router.Initialize().ServeHTTP(rec, req)
			require.Equal(t, http.StatusUnprocessableEntity, rec.Code)

			responseBody, err := io.ReadAll(rec.Body)
			require.NoError(t, err)

			var response map[string]any
			err = json.Unmarshal(responseBody, &response)

			require.NoError(t, err)
			require.Equal(t, "Unprocessable Entity", response["message"])

			year := strconv.Itoa(time.Now().Year())
			expectedInvalidFields := map[string]interface{}{
				"Name":        "The name is required",
				"Prolang":     "The programming language is required",
				"Slug":        "The slug is required",
				"ReleaseYear": "The release year is invalid. Valid release years are between 1940 and " + year,
				"Website":     "The website is invalid. Valid websites must be a valid URL",
				"Github":      "The github is invalid. Valid github must be a valid URL",
			}
			expectedDetails := map[string]interface{}{
				"Fields": expectedInvalidFields,
			}
			require.Equal(t, expectedDetails, response["details"])
		})

		t.Run("Create tool with valid input", func(t *testing.T) {
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
			router.Initialize().ServeHTTP(rec, req)

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
			require.Equal(t, "active", response["dev_status"])
			require.Equal(t, "Go", response["prolang"])
			require.Equal(t, float64(2009), response["release_year"])
			require.Equal(t, "Test Details", response["details"])
			require.Equal(t, []interface{}{"api", "backend"}, response["use_cases"])
			require.Equal(t, []interface{}{"rest api", "server side", "cli"}, response["tags"])
			require.Equal(t, "https://go.dev", response["website"])
			require.Equal(t, "https://github.com/golang/go", response["github"])
			require.Contains(t, response, "created_at")
			require.Contains(t, response, "updated_at")

			deleteTools(t, db, int(response["id"].(float64)))
		})

		t.Run("return error 400 when the tool already exists", func(t *testing.T) {
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
			router.Initialize().ServeHTTP(recOne, reqOne)

			require.Equal(t, http.StatusCreated, recOne.Code)

			responseBodyOne, err := io.ReadAll(recOne.Body)

			require.NoError(t, err)

			var responseOne map[string]any
			err = json.Unmarshal(responseBodyOne, &responseOne)

			require.NoError(t, err)

			toolIdCreated := responseOne["id"].(float64)

			// Send the same request again
			bodyTwo := bytes.NewBufferString(bodyString)

			reqTwo := httptest.NewRequest(http.MethodPost, "/tools", bodyTwo)
			reqTwo.Header.Set("Content-Type", "application/json")
			recTwo := httptest.NewRecorder()
			router.Initialize().ServeHTTP(recTwo, reqTwo)
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

			deleteTools(t, db, int(toolIdCreated))
		})
	})

	t.Run("[PUT] /tools/{id}", func(t *testing.T) {
		db := testutil.SetupTestDB(t)
		router := buildRouter(t, db)
		createdToolId := createTool(t, db)
		t.Cleanup(func() {
			deleteTools(t, db, createdToolId)
			db.Close()
		})
		t.Run("return error 400 when the request body is an invalid JSON", func(t *testing.T) {
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
			router.Initialize().ServeHTTP(rec, req)
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
			body := bytes.NewBufferString(`{
				"name": "Golang",
				"slug": "golang",
				"category": "invalid"
			}`)
			req := httptest.NewRequest(http.MethodPut, "/tools/golang", body)
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
			expectedDetails := map[string]interface{}{
				"Category": "invalid",
				"Message":  "The category is invalid. Valid categories are: language, framework, library",
			}
			require.Equal(t, expectedDetails, response["details"])
		})

		t.Run("return error 400 when the sub type is invalid", func(t *testing.T) {
			body := bytes.NewBufferString(`{
				"name": "Golang",
				"slug": "golang",
				"category": "language",
				"subType": "invalid"
			}`)
			req := httptest.NewRequest(http.MethodPut, "/tools/golang", body)
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
			expectedDetails := map[string]interface{}{
				"SubType": "invalid",
				"Message": "The sub type is invalid. Valid sub types are: backend, frontend, fullstack, mobile, desktop, game, other",
			}
			require.Equal(t, expectedDetails, response["details"])
		})

		t.Run("return error 400 when the dev status is invalid", func(t *testing.T) {
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
			router.Initialize().ServeHTTP(rec, req)

			require.Equal(t, http.StatusBadRequest, rec.Code)
		})

		t.Run("return error 422 when the fields are invalid", func(t *testing.T) {
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
			router.Initialize().ServeHTTP(rec, req)

			require.Equal(t, http.StatusUnprocessableEntity, rec.Code)

			responseBody, err := io.ReadAll(rec.Body)

			require.NoError(t, err)

			var response map[string]any
			err = json.Unmarshal(responseBody, &response)

			require.NoError(t, err)
			require.Equal(t, "Unprocessable Entity", response["message"])

			year := strconv.Itoa(time.Now().Year())
			expectedInvalidFields := map[string]interface{}{
				"Name":        "The name is required",
				"Prolang":     "The programming language is required",
				"ReleaseYear": "The release year is invalid. Valid release years are between 1940 and " + year,
				"Website":     "The website is invalid. Valid websites must be a valid URL",
				"Github":      "The github is invalid. Valid github must be a valid URL",
			}
			expectedDetails := map[string]interface{}{
				"Fields": expectedInvalidFields,
			}
			require.Equal(t, expectedDetails, response["details"])
		})

		t.Run("Update tool with valid input", func(t *testing.T) {
			updateBody := bytes.NewBufferString(`{
				"name": "GoLang",
				"slug": "nodejs",
				"category": "language",
				"subType": "backend",
				"prolang": "Golang",
				"releaseYear": 2009,
				"devStatus": "deprecated",
				"details": "Golang details",
				"usecases": ["api", "backend"],
				"tags": ["web", "api", "cli"],
				"website": "https://golang-updated.org",
				"github": "https://github.com/golang/go-updated"
			}`)

			req := httptest.NewRequest(http.MethodPut, "/tools/golang", updateBody)
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			router.Initialize().ServeHTTP(rec, req)

			require.Equal(t, http.StatusOK, rec.Code)

			responseBody, err := io.ReadAll(rec.Body)

			require.NoError(t, err)

			var response map[string]any
			err = json.Unmarshal(responseBody, &response)

			require.NoError(t, err)

			require.Equal(t, float64(createdToolId), response["id"])
			require.Equal(t, "GoLang", response["name"])
			require.Equal(t, "golang", response["slug"]) // slug should not be updated
			require.Equal(t, "language", response["category"])
			require.Equal(t, "backend", response["sub_type"])
			require.Equal(t, "deprecated", response["dev_status"])
			require.Equal(t, "Golang", response["prolang"])
			require.Equal(t, float64(2009), response["release_year"])
			require.Equal(t, "Golang details", response["details"])
			require.Equal(t, []interface{}{"api", "backend"}, response["use_cases"])
			require.Equal(t, []interface{}{"web", "api", "cli"}, response["tags"])
			require.Equal(t, "https://golang-updated.org", response["website"])
			require.Equal(t, "https://github.com/golang/go-updated", response["github"])
			require.Contains(t, response, "created_at")
			require.Contains(t, response, "updated_at")
		})
	})

	t.Run("[DELETE] /tools/{id}", func(t *testing.T) {
		db := testutil.SetupTestDB(t)
		router := buildRouter(t, db)
		createTool(t, db)

		t.Cleanup(func() {
			db.Close()
		})

		t.Run("Do not nothing when the tool is not found", func(t *testing.T) {
			req := httptest.NewRequest(http.MethodDelete, "/tools/not-found", nil)
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			router.Initialize().ServeHTTP(rec, req)

			require.Equal(t, http.StatusNoContent, rec.Code)
		})

		t.Run("delete tool successfully", func(t *testing.T) {
			req := httptest.NewRequest(http.MethodDelete, "/tools/golang", nil)
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			router.Initialize().ServeHTTP(rec, req)

			require.Equal(t, http.StatusNoContent, rec.Code)
		})
	})

	t.Run("[GET] /tools/{id}", func(t *testing.T) {
		db := testutil.SetupTestDB(t)
		router := buildRouter(t, db)
		createdToolId := createTool(t, db)
		t.Cleanup(func() {
			deleteTools(t, db, createdToolId)
			db.Close()
		})

		t.Run("return tool successfully", func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/tools/golang", nil)
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			router.Initialize().ServeHTTP(rec, req)
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
			require.Equal(t, "active", response["dev_status"])
			require.Equal(t, "Test Details", response["details"])
			require.Equal(t, []interface{}{"api", "backend"}, response["use_cases"])
			require.Equal(t, []interface{}{"rest api", "server side", "cli"}, response["tags"])
			require.Equal(t, "https://go.dev", response["website"])
			require.Equal(t, "https://github.com/golang/go", response["github"])
			require.Contains(t, response, "created_at")
			require.Contains(t, response, "updated_at")
		})

		t.Run("return error 404 when the tool is not found", func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/tools/not-found", nil)
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			router.Initialize().ServeHTTP(rec, req)
			require.Equal(t, http.StatusNotFound, rec.Code)
		})
	})

	t.Run("[GET] /tools/{id}/alternatives", func(t *testing.T) {
		db := testutil.SetupTestDB(t)
		router := buildRouter(t, db)
		toolSlugs := populateAlternatives(t, db)

		t.Cleanup(func() {
			cleanAlternatives(t, db, toolSlugs)
			db.Close()
		})

		t.Run("return tool alternatives successfully", func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/tools/golang/alternatives", nil)
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			router.Initialize().ServeHTTP(rec, req)
			require.Equal(t, http.StatusOK, rec.Code)

			responseBody, err := io.ReadAll(rec.Body)
			require.NoError(t, err)

			var response []map[string]any
			err = json.Unmarshal(responseBody, &response)
			require.NoError(t, err)

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
			require.Equal(t, float64(1995), response[1]["release_year"])
			require.Equal(t, "active", response[1]["dev_status"])
			require.Equal(t, "Python details", response[1]["details"])
			require.Equal(t, []interface{}{"api", "backend"}, response[1]["use_cases"])
			require.Equal(t, []interface{}{"rest api", "server side", "cli"}, response[1]["tags"])
			require.Equal(t, "https://python.org", response[1]["website"])
			require.Equal(t, "https://github.com/python/python", response[1]["github"])
			require.Equal(t, map[string]interface{}{"reason": "test relationship 3"}, response[1]["metadata"])
		})

		t.Run("return error 404 when the tool is not found", func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/tools/not-found/alternatives", nil)
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			router.Initialize().ServeHTTP(rec, req)

			require.Equal(t, http.StatusNotFound, rec.Code)
		})
	})
}
