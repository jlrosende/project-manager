# Tasks: TUI environments column

**Input**: Design documents from `/specs/001-a-project-manager/`
**Prerequisites**: plan.md (required), research.md, data-model.md, contracts/

## Execution Flow (main)
```
Order: Setup → Tests (failing) → Core → Integration → Polish
Apply [P] for independent files; omit [P] when same file or dependent
```

## Phase 3.1: Setup
- [x] T001 Ensure dependencies in go.mod: bubbletea, bubbles, lipgloss in /home/jlrosende/personal/project-manager/go.mod
- [x] T002 Create pkg/ui/list and pkg/ui/styles scaffolding in /home/jlrosende/personal/project-manager/pkg/ui
- [x] T003 [P] Add basic README usage notes in docs/cli (skip if already present)

## Phase 3.2: Tests First (TDD)
- [x] T004 [P] Contract test for ProjectsLoadedMsg handling: derive envs list in tests/integration/tui_projects_loaded_test.go
- [x] T005 [P] Contract test for ProjectSelectedMsg handling: updates envs list in tests/integration/tui_project_selected_test.go
- [x] T006 [P] Integration test for collapse behavior on narrow terminal in tests/integration/tui_collapse_envs_test.go
- [x] T007 [P] Integration test for view rendering envs column when width sufficient in tests/integration/tui_view_envs_test.go

## Phase 3.3: Core Implementation
- [x] T008 Add UI model fields to internal/adapters/handlers/tui/window.go: selected index, envsForSelected, width/height
- [x] T009 Implement Update() cases for ProjectsLoadedMsg and ProjectSelectedMsg in internal/adapters/handlers/tui/window.go
- [x] T010 Implement WindowSizeMsg handling and collapse logic in internal/adapters/handlers/tui/window.go
- [x] T011 Implement environments list component in /home/jlrosende/personal/project-manager/pkg/ui/list/envs_list.go
- [x] T012 [P] Implement styles for panels and titles in /home/jlrosende/personal/project-manager/pkg/ui/styles/styles.go
- [x] T013 Wire environments column into Window.viewPanelRender() in internal/adapters/handlers/tui/window.go
- [x] T014 Replace hardcoded lorem ipsum with project.Description from services in internal/adapters/handlers/tui/window.go

## Phase 3.4: Integration
- [x] T015 Use internal/core/services/ProjectService to load projects, map to UI model in internal/adapters/handlers/tui/window.go
- [x] T016 Ensure no secrets or credential values are rendered or logged (defensive check) across UI rendering functions

## Phase 3.5: Polish
- [x] T017 [P] Unit tests for envs list component in tests/unit/pkg_ui_envs_list_test.go
- [x] T018 [P] Unit tests for styles helpers in tests/unit/pkg_ui_styles_test.go
- [x] T019 Verify performance: render update on selection <50ms (manual measurement notes in docs/rest/perf-notes.md)
- [x] T020 Accessibility/UX: ensure clear focus/selection visuals for lists
- [x] T021 Run Quickstart to manually verify envs column behavior in /home/jlrosende/personal/project-manager/.spec-kit/specs/001-a-project-manager/quickstart.md

## Dependencies
- T004–T007 before T008–T014
- T008 blocks T009–T013
- T011 independent of T009–T010 but required by T013
- T012 independent; can run parallel with T011
- T015 depends on services being available; runs after T008

## Parallel Example
```
# Run these in parallel once Setup done and before Core:
Task: "Contract test ProjectsLoadedMsg" (T004)
Task: "Contract test ProjectSelectedMsg" (T005)
Task: "Integration test collapse behavior" (T006)
Task: "Integration test envs view" (T007)

# Core parallel chunk:
Task: "Implement envs list component" (T011)
Task: "Implement styles helpers" (T012)
```

## Notes
- Keep Update() pure and return new state + Cmd
- Avoid logging sensitive data; render only environment names
- Exact file paths included above for all edits and new files