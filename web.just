set working-directory := 'apps/frontend'

# Start Vite dev server on port 8800
[group('web')]
dev:
  yarn dev

# Production build
[group('web')]
build:
  yarn build

# Serve the production build locally
[group('web')]
preview:
  yarn preview

# Regenerate API client from OpenAPI spec (Orval)
[group('web')]
api-generate:
  yarn api:generate

# Regenerate TanStack Router route tree
[group('web')]
generate-routes:
  yarn generate-routes

# Install TanStack Intent skills
[group('web')]
intent-install:
  yarn intent:install

# List available TanStack Intent skills
[group('web')]
intent-list:
  yarn intent:list

# Run ESLint and Biome
[group('web')]
lint:
  yarn lint

# Run Biome only
[group('web')]
lint-biome:
  yarn lint:biome

# Run ESLint only
[group('web')]
lint-eslint:
  yarn lint:eslint

# Auto-fix ESLint and Biome issues
[group('web')]
lint-fix:
  yarn lint:fix

# Lint staged files (pre-commit)
[group('web')]
lint-staged:
  yarn lint:staged

# Scaffold a new TanStack app
[group('web')]
tanstack-create:
  yarn tanstack:create

# Run Vitest once (append -- file to filter)
[group('web')]
test *args:
  yarn test {{args}}
