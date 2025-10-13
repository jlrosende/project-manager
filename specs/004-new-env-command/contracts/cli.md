# CLI Contract: `pm new`

## Positional Arguments
1. `name` (required): Project identifier.
2. `arg2` (optional):
   - Treated as project path when the project does not yet exist (registry miss or `.project.hcl` absent).
   - Treated as environment name when the project exists and `.project.hcl` is present.

## Flags (additions)
- `--env-file <path>`: Overrides environment vars file name (defaults to `.<slug>.env`).
- `--env-mode <merge|replace>`: Sets env vars merge strategy (default `merge`).
- `--env-color <value>`: Optional color metadata for downstream tooling.
- `--env-var KEY=VALUE`: Repeatable flag appending/overriding environment variables.

## Behavior Matrix

| Scenario | Inputs | Outcome |
|----------|--------|---------|
| First-time project creation | `pm new demo`, optional `--here` or path | Creates project, no environment processing. |
| Repeat invocation without env | `pm new demo` | Returns success with message “project already exists”; no mutations. |
| Environment addition | `pm new demo staging [flags]` | Validates existence, merges environment inputs, adds environment (force overwrites metadata/file when `--force`). |
| Environment attempt before project exists | `pm new demo staging` | Fails with project-not-found style error (same for `--dry-run`). |
| Dry-run environment addition | `pm new demo staging --dry-run` | Outputs merged project + environment details in text/JSON without writing files. |
| Legacy config with `--allow-unknown` | Config includes deprecated `environments` map | Legacy map ignored silently; command continues using new `environment` object if present. |

## Output Formats
- **Text dry-run**: Lists `name`, `path`, `env vars` count, and environment details when applicable.
- **JSON dry-run**: Structured object containing project fields and nested environment payload (name, file, mode, color, env_vars map).
- **Errors**: Include human-readable error codes for duplicates, validation failures, or missing project artifacts.
