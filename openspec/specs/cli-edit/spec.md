# cli-edit Specification

## Purpose
Define the behavior for the `pm edit` command so users can safely update project-level or environment-specific configuration with validation and clear feedback.
## Requirements
### Requirement: Argument Handling and Scope Selection
The CLI SHALL require a project argument and determine edit scope based on the presence of an environment argument.

#### Scenario: Project-level editing
- **WHEN** a user runs `pm edit <project>`
- **THEN** the command loads the project configuration for editing, applies confirmed changes to shared metadata, and leaves environment overrides untouched.

#### Scenario: Environment-level editing
- **WHEN** a user runs `pm edit <project> <env>`
- **THEN** the command isolates edits to the specified environment and preserves other environments.

#### Scenario: Unknown project or environment
- **WHEN** the provided project or environment does not exist
- **THEN** the command reports the missing target, performs no updates, and exits with a non-zero status.

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

### Requirement: User-Controlled Outcomes
The CLI SHALL allow users to cancel edits and avoid partial writes.

#### Scenario: Cancel edit flow
- **WHEN** a user exits the edit operation without confirming changes
- **THEN** the command reports that no updates were applied and leaves configuration untouched.

#### Scenario: Handle concurrent edits
- **WHEN** the target configuration is locked or another edit is in progress
- **THEN** the command fails fast with retry guidance and performs no writes.

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

