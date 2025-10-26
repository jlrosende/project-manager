# `pm new` Enhanced Behavior Specification

## User Stories & Priorities
- **US1 (P1)**: As a CLI user, I need `pm new` to detect existing projects and avoid destructive reinitialization so repeated runs are safe.
- **US2 (P1)**: As a CLI user, I can add environments via `pm new <name> <env>` with customization flags/config once the project exists.
- **US3 (P2)**: As a CLI user, I can preview environment changes and generate skeleton configs reflecting the new schema.

## Project Detection & Routing
- Treat a project as “existing” only when both of these are true: the project is present in the registry and its `.project.hcl` file exists on disk (Q1=C). Missing metadata means the invocation behaves like a fresh project creation attempt (Q22=B).
- The second positional argument is always interpreted as a filesystem path when the project does **not** exist; once the project is registered, that argument is treated as an environment name, even if it resembles a path (Q10=A).
- Running `pm new <name>` with no second argument behaves differently based on existence: it scaffolds the project during first-run; on subsequent runs it returns a no-op success explaining that the project already exists (Q9=B).
- `--here` continues to forbid any positional path or environment argument (existing behaviour, restated for clarity).

## Environment Shortcut Semantics
- Environment creation is blocked during the same run as project creation. If environment flags or config data are supplied while the project is missing, the command fails with guidance to create the project first (Q13=B, Q18=A, Q28=A).
- Calling `pm new <project> <env>` on an existing project invokes `ProjectService.AddEnvironment`. Duplicate environment names produce a clear error; `--force` overwrites metadata and rewrites the env file without backups (Q6=A, Q20=A, Q25=A).
- `--environment-color`, `--environment-env-file`, `--environment-mode`, and `--environment-env-var KEY=VALUE` form the new CLI flag surface for environment customization; the positional argument remains the sole way to set the environment name (Q2, Q16).
- For ambiguity, the CLI always prefers environment mode once the project is known; users must use explicit flags for paths if they intend to re-target creation (Q10=A).

## Configuration & Skeletons
- CLI input files now expose a structured `environment` object with fields mirroring the CLI flags plus an `env_vars` map for key/value pairs (Q3=B, Q31=A).
- Skeleton JSON/YAML output includes that `environment` object with placeholders for `env_vars_file`, `env_vars_mode`, `color`, and an empty `env_vars` map (Q5=B, Q11=D).
- Legacy top-level `environments` maps are rejected unless the user passes `--allow-unknown`, in which case they are silently ignored and treated as unknown keys (Q17=B, Q34=D, Q36). No warning or logging is emitted when ignored (Q36 answer + Q37=A).

## Environment Data Merging
- Environment values can come from three sources: positional name, CLI flags (`--environment-env-file`, `--environment-mode`, `--environment-color`, `--environment-env-var KEY=VALUE`), and the config `environment` object.
- Precedence rules:
  - Positional argument sets the environment name. If omitted and config supplies a name, the config name is used (Q12=A).
  - CLI flags override matching fields in the config object (Q24=A).
  - CLI `--environment-env-var` entries override individual keys from the config `env_vars` map; collisions favour the flag-supplied value (Q30=B, Q32=A).
- Dry-run and execution consume the merged environment definition; no attempt is made to stage environment creation for later runs.

## Dry-Run Behaviour
- Dry-run output always includes both project and environment data when an environment would be created, regardless of whether the name came from the positional argument or the config file (Q14=A, Q19=A, Q35=C).
- If the project does not yet exist but the config file supplies environment data, dry run fails with the same error as the real run (Q23=A).
- Text dry runs list project name/path and total environment vars after merging. JSON dry runs include the complete environment structure (name, mode, file, color, env_vars) with CLI overrides applied.

## Registry & Uniqueness Rules
- Global uniqueness is enforced via a shared helper used by CLI and services (Q8=D). Entries lacking a `.project.hcl` file are ignored during uniqueness checks so they do not block recreation (Q27=A).
- When recreating a project whose registry entry still exists, `pm new` skips adding a duplicate includeIf entry if one already references the resolved path (Q33=C) but does not clean up the stale entry automatically (Q26=B, Q29=B).
- The CLI continues to require unique project names and paths; rename operations must still reject conflicts to maintain a 1:1 mapping.

## Error Handling & Messages
- Missing projects: executing `pm new <name> <env>` when the project cannot be resolved produces the existing “project not found” error. With dry-run, the same error is surfaced immediately (Q23=A).
- Duplicate environment names, path collisions, and validation failures continue to return CLI-friendly errors augmented with human-readable codes (Q7=C, Q21=A).
- Legacy maps ignored under `--allow-unknown` quiet mode generate no user-visible feedback (Q36 custom answer, confirmed by Q37=A).

## Summary of Feature Changes
1. Dual-mode positional handling with strict existence detection and graceful no-op messaging for repeat runs.
2. Environment shortcut enhancements with dedicated flags, config schema, merge precedence, and dry-run visibility.
3. Updated skeleton templates and config parsing that reject legacy maps (unless ignored) and support env var payloads.
4. Registry and uniqueness refinements that tolerate missing `.project.hcl` files while avoiding duplicate includeIf entries.
5. Simplified force semantics and dry-run parity, ensuring users always see the exact environment that will be created.
