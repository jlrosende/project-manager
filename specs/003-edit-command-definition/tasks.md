# Tasks: pm edit flag-driven updates

**Input**: Design documents from `/specs/003-edit-command-definition/`
**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/

**Tests**: Add targeted unit, contract, and integration coverage alongside functional changes to uphold repository standards.

**Organization**: Tasks are grouped by user story to enable independent implementation and validation.

## Format: `[ID] [P?] [Story] Description`
- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Project initialization and basic structure

No additional setup work is required; existing CLI infrastructure and tooling remain unchanged.

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core infrastructure that MUST be complete before ANY user story can be implemented

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [x] T001 [Foundation] Add edit change-set domain structs capturing optional project/environment mutations in `internal/core/domain/edit_changeset.go`.
- [x] T002 [Foundation] Extend project port interfaces with scoped load/save and locking methods in `internal/core/ports/project_port.go`, documenting that the edit command continues to require a `<project>` argument via the shared Cobra guard.
- [x] T003 [Foundation] Implement new port methods with atomic writes and lock acquisition in `internal/adapters/repositories/project_repository.go` (depends on T002).
- [x] T004 [Foundation] Wire change-set and locking dependencies through the bootstrap container in `internal/bootstrap/bootstrap.go` (depends on T001–T003).

**Checkpoint**: Foundation ready - user story implementation can now begin in parallel

---

## Phase 3: User Story 1 - Project configuration editing (Priority: P1) 🎯 MVP

**Goal**: Allow maintainers to modify shared project settings using flags or CLI input files without affecting environments.

**Independent Test**: Run `pm edit <project>` with field updates, confirm project metadata changes and environments remain untouched.

### Implementation for User Story 1

- [x] T005 [US1] Register project-level flags, `--cli-input`, and skeleton options in `internal/adapters/handlers/cli/edit/edit.go` (depends on T004).
- [x] T006 [US1] Implement project-scope CLI input parsing, skeleton generation, and change application in `internal/core/services/project_service.go` (depends on T001, T004).
- [x] T007 [US1] Emit project edit summaries and dry-run output in `internal/adapters/handlers/cli/edit/edit.go` (depends on T005, T006).
- [x] T017 [US1] Guarantee cancel/no-save exits leave project configuration untouched and surface "no changes applied" messaging in `internal/adapters/handlers/cli/edit/edit.go` and related service logic (depends on T006, T007).
- [x] T018 [US1] Add contract or integration tests covering project-only flag/file edits, including cancel/no-save flows, in `tests/integration/cli_edit_project_*` (depends on T005–T007, T017).

**Checkpoint**: User Story 1 is functional and independently verifiable

---

## Phase 4: User Story 2 - Environment-specific editing (Priority: P2)

**Goal**: Allow release engineers to target a single environment with flag or file-driven overrides while leaving other environments unchanged.

**Independent Test**: Run `pm edit <project> <env>` updating environment overrides and verify no other environments change.

### Implementation for User Story 2

- [x] T008 [US2] Extend change application logic to support environment-scoped updates in `internal/core/services/project_service.go` (depends on T006).
- [x] T009 [US2] Expose environment flag registration and scope routing in `internal/adapters/handlers/cli/edit/edit.go` (depends on T005, T008).
- [x] T010 [US2] Persist environment-specific configuration without touching other environments in `internal/adapters/repositories/project_repository.go` (depends on T003, T008).
- [x] T019 [US2] Emit environment confirmation output summarizing updated overrides and unchanged scopes in `internal/adapters/handlers/cli/edit/edit.go` (depends on T009, T010).
- [x] T020 [US2] Add integration tests covering environment-scoped edits and confirmation messaging in `tests/integration/cli_edit_environment_*` (depends on T008–T010, T019).

**Checkpoint**: User Stories 1 and 2 are both functional and independently verifiable

---

## Phase 5: User Story 3 - Validation feedback (Priority: P3)

**Goal**: Guarantee invalid edits are rejected with actionable feedback and no partial persistence.

**Independent Test**: Attempt an invalid edit (e.g., clearing required field) and confirm the command reports the issue and leaves files unchanged.

### Implementation for User Story 3

- [x] T011 [US3] Aggregate validation errors and enforce immutability checks in `internal/core/services/project_service.go` (depends on T006, T008).
- [x] T012 [US3] Map domain errors to exit codes and user-facing messages in `internal/adapters/handlers/cli/edit/edit.go` (depends on T005, T011).
- [x] T013 [US3] Guard repository writes so failed validations leave filesystem state untouched in `internal/adapters/repositories/project_repository.go` (depends on T003, T011).
- [x] T021 [US3] Add integration and unit tests verifying validation rejection paths keep state unchanged and report detailed errors in `tests/integration/cli_edit_validation_*` (depends on T011–T013).

**Checkpoint**: All user stories are functional with robust validation feedback

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Improvements that affect multiple user stories

- [x] T014 [P] Regenerate CLI docs via `make gendocs` and verify updates in `docs/cli/pm_edit.md` and `man/pm-edit.1` reflect new flags (depends on T005–T021).
- [x] T015 [P] Re-run quickstart steps and adjust guidance in `specs/003-edit-command-definition/quickstart.md` if behavior changed (depends on T005–T021).
- [x] T016 Run `go test ./...` to confirm repository health after changes (depends on T005–T021).
- [x] T022 [P] Collect first-attempt success feedback from maintainers and document findings in `specs/003-edit-command-definition/quickstart.md` or a linked follow-up note (depends on T005–T021).

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No tasks – existing infrastructure reused
- **Foundational (Phase 2)**: Must complete before User Stories 1–3
- **User Story Phases (3–5)**: Can start after Phase 2; proceed in priority order (P1 → P2 → P3) for MVP delivery
- **Polish (Phase 6)**: Runs after desired user stories are complete

### User Story Dependencies

- **US1** depends on Phase 2
- **US2** depends on Phase 2 and completion of US1 service groundwork (T006)
- **US3** depends on Phases 2–4 (needs validation hooks from earlier work)

### Within Each User Story

- Complete CLI updates before service or repository adjustments that rely on them
- Apply service changes before repository persistence tweaks within the same story
- Each story culminates in a checkpoint ensuring independent validation

### Parallel Opportunities

- T001 and T002 touch different files and can run in parallel once design is settled; T002 must finish before T003
- During US2, T009 and T010 operate on distinct files once T008 completes and may proceed in parallel
- Polish tasks T014 and T015 can run concurrently after functional work is done

---

## Parallel Example: User Story 2

```bash
# After T008 completes:
Task: "T009 [US2] Expose environment flag registration and scope routing in internal/adapters/handlers/cli/edit/edit.go"
Task: "T010 [US2] Persist environment-specific configuration without touching other environments in internal/adapters/repositories/project_repository.go"
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)
1. Complete Phase 2 foundational tasks
2. Deliver Phase 3 (US1) functionality and verify independently
3. Optionally ship MVP after confirming validations for project-level edits

### Incremental Delivery
1. Phase 2 foundation → Phase 3 (US1) → Validate → demo
2. Phase 4 (US2) → Validate in isolation
3. Phase 5 (US3) → Validate failure paths
4. Phase 6 polish → regenerate docs, rerun tests

### Parallel Team Strategy
- Developer A: Lead Phase 2 tasks T001–T004
- Developer A continues with US1 tasks T005–T018 while Developer B handles US2 tasks T008–T020 once US1 groundwork lands
- Developer C focuses on US3 tasks T011–T021 after US1/US2 dependencies are ready
- Team coordinates final polish and feedback tasks T014–T016 and T022 together
