# Implementation Plan: `pm new` Environment Enhancements

**Branch**: `004-new-env-command` | **Date**: 2025-10-13 | **Spec**: [/specs/004-new-env-command/spec.md](/specs/004-new-env-command/spec.md)
**Input**: Feature specification from `/specs/004-new-env-command/spec.md`

**Note**: This template is filled in by the `/speckit.plan` command. See `.specify/templates/commands/plan.md` for the execution workflow.

## Summary

Implement dual-mode `pm new` behavior that distinguishes project creation from environment addition, introduce environment-customization flags and config schema, ensure dry-run parity, refine uniqueness/registry handling, and update skeletons/tests to reflect the new workflow. The approach extends CLI parsing and shared merge helpers, adds an `EnvironmentInput` data model, adjusts services/repositories for clarified rules, and refreshes docs/tests to keep users aligned with the new experience.

## Technical Context

**Language/Version**: Go 1.25  
**Primary Dependencies**: Cobra CLI (`spf13/cobra`), internal hexagonal service/repository abstractions, stdlib (`os`, `filepath`, `encoding/json`, `gopkg.in/yaml.v3`)  
**Storage**: Filesystem-backed project metadata (`.project.hcl`, includeIf git config)  
**Testing**: `go test` suites (unit, integration, contract), make targets (`make unit`, `make integration`)  
**Target Platform**: Cross-platform CLI (Linux/macOS/WSL) running in terminal environments  
**Project Type**: Monolithic Go CLI with hexagonal architecture  
**Performance Goals**: No new performance targets beyond existing CLI responsiveness (<100ms command overhead)  
**Constraints**: Must preserve backward compatibility for project creation, avoid breaking existing scripts, and keep skeleton generation idempotent  
**Scale/Scope**: Single CLI command enhancement affecting relevant services/repos and associated tests/documentation

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

The constitution file contains placeholder sections with no enforceable principles or gates. There are no explicit governance constraints to evaluate, so the plan proceeds with standard repository guidelines. Post-design review confirms no additional governance content was introduced during this planning pass.

## Project Structure

### Documentation (this feature)

```
specs/004-new-env-command/
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
├── contracts/
└── tasks.md
```

### Source Code (repository root)

```
internal/
├── adapters/
│   ├── handlers/
│   │   └── cli/
│   │       └── new/
│   └── repositories/
│       └── project_repository.go
├── bootstrap/
│   └── bootstrap.go
└── core/
    ├── domain/
    ├── ports/
    └── services/

configs/
└── config.example.hcl

docs/
└── cli/
    └── pm_new.md

tests/
├── contract/
│   └── cli_new_flags_test.go
├── integration/
│   └── cli_new_*.go
└── unit/
    └── project_service_*.go
```

**Structure Decision**: Work occurs within the existing Go CLI layout under `internal/` (handlers, services, repositories, domain) with supporting updates to configs/docs/tests. No new top-level packages are introduced; all changes integrate with the established hexagonal structure.

## Complexity Tracking

*Fill ONLY if Constitution Check has violations that must be justified*

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| _None_ | — | — |
