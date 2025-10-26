# Task Plan: `pm new` Environment Enhancements

## Phase 1 – Setup

- [x] **T001** – Checkout branch `004-new-env-command` and ensure Go 1.25 toolchain & `make` targets are ready.  
  _Files_: git workspace, local toolchain  
  _Notes_: Required for all subsequent phases.

## Phase 2 – Foundational Tasks

- [x] **T002** – Update `internal/core/domain/project_definition.go` and related domain types to add `Environment *EnvironmentInput` pointer and define the new `EnvironmentInput` struct with validation tags.  
  _Notes_: Enables environment data flow for later stories.
- [x] **T003** – Extend `internal/core/domain/project_definition.go` and `internal/core/domain/project.go` to remove the legacy `Environments map[string]string` usage and wire the new environment pointer through domain constructors.  
  _Notes_: Cleans legacy schema before story work begins.

## Phase 3 – User Story US1 (P1)
**Goal**: As a CLI user, I need `pm new` to detect existing projects reliably and avoid destructive reinitialization.

**Independent Test Criteria**: Running `pm new demo` twice should report the second run as a no-op, while recreating after removing `.project.hcl` should succeed without duplicate include entries.

- [x] **T004 [Story US1]** – Implement a project existence probe in `internal/core/services/project_service.go` that verifies registry membership and `.project.hcl` presence.  
  _Notes_: Provides shared detection logic.  
- [x] **T005 [Story US1]** – Refactor `internal/adapters/handlers/cli/new/new.go` to use the new existence probe, treating the second positional argument as path vs environment based on probe result and emitting the no-op message for repeat invocations.  
  _Depends on_: T004.  
- [x] **T006 [Story US1]** – Update uniqueness helper logic in `internal/core/services/project_options.go` to ignore registry entries missing `.project.hcl` and return explicit error codes for real conflicts.  
  _Notes_: Must run before repository updates.  
- [x] **T007 [Story US1]** – Adjust `internal/adapters/repositories/project_repository.go` to skip writing duplicate includeIf entries when the same path is already registered and keep stale entries untouched.  
  _Depends on_: T006.  
- [x] **T008 [Story US1]** – Refresh integration tests in `tests/integration/cli_new_*.go` and unit coverage for project existence helpers to assert no-op messaging and stale-entry tolerance.  
- [x] **T009 [Story US1]** – Extend `tests/integration/cli_new_here_test.go` (and related suites) to confirm `pm new --here` still rejects additional positional arguments or environment names.  

  _Depends on_: T005.  

_Checkpoint_: US1 completed when rerun behavior, stale registry handling, and `--here` guard verification are all in place.

## Phase 4 – User Story US2 (P1)
**Goal**: As a CLI user, I can add environments via `pm new <name> <env>` with new flags and config schema, ensuring the project already exists.

**Independent Test Criteria**: Attempting to add an environment before project creation fails; adding with flags/config succeeds; `--force` overwrites env metadata/file; legacy maps ignored silently under `--allow-unknown`.

- [x] **T010 [Story US2]** – Introduce new flags (`--environment-env-file`, `--environment-mode`, `--environment-color`, `--environment-env-var`) in `internal/adapters/handlers/cli/new/new.go` and update flag help text.  
- [x] **T011 [Story US2]** – Extend merge logic in `internal/core/services/project_options.go` to assemble an `EnvironmentInput` from config object and CLI flags, including precedence for `--environment-env-var`.  
- [x] **T012 [Story US2]** – Enforce project-first rule in `internal/adapters/handlers/cli/new/new.go`, returning a descriptive error when environment input is provided but the existence probe fails.  
- [x] **T013 [Story US2]** – Update `internal/core/services/project_service.go` to consume `EnvironmentInput`, validate duplicates, and honor `--force` by rewriting metadata/env files.  
- [x] **T014 [Story US2]** – Modify config loading in `internal/core/services/project_options.go` to drop legacy `environments` map, ignore it when `--allow-unknown` is set, and surface errors otherwise.  
- [x] **T015 [Story US2]** – Update unit and integration tests covering environment creation (`tests/unit/project_service_*`, `tests/integration/cli_new_*`) to cover flag parsing, `--force`, legacy map handling, and config precedence.  

  _Depends on_: T010–T014.  

_Checkpoint_: US2 completed when environment creation behaves as specified and tests pass.

## Phase 5 – User Story US3 (P2)
**Goal**: As a CLI user, I can preview environment changes and scaffold configs that mirror the new schema.

**Independent Test Criteria**: Dry-run output shows project + environment details in text/JSON; skeleton generators emit the new `environment` object; documentation reflects workflow.

- [x] **T016 [Story US3]** – Enhance dry-run rendering in `internal/adapters/handlers/cli/new/new.go` to output merged environment fields in both text and JSON modes.  
- [x] **T017 [Story US3]** – Revise skeleton generation in `internal/core/services/project_options.go` to produce the new environment object and remove the legacy map.  
- [x] **T018 [Story US3]** – Update CLI docs (`docs/cli/pm_new.md`) and quickstart examples (`specs/004-new-env-command/quickstart.md`) to describe dry-run output, new flags, and skeleton structure.  
- [x] **T019 [Story US3]** – Adjust skeleton and dry-run related tests (`tests/integration/cli_new_dry_run_test.go`, `tests/integration/cli_new_skeleton_test.go`) to validate the new outputs.  

  _Depends on_: T016, T017.  

_Checkpoint_: US3 completed when previews, skeletons, and docs align with the new environment schema.

## Phase 6 – Polish & Cross-Cutting

- [ ] **T020** – Run `go test ./...`, `make lint`, and apply gofmt/gofumpt across touched files; fix any issues surfaced by linters or tests.  
- [x] **T021** – Perform final documentation audit (README.md, release notes if applicable) to ensure messaging matches new behavior before handoff.  

  _Depends on_: T018.  

## Dependencies Overview

1. Phase 1 → Phase 2 → US1 → US2 → US3 → Polish.  
2. US1 (T004–T009) must complete before US2 tasks begin; US2 (T010–T015) completion is required for US3 (T016–T019).  
3. Foundational tasks (T002–T003) block both US2 and US3 since they introduce the shared data model.

## Parallel Execution Examples

- Within US1, after T006 completes, T007 (repository adjustments) and T008 (existence tests) can proceed in parallel; T009 (here regression) remains focused on integration coverage.  
- Within US2, T010 (flag wiring) and T011 (merge logic) touch different files and may run concurrently once dependencies resolve.  
- Within US3, T016 (handler dry-run) and T017 (skeleton generator) operate on separate files and can execute in parallel.

## Implementation Strategy

- **MVP Scope**: Complete US1 to stabilize project detection and rerun safety; this provides immediate user value and unblocks environment enhancements.  
- **Incremental Delivery**: Ship US2 next to unlock environment customization, followed by US3 for user-facing previews and documentation polish.  
- **Risk Mitigation**: Prioritize core CLI changes and domain model updates before touching docs/tests; run the full test suite (T020) after each major story completion.
