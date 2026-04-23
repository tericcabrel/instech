package testutil

import (
	"context"
	"database/sql"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/require"

	_ "modernc.org/sqlite"
)

func ModuleRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("go.mod not found from %s upward", dir)
		}
		dir = parent
	}
}

func MigrationsDir() (string, error) {
	root, err := ModuleRoot()
	if err != nil {
		return "", err
	}
	return filepath.Join(root, "db", "migrations"), nil
}

func SetupTestDB(t *testing.T) *sql.DB {
	t.Helper()
	dir := t.TempDir()
	dsn := "file:" + filepath.Join(dir, "test.db") //+ "?_pragma=foreign_keys(1)"
	db, err := sql.Open("sqlite", dsn)

	require.NoError(t, err)

	t.Cleanup(func() { _ = db.Close() })

	migrationsDir, err := MigrationsDir()
	require.NoError(t, err)

	// run migrations (goose or raw SQL)
	migrations, err := goose.NewProvider(goose.DialectSQLite3, db, fs.FS(os.DirFS(migrationsDir)))
	require.NoError(t, err)

	res, err := migrations.Up(context.Background())

	require.NoError(t, err)
	fmt.Printf("Migrations applied: %v", res)
	return db
}
