mod web
mod api

default:
  @just --list --list-submodules

# Start web and api dev servers together
[group('repo')]
dev:
  [parallel]
  just api dev
  just web dev

# Run full CI checks locally (lint, tests, OpenAPI)
[group('repo')]
ci:
  just api lint
  just api test-integration
  just api openapi-lint
  just web lint
  just web test
