# Implementation Plan: CLI `pm delete`

**Branch**: `002-delete-command-overview` | **Date**: 2025-10-05 | **Spec**: `/workspaces/project-manager/specs/002-delete-command-overview/spec.md`
**Input**: Feature specification from `/workspaces/project-manager/specs/002-delete-command-overview/spec.md`

## Execution Flow (/plan command scope)
```
1. Load feature spec from Input path
   → If not found: ERROR "No feature spec at {path}"
2. Fill Technical Context (scan for NEEDS CLARIFICATION)
   → Detect Project Type from file system structure or context (web=frontend+backend, mobile=app+api)
   → Set Structure Decision based on project type
3. Fill the Constitution Check section based on the content of the constitution document.
4. Evaluate Constitution Check section below
   → If violations exist: Document in Complexity Tracking
   → If no justification possible: ERROR "Simplify approach first"
   → Update Progress Tracking: Initial Constitution Check
5. Execute Phase 0 → research.md
   → If NEEDS CLARIFICATION remain: ERROR "Resolve unknowns"
6. Execute Phase 1 → contracts, data-model.md, quickstart.md, agent-specific template file
7. Re-evaluate Constitution Check section
   → If new violations: Refactor design, return to Phase 1
   → Update Progress Tracking: Post-Design Constitution Check
8. Plan Phase 2 → Describe task generation approach (DO NOT create tasks.md)
9. STOP - Ready for /tasks command
```

## Summary
Implement the `pm delete` CLI workflow that resolves a target project by name or path, validates mutually exclusive scope flags, and drives a new `DeleteProject` service operation that removes registry metadata, environment data, filesystem assets, git hooks, and optional backups while respecting dry-run, force, and confirmation safeguards.

## Technical Context
**Language/Version**: Go 1.25
**Primary Dependencies**: `spf13/cobra` for CLI parsing, `bubbletea` UI utilities, `go-git` for git hooks, project-local filesystem/env services
**Storage**: Local filesystem directories, project registry persistence, env var store
**Testing**: `go test`, `make unit`, `make integration`
**Target Platform**: Cross-platform terminal environments (Linux/macOS/Windows)
**Project Type**: Single CLI application with hexagonal architecture
**Performance Goals**: CLI operations should complete within a few seconds and stream progress logs
**Constraints**: Preserve user data unless `--all` is selected, respect project locks/read-only flags, abort on backup failure
**Scale/Scope**: Single-project operations; deletion affects one resolved project per invocation

## Constitution Check
- **Hexagonal Integrity**: Domain logic stays within `internal/core`; adapters only orchestrate CLI parsing, repositories, and logging. New delete flow must route through ports instead of reaching adapters directly.
- **Port-Driven Delivery**: Extend filesystem, env vars, git, and project ports with deletion/backup methods before implementing repositories. No adapter should bypass these interfaces.
- **Reproducible Workflows**: Add unit and integration tests covering delete paths; ensure `make lint`, `make unit`, and `make integration` remain green.
- **Configuration Contracts**: Document CLI usage and safety semantics across docs, man page, and README; backups live under `~/.pm/backups` with consistent naming.

## Project Structure
```
specs/002-delete-command-overview/
├── plan.md
├── research.md          # Phase 0 output
├── data-model.md        # Phase 1 output
├── quickstart.md        # Phase 1 output
├── contracts/           # Phase 1 output
└── tasks.md             # Phase 2 output

internal/
├── adapters/
│   ├── handlers/cli/
│   │   ├── root.go
│   │   ├── delete/          # new command wiring
│   │   ├── new/
│   │   ├── edit/
│   │   └── list/
│   └── repositories/
│       ├── filesystem.go
│       ├── env_vars_repository.go
│       ├── git_repository.go
│       └── project_repository.go
├── core/
│   ├── domain/
│   ├── ports/
│   │   ├── filesystem_port.go
│   │   ├── env_vars_port.go
│   │   ├── git_port.go
│   │   └── project_port.go
│   └── services/
│       ├── project_service.go   # add DeleteProject flow
│       └── project_options.go   # extend with delete options
└── bootstrap/
    └── bootstrap.go

tests/
├── integration/
│   ├── cli_new_*.go
│   └── edit_project_rename_refresh_test.go
├── unit/
│   ├── project_service_*_test.go
│   └── project_repository_*_test.go
└── contract/
    └── cli_new_flags_test.go

docs/
├── cli/
│   ├── pm.md
│   ├── pm_new.md
│   ├── pm_edit.md
│   └── pm_list.md
└── rest/
    └── ... (RST mirrors of CLI docs)

man/
└── pm.1, pm-new.1, ... (add pm-delete.1)
```

