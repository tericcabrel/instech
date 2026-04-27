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
sqlite3 db/data/instech.db "VACUUM;"
goose up
```

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

Configuration lives in [`.golangci.yml`](.golangci.yml). CI runs the same check via [`.github/workflows/golangci-lint.yml`](../../.github/workflows/golangci-lint.yml) with `working-directory: apps/backend`.

## Pre-commit hook (optional)

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