# Feature Specification: CLI `new` Project Command

**Feature Branch**: `001-cli-new-project`
**Created**: 2025-09-22
**Status**: Draft
**Input**: User description: "CLI new project command, i want to create a terminal command 'new' to create a new project. This command recive flags to configure all the required parameters of a project, and optinal parameters tu configure the other project values. The inputs are validated. Also the command mus have a flag to sent a json or yaml file to configure all the inputs. If a flag path is set the project is created in the path, but if a flag --here is set, the project is initialized in te currect directory."

## Execution Flow (main)
```
1. Parse user description from Input
   � If empty: ERROR "No feature description provided"
2. Extract key concepts from description
   � Identify: actors, actions, data, constraints
3. For each unclear aspect:
   � Mark with [NEEDS CLARIFICATION: specific question]
4. Fill User Scenarios & Testing section
   � If no clear user flow: ERROR "Cannot determine user scenarios"
5. Generate Functional Requirements
   � Each requirement must be testable
   � Mark ambiguous requirements
6. Identify Key Entities (if data involved)
7. Run Review Checklist
   � If any [NEEDS CLARIFICATION]: WARN "Spec has uncertainties"
   � If implementation details found: ERROR "Remove tech details"
8. Return: SUCCESS (spec ready for planning)
```

---

## � Quick Guidelines
-  Focus on WHAT users need and WHY
- L Avoid HOW to implement (no tech stack, APIs, code structure)
- =e Written for business stakeholders, not developers

### Section Requirements
- **Mandatory sections**: Must be completed for every feature
- **Optional sections**: Include only when relevant to the feature
- When a section doesn't apply, remove it entirely (don't leave as "N/A")

### For AI Generation
When creating this spec from a user prompt:
1. **Mark all ambiguities**: Use [NEEDS CLARIFICATION: specific question] for any assumption you'd need to make
2. **Don't guess**: If the prompt doesn't specify something (e.g., "login system" without auth method), mark it
3. **Think like a tester**: Every vague requirement should fail the "testable and unambiguous" checklist item
4. **Common underspecified areas**:
   - User types and permissions
   - Data retention/deletion policies
   - Performance targets and scale
   - Error handling behaviors
   - Integration requirements
   - Security/compliance needs

---

## User Scenarios & Testing (mandatory)

### Primary User Story
As a developer managing multiple projects, I want a `pm new` command that creates or initializes a project folder with a `.project.hcl` and `.env`, so I can quickly set up consistent project scaffolding with required and optional properties via flags or a config file.

### Acceptance Scenarios
1. Given no existing folder and a target path argument, When I run `pm new NAME /abs/dir`, Then a new folder is created with `.project.hcl` and `.env` populated from arguments and validation passes.
2. Given an existing folder and `--here`, When I run `pm new NAME --here`, Then the current directory is initialized with `.project.hcl` and `.env` without overwriting unrelated files.
3. Given a CLI input file `project.yaml`, When I run `pm new NAME --cli-input project.yaml`, Then required and optional inputs are read from the file and validated, and the project is created/initialized accordingly.
4. Given both arguments and a CLI input file, When I run `pm new override /abs/dir --cli-input project.json`, Then positional arguments override values from the file with validation applied to the final merged inputs.
5. Given a request to scaffold a CLI input file, When I run `pm new --generate-cli-skeleton-json project.json`, Then a JSON skeleton with the expected fields is produced (or printed to stdout when no path is supplied).
6. Given invalid inputs (e.g., missing required name, invalid path, conflicting `--here` with a path argument), When I run the command, Then the command fails with actionable error messages and no partial project is left behind.
7. Given request to create per-project git config, When I run the command, Then a per-project git config entry is prepared according to project policy without exposing secrets.

