<!-- OPENSPEC:START -->
# OpenSpec Instructions

These instructions are for AI assistants working in this project.

Always open `@/openspec/AGENTS.md` when the request:
- Mentions planning or proposals (words like proposal, spec, change, plan)
- Introduces new capabilities, breaking changes, architecture shifts, or big performance/security work
- Sounds ambiguous and you need the authoritative spec before coding

Use `@/openspec/AGENTS.md` to learn:
- How to create and apply change proposals
- Spec format and conventions
- Project structure and guidelines

Keep this managed block so 'openspec update' can refresh the instructions.

<!-- OPENSPEC:END -->

# Agent Quickstart
## Build / Lint / Test
- Build snapshot: `make build` (goreleaser snapshot).
- Run CLI locally: `make run args="pm list"`.
- Lint or autofix: `make lint` / `make lint-fix`.
- Full tests: `make test` (`go test ./...`).
- Unit suite: `make unit` (`go test ./tests/... -tags unit`).
- Integration suite: `make integration` (`go test ./tests/... -tags integration`).
- Single test: `go test ./tests/<pkg> -run TestName -tags unit` (swap tag).
- Mock generation: `make generate`.
## Formatting & Imports
- Apply `gofmt`, `gofumpt --extra-rules`, and `golines --max-len 120`.
- Keep imports ordered via gci: std, third-party, project prefix.
## Code Style
- Go 1.25; write small, cohesive packages with clear interfaces.
- Exported identifiers need comments; avoid stuttered names.
- Errors must be checked, wrapped with context, and surfaced upward.
- Prefer context-aware APIs; no global state; keep functions pure when possible.
- Watch for security lint (gosec) and duplication (dupl); refactor early.
## Hexagonal Architecture
- Keep `internal/core` domain-centric with pure types, services, and ports (interfaces) only.
- Implement adapters under `internal/adapters` to fulfill ports; avoid domain imports from adapters other than interfaces.
- Route all infrastructure wiring through `internal/bootstrap`; inject dependencies via constructors.
- Drive features from adapters into domain via ports, keeping business logic in services and repositories.
- Favor small, focused packages; add new ports before introducing new adapters.

## Misc
- Docs via `make gendocs` (update CLI docs and man pages using this command only); no Cursor/Copilot rules, follow this guide.
