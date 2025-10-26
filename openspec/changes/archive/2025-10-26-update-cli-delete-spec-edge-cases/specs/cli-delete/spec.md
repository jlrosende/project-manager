## MODIFIED Requirements
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

#### Scenario: Force skips confirmation
- **WHEN** a user runs `pm delete <name> --force`
- **THEN** the command bypasses the confirmation prompt, applies the requested scope immediately, and still reports the outcome.

#### Scenario: Dry run with backup preview
- **WHEN** a user runs `pm delete <name> --dry-run --backup`
- **THEN** the command reports the backup destination that would be used without creating any archive or modifying files.

### Requirement: Backup Handling
The CLI SHALL provide optional backups that gate destructive actions.

#### Scenario: Successful backup before deletion
- **WHEN** a user runs `pm delete <name> --backup`
- **THEN** the command creates a timestamped archive of the project configuration (and user files when combined with `--all`) before deletion and reports the backup path.

#### Scenario: Abort on backup failure
- **WHEN** backup creation fails for any reason
- **THEN** the command aborts the deletion, reports the failure, and leaves all project assets untouched.

#### Scenario: Backup destination validation
- **WHEN** a user supplies `--backup-destination` without `--backup`
- **THEN** the command returns an error stating that the destination flag requires `--backup` and performs no changes.
