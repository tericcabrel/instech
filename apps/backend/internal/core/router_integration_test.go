//go:build integration

package core_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"tericcabrel/instech/internal/core"
	"tericcabrel/instech/internal/repository"
	"tericcabrel/instech/testutil"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHTTPRouter_GetRoot(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer db.Close()

	toolRepository := repository.NewToolRepository(db)
	relationshipRepository := repository.NewRelationshipRepository(db)
	router := core.HTTPRouter{
		ToolRepository:         toolRepository,
		RelationshipRepository: relationshipRepository,
	}
	h := router.Initialize()

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	body, err := io.ReadAll(rec.Body)
	require.NoError(t, err)
	require.Equal(t, "{\"message\":\"Hello from Instech\"}\n", string(body))
}
