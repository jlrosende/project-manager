# Feature Specification: CLI `delete` Command

**Feature Branch**: `002-delete-command-overview`  
**Created**: 2025-10-05  
**Status**: Draft  
**Input**: User description: "Delete Command Overview • Invocation: Run pm delete with either a project name or an explicit filesystem path. The command first resolves the project context: names map through the project registry, while paths must correspond to a registered project root; otherwise the command aborts with a “project not found at path” error. • Confirmation Flow: By default the user is prompted to confirm the deletion to prevent accidental loss. Passing --force skips this prompt and proceeds immediately. If the project is already undergoing another operation, or the workspace is read-only, the command stops before making changes. • Functional Scopes: • --all: Purges everything linked to the project, including registry entry, configuration files, generated skeleton, caches, and git hooks if applicable. • --keep-files: Removes only project metadata (registry entry, config, env vars) while leaving the working directory and user-authored files intact. • --only-env: Drops stored environment variables for the project but preserves the project definition and files. • Scope flags combine with path-based targeting; the chosen scope determines exactly which artifacts are removed once the project is resolved. • Dry Run Mode: Supplying --dry-run runs all resolution and validation logic, then prints the list of components that would be deleted under the chosen scope. No changes are committed during a dry run. • Backup Handling: When --backup is present, the command packages the project’s configuration and environment data into a timestamped archive (default location ~/.pm/backups/). Deletion only proceeds after the archive is safely written; on backup failure the entire operation is aborted to avoid data loss. • Execution Feedback: As deletion progresses, each component (config file, env vars, skeleton directory, etc.) is logged. On success the command summarizes what was removed and where the optional backup lives. Any step that cannot be completed cleanly emits a descriptive error so the user knows what remains."

## Execution Flow (main)
```
1. Resolve project target from provided name or path
   → Name: look up project registry; Path: ensure it matches a registered project root or abort with "project not found at path"
2. Verify workspace readiness and project availability
   → Abort with messaging if project is locked by another operation or the workspace is read-only
3. Determine requested deletion scope
   → Interpret `--all`, `--keep-files`, or `--only-env`; reject conflicting scope flags and default to safe metadata-only removal when unspecified
4. If `--backup` supplied, create timestamped archive of project configuration and environment data under `~/.pm/backups/`
   → Abort entire command if archive creation fails
5. If `--dry-run` supplied, produce detailed list of components slated for removal based on scope and exit without changes
6. Request user confirmation unless `--force` is provided
   → On decline, abort with no changes
7. Execute scoped deletions with step-by-step logging; halt and report if any artifact cannot be removed cleanly
8. Summarize completed deletions and reference backup location when applicable
```

---

## ⚡ Quick Guidelines
- Prioritize safeguards that prevent accidental loss while keeping cleanup fast for confident users
- Communicate decisions clearly: every abort or deletion step must have actionable messaging
- Preserve user-authored assets unless the user explicitly opts in to full removal

## User Scenarios & Testing *(mandatory)*

### Primary User Story
As a developer maintaining multiple projects, I want a `pm delete` command that safely removes a project's records and optional assets, so I can decommission workspaces without leaving stale configuration or risking unintended data loss.

