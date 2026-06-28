<!-- intent-skills:start -->
## Skill Loading
If there is network permission issues, request the network permission from the user.

Before editing files for a substantial task:
- Run `yarn dlx @tanstack/intent@latest list` from the workspace root to see available local skills. 
- If a listed skill matches the task, run `yarn dlx @tanstack/intent@latest load <package>#<skill>` before changing files.
- Use the loaded `SKILL.md` guidance while making the change.
- Monorepos: when working across packages, run the skill check from the workspace root and prefer the local skill for the package being changed.
- Multiple matches: prefer the most specific local skill for the package or concern you are changing; load additional skills only when the task spans multiple packages or concerns.
<!-- intent-skills:end -->

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


## Architecture Decisions

- Preserve TanStack CLI generated structure as the source of truth.
- Keep the app intentionally minimal ("blank-style"), while retaining Query and Form demo routes to demonstrate integrations.
- Keep `src/router.tsx` SSR query wiring (`setupRouterSsrQueryIntegration`) intact.
- Keep TanStack Devtools wiring in root layout for local debugging.

## Feature Folder Conventions

- Keep route files thin; each route component should import and render one or more feature-level container components.
- Place feature modules under `src/features/<feature-name>`.
- Keep all container components at feature root using `*.container.tsx` naming.
- Place container-local subcomponents in dedicated subfolders inside the same feature folder.
- If a component is reused by multiple containers in one feature, move it to `src/features/<feature-name>/shared`.
- If a component is reused across multiple features, move it to `src/components`.
- A feature can contain multiple containers. Example:
  - `src/features/tool-graph/tool-graph-home.container.tsx`
  - `src/features/tool-graph/tool-graph-detail.container.tsx`
  - `src/features/tool-graph/tool-graph-home/`
  - `src/features/tool-graph/tool-graph-detail/`
  - `src/features/tool-graph/shared/`

## Dependency Guidelines

- Keep all direct `dependencies` and `devDependencies` pinned to exact versions (no `^`, `~`, `latest`, or ranges).
- Add or update packages with Yarn, then commit the synchronized pair: `package.json` and `yarn.lock`.
- Prefer updating TanStack packages as a coordinated set to avoid cross-version mismatches.
- After dependency changes, run `yarn install`, `yarn build`, and `yarn lint` to verify lockfile and compile health.

## How to Use Shadcn Components

- [Shadcn Components](docs/guidelines/shadcn-components.md)

## Styling Guidelines

- Do not create custom CSS classes in styles.css or inline styles. Use Tailwind CSS classes instead.

## Known Gotchas

- TanStack CLI currently reports `--tailwind` as deprecated/ignored because Tailwind is always enabled in Start scaffolds.
- If Intent list initially returns no skills, ensure `@tanstack/intent` is installed in the project and rerun `npx @tanstack/intent@latest list`.
- Do not edit route path strings in `createFileRoute(...)` manually; they must match route file paths.

## Execution and Verification Guidelines
1. After you have executed a plan, run: `yarn build` and `yarn lint` to verify the project is working.
2. If there are any linting errors that can be fixed safely, run `yarn lint:fix` to fix them.
3. To run tests: `yarn test`
4. To run test on specific file, append the filename after the command. For example: `yarn test -- <filename>`; example: `yarn test -- tool-graph.test.ts`
5. Never try to start a dev server or curl a local endpoint.
