# Agent Quickstart
## Build / Lint / Test
- Build snapshot: `make build` (goreleaser snapshot).
- Run CLI locally: `make run args="pm list"`.
- Lint or autofix: `make lint` / `make lint-fix`.
- Full tests: `go test ./...`.
- Unit suite: `make unit` (`go test ./tests/... -tags unit`).
- Integration suite: `make integration` (`go test ./tests/... -tags integration`).
- Single test: `go test ./tests/<pkg> -run TestName -tags unit` (swap tag).
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
- Docs via `make gendocs`; no Cursor/Copilot rules, follow this guide.