### Edge Cases
- Per-project git config uses includeIf include.path in $HOME/.gitconfig referencing project path; ensure no secrets are written.
- `--here` used in a non-empty directory containing an existing `.project.hcl`: default FAIL unless `--force`; with `--force`, overwrite fields per schema without exposing secrets.
- Path arguments may be relative: must resolve to absolute path for creation; conflicting or inaccessible parent dirs must be handled.
- CLI input files contain unknown fields: default FAIL; allow via `--allow-unknown`.
- Skeleton generation flags may omit a destination path to emit the template to stdout.
- Simultaneous use of a path argument and `--here`: must be rejected with a clear error.
- YAML vs JSON: both supported; malformed files produce clear validation errors.
- Name collisions: check filesystem only; no global registry involved.
- When `--force` is used, overwrite entire `.env` with new values; backup is not created.

## Clarifications

### Session 2025-09-23
- Q: How should `.env` be handled when `--force` is used? → A: Overwrite entire `.env` with new values
- Q: How should per-project git config be set up? → A: Add includeIf in $HOME/.gitconfig for project path
- Q: What should `--dry-run` output by default? → A: Selectable via --output=json|text (default text)

## Requirements (mandatory)

### Functional Requirements
- FR-001: Users MUST be able to create a new project via `pm new` using flags only.
- FR-002: Users MUST be able to initialize the current directory via `pm new --here`.
- FR-003: Users MUST be able to specify a destination via the optional path argument to create the project there.
- FR-004: Users MUST be able to provide a CLI input file via `--cli-input` supporting JSON and YAML.
- FR-005: System MUST validate required inputs (e.g., project name) and optional inputs according to schema.
- FR-006: System MUST merge config file values with CLI flags; flags override file values.
- FR-007: System MUST create `.project.hcl` and `.env` with provided values and defaults; when `--force` is used, `.env` is overwritten entirely.
- FR-008: System MUST handle existing directories: create if missing, initialize if present without destructive changes.
- FR-009: System MUST provide clear error messages and avoid leaving partial state on failure.
- FR-010: System MUST support a `--dry-run` flag to preview actions without modifying the filesystem; output selectable via `--output=json|text` with default `text`.
- FR-011: System MUST support per-project git configuration setup consistent with project governance; use includeIf in $HOME/.gitconfig for project path.
- FR-012: System MUST log actions with structured, user-friendly messages without exposing secrets.
- FR-013: Users MUST be able to generate CLI input skeleton files via `--generate-cli-skeleton-json` and `--generate-cli-skeleton-yaml`.

- FR-013: System MUST reject conflicting options (e.g., providing a path argument together with `--here`).
- FR-014: System MUST support defaults for unspecified optional values as documented.
- FR-015: System MUST ensure created files respect .gitignore and do not commit secrets.

- FR-016: System MUST return non-zero exit code on validation or creation errors.

- FR-017: System MUST provide help/usage output describing flags and examples.

- FR-018: System MUST be idempotent: re-run initializes missing files; no overwrites occur unless `--force` is provided.

### Non-Functional Requirements
- NFR-001: Startup under 200ms; file operations non-blocking.
- NFR-002: No secret logging; sensitive values only in `.env`; respect `.gitignore`.
- NFR-003: CI quality gates: lint, unit, integration must pass; reproducible builds required.
- NFR-004: Structured, user-friendly logging and deterministic behavior.

### Key Entities (include if feature involves data)
- Project Definition: name (required), path or here (one required), environments (optional), metadata (optional). Persisted in `.project.hcl`; sensitive values stored in `.env` only.
- Config Input: a JSON or YAML document containing required and optional project properties; validated against a schema.

---

## Review & Acceptance Checklist

### Content Quality
- [ ] No implementation details (languages, frameworks, APIs)
- [x] Focused on user value and business needs
- [x] Written for non-technical stakeholders
- [x] All mandatory sections completed

### Requirement Completeness
- [x] No [NEEDS CLARIFICATION] markers remain
- [x] Requirements are testable and unambiguous
- [x] Success criteria are measurable
- [x] Scope is clearly bounded
- [x] Dependencies and assumptions identified

---

## Execution Status

- [x] User description parsed
- [x] Key concepts extracted
- [x] Ambiguities marked
- [x] User scenarios defined
- [x] Requirements generated
- [x] Entities identified
- [ ] Review checklist passed

---
