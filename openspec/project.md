# Project Context

## Purpose
Project Manager (`pm`) is a terminal-first CLI/TUI that standardizes how developers register, launch, and retire local projects. It creates opinionated scaffolding for project metadata, environment variables, and Git configuration so teams can spin up consistent workspaces, switch between environments quickly, and remove stale projects without leaving residue on disk or in global config.

## Tech Stack
- Go 1.25 as the primary language targeting cross-platform builds
- Bubble Tea/Bubbles/Lipgloss for the interactive terminal UI
- Cobra + Viper for CLI command discovery, configuration, and flag parsing
- HashiCorp HCL and `godotenv` for project definitions and environment management
- GoReleaser-driven distribution backed by GitHub Actions for CI/CD automation

## Project Conventions

### Code Style
All Go code is formatted with `gofmt`, `gofumpt --extra-rules`, and `golines --max-len 120`, with imports grouped via `gci` (stdlib, third-party, module). Packages stay small and focused, exported identifiers are documented, and errors are wrapped with contextual messages before returning. Services accept interfaces (ports) rather than concrete implementations, avoid global state, and follow test-friendly patterns.

### Architecture Patterns
The application follows a strict hexagonal architecture:
- `internal/core` owns domain types, validation rules, and service logic expressed against ports.
- `internal/adapters` provides driving adapters (CLI/TUI handlers) and driven adapters (filesystem, Git, logger, shell) that satisfy those ports.
- `internal/bootstrap` wires adapters into services, ensuring dependencies are injected at the edges.
This separation lets new frontends or infrastructure components plug in without leaking platform details into the domain.

### Testing Strategy
Tests live outside production packages under `tests/`. Unit suites (`make unit`) target domain services and ports using mocks, while integration suites (`make integration`) exercise CLI/TUI flows end-to-end against the filesystem, Git, and environment adapters. `make test` runs both suites, and CI additionally executes `go test -race -cover ./...` plus integration runs. Contributors should add or update unit tests alongside domain changes and prefer integration tests for user-visible CLI behavior.

### Git Workflow
Work happens on short-lived feature branches that open pull requests against `main`. Every PR runs lint, unit, integration, and snapshot build jobs in GitHub Actions before merge. Releases are cut by pushing semantic tags (`v*`) which trigger the GoReleaser workflow to publish binaries. Reference relevant OpenSpec change IDs in PR descriptions to keep specs and implementation in sync, and write descriptive commit messages that explain the motivation behind the change.

## Domain Context
`pm` maintains a registry of local projects backed by `.project.hcl` files, per-environment `.env` manifests, and optional Git overrides stored via `includeIf` entries. The TUI lists projects alongside their environments, supports creation/edit flows with keyboard-driven forms, and provides shortcuts to start sessions. CLI subcommands (`pm new`, `pm list`, `pm edit`, `pm delete`) mirror those capabilities for scripting and CI usage, including safe deletion with dry-run previews, backup destinations, and file preservation flags.

## Important Constraints
- Git must be available on the host, and modifying the user’s global `~/.gitconfig` is required for per-project include entries.
- Project creation validates target paths: parents must exist and destinations must be empty unless initializing an existing folder.
- Environment files are never committed by default; `.gitignore` entries are enforced to keep secrets local.
- Destructive operations (deleting workspaces) require explicit flags (`--all`, `--force`) so defaults remain non-destructive.
- The CLI/TUI assumes ANSI-capable terminals; degraded behavior may occur on minimal terminals without alternate-screen support.

## External Dependencies
- System Git for include rules, status checks, and repository operations
- Local filesystem access for creating project directories, HCL manifests, and environment files
- Optional GPG tooling when users enable commit or tag signing during project setup
- Go toolchain, GoReleaser, and GitHub Actions for builds, linting, and release automation
