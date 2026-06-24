# Backend (`tericcabrel/instech`)

Instech is a tool for finding alternatives to popular tools.
## Tech stack
- Go
- Chi
- SQLite via `database/sql`
- Goose
- SQLC
- Air
- Golangci-lint
- OpenAPI spec ([`openapi.yaml`](openapi.yaml), linted with `@redocly/cli`)
- Pre-commit hook
- GitHub Actions

## Run locally

### Install dependencies:

```bash
go mod tidy
go mod download
```

### Create the database and run migrations:

```bash
mkdir -p db/data

sqlite3 db/data/instech.db "VACUUM;"                        
sqlite3 db/data/instech.db "PRAGMA journal_mode=WAL;"          

# Apply migrations
goose up
```
VACUUM is the process of compacting the database file to remove deleted rows and free up space.
WAL is the Write-Ahead Logging mode, it allows readers to proceed without blocking writers.

### Set environment variables
Create the `.env` file by copying the `.env.example` file

```bash
cp .env.example .env
```
Set the environment variables in the `.env` file.

### Run the server
```bash
go run . [-debug]
```

With watch mode (using [air](https://github.com/air-verse/air)):
```bash
air
```

The server listens on `:8800`.

## Lint

Install [golangci-lint](https://golangci-lint.run/welcome/install/) v2.x, then from this directory:

```bash
golangci-lint run ./...
```

To run the linting and fix the files, run the following command:
```bash
golangci-lint run ./... --fix
```

Configuration lives in [`.golangci.yml`](.golangci.yml). CI runs the same check via [`.github/workflows/build.yml`](../../.github/workflows/build.yml) with `working-directory: apps/backend`.

### OpenAPI spec

The API contract lives in [`openapi.yaml`](openapi.yaml). Lint it with Redocly (via `npx` — nothing is added to `go.mod` or `package.json`):

```bash
npx --yes @redocly/cli@2.34.0 lint openapi.yaml
```

CI runs the same command in the `lint` job.

### Pre-commit hook

From the **repository root** (not `apps/backend`), enable the shared hooks so commits run lint:

```bash
chmod +x .githooks/pre-commit
git config core.hooksPath .githooks
```

The hook runs `golangci-lint run ./...` inside `apps/backend`. Skip when needed with `git commit --no-verify`.

## How to add new database migrations

1. Create a new migration file in the `db/migrations` directory.
```bash
goose create <migration_name> sql
```
The migration file will be created in the `db/migrations` directory; edit the file to add the SQL queries.

Example:
```bash
goose create add_website_column sql
```

2. Run the migrations:

```bash
goose up
```

3. Add SQL queries to the `db/sqlc` directory if needed and run the following command to generate the SQLC code:

```bash
sqlc generate
```
The SQLC code will be generated in the `db/queries` directory. Edit the `sqlc.yaml` file to update the configuration for generating the SQLC code.