# Project Manager Refactor Plan (Simplified Hexagonal)

Goal: Reduce repository responsibilities, keep clear boundaries, add tests to ensure behavior parity.

## Architecture targets
- Repositories:
  - ProjectRepository: only read/write .project.hcl and in-memory project aggregate
  - EnvVarsRepository: read/write .env files
  - GitRepository: read/write per-project .gitconfig and manage global includeIf
- Services: ProjectService orchestrates validations, defaults, ordering, rollback on failure
- Utilities: small helpers in internal/tools for path ops and dir checks (no new ports)

## Current coupling hotspots (references)
- Env file IO and naming in ProjectRepository:
  - Create: internal/adapters/repositories/project_repository.go:228-254, 255-283
  - AddEnvironment: internal/adapters/repositories/project_repository.go:410-474
  - UpdateEnvironment (rename): internal/adapters/repositories/project_repository.go:358-384
- Git global includes + per-project .gitconfig in ProjectRepository:
  - Create: internal/adapters/repositories/project_repository.go:163-195, 196-226
- Project persistence (keep here):
  - Load .project.hcl: internal/adapters/repositories/project_repository.go:478-499
  - Save .project.hcl: internal/adapters/repositories/project_repository.go:320-327

## Milestones and tasks

### 0) Prep and safety
- [x] T001: Add internal/tools/path.go with: ExpandHome, AbsPath, Join, IsAbs
- [x] T002: Add internal/tools/fs.go with: EnsureDir(path, mode), IsDirEmpty(path), Rename(old,new), WriteFile(path, bytes, mode)
- [x] T003: Move IsDirEmpty from project_repository.go into internal/tools/fs.go and update references
- [x] T004: Test: unit tests for tools/fs IsDirEmpty and EnsureDir using temp dirs

### 1) Extract Env vars responsibilities
- [x] T005: ProjectRepository.Create: delegate env file creation to EnvVarsRepository.Save (refs: 228-254)
- [x] T006: ProjectRepository.AddEnvironment: delegate env file creation and path resolution to EnvVarsRepository (refs: 410-441)
- [x] T007: ProjectRepository.UpdateEnvironment: delegate env file rename to tools/fs.Rename after resolving paths (refs: 358-384)
- [x] T008: Ensure ProjectService computes env file default name (slug) instead of repository (refs: 410-418)
- [x] T009: Tests:
  - [x] T010: EnvVarsRepository unit tests: Load/Save round-trip
  - [x] T011: ProjectService AddEnvironment: creates file with defaults, idempotency on name conflicts

### 2) Extract Git responsibilities
- [x] T012: Extend ports.GitRepository with global config helpers:
  - [x] T013: LoadGlobal() error
  - [x] T014: UpdateIncludeIf(gitdir, perProjectPath, subproject string) error
  - [x] T015: SaveGlobal(home string) error
- [x] T016: Implement in internal/adapters/repositories/git_repository.go using go-git config (reuse existing Load/Save for per-project files)
- [x] T017: ProjectRepository.Create: replace direct includeIf + ~/.gitconfig writes with GitRepository calls (refs: 163-195)
- [x] T018: ProjectRepository.Create: move per-project .gitconfig write to GitRepository.Save (refs: 196-226)
- [x] T019: Tests:
  - [x] T020: GitRepository per-project Save/Load round-trip
  - [x] T021: GitRepository global includeIf end-to-end using temp HOME (override HOME env var)
  - [x] T022: ProjectService.Create integration: ensures includeIf points to per-project path

### 3) Narrow ProjectRepository to .project.hcl persistence
- [x] T023: Keep HCL encoding/decoding but remove unrelated os.*, logging, and path expansion (use tools)
- [x] T024: Factor loadDotProject into a method and use tools.ExpandHome (refs: 478-499)
- [x] T025: Tests: ProjectRepository List/Get/UpdateProject using temp dirs and sample .project.hcl

### 4) Move validations and orchestration to ProjectService
- [x] T026: Create path from user input: expand ~, normalize abs (refs currently in Create 121-146)
- [x] T027: Ensure directory exists, require empty dir for new project
- [x] T028: Validate env uniqueness and default EnvVarsMode
- [x] T029: Orchestrate order: fs prepare -> env file (optional) -> git includes -> dotproject write
- [x] T030: Rollback on failure: track created files and revert renames
- [x] T031: Tests:
  - [x] T032: Create success path: project dir empty -> all artifacts created
  - [x] T033: Create failure during git step: env file rolled back, no partial artifacts left
  - [x] T034: UpdateEnvironment rename: handles relative/absolute paths

### 5) Logging and side effects
- [x] T035: Remove slog/log prints from repositories (refs: 43-48, 96-99, 147-175)
- [x] T036: Centralize logs in ProjectService with context

### 6) Handlers use services only
- [x] T037: Ensure TUI/CLI handlers depend on ports.ProjectService; no direct repository usage
- [x] T038: Smoke tests: commands still work via make test; add integration test where feasible

### 7) Naming and structure (optional, low risk)
- [x] T039: Keep adapters layout; skip renaming to outbound/inbound to avoid churn

## Acceptance criteria
- ProjectRepository contains only .project.hcl load/list/save logic; no env or git writes
- Env/ Git operations are delegated to their repositories
- All unit and service tests pass; no regressions in CLI/TUI
- No direct os or slog usage in repositories beyond necessary file IO already encapsulated or minimized

## Test matrix (summary)
- EnvVarsRepository: Load/Save; handles ~ expansion
- GitRepository: per-project Save/Load; global includeIf round-trip
- ProjectRepository: List/Get/UpdateProject with sample HCL
- ProjectService: Create/Add/UpdateEnvironment happy paths; rollback on injected failures

## Notes
- Use temp directories and HOME overrides in tests to avoid touching user files
- Keep ports count minimal; only extend GitRepository methods as above
