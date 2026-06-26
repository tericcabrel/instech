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

func TestToolRouter_AllToolEndpoints(t *testing.T) {
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
			require.Equal(t, "backend", response["subType"])
			require.Equal(t, "active", response["devStatus"])
			require.Equal(t, "Go", response["prolang"])
			require.Equal(t, float64(2009), response["releaseYear"])
			require.Equal(t, "Test Details", response["details"])
			require.Equal(t, []interface{}{"api", "backend"}, response["useCases"])
			require.Equal(t, []interface{}{"rest api", "server side", "cli"}, response["tags"])
			require.Equal(t, "https://go.dev", response["website"])
			require.Equal(t, "https://github.com/golang/go", response["github"])
			require.Contains(t, response, "createdAt")
			require.Contains(t, response, "updatedAt")

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

	t.Run("[PATCH] /tools/{id}", func(t *testing.T) {
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
				"devStatus": "activ
			}`)
			req := httptest.NewRequest(http.MethodPatch, "/tools/golang", body)
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

		t.Run("return error 404 when tool is not found", func(t *testing.T) {
			body := bytes.NewBufferString(`{"name":"GoLang"}`)
			req := httptest.NewRequest(http.MethodPatch, "/tools/not-found", body)
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			router.Initialize().ServeHTTP(rec, req)
			require.Equal(t, http.StatusNotFound, rec.Code)
		})

		t.Run("patch tool with valid full payload", func(t *testing.T) {
			body := bytes.NewBufferString(`{
				"name": "GoLang",
				"category": "language",
				"subType": "backend",
				"prolang": "Go",
				"releaseYear": 2009,
				"devStatus": "deprecated",
				"details": "Golang details",
				"useCases": ["api", "backend"],
				"tags": ["web", "api", "cli"],
				"website": "https://golang-updated.org",
				"github": "https://github.com/golang/go-updated"
			}`)

			req := httptest.NewRequest(http.MethodPatch, "/tools/golang", body)
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
			require.Equal(t, "golang", response["slug"])
			require.Equal(t, "deprecated", response["devStatus"])
			require.Equal(t, "Golang details", response["details"])
			require.Equal(t, "https://golang-updated.org", response["website"])
		})

		t.Run("updates enum fields when valid", func(t *testing.T) {
			body := bytes.NewBufferString(`{
				"category": "framework",
				"subType": "frontend",
				"devStatus": "deprecated"
			}`)
			req := httptest.NewRequest(http.MethodPatch, "/tools/golang", body)
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			router.Initialize().ServeHTTP(rec, req)
			require.Equal(t, http.StatusOK, rec.Code)

			responseBody, err := io.ReadAll(rec.Body)
			require.NoError(t, err)
			var response map[string]any
			err = json.Unmarshal(responseBody, &response)
			require.NoError(t, err)
			require.Equal(t, "framework", response["category"])
			require.Equal(t, "frontend", response["subType"])
			require.Equal(t, "deprecated", response["devStatus"])
		})

		t.Run("updates only one scalar field and preserves others", func(t *testing.T) {
			body := bytes.NewBufferString(`{"name":"Go patched"}`)
			req := httptest.NewRequest(http.MethodPatch, "/tools/golang", body)
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			router.Initialize().ServeHTTP(rec, req)
			require.Equal(t, http.StatusOK, rec.Code)

			responseBody, err := io.ReadAll(rec.Body)
			require.NoError(t, err)
			var response map[string]any
			err = json.Unmarshal(responseBody, &response)
			require.NoError(t, err)
			require.Equal(t, "Go patched", response["name"])
			require.Equal(t, "framework", response["category"])
			require.Equal(t, "frontend", response["subType"])
			require.Equal(t, "https://golang-updated.org", response["website"])
		})

		t.Run("return error 400 when enum value is invalid", func(t *testing.T) {
			body := bytes.NewBufferString(`{"category":"invalid"}`)
			req := httptest.NewRequest(http.MethodPatch, "/tools/golang", body)
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			router.Initialize().ServeHTTP(rec, req)
			require.Equal(t, http.StatusBadRequest, rec.Code)
		})

		t.Run("return error 422 when fields are invalid", func(t *testing.T) {
			body := bytes.NewBufferString(`{
				"name": "",
				"website": "invalid"
			}`)
			req := httptest.NewRequest(http.MethodPatch, "/tools/golang", body)
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			router.Initialize().ServeHTTP(rec, req)
			require.Equal(t, http.StatusUnprocessableEntity, rec.Code)
		})

		t.Run("clear details and website explicitly with null", func(t *testing.T) {
			body := bytes.NewBufferString(`{"details":null, "website":null}`)
			req := httptest.NewRequest(http.MethodPatch, "/tools/golang", body)
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			router.Initialize().ServeHTTP(rec, req)
			require.Equal(t, http.StatusOK, rec.Code)
			responseBody, err := io.ReadAll(rec.Body)
			require.NoError(t, err)
			var response map[string]any
			err = json.Unmarshal(responseBody, &response)
			require.NoError(t, err)
			require.NotContains(t, response, "details")
			require.NotContains(t, response, "website")
		})

		t.Run("omitted optional fields are preserved", func(t *testing.T) {
			body := bytes.NewBufferString(`{"name":"Go final"}`)
			req := httptest.NewRequest(http.MethodPatch, "/tools/golang", body)
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			router.Initialize().ServeHTTP(rec, req)
			require.Equal(t, http.StatusOK, rec.Code)
			responseBody, err := io.ReadAll(rec.Body)
			require.NoError(t, err)
			var response map[string]any
			err = json.Unmarshal(responseBody, &response)
			require.NoError(t, err)
			require.Equal(t, "Go final", response["name"])
			require.NotContains(t, response, "website")
			require.NotContains(t, response, "details")
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

	t.Run("[GET] /tools/query", func(t *testing.T) {
		db := testutil.SetupTestDB(t)
		router := buildRouter(t, db)
		toolSlugs := populateAlternatives(t, db)

		t.Cleanup(func() {
			cleanAlternatives(t, db, toolSlugs)
			db.Close()
		})

		readSearchResponse := func(t *testing.T, query string) []map[string]any {
			t.Helper()

			req := httptest.NewRequest(http.MethodGet, "/tools/query"+query, nil)
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			router.Initialize().ServeHTTP(rec, req)
			require.Equal(t, http.StatusOK, rec.Code)

			responseBody, err := io.ReadAll(rec.Body)
			require.NoError(t, err)

			var response []map[string]any
			err = json.Unmarshal(responseBody, &response)
			require.NoError(t, err)

			return response
		}

		t.Run("missing q returns 200 with empty array", func(t *testing.T) {
			response := readSearchResponse(t, "")
			require.Empty(t, response)
		})

		t.Run("empty q returns 200 with empty array", func(t *testing.T) {
			response := readSearchResponse(t, "?q=")
			require.Empty(t, response)
		})

		t.Run("whitespace q returns 200 with empty array", func(t *testing.T) {
			response := readSearchResponse(t, "?q=%20%20")
			require.Empty(t, response)
		})

		t.Run("search success covers substring, case-insensitive, trimmed input, ordered ascending, and payload shape", func(t *testing.T) {
			response := readSearchResponse(t, "?q=%20O%20")

			require.Len(t, response, 3)

			require.Equal(t, "Golang", response[0]["name"])
			require.Equal(t, "Node.js", response[1]["name"])
			require.Equal(t, "Python", response[2]["name"])

			require.Equal(t, "golang", response[0]["slug"])
			require.Equal(t, "nodejs", response[1]["slug"])
			require.Equal(t, "python", response[2]["slug"])

			require.Equal(t, "language", response[0]["category"])
			require.Equal(t, "language", response[1]["category"])
			require.Equal(t, "language", response[2]["category"])

			require.Equal(t, "backend", response[0]["subType"])
			require.Equal(t, "fullstack", response[1]["subType"])
			require.Equal(t, "backend", response[2]["subType"])
		})

		t.Run("no match returns empty array", func(t *testing.T) {
			response := readSearchResponse(t, "?q=nonexistent-value")
			require.Empty(t, response)
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
			require.Equal(t, "backend", response["subType"])
			require.Equal(t, "Go", response["prolang"])
			require.Equal(t, float64(2009), response["releaseYear"])
			require.Equal(t, "active", response["devStatus"])
			require.Equal(t, "Test Details", response["details"])
			require.Equal(t, []interface{}{"api", "backend"}, response["useCases"])
			require.Equal(t, []interface{}{"rest api", "server side", "cli"}, response["tags"])
			require.Equal(t, "https://go.dev", response["website"])
			require.Equal(t, "https://github.com/golang/go", response["github"])
			require.Contains(t, response, "createdAt")
			require.Contains(t, response, "updatedAt")
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
			require.Equal(t, "fullstack", response[0]["subType"])
			require.Equal(t, "JavaScript", response[0]["prolang"])
			require.Equal(t, float64(2009), response[0]["releaseYear"])
			require.Equal(t, "active", response[0]["devStatus"])
			require.Equal(t, "Node.js details", response[0]["details"])
			require.Equal(t, []interface{}{"api", "frontend", "fullstack"}, response[0]["useCases"])
			require.Equal(t, []interface{}{"web", "api", "frontend"}, response[0]["tags"])
			require.Equal(t, "https://nodejs.org", response[0]["website"])
			require.Equal(t, "https://github.com/nodejs/node", response[0]["github"])
			require.Equal(t, map[string]interface{}{"reason": "test relationship 1"}, response[0]["metadata"])

			require.Equal(t, "python", response[1]["id"])
			require.Equal(t, "Python", response[1]["name"])
			require.Equal(t, "language", response[1]["category"])
			require.Equal(t, "backend", response[1]["subType"])
			require.Equal(t, "Python", response[1]["prolang"])
			require.Equal(t, float64(1995), response[1]["releaseYear"])
			require.Equal(t, "active", response[1]["devStatus"])
			require.Equal(t, "Python details", response[1]["details"])
			require.Equal(t, []interface{}{"api", "backend"}, response[1]["useCases"])
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

	t.Run("[GET] /tools/{id}/graph", func(t *testing.T) {
		db := testutil.SetupTestDB(t)
		router := buildRouter(t, db)

		toolRepository := repository.NewToolRepository(db)
		relationshipRepository := repository.NewRelationshipRepository(db)
		createTestTool := func(t *testing.T, input domain.CreateToolInput) domain.Tool {
			t.Helper()
			tool, err := domain.CreateTool(input)
			require.NoError(t, err)
			createdTool, err := toolRepository.CreateTool(context.Background(), tool)
			require.NoError(t, err)

			return createdTool
		}
		createTestRelationship := func(t *testing.T, input domain.CreateRelationshipInput) {
			t.Helper()
			relationship, err := domain.CreateRelationship(input)
			require.NoError(t, err)
			_, err = relationshipRepository.CreateRelationship(context.Background(), relationship)
			require.NoError(t, err)
		}

		golang := createTestTool(t, domain.CreateToolInput{
			Name:        "Golang",
			Slug:        "golang",
			Category:    "language",
			SubType:     "backend",
			DevStatus:   "active",
			Prolang:     new("Go"),
			ReleaseYear: 2009,
			UseCases:    []string{"api"},
			Tags:        []string{"server"},
		})
		nodejs := createTestTool(t, domain.CreateToolInput{
			Name:        "Node.js",
			Slug:        "nodejs",
			Category:    "language",
			SubType:     "fullstack",
			DevStatus:   "active",
			Prolang:     new("JavaScript"),
			ReleaseYear: 2009,
			UseCases:    []string{"api"},
			Tags:        []string{"web"},
		})
		python := createTestTool(t, domain.CreateToolInput{
			Name:        "Python",
			Slug:        "python",
			Category:    "language",
			SubType:     "backend",
			DevStatus:   "active",
			Prolang:     new("Python"),
			ReleaseYear: 1995,
			UseCases:    []string{"api"},
			Tags:        []string{"script"},
		})
		rust := createTestTool(t, domain.CreateToolInput{
			Name:        "Rust",
			Slug:        "rust",
			Category:    "language",
			SubType:     "backend",
			DevStatus:   "active",
			Prolang:     new("Rust"),
			ReleaseYear: 2010,
			UseCases:    []string{"systems"},
			Tags:        []string{"performance"},
		})

		createTestRelationship(t, domain.CreateRelationshipInput{
			FromToolID: golang.Id,
			ToToolID:   nodejs.Id,
			Kind:       "alternative_to",
			Reason:     "rel 1",
		})
		createTestRelationship(t, domain.CreateRelationshipInput{
			FromToolID: golang.Id,
			ToToolID:   python.Id,
			Kind:       "used_with",
			Reason:     "rel 2",
		})
		createTestRelationship(t, domain.CreateRelationshipInput{
			FromToolID: nodejs.Id,
			ToToolID:   rust.Id,
			Kind:       "alternative_to",
			Reason:     "rel 3",
		})

		t.Cleanup(func() {
			cleanAlternatives(t, db, []string{"golang", "nodejs", "python", "rust"})
			db.Close()
		})

		t.Run("returns graph with default params", func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/tools/golang/graph", nil)
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			router.Initialize().ServeHTTP(rec, req)
			require.Equal(t, http.StatusOK, rec.Code)

			responseBody, err := io.ReadAll(rec.Body)
			require.NoError(t, err)

			var response map[string]any
			err = json.Unmarshal(responseBody, &response)
			require.NoError(t, err)

			require.Equal(t, "golang", response["focusNodeId"])
			require.Equal(t, map[string]any{
				"depth":      float64(1),
				"layoutMode": "chronological",
				"totalLinks": float64(2),
				"totalNodes": float64(3),
				"kindsApplied": []any{
					"built_on",
					"inspired_by",
					"alternative_to",
					"replaced_by",
					"used_with",
				},
			}, response["meta"])
		})

		t.Run("returns graph with depth and repeated kinds filters", func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/tools/golang/graph?depth=2&kinds=alternative_to&kinds=used_with&layoutMode=force", nil)
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			router.Initialize().ServeHTTP(rec, req)
			require.Equal(t, http.StatusOK, rec.Code)

			responseBody, err := io.ReadAll(rec.Body)
			require.NoError(t, err)

			var response map[string]any
			err = json.Unmarshal(responseBody, &response)
			require.NoError(t, err)

			meta := response["meta"].(map[string]any)

			require.Equal(t, float64(2), meta["depth"])
			require.Equal(t, "force", meta["layoutMode"])
			require.Equal(t, float64(4), meta["totalNodes"])
			require.Equal(t, float64(3), meta["totalLinks"])
		})

		t.Run("returns error 404 when tool is not found", func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/tools/unknown/graph", nil)
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			router.Initialize().ServeHTTP(rec, req)
			require.Equal(t, http.StatusNotFound, rec.Code)
		})
	})
}
