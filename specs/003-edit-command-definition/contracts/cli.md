# CLI Contract: pm edit

## Synopsis
```
pm edit <project> [env] [flags]
```

- `<project>` (string, required): Name or path identifier of the project to edit.
- `<env>` (string, optional): Name of environment scoped to the project. When omitted, edits apply to project-level configuration.

## Flags
| Flag | Type | Scope | Description |
|------|------|-------|-------------|
| `--cli-input` | path | project/env | Load updates from JSON or YAML file. Missing keys leave values unchanged; explicit `null`/`""` clear fields. |
| `--generate-cli-skeleton-json [path]` | optional path (`-` for stdout) | project/env | Emit JSON representation of current configuration and exit without modifications. |
| `--generate-cli-skeleton-yaml [path]` | optional path | project/env | Emit YAML representation of current configuration and exit without modifications. |
| `--allow-unknown` | boolean | project/env | Ignore extra keys found in `--cli-input` payload instead of failing validation. |
| `--dry-run` | boolean | project/env | (Inherited) Preview changes/diff without persisting. |
| `--output` | enum (`text`,`json`) | project/env | (Inherited) Select output format for command response. |

> Note: Additional field-specific flags (e.g., `--project-description`, `--env-timezone`) are generated in alignment with configuration schema; flags corresponding to immutable fields are not exposed.

## Behavior
- **Skeleton Generation**: If either skeleton flag is provided, the command outputs the serialized configuration for the given scope and exits with status 0.
- **Update Flow**: Applies `--cli-input` payload first, then overrides specified via explicit flags. Immutable fields trigger immediate error.
- **Validation**: Runs full configuration validation before persistence. Failures keep disk state unchanged and return non-zero exit code.
- **Persistence**: Writes project or environment configuration atomically using repository services; emits summary of changed fields on success.
- **Locking**: Acquires per-project lock for the full read/validate/write cycle. On contention, returns lock error with retry guidance.

## Exit Codes
| Code | Meaning |
|------|---------|
| `0` | Success (update performed or skeleton generated). |
| `2` | Validation failure or attempt to clear required field; no changes persisted. |
| `3` | Immutable field modification attempted or unknown environment specified. |
| `4` | Lock acquisition failure (another operation in progress). |
| `1` | Unexpected error (I/O failure, parse error, etc.). |

## Examples
- Generate editable YAML skeleton for project-level configuration:  
  `pm edit demo --generate-cli-skeleton-yaml ./demo-edit.yaml`
- Apply updates from file and clear description:  
  `pm edit demo --cli-input ./updates.yaml --project-description ""`
- Update staging environment config with new feature flags:  
  `pm edit demo staging --cli-input ./staging.json`
