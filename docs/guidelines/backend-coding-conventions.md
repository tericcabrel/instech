# Backend Coding Conventions

This document captures the most important coding patterns for `apps/backend`.
Keep it short and use the existing code as the final source of truth.

## Architecture

The backend follows a small layered Go structure:

- `internal/domain`: business entities, input structs, validation rules, and typed domain errors.
- `internal/feature/<name>/usecase`: application actions that orchestrate domain logic and repositories.
- `internal/feature/<name>/http`: chi route wiring, request decoding, HTTP error mapping, and response writing.
- `internal/repository`: persistence interfaces, SQLC calls, and database-record-to-domain mapping.
- `internal/infra`: shared infrastructure helpers such as HTTP response helpers.
- `internal/core`: top-level router composition and middleware.

Keep new behavior in the lowest layer that owns it. Validation that defines whether an entity is valid belongs in `domain`; request parsing belongs in `http`; SQLC details belong in `repository`.

## Domain

Domain code should be usable without HTTP or database dependencies.

- Create entities through constructor-style functions such as `CreateTool` and `CreateRelationship`.
- Put entity updates on pointer receiver methods such as `(*Tool).Update`.
- Return typed errors for known invalid states, for example `ErrInvalidField` or `ErrInvalidToolCategory`.
- Aggregate multiple field validation failures into `ErrInvalidField{Fields: map[string]string{...}}`.
- Keep allowed values near the entity they validate, with small helper functions such as `IsCategoryValid`.
- Prefer explicit input structs (`CreateToolInput`, `UpdateRelationshipInput`) over passing large entity structs into domain creation.

## Use Cases

Use cases are thin orchestration units.

- Name them `<Action><Entity>UseCase` with an `Execute(...)` method.
- Inject repository interfaces as exported struct fields.
- Build or update domain entities before calling repositories.
- Convert storage-specific misses, such as `sql.ErrNoRows`, into common application errors such as `common.ErrResourceNotFound`.
- Return domain/common errors unchanged so HTTP handlers can choose the response status.
- Keep use cases independent from chi, `http.ResponseWriter`, JSON decoding, and response formatting.

## HTTP Routers

HTTP packages adapt the API boundary to use cases.

- Use chi routers with `Initialize()` methods on dependency structs such as `ToolRouter`.
- Decode JSON into usecase input structs.
- Use `httprouter` helpers for every JSON response and error response.
- Map expected domain/common errors explicitly to HTTP status codes.
- For unexpected errors, call `httprouter.InternalServerError(w, err, "<UseCaseName>")`.
- Keep handler logic focused on decode, execute, map error, respond. Move business decisions into domain or usecase code.

## Repositories

Repositories hide SQLC and database representation details from the rest of the app.

- Define a repository interface next to its concrete implementation.
- Use SQLC-generated query types only inside `internal/repository`.
- Convert SQLC records to domain entities with mapper functions.
- Keep JSON marshaling/unmarshaling for persisted JSON fields in repository/mappers.
- Return empty domain values with errors when persistence or mapping fails.

## Errors

Use typed errors as part of the application contract.

- Domain errors describe invalid business input.
- Common errors describe app-level resource states such as not found or already exists.
- HTTP handlers are responsible for translating known error types to status codes.
- Internal server errors should be logged through `httprouter.InternalServerError` with a useful source string.

## Tests

Follow `docs/backend-testing-guidelines.md` for the full testing policy.

## Style

- Run `gofmt`/`goimports`; imports stay at the top.
- Keep packages small and named by their layer or feature role.
- Prefer clear struct literals and explicit field names.
- Keep constants and validation messages close to the code that uses them.
- Do not add abstraction unless it removes repeated behavior already visible in the backend.
- One endpoint = one use case. When a backend route use case needs multiple helper files, create a named folder under the `usecase` folder, put the main use case file inside that folder, and keep route-specific helper modules scoped to the same folder unless they are reused by another route. Example: `feature/tool/usecase/create-tool/create-tool.go`.
- A use case must not import or call another use case. If a flow requires several use cases (e.g. validate then write), compose them in the router. Shared, non-use-case logic goes in a `lib/` helper (scoped to the route folder if route-specific, or to the domain's `lib/` if shared across routes).
