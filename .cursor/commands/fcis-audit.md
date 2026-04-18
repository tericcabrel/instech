---
description: Structure-agnostic Functional Core / Imperative Shell architecture review
---

You are performing a Functional Core / Imperative Shell (FC/IS) architecture review.

## Scope
- Argument received: $ARGUMENTS
- If `$ARGUMENTS` is non-empty, review ONLY that folder/file.
- If `$ARGUMENTS` is empty, review the entire repository.
- Never fail just because structure/naming differs from common conventions.

## How to Review (Structure-Agnostic)
1. Discover layers by behavior (not folder names):
   - Functional core candidates: pure business rules, validation, transformations, decision logic.
   - Imperative shell candidates: HTTP/controllers/routes, DB access, filesystem/env, network, external SDKs, logging.
   - Orchestration candidates: use cases/services coordinating core + shell.

2. Evaluate FC/IS boundaries:
   - Flag business logic mixed with side effects.
   - Flag transport/infra making domain decisions.
   - Flag orchestration bypassing domain rules.
   - Flag inconsistent error-to-HTTP mapping across similar endpoints.
   - Flag pagination/list metadata inconsistencies after filtering/skipping.
   - Highlight strong FC/IS patterns where present.

3. Robustness:
   - If uncertain, state assumptions explicitly.
   - Use "possible issue" wording when confidence is low.
   - Continue with best-effort analysis even if parts are missing.

## Output Format (Strict)
1. Findings (High -> Medium -> Low), each with:
   - Severity
   - Why it matters
   - File path(s)
   - Short code reference/snippet
2. Open questions / assumptions
3. What is good
4. Top 3 practical next improvements
