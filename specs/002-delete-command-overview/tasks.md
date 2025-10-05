# Tasks: CLI `pm delete`

**Input**: Design documents from `/workspaces/project-manager/specs/002-delete-command-overview/`
**Prerequisites**: plan.md (required), research.md, data-model.md, contracts/, quickstart.md

## Phase 3.1: Setup
- [X] T001 Verify branch `002-delete-command-overview` is checked out and workspace clean (`/workspaces/project-manager`).
- [X] T002 Ensure Go toolchain 1.25 and project linters are ready by running `make lint` dry-run (`/workspaces/project-manager/Makefile`).
- [X] T003 Create CLI delete command package directory `internal/adapters/handlers/cli/delete` with placeholder doc.go.

## Phase 3.2: Tests First (TDD)
- [X] T004 [P] Add contract test covering `pm delete` help/flag synopsis in `/workspaces/project-manager/tests/contract/cli_delete_flags_test.go`.
- [X] T005 [P] Add unit tests for flag parsing and scope validation in `/workspaces/project-manager/tests/unit/cli_delete_flags_test.go`.
- [X] T006 [P] Add unit tests for `ProjectService.DeleteProject` orchestration using mocks in `/workspaces/project-manager/tests/unit/project_service_delete_test.go`.
- [X] T007 [P] Add unit tests for filesystem backup planner, covering dry-run behaviour that MUST avoid archive creation, in `/workspaces/project-manager/tests/unit/filesystem_delete_test.go`.
- [X] T008 [P] Add unit tests asserting CLI exit codes for success, dry-run, and failure flows in `/workspaces/project-manager/tests/unit/cli_delete_exit_codes_test.go`.
- [X] T009 [P] Add unit tests ensuring missing-target requests surface not-found errors without invoking destructive ports in `/workspaces/project-manager/tests/unit/project_service_delete_notfound_test.go`.
- [X] T010 [P] Add integration test for name-based deletion success in `/workspaces/project-manager/tests/integration/cli_delete_smoke_test.go`.
- [X] T011 [P] Add integration test for path-based resolution in `/workspaces/project-manager/tests/integration/cli_delete_path_test.go`.
- [X] T012 [P] Add integration test for `--dry-run --all` preview output in `/workspaces/project-manager/tests/integration/cli_delete_dry_run_test.go`.
- [X] T013 [P] Add integration test for `--keep-files` scope in `/workspaces/project-manager/tests/integration/cli_delete_keep_files_test.go`.
- [X] T014 [P] Add integration test for `--only-env` scope in `/workspaces/project-manager/tests/integration/cli_delete_only_env_test.go`.
- [X] T015 [P] Add integration test for `--all --backup --force` success in `/workspaces/project-manager/tests/integration/cli_delete_backup_success_test.go`.
- [X] T016 [P] Add integration test for backup failure aborting deletion in `/workspaces/project-manager/tests/integration/cli_delete_backup_failure_test.go`.
- [X] T017 [P] Add integration test for `--dry-run --backup` ensuring no archive file is created while reporting the destination in `/workspaces/project-manager/tests/integration/cli_delete_dry_run_backup_test.go`.
- [X] T018 [P] Add integration test for read-only/locked workspace rejection in `/workspaces/project-manager/tests/integration/cli_delete_readonly_test.go`.
- [X] T019 [P] Add integration test for partial cleanup refusal (filesystem failure) in `/workspaces/project-manager/tests/integration/cli_delete_cleanup_failure_test.go`.
- [X] T020 [P] Add integration test for confirmation bypass via `--force` in `/workspaces/project-manager/tests/integration/cli_delete_force_test.go`.
- [X] T021 [P] Add integration test for conflicting scope flag error messaging in `/workspaces/project-manager/tests/integration/cli_delete_conflict_test.go`.
- [X] T022 [P] Add integration test for missing-target inputs returning non-zero exit codes and descriptive errors in `/workspaces/project-manager/tests/integration/cli_delete_missing_target_test.go`.
- [X] T023 [P] Add integration test verifying CLI exit codes for success, dry-run, and failure scenarios in `/workspaces/project-manager/tests/integration/cli_delete_exit_codes_test.go`.

