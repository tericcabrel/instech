# Project Guidelines

All rules are defined in [`docs/references/`](docs/references/) and apply to any AI agent or IDE.

## Always read before making changes

- [General Coding Conventions](docs/general-coding-conventions.md)
- [Backend Coding Conventions](docs/backend-coding-conventions.md)

## When working on TypeScript or TSX files

- [TypeScript Standards](docs/references/typescript.md)


## Plan Mode

- Make the plan extremely concise. Sacrifice grammar for the sake of concision.
- At the end of each plan, give me a list of unresolved questions to answer, if any.
- Don't write tests for what the type system already guarantees

## Plan-mode safety checklist

Before finalising any plan, audit for bypass / integrity flaws and call them out explicitly in a "Risks / flaws" section:

1. Every client-side validation, gating, or rate-limit must have an equivalent server-side enforcement; otherwise an attacker can bypass it by calling the API directly (curl, scripts, etc.).
2. Identify which checks are authoritative (server) versus UX-only (client).
3. If a check requires data only the server can fetch (third-party probe, secret, signed token), express it as router orchestration of use cases on the backend, not as a client round-trip the user is expected to make.
4. Flag remaining trade-offs in the plan's "Risks / flaws" section so the user can confirm.