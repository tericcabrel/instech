# OpenAPI + Orval contract workflow (current, docs-first)

This project currently uses a manual contract synchronization workflow between backend and frontend.

## Scope and current constraints

- `apps/backend/openapi.yaml` is the API contract source of truth.
- Frontend API client code is generated with Orval from that contract.
- `apps/frontend/src/api/generated/` is ignored from versioning in the current deferred setup.
- There is no CI gate yet that enforces OpenAPI updates or generated-client freshness.

## Required flow for backend API changes

When changing backend API behavior (new endpoint, request body shape, response field rename/removal, or route signature change), run this flow before opening a PR:

1. Update `apps/backend/openapi.yaml` to match backend behavior.
2. Regenerate frontend API artifacts:

```bash
yarn --cwd apps/frontend api:generate
```

3. Verify frontend compiles with the new contract:

```bash
yarn --cwd apps/frontend tsc --noEmit
```

4. If TypeScript errors appear in frontend usage, fix those call sites to match the new contract.

## Practical implications

- Because generated artifacts are ignored, each developer and CI environment must run `api:generate` locally before validating frontend compilation.
- Contract mismatches are currently detected through local regeneration and frontend type checking, not through a dedicated contract-sync CI job.
- A future phase may reintroduce stricter CI enforcement once the team is ready.
