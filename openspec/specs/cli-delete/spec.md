# cli-delete Specification

## Purpose
Define the required behavior for the `pm delete` command so users can safely remove project assets with clear confirmations, backups, and scoped operations.

## Requirements
### Requirement: Target Resolution and Safeguards
The CLI SHALL resolve deletion targets by project name or registered path and refuse operations when safety checks fail.

#### Scenario: Resolve project by name
- **WHEN** a user runs `pm delete <name>` for a registered project
- **THEN** the command resolves the registry entry, prompts for confirmation, and proceeds only after acceptance or `--force`.

#### Scenario: Resolve project by path
- **WHEN** a user runs `pm delete /path/to/project`
- **THEN** the command verifies the path matches a registered project root and continues as with a name-based invocation.

#### Scenario: Abort when project unavailable
- **WHEN** the supplied name or path does not match a registered project or the workspace is read-only
- **THEN** the command aborts before deletion and returns a non-zero exit code with guidance.

### Requirement: Scoped Deletion Modes
The CLI SHALL support mutually exclusive scopes controlling which assets are removed.

#### Scenario: Delete everything with --all
- **WHEN** a user runs `pm delete <name> --all`
- **THEN** the command removes registry data, configs, environment files, generated assets, caches, git hooks, and user-authored files while logging each removal.

#### Scenario: Keep working files
- **WHEN** a user runs `pm delete <name> --keep-files`
- **THEN** the command removes only registry, config, and environment data and leaves the working directory untouched.

#### Scenario: Delete only environment data
- **WHEN** a user runs `pm delete <name> --only-env`
- **THEN** the command deletes stored environment data while preserving project definitions and files.

#### Scenario: Reject conflicting scope flags
- **WHEN** a user supplies more than one scope flag in the same invocation
- **THEN** the command fails with an error explaining the conflict and performs no deletions.

### Requirement: Confirmation, Dry Run, and Feedback
The CLI SHALL prevent accidental destruction through confirmations, dry-run previews, and explicit reporting.

#### Scenario: Confirmation prompt by default
- **WHEN** a user runs `pm delete <name>` without `--force`
- **THEN** the command prompts for confirmation and proceeds only if the user agrees.

#### Scenario: Dry run preview
- **WHEN** a user runs `pm delete <name> --dry-run [--all|--keep-files|--only-env]`
- **THEN** the command outputs the artifacts that would be removed for the chosen scope without making changes.

#### Scenario: Operation summary
- **WHEN** the command completes a deletion
- **THEN** it summarizes which artifacts were removed or skipped and reports any errors encountered.

### Requirement: Backup Handling
The CLI SHALL provide optional backups that gate destructive actions.

#### Scenario: Successful backup before deletion
- **WHEN** a user runs `pm delete <name> --backup`
- **THEN** the command creates a timestamped archive of the project configuration (and user files when combined with `--all`) before deletion and reports the backup path.

#### Scenario: Abort on backup failure
- **WHEN** backup creation fails for any reason
- **THEN** the command aborts the deletion, reports the failure, and leaves all project assets untouched.

### Requirement: Error Handling and Exit Codes
The CLI SHALL surface descriptive errors and exit codes.

#### Scenario: Partial deletion error reporting
- **WHEN** a deletion step fails due to permissions or other issues
- **THEN** the command reports the failing artifact, stops further removal, and exits with a non-zero code.

#### Scenario: Successful operation exit
- **WHEN** the command completes without errors or user cancellation
- **THEN** it exits with code 0.
