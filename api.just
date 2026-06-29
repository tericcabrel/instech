set working-directory := 'apps/backend'

# Start API server with hot reload (air, port 8801)
[group('api')]
dev:
  air

# Run golangci-lint
[group('api')]
lint:
  golangci-lint run ./...

# Run golangci-lint with auto-fix
[group('api')]
lint-fix:
  golangci-lint run ./... --fix

# Run unit tests (no integration tag)
[group('api')]
test:
  gotestsum -- ./...

# Run all tests including integration
[group('api')]
test-integration:
  gotestsum -- -tags=integration ./...

# Run integration tests and write cover.out
[group('api')]
test-coverage:
  gotestsum -- -tags=integration -coverprofile=cover.out ./...

# Run tests in one package
[group('api')]
test-pkg path:
  gotestsum -- {{path}}

# Run one test by name
[group('api')]
test-run path name:
  gotestsum -- {{path}} -run {{name}}

# Re-run package tests on file changes
[group('api')]
test-watch path:
  gotestsum --watch -- {{path}}

# Re-run integration tests in a package on changes
[group('api')]
test-watch-integration path:
  gotestsum --watch -- -tags=integration {{path}}

# Re-run one integration test on changes
[group('api')]
test-watch-run path name:
  gotestsum --watch -- -tags=integration -run {{name}} {{path}}

# Apply pending database migrations
[group('api')]
db-up:
  goose up

# Create a new SQL migration file
[group('api')]
db-migrate name:
  goose create {{name}} sql

# Regenerate type-safe query code from SQL
[group('api')]
sqlc:
  sqlc generate

# Regenerate testify mocks in testutil/testmocks.go
[group('api')]
mockery:
  mockery

# Lint the OpenAPI spec (Redocly)
[group('api')]
openapi-lint:
  npx --yes @redocly/cli@2.34.0 lint openapi.yaml
