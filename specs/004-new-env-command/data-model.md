# Data Model

## Entities

### ProjectDefinition (existing, updated usage)
- **Fields**:
  - `Name` (string) – project identifier, globally unique.
  - `Here` (bool) – indicates in-place initialization.
  - `Path` (string) – absolute target path (required when `Here` is false).
  - `Metadata` (map[string]string) – auxiliary properties (unchanged).
  - `Environment` (*EnvironmentInput) – optional pointer populated when adding or previewing an environment (new linkage).
- **Relationships**: Links to `EnvironmentInput` when environment creation shortcut is in effect.
- **Validation Rules**:
  - Either `Here` or `Path` must be provided (existing rule).
  - `Name` must be non-empty and remain globally unique.
  - When `Environment` is provided, the project must already exist on disk/registry.

### EnvironmentInput (new)
- **Fields**:
  - `Name` (string) – required environment identifier (from positional argument or config).
  - `Color` (string, optional) – CLI `--env-color` or config `environment.color`.
  - `EnvVarsMode` (string) – required; accepts `merge` (default) or `replace`.
  - `EnvVarsFile` (string) – required path relative to project root (default `.<slug>.env`).
  - `EnvVars` (map[string]string) – environment key/value pairs combined from config + CLI `--env-var` flags.
- **Validation Rules**:
  - `Name` must be unique within project and match allowed slug pattern.
  - `EnvVarsMode` must be one of `merge` or `replace`.
  - `EnvVarsFile` must resolve within project directory (no traversal outside workspace).
  - `EnvVars` may be empty but must not contain duplicate keys post-merge.

### ProjectExistenceProbe (conceptual service helper)
- **Fields**:
  - `RegistryHit` (bool) – indicates project entry exists in repository.
  - `ProjectFilePresent` (bool) – indicates `.project.hcl` file exists on disk.
- **Usage**: Consumed by CLI/service to decide between creation vs environment mode.

## State Transitions

1. **Project Creation**
   - Pre-state: `ProjectExistenceProbe` reports `false/false`.
   - Action: CLI merges inputs without `Environment`; creator persists `.project.hcl` and registry.
   - Post-state: Project entry exists, `.project.hcl` present.

2. **Environment Addition**
   - Pre-state: `ProjectExistenceProbe` reports `true/true`.
   - Action: CLI populates `EnvironmentInput`; service validates uniqueness, applies force semantics if requested, writes env file/metadata.
   - Post-state: Environment appended; dry run mirrors expected output.

3. **Dry Run Failure for Missing Project**
   - Pre-state: `ProjectExistenceProbe` reports `false/*` while `EnvironmentInput` populated.
   - Action: CLI aborts before generating preview, returning “project not found” style error.
