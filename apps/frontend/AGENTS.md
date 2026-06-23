<!-- intent-skills:start -->
## Skill Loading

Before editing files for a substantial task:
- Run `yarn dlx @tanstack/intent@latest list` from the workspace root to see available local skills.
- If a listed skill matches the task, run `yarn dlx @tanstack/intent@latest load <package>#<skill>` before changing files.
- Use the loaded `SKILL.md` guidance while making the change.
- Monorepos: when working across packages, run the skill check from the workspace root and prefer the local skill for the package being changed.
- Multiple matches: prefer the most specific local skill for the package or concern you are changing; load additional skills only when the task spans multiple packages or concerns.
<!-- intent-skills:end -->

## Frontend Project Context

- Scaffold source command:
  - `npx @tanstack/cli@latest create my-tanstack-app --agent --package-manager yarn --tailwind --add-ons tanstack-query,form`
- Follow-up TanStack Intent commands executed:
  - `npx @tanstack/intent@latest install`
  - `npx @tanstack/intent@latest list`
  - Loaded guidance before library-structure edits:
    - `npx @tanstack/intent@latest load @tanstack/cli#create-app-scaffold`
    - `npx @tanstack/intent@latest load @tanstack/start-client-core#start-core`
    - `npx @tanstack/intent@latest load @tanstack/router-core#router-core`

## Chosen Stack And Integrations

- Framework: React + TanStack Start
- Router: TanStack Router (file-based route generation)
- Data fetching: TanStack Query (`@tanstack/react-query`, SSR query integration enabled)
- Forms: TanStack Form (`@tanstack/react-form`) demo routes included
- Tooling: Yarn, Vite, TypeScript, Tailwind (default TanStack Start scaffold toolchain)
- Developer tooling: TanStack CLI and TanStack Intent installed in `devDependencies`

## Environment Variables

- No required runtime environment variables for local development in this scaffold.
- If adding client-side env vars, use `VITE_*` prefix (for example `VITE_API_BASE_URL`).
- Keep secrets server-only and access them through server-side code paths.

## Deployment Notes

- Build command: `yarn build`
- Preview production build: `yarn preview`
- Generated output follows TanStack Start defaults for Vite-based deployments.

## Architecture Decisions

- Preserve TanStack CLI generated structure as the source of truth.
- Keep the app intentionally minimal ("blank-style"), while retaining Query and Form demo routes to demonstrate integrations.
- Keep `src/router.tsx` SSR query wiring (`setupRouterSsrQueryIntegration`) intact.
- Keep TanStack Devtools wiring in root layout for local debugging.

## Dependency Guidelines

- Keep all direct `dependencies` and `devDependencies` pinned to exact versions (no `^`, `~`, `latest`, or ranges).
- Add or update packages with Yarn, then commit the synchronized pair: `package.json` and `yarn.lock`.
- Prefer updating TanStack packages as a coordinated set to avoid cross-version mismatches.
- After dependency changes, run `yarn install` and `yarn build` to verify lockfile and compile health.

## Known Gotchas

- TanStack CLI currently reports `--tailwind` as deprecated/ignored because Tailwind is always enabled in Start scaffolds.
- If Intent list initially returns no skills, ensure `@tanstack/intent` is installed in the project and rerun `npx @tanstack/intent@latest list`.
- Do not edit route path strings in `createFileRoute(...)` manually; they must match route file paths.

## Next Steps

- Run `yarn dev` and verify:
  - `/` renders the minimal home page
  - `/demo/tanstack-query` shows query data
  - `/demo/form/simple` and `/demo/form/address` render form demos
- Add app-specific routes and server functions as needed.