## Phase 3.3: Core Implementation (after tests fail)
- [X] T024 [P] Implement `ProjectDeleteOptions` struct and helpers in `/workspaces/project-manager/internal/core/domain/project_delete_options.go`.
- [X] T025 [P] Implement `ProjectIdentifier` type with lock/read-only metadata in `/workspaces/project-manager/internal/core/domain/project_identifier.go`.
- [X] T026 [P] Implement `DeleteScope` enum and validation utilities in `/workspaces/project-manager/internal/core/domain/delete_scope.go`.
- [X] T027 [P] Implement `ProjectDeletePlan` builder structures in `/workspaces/project-manager/internal/core/domain/project_delete_plan.go`.
- [X] T028 [P] Implement `DeletionArtifact` definitions in `/workspaces/project-manager/internal/core/domain/deletion_artifact.go`.
- [X] T029 [P] Implement `BackupArtifact` struct with staging metadata in `/workspaces/project-manager/internal/core/domain/backup_artifact.go`.
- [X] T030 [P] Implement `ProjectDeleteResult` summary struct in `/workspaces/project-manager/internal/core/domain/project_delete_result.go`.
- [X] T031 [P] Extend filesystem port with delete/backup methods in `/workspaces/project-manager/internal/core/ports/filesystem_port.go`.
- [X] T032 [P] Extend env vars port with delete method in `/workspaces/project-manager/internal/core/ports/env_vars_port.go`.
- [X] T033 [P] Extend git port with hook cleanup support in `/workspaces/project-manager/internal/core/ports/git_port.go`.
- [X] T034 [P] Extend project port with delete/lock status APIs in `/workspaces/project-manager/internal/core/ports/project_port.go`.
- [X] T035 Update project service interface to expose `DeleteProject` in `/workspaces/project-manager/internal/core/ports/project_port.go`.
- [X] T036 Implement option parsing helpers in `/workspaces/project-manager/internal/core/services/project_options.go`.
- [X] T037 Implement `DeleteProject` orchestration, dry-run planning, and logging in `/workspaces/project-manager/internal/core/services/project_service.go`.
- [X] T038 Add backup writer utility in `/workspaces/project-manager/internal/core/services/project_service.go` or companion file.
- [X] T039 Update filesystem repository with recursive deletion and archive staging in `/workspaces/project-manager/internal/adapters/repositories/filesystem.go`.
- [X] T040 Update env vars repository with delete operations in `/workspaces/project-manager/internal/adapters/repositories/env_vars_repository.go`.
- [X] T041 Update git repository to remove include-if and hooks in `/workspaces/project-manager/internal/adapters/repositories/git_repository.go`.
- [X] T042 Update project repository to drop registry entries and set status flags in `/workspaces/project-manager/internal/adapters/repositories/project_repository.go`.
- [ ] T043 Regenerate mocks for updated ports (`/workspaces/project-manager/mocks/*`).
- [X] T044 Implement CLI confirmation helper in `/workspaces/project-manager/internal/adapters/handlers/cli/delete/confirm.go`.
- [X] T045 Implement Cobra command with flag setup and validation in `/workspaces/project-manager/internal/adapters/handlers/cli/delete/command.go`.
- [X] T046 Map CLI options to service call and render summaries in `/workspaces/project-manager/internal/adapters/handlers/cli/delete/run.go`.
- [X] T047 Register delete command in root command in `/workspaces/project-manager/internal/adapters/handlers/cli/root.go`.
- [X] T048 Wire new dependencies in bootstrap container in `/workspaces/project-manager/internal/bootstrap/bootstrap.go`.

## Phase 3.4: Integration
- [ ] T049 Ensure logger adapter emits structured delete events in `/workspaces/project-manager/internal/adapters/repositories/logger.go` if required.
- [ ] T050 Add backup directory initialization to configuration defaults in `/workspaces/project-manager/configs/config.default.hcl` or applicable config.
- [ ] T051 Verify doc generator picks up new command by updating `/workspaces/project-manager/tools/docgen/main.go` if needed.

## Phase 3.5: Polish
- [ ] T052 [P] Add automated test coverage confirming logger emits per-artifact entries in `/workspaces/project-manager/tests/integration/cli_delete_logging_test.go`.
- [ ] T053 [P] Add unit tests for new domain validation helpers in `/workspaces/project-manager/tests/unit/domain_delete_validation_test.go`.
- [ ] T054 [P] Update CLI docs with delete command usage in `/workspaces/project-manager/docs/cli/pm_delete.md` and `/workspaces/project-manager/docs/rest/pm_delete.rst`.
- [ ] T055 [P] Add `pm delete` man page in `/workspaces/project-manager/man/pm-delete.1` and link from `/workspaces/project-manager/man/pm.1`.
- [ ] T056 [P] Update README and changelog entries referencing delete command in `/workspaces/project-manager/README.md` and `/workspaces/project-manager/docs/cli/pm.md`.
- [ ] T057 Run `make lint`, `make unit`, and `make integration` to confirm all suites pass (`/workspaces/project-manager`).
- [ ] T058 Capture release notes snippet for delete command in `/workspaces/project-manager/docs/refactor_plan.md` or relevant summary.

## Dependencies
- T001 → T002 (workspace readiness before lint check).
- T003 depends on setup; T004-T023 require T003 so the CLI delete package exists for test scaffolding.
- Tests (T004-T023) must complete before implementations (T024-T048).
- Domain types (T024-T030) precede port/service updates (T031-T038).
- Port extensions (T031-T035) precede mock regeneration (T043) and repository updates (T039-T042).
- Service orchestration (T036-T038) and repository tasks (T039-T042) must finish before CLI adapter tasks (T044-T047).
- Bootstrap wiring (T048) follows CLI command integration.
- Integration polish (T049-T051) depends on core implementation completion (T024-T048).
- Logging verification test (T052) depends on service logging (T037) and logger adapter work (T049).
- Remaining polish tasks (T053-T058) run after integration tasks.

## Parallel Execution Examples
```
# Kick off contract/unit scaffolding together once setup done
/task run --ids T004,T005,T006,T007,T008,T009

# Execute independent integration specs concurrently
/task run --ids T010,T011,T012,T013
/task run --ids T014,T015,T016,T017
/task run --ids T018,T019,T020,T021,T022,T023

# After domain scaffolding, generate ports in parallel
/task run --ids T024,T025,T026,T027,T028,T029,T030
```

## Notes
- [P] tasks operate on distinct files or generated artifacts; ensure no shared edits before parallel execution.
- Maintain TDD discipline: run tests after each implementation milestone to confirm previously failing tests now pass.
- Regenerate mocks (T038) after port signatures change to keep unit tests compiling.
- Record manual confirmation flows while implementing CLI to aid documentation updates later.