### Acceptance Scenarios
1. **Given** a registered project referenced by name, **When** I run `pm delete my-app`, **Then** I am prompted to confirm and, after agreeing, the registry entry, config, and env data are removed with a success summary.
2. **Given** a registered project located at `/work/repos/my-app`, **When** I run `pm delete /work/repos/my-app`, **Then** the command resolves the path, confirms deletion, and removes the same artifacts as a name-based invocation.
3. **Given** a registered project, **When** I run `pm delete my-app --force --all`, **Then** the command skips confirmation and removes metadata, skeleton files, caches, git hooks, and user-authored files while logging each step.
4. **Given** a registered project with a populated working directory, **When** I run `pm delete my-app --keep-files`, **Then** only registry, config, and env metadata are removed while the working directory remains untouched.
5. **Given** a registered project with stored environment variables, **When** I run `pm delete my-app --only-env`, **Then** the environment data is removed and other project assets remain available.
6. **Given** a registered project, **When** I run `pm delete my-app --dry-run --all`, **Then** the command lists all components that would be removed and exits without modifying anything.
7. **Given** a registered project, **When** I run `pm delete my-app --backup --all`, **Then** the command saves a timestamped archive of the entire project directory (including user-authored files) before removing artifacts and reports the archive path.
8. **Given** no registered project matches the supplied name or path, **When** I run `pm delete ghost-app`, **Then** the command aborts before any deletion, returns a non-zero exit code, and surfaces a "project not found" error message.
9. **Given** a registered project, **When** I run `pm delete my-app --dry-run --backup`, **Then** the command reports the intended archive location without creating an archive file and exits without changes while returning success.
 
 ### Edge Cases
 - Unregistered name or path triggers a "project not found" error before any destructive action and MUST return a non-zero exit code.
 - Workspace read-only or project currently locked causes the command to abort with instructions for resolving the conflict.
 - Scope flags supplied together (e.g., `--all` with `--keep-files`) are rejected with guidance to choose one.
 - Backup creation failure halts the operation and leaves the project untouched.
 - Dry run combined with backup reports intended archive path but skips file creation to avoid side effects and MUST avoid touching the filesystem.
 - File system permission issues or partial deletions surface descriptive errors indicating what remains and why.
 - `--backup --all` ensures user-authored files are captured in the archive so users can restore their workspace after full deletion.


## Clarifications

### Session 2025-10-05
- Q: When `pm delete TARGET --all --backup` runs, what should go into the backup archive before deletion? → A: Entire project directory, including user-authored files

## Requirements *(mandatory)*

### Functional Requirements
- **FR-001**: CLI MUST resolve the target project by name through the registry or by validating that a provided path corresponds to a registered project root.
- **FR-002**: CLI MUST refuse to proceed when the project is already under another operation or the workspace is marked read-only, providing actionable messaging.
- **FR-003**: CLI MUST prompt for confirmation before deleting anything, unless `--force` is provided.
- **FR-004**: CLI MUST support `--all` to delete registry entries, configuration files, environment data, generated skeletons, caches, related git hooks, and user-authored files.
- **FR-005**: CLI MUST support `--keep-files` to remove project metadata (registry, config, env data) while leaving the working directory untouched.
- **FR-006**: CLI MUST support `--only-env` to delete stored environment variables while preserving project definitions and files.
- **FR-007**: CLI MUST validate that scope flags are mutually exclusive and return a descriptive error when conflicting options are combined.
- **FR-008**: CLI MUST support `--dry-run` that executes all validations, computes the scope, and outputs the items that would be deleted without performing any changes.
- **FR-009**: CLI MUST support `--backup`, creating a timestamped archive of the entire project directory contents (including user-authored files) prior to scoped deletions, and abort entirely on backup failure.
- **FR-010**: CLI MUST log each deletion action and conclude with a summary describing what was removed and, when applicable, where the backup resides.
- **FR-011**: CLI MUST surface descriptive errors when a requested artifact cannot be removed and leave the remaining project assets untouched.
- **FR-012**: CLI MUST return a non-zero exit code on failures or aborted operations and zero on successful completion.

### Key Entities *(include if feature involves data)*
- **Project Registry Entry**: Represents the registered project, including name, root path, metadata references, and status flags indicating active operations.
- **Project Configuration Assets**: User-facing configuration file(s) and stored environment variables tied to the project; subject to backup and deletion scope.
- **Project Backup Archive**: Timestamped package containing the entire project directory contents captured when `--backup` is requested, ensuring configuration, environment data, and user-authored files can be restored if needed.

---

## Review & Acceptance Checklist
*GATE: Automated checks run during main() execution*

### Content Quality
- [x] No implementation details (languages, frameworks, APIs)
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
*Updated by main() during processing*

- [x] User description parsed
- [x] Key concepts extracted
- [x] Ambiguities marked
- [x] User scenarios defined
- [x] Requirements generated
- [x] Entities identified
- [ ] Review checklist passed

---
