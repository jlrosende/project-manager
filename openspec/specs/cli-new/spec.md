# cli-new Specification

## Purpose
Define the required behavior for the `pm new` command when creating or updating local projects, including environment shortcuts, configuration inputs, and dry-run reporting.
## Requirements
### Requirement: Project Creation Workflow
The CLI SHALL create or initialize a project directory when invoked with a project name and optional path while enforcing safe defaults.

#### Scenario: Create project from arguments
- **WHEN** a user runs `pm new <name> [path]` and passes valid inputs
- **THEN** the command creates or initializes the target directory with `.project.hcl` and `.env`, registers the project, and reports success.

#### Scenario: Initialize current directory
- **WHEN** a user runs `pm new <name> --here`
- **THEN** the command initializes the current working directory without creating a new folder and preserves unrelated files.

#### Scenario: Reject conflicting options
- **WHEN** a user supplies both a positional path and `--here`
- **THEN** the command fails with an actionable error explaining the conflict and performs no writes.

#### Scenario: Repeat invocation no-ops
- **WHEN** a user reruns `pm new <name>` after the project is already registered with an existing `.project.hcl`
- **THEN** the command exits successfully without recreating artifacts and prints a message explaining that the project already exists.

### Requirement: External Configuration Inputs
The CLI SHALL merge configuration from CLI flags, positional arguments, and optional JSON/YAML input files with flags taking precedence.

#### Scenario: Consume CLI input file
- **WHEN** a user runs `pm new <name> --cli-input config.yaml`
- **THEN** the command parses the file, validates required fields, merges values with CLI flags, and creates the project using the merged configuration.

#### Scenario: Generate skeleton files
- **WHEN** a user runs `pm new --generate-cli-skeleton-json` or `--generate-cli-skeleton-yaml`
- **THEN** the command emits a structured skeleton to the provided path (or stdout) that contains project and environment placeholders.

#### Scenario: Handle unknown fields
- **WHEN** a CLI input file contains unknown fields and the user omits `--allow-unknown`
- **THEN** the command fails validation with a clear error identifying the unsupported keys.

#### Scenario: Allow unknown fields when requested
- **WHEN** a CLI input file contains unknown fields and the user passes `--allow-unknown`
- **THEN** the command accepts the file, ignores the unknown keys without emitting warnings, and proceeds with the known configuration values.

### Requirement: Dry Run and Reporting
The CLI SHALL support a dry-run mode that previews planned actions with deterministic output formats.

#### Scenario: Text dry run output
- **WHEN** a user runs `pm new <name> --dry-run`
- **THEN** the command prints a text summary describing the planned project path, created files, and environment details without modifying the filesystem.

#### Scenario: JSON dry run output
- **WHEN** a user runs `pm new <name> --dry-run --output=json`
- **THEN** the command returns structured JSON containing the project and environment artifacts that would be produced.

#### Scenario: Dry run for invalid state
- **WHEN** a user runs `pm new <name> <env> --dry-run` for a project that is not yet registered
- **THEN** the command fails with the same error as execution, indicating the project must exist before adding an environment.

### Requirement: Environment Shortcut Workflow
The CLI SHALL add or update environments for existing projects using the positional environment argument and dedicated flags.

#### Scenario: Add environment to existing project
- **WHEN** a user runs `pm new <project> <env>` for a registered project
- **THEN** the command delegates to the environment service to create the environment file with merged configuration and confirms success.

#### Scenario: Duplicate environment name
- **WHEN** a user runs `pm new <project> <env>` for an environment that already exists without `--force`
- **THEN** the command fails with a descriptive error and leaves existing environment files untouched.

#### Scenario: Force overwrite environment
- **WHEN** a user runs `pm new <project> <env> --force`
- **THEN** the command overwrites the environment metadata and env file without creating backups and reports the overwrite.

#### Scenario: Merge environment configuration sources
- **WHEN** a user combines positional arguments, `--environment-*` flags, and an `environment` object in a CLI input file
- **THEN** the command applies precedence where CLI flags override file values and reports the merged environment in dry-run and success messages.

#### Scenario: Missing project blocks environment creation
- **WHEN** a user runs `pm new <project> <env>` but the project is not registered
- **THEN** the command fails with an error instructing the user to create the project first and makes no changes.

#### Scenario: Environment name from configuration file
- **WHEN** a user omits the positional environment argument but supplies an `environment.name` value in a CLI input file for an existing project
- **THEN** the command uses the configuration value as the environment name unless a CLI flag overrides it.

