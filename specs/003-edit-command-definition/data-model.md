# Data Model: pm edit flag-driven updates

## Project Configuration
- **Identifier**: Project name (string, unique across workspace)
- **Path**: Absolute or workspace-relative filesystem path (immutable for edit command)
- **Description**: Optional string; cleared only via explicit flag or CLI input set to empty
- **Owner / Maintainers**: Optional contact metadata; when provided must include non-empty name
- **Default Environment**: Name of environment that loads by default; must reference existing environment
- **Templates / Tooling Defaults**: Structured fields for skeletons, build commands, git integration (retained when untouched)
- **Environment List**: Collection of Environment Configuration entries linked by name

## Environment Configuration
- **Name**: Unique per project; immutable during edit
- **Overrides**: Key/value pairs for environment-specific metadata (URLs, shells, workspace paths)
- **Feature Flags**: Map of flag name → boolean/string values; optional
- **Secrets / Env Vars File**: Path to environment variables file; must remain inside project scope; may be cleared explicitly if validation allows
- **Defaults / Modes**: Fields like `env_vars_mode` (e.g., `merge` or `replace`); must match supported enumeration

## Change Set Representation
- **Scope**: Project or specific environment name
- **Mutations**: Map of field identifier → pointer to value (nil when not supplied, empty string when clearing)
- **Source Metadata**: Tracks whether change originated from CLI flag or CLI input file (for error messaging)
- **Validation Impact**: Clearing a field triggers re-validation for required constraints before persistence

## Lock Metadata
- **Project Lock File**: Path `<project>/._pm.edit.lock`, containing process identifier and timestamp while held
- **Timeouts**: Lock acquisition retries for up to ~2 s before surfacing contention error
