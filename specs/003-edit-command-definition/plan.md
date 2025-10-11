# Implementation Plan: pm edit flag-driven updates

**Branch**: `003-edit-command-definition` | **Date**: 2025-10-11 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/003-edit-command-definition/spec.md`

## Summary

Implement flag-driven project and environment editing for the `pm edit` command, adding CLI input file support and skeleton generation so configuration updates can be scripted, validated, and persisted without launching an editor.

## Technical Context

**Language/Version**: Go 1.25 (module `github.com/jlrosende/project-manager`)  
**Primary Dependencies**: Cobra CLI framework, internal bootstrap/services, filesystem repositories  
**Storage**: Local filesystem-backed configuration stores  
**Testing**: `go test`, `make unit`, `make integration`  
**Target Platform**: Cross-platform CLI (Linux/macOS)  
**Project Type**: Single CLI application with hexagonal architecture  
**Performance Goals**: p95 ≤750 ms (≤200 lines), ≤1.4 s (≤800 lines), ≤2.2 s (>800 lines); max ≤3 s; support ≥20 edits/min sequential  
**Constraints**: Enforce project/environment immutability rules, per-project file locking, atomic writes, explicit clears for metadata, validation before persistence  
**Scale/Scope**: Configuration management for local projects (small-to-medium project lists)

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

- Constitution placeholders contain no enforceable principles; treating as no blocking gates.  
- Proceed with standard repository practices (tests required for new behavior, documentation auto-generation).  
- Post-Phase 1 review: no new constraints introduced; plan remains in compliance.

## Project Structure

### Documentation (this feature)

```
specs/003-edit-command-definition/
├── plan.md              # This file (/speckit.plan output)
├── research.md          # Phase 0 output
├── data-model.md        # Phase 1 output
├── quickstart.md        # Phase 1 output
├── contracts/           # Phase 1 output
└── tasks.md             # Phase 2 output (/speckit.tasks)
```

### Source Code (repository root)

```
cmd/
└── main.go

internal/
├── adapters/
│   ├── handlers/
│   │   └── cli/
│   │       ├── edit/
│   │       ├── new/
│   │       └── list/
│   ├── repositories/
│   │   └── project_repository.go
│   └── ...
├── core/
│   ├── domain/
│   ├── ports/
│   └── services/
└── bootstrap/
    └── bootstrap.go

tests/
├── contract/
├── integration/
└── unit/
```

**Structure Decision**: Single Go CLI application leveraging hexagonal layers (`internal/core` services + `internal/adapters` handlers/repositories); existing directories reused.

## Complexity Tracking

*N/A (no constitution violations identified)*

## Deferred Work

- Backlog item: "pm edit performance timing" to implement response-time instrumentation that validates SC-004 after initial delivery.