**Structure Decision**: Single Go CLI with shared domain/service core; new delete command touches CLI handler, repositories, services, and documentation/test suites listed above.

## Phase 0: Outline & Research
1. Confirm existing project resolution, locking, and confirmation utilities to understand integration points for delete.
2. Investigate current filesystem/env/git repository implementations for reusable patterns (e.g., skeleton removal, env cleanup) and identify gaps for deletion/backup behaviors.
3. Research Go archival approach (standard `archive/zip` vs `compress/gzip` + tar) and determine consistent naming/location under `~/.pm/backups` with atomic rename semantics.
4. Document safety expectations for dry-run and backup flows so service/CLI layers share the same contract.

**Output**: `/workspaces/project-manager/specs/002-delete-command-overview/research.md` capturing decisions, rationale, and alternatives.

## Phase 1: Design & Contracts
1. Translate project entities and deletion scopes into `/specs/002-delete-command-overview/data-model.md`, including new `ProjectDeleteOptions`, scope enum, and backup metadata.
2. Define CLI contract behavior (flags, prompts, logging expectations) under `/specs/002-delete-command-overview/contracts/cli.md`, mirroring documentation requirements.
3. Outline integration and unit test scenarios in `/specs/002-delete-command-overview/quickstart.md`, aligning with acceptance tests from the spec.
4. Extend repository and port contracts conceptually within the documentation to ensure new methods align with hexagonal boundaries before coding.
5. Document explicit exit-code outcomes (success, dry-run, failure) and missing-target handling so tests can assert CLI return values and error messaging.
6. Capture dry-run + backup behavior expectations, including “no archive written” guarantees, to drive matching tests and implementation safeguards.
7. Update agent context by running `.specify/scripts/bash/update-agent-context.sh opencode` after documenting new technologies or decisions.

**Output**: Refresh plan.md, produce data-model.md, contracts, quickstart, and update agent file to reflect new context.

### Validation Hooks
- Add a service validation checklist ensuring logging occurs for every deletion artifact (FR-010) and is verifiable via automated tests.
- Record planned test coverage for non-zero exit codes on error paths and zero exit code for successful dry runs.
- Confirm dry-run + backup flow explicitly avoids filesystem writes while still reporting archive destinations.

## Phase 2: Task Planning Approach
- Generate `tasks.md` summarizing TDD workflow: add deletion tests (unit/integration/contract) before implementation, then CLI/service/adapter changes, then docs and backups.
- Organize tasks by scope: option parsing, service orchestration, port extensions, adapter implementations, backups, dry-run, confirmations, logging, docs, exit-code handling, and tests.
- Flag parallelizable documentation/test tasks; ensure destructive flows are completed sequentially after tests verify coverage.

**Estimated Output**: `tasks.md` with ~25 ordered tasks spanning tests, implementation, backups, and documentation updates.

## Complexity Tracking
| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| *(none)* | | |

## Progress Tracking
**Phase Status**:
- [x] Phase 0: Research complete (/plan command)
- [x] Phase 1: Design complete (/plan command)
- [x] Phase 2: Task planning complete (/plan command - describe approach only)
- [ ] Phase 3: Tasks generated (/tasks command)
- [ ] Phase 4: Implementation complete
- [ ] Phase 5: Validation passed

**Gate Status**:
- [x] Initial Constitution Check: PASS
- [x] Post-Design Constitution Check: PASS
- [x] All NEEDS CLARIFICATION resolved
- [x] Complexity deviations documented

---
*Based on Constitution v1.0.0 - See `/workspaces/project-manager/.specify/memory/constitution.md`*
