# Backend Testing Guidelines

This document defines the testing strategy and conventions for `apps/backend`.

## Goals

- Keep feedback fast during development.
- Validate business logic in isolation.
- Validate real behavior for DB and HTTP flows.
- Keep CI deterministic and easy to maintain.

## Test Layers

### 1) Unit Tests (default)

Use unit tests for:

- Domain validation and entity behavior (`internal/domain`).
- Use-case orchestration and error mapping (`internal/feature/**/usecase`).
- Small pure helpers and serialization logic.

Guidelines:

- Prefer table-driven tests.
- Use explicit fixture values (no unnecessary randomness).
- Mock interfaces with `mockery`-generated mocks or lightweight fakes.
- Assert behavior and error types, not implementation details.

Run:

```bash
cd apps/backend
gotestsum -- ./...
```

### 2) Integration Tests (`integration` build tag)

Use integration tests for:

- Repository tests with real SQLite.
- HTTP router tests with real chi wiring + repositories.
- End-to-end request/response behavior (status + payload + DB effects).

Conventions:

- File naming: `*_integration_test.go`
- Build constraint at top of file:

```go
//go:build integration
```

- Use `testutil.SetupTestDB(t)` for temporary migrated SQLite DB.
- Keep each test isolated (prefer fresh setup per test/subtest when practical).
- Cleanup should use captured IDs/slugs from setup, not hardcoded IDs.

Run:

```bash
cd apps/backend
gotestsum -- -tags=integration ./...
```

## Local Developer Workflow (gotestsum)

Use `gotestsum` as the default local test runner for better readability and watch mode support.

Install:

```bash
go install gotest.tools/gotestsum@latest
```

If `gotestsum` is not found after install, ensure `$(go env GOPATH)/bin` is on your `PATH`.

Recommended commands:

```bash
cd apps/backend

# Full default suite (unit/default tests)
gotestsum -- ./...

# Integration suite
gotestsum -- -tags=integration ./...

# Focus one package while developing
gotestsum -- ./internal/feature/tool/http

# Focus one test
gotestsum -- -run TestToolRouter_CreateTool ./internal/feature/tool/http
```

Watch mode (great local dev experience):

```bash
cd apps/backend

# Re-run focused test/package on file changes
gotestsum --watch -- -run TestToolRouter_CreateTool ./internal/feature/tool/http
```

## Database Test Setup

- Integration tests use SQLite temp DB from `t.TempDir()`.
- Migrations are applied in test setup via goose.
- SQLite driver must be imported with a blank import in setup code:

```go
_ "modernc.org/sqlite"
```

## HTTP Integration Testing Conventions

- Build requests with `httptest.NewRequest`.
- Execute with `router.Initialize().ServeHTTP(rec, req)` or shared root router.
- Set `Content-Type: application/json` for JSON endpoints.
- Assert:
  - HTTP status code.
  - Response shape and key fields.
  - Important side effects in DB when relevant.

Prefer resilient assertions:

- Avoid over-asserting exact parser error strings when not needed.
- Assert stable error structure/messages and codes.

## Coverage Conventions

CI computes coverage from merged test execution (`-tags=integration`) and excludes noisy/generated/helper paths from reporting:

- Excluded: `db`, `testutil`

Local coverage command:

```bash
cd apps/backend
gotestsum -- -tags=integration -covermode=atomic -coverprofile=cover.out ./...
grep -vE '/db/|/testutil/' cover.out > cover.filtered.out
go tool cover -func=cover.filtered.out
```

Optional HTML view:

```bash
go tool cover -html=cover.filtered.out
```

## CI Reference

The GitHub Actions workflow at `.github/workflows/build.yml` currently runs:

- `golangci-lint`
- Backend tests with integration tag and filtered coverage report

Any change to test tags, coverage exclusions, or commands must be reflected in this workflow.

## Naming and Structure Guidelines

- Keep test names behavior-oriented (`return error 404 when ...` / `create ... successfully`).
- Group scenarios per endpoint/use-case.
- Share setup helpers in `apps/backend/testutil` when reused across packages.
- Avoid print/debug statements in committed tests.

## When to Add Which Test

- Add/modify **unit tests** for domain and use-case logic changes.
- Add/modify **integration tests** when:
  - SQL/repository behavior changes.
  - HTTP routing, middleware, or response contracts change.
  - Cross-layer behavior (HTTP -> use-case -> repository) changes.

