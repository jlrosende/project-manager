## MODIFIED Requirements
### Requirement: Validation and Persistence
The CLI SHALL validate edits before persistence and ensure atomic updates.

#### Scenario: Reject invalid edits
- **WHEN** a user attempts to save changes that violate validation rules
- **THEN** the command rejects the update, surfaces all validation errors, and restores the previous configuration.

#### Scenario: Confirm successful edit
- **WHEN** a user submits valid changes
- **THEN** the command persists the updates atomically and reports which scope (project or environment) was modified.

#### Scenario: Confirmation highlights scope details
- **WHEN** a user completes an edit
- **THEN** the confirmation output explicitly identifies whether the project or a specific environment changed and summarizes key fields that were updated.

## ADDED Requirements
### Requirement: Dry-Run Preview
The CLI SHALL provide a dry-run mode that validates edits without persisting them while reporting intended changes.

#### Scenario: Project-level dry run
- **WHEN** a user runs `pm edit <project> --dry-run`
- **THEN** the command validates the proposed changes, reports the differences it would apply to the project, and exits without modifying files.

#### Scenario: Environment-level dry run
- **WHEN** a user runs `pm edit <project> <env> --dry-run`
- **THEN** the command reports the environment-specific changes it would perform, leaves all configuration untouched, and returns success when validation passes.

#### Scenario: Dry run failure
- **WHEN** a user runs `pm edit ... --dry-run` and validation fails
- **THEN** the command reports the blocking validation errors and exits with a non-zero status while preserving the original configuration.
