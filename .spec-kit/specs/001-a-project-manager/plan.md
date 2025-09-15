# Implementation Plan: TUI environments column

**Branch**: `001-a-project-manager` | **Date**: 2025-09-14 | **Spec**: /home/jlrosende/personal/project-manager/.spec-kit/specs/001-a-project-manager/spec.md
**Input**: Feature specification from `.spec-kit/specs/001-a-project-manager/spec.md`

## Execution Flow (/plan command scope)

```
1. Load feature spec from Input path
   → Found: spec.md
2. Fill Technical Context (scan for NEEDS CLARIFICATION)
   → No new unknowns for this TUI-only change beyond existing global clarifications
3. Evaluate Constitution Check section below
   → No violations for a small UI enhancement
   → Update Progress Tracking: Initial Constitution Check
4. Execute Phase 0 → research.md
   → All TUI unknowns resolved
5. Execute Phase 1 → contracts, data-model.md, quickstart.md
6. Re-evaluate Constitution Check section
   → No new violations
   → Update Progress Tracking: Post-Design Constitution Check
7. Plan Phase 2 → Describe task generation approach (DO NOT create tasks.md)
8. STOP - Ready for /tasks command
```

## Summary

Add an environments column/pane to the TUI that displays the available environments for the currently selected project. When the project selection changes, the environments list updates. Styling uses lipgloss, and state management follows the Elm architecture using bubbletea. No secrets are displayed; only environment names.

## Technical Context

**Language/Version**: Go 1.25.1
**Primary Dependencies**: github.com/charmbracelet/bubbletea v1.3.9, github.com/charmbracelet/bubbles v0.21.0, github.com/charmbracelet/lipgloss v1.1.0
**Storage**: N/A
**Testing**: go test (NEEDS CLARIFICATION: test layout and commands)
**Target Platform**: Linux/macOS terminal
**Project Type**: single
**Performance Goals**: Instant UI update on selection (<50ms perceived)
**Constraints**: Do not display or log secrets; keep UI responsive in small terminals
**Scale/Scope**: Dozens of projects, a handful of environments each

Implementation placement and constraints from user input:

-   TUI handlers in `internal/adapters/handlers`
-   Custom UI components in `pkg/ui`
-   Elm architecture with bubbletea; styles via lipgloss

## Constitution Check

**Simplicity**:

-   Projects: 1 (cli)
-   Using framework directly? Yes (bubbletea/bubbles)
-   Single data model? Yes (UI model augmented with envs list)
-   Avoiding patterns? Yes (no extra abstraction)

**Architecture**:

-   Feature shipped within existing app; no new library required
-   Libraries listed: bubbletea/bubbles (UI), lipgloss (styles)
-   CLI per library: N/A
-   Library docs: N/A

**Testing (NON-NEGOTIABLE)**:

-   RED-GREEN planned for UI message handling and rendering helpers
-   Order: contract (messages) → integration (model update/render) → unit (helpers)
-   Real deps: yes (bubbletea)

**Observability**:

-   Use existing structured logging if present; avoid logging secrets

**Versioning**:

-   Folded into current feature branch; no API breakage

## Project Structure

### Documentation (this feature)

```
specs/001-a-project-manager/
├── plan.md              # This file (/plan command output)
├── research.md          # Phase 0 output (/plan command)
├── data-model.md        # Phase 1 output (/plan command)
├── quickstart.md        # Phase 1 output (/plan command)
└── contracts/           # Phase 1 output (/plan command)
```

### Source Code (repository root)

```
# Single project
internal/adapters/handlers        # TUI handlers (update/view)
pkg/ui                            # UI components (lists, columns, styles)
```

**Structure Decision**: Option 1 (single project)

## Phase 0: Outline & Research

1. Unknowns identified and resolved
    - Layout approach: split view with projects list (left) and environments list (right)
    - Component choice: bubbles list or table; choose list for environments for simplicity
    - State updates: on project selection change, derive envs for selected project
    - Sizing behavior: responsive split using terminal width; collapse envs when width too small
2. Best practices
    - Use tea.Msg for selection changes; keep rendering stateless from model
    - Avoid heavy computation in View; precompute styles
3. Output: research.md with decisions and rationale

**Output**: research.md with all TUI unknowns resolved

## Phase 1: Design & Contracts

1. Entities → `data-model.md`
    - Project(name, environments[]string)
    - UI Model: selectedProjectIdx int, projects []Project, envsForSelected []string, width/height, styles
    - Messages: ProjectsLoadedMsg, ProjectSelectedMsg, WindowSizeMsg
2. Contracts → `/contracts/`
    - Message/event contracts for selection and rendering
3. Contract tests (planned)
    - Ensure envs list updates when selection changes
    - Ensure collapse behavior under narrow width
4. Test scenarios
    - Validate user can see environments of selected project; switching selection updates list
5. Agent file: N/A

**Output**: data-model.md, /contracts/\*, quickstart.md

## Phase 2: Task Planning Approach

**Task Generation Strategy**:

-   Derive tasks from data-model and contracts: add model fields, implement Update cases, implement View split, add components in pkg/ui

**Ordering Strategy**:

-   TDD: write contract tests for message handling before implementation
-   Dependency: model before view; styles last
-   Parallel: style constants and component skeletons

**Estimated Output**: Will be produced by /tasks

## Complexity Tracking

| Violation | Why Needed | Simpler Alternative Rejected Because |
| --------- | ---------- | ------------------------------------ |
| —         | —          | —                                    |

## Progress Tracking

**Phase Status**:

-   [x] Phase 0: Research complete (/plan command)
-   [x] Phase 1: Design complete (/plan command)
-   [x] Phase 2: Task planning complete (/plan command - describe approach only)
-   [ ] Phase 3: Tasks generated (/tasks command)
-   [ ] Phase 4: Implementation complete
-   [ ] Phase 5: Validation passed

**Gate Status**:

-   [x] Initial Constitution Check: PASS
-   [x] Post-Design Constitution Check: PASS
-   [x] All NEEDS CLARIFICATION resolved (for this TUI scope)
-   [ ] Complexity deviations documented (N/A)

---

_Based on Constitution v2.1.1 - See `/memory/constitution.md`_
