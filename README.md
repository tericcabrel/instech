# Instech

Run project commands from the repository root with [Just](https://github.com/casey/just).

## Setup

```bash
brew install just
```

Optional Go tools for backend recipes (`GOPATH/bin` on `PATH`):

```bash
go install gotest.tools/gotestsum@latest
go install github.com/vektra/mockery/v3@latest
```

## Commands

List all recipes (including `web` and `api` modules):

```bash
just
```

List one module: `just --list web` or `just --list api`.

Common recipes:

| Command | Description |
|---------|-------------|
| `just dev` | Start web (:8800) and api (:8801) dev servers |
| `just ci` | Run full local CI checks |
| `just web dev` | Frontend dev server |
| `just web build` | Production frontend build |
| `just web lint` | ESLint + Biome |
| `just web test` | Vitest |
| `just web api-generate` | Regenerate Orval client from OpenAPI |
| `just api dev` | Backend with hot reload (air) |
| `just api lint` | golangci-lint |
| `just api test` | Unit tests (gotestsum) |
| `just api test-integration` | All tests including integration |
| `just api mockery` | Regenerate testify mocks |

Use `just` to list all commands. Filter with `just --list web` or `just --list api`.

## App docs

- Web: [apps/frontend/README.md](apps/frontend/README.md)
- Api: [apps/backend/README.md](apps/backend/README.md)
