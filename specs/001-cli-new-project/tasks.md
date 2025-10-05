# Tasks: CLI `new` Project Command

**Input**: Design documents from `/specs/001-cli-new-project/`
**Prerequisites**: plan.md (required), research.md, data-model.md, contracts/

## Phase 3.1: Setup
- [X] T001 Create internal/adapters/handlers/cli structure and wire Cobra root in cmd/main.go (no business logic)
- [X] T002 Initialize Cobra dependency and update go.mod; ensure make lint/test targets cover CLI
- [X] T003 [P] Configure structured logging and ensure no secret logging in CLI adapter

## Phase 3.2: Tests First (TDD)
- [X] T004 [P] Contract test for CLI flags and precedence in tests/contract/cli_new_flags_test.go using contracts/cli.md
- [X] T005 [P] Integration test quickstart path creation in tests/integration/cli_new_path_test.go based on quickstart.md
- [X] T006 [P] Integration test init here in tests/integration/cli_new_here_test.go
- [X] T007 [P] Integration test config merge and overrides in tests/integration/cli_new_config_merge_test.go
- [X] T008 [P] Integration test dry-run behavior in tests/integration/cli_new_dry_run_test.go
- [X] T009 [P] Integration test idempotent re-run with and without --force in tests/integration/cli_new_idempotency_test.go
- [X] T010 [P] Integration test failure paths ensure non-zero exit and no partial state in tests/integration/cli_new_failure_paths_test.go
- [X] T011 [P] Integration test per-project git config includeIf entries in tests/integration/cli_new_gitconfig_test.go
- [X] T012 [P] Unit test for defaults in tests/unit/defaults_test.go (ensure documented defaults applied)
- [X] T013 [P] Integration test .gitignore/no secrets committed in tests/integration/cli_new_gitignore_test.go

## Phase 3.3: Core Implementation
- [X] T014 [P] Define ProjectDefinition and ConfigInput domain types in internal/core/domain/project_definition.go per data-model.md
- [X] T015 [P] Implement validation rules in internal/core/domain/project_validation.go (mutual exclusivity, name pattern, absolute path)
- [X] T016 Implement merge logic (file + flags, flags override) in internal/core/services/project_options.go
- [X] T017 Implement filesystem service to create `.project.hcl` and `.env` in internal/core/services/file_service.go
- [X] T018 Implement CLI adapter command `new` in internal/adapters/handlers/cli/new.go using Cobra; map flags to domain
- [X] T019 Wire command into root and ensure help/usage in internal/adapters/handlers/cli/root.go
- [X] T020 Error handling and messages consistent with constitution (user-centered, no secrets)

## Phase 3.4: Integration
- [X] T021 [P] Add `--dry-run`, `--force`, and `--allow-unknown` flag handling in internal/adapters/handlers/cli/new.go
- [X] T022 [P] Per-project git config setup stub aligned with governance in internal/gitconfig/service.go
- [X] T023 Ensure .gitignore respects .env and no secrets committed; update docs if needed

## Phase 3.5: Polish
- [X] T024 [P] Unit tests for validation rules in tests/unit/validate_test.go
- [X] T025 [P] Update quickstart with examples for new flags in specs/001-cli-new-project/quickstart.md
- [X] T026 [P] Man page snippet for `pm new` in man/pm-new.1.md
- [X] T027 Performance check: command startup under 200ms and non-blocking file ops
- [X] T028 Run make lint, make test, and ensure reproducible build

## Dependencies
- Tests (T004-T013) before implementation (T014-T020)
- T014 blocks T015, T016
- T016 blocks T017, T018
- CLI wiring (T018, T019) after validation and merge implemented
- Integrations after core

## Parallel Example
```
# Launch independent tests together:
Task: "Contract test for CLI flags and precedence in tests/contract/cli_new_flags_test.go"
Task: "Integration test quickstart path creation in tests/integration/cli_new_path_test.go"
Task: "Integration test init here in tests/integration/cli_new_here_test.go"
Task: "Integration test config merge and overrides in tests/integration/cli_new_config_merge_test.go"
Task: "Integration test dry-run behavior in tests/integration/cli_new_dry_run_test.go"
Task: "Integration test idempotent re-run in tests/integration/cli_new_idempotency_test.go"
Task: "Integration test failure paths in tests/integration/cli_new_failure_paths_test.go"
Task: "Integration test per-project git config in tests/integration/cli_new_gitconfig_test.go"
```

## Validation Checklist
- [ ] All contracts have corresponding tests (T004)
- [ ] All entities have model tasks (T014)
- [ ] All tests come before implementation
- [ ] Parallel tasks truly independent
- [ ] Each task specifies exact file path
- [ ] No task modifies same file as another [P] task


---
