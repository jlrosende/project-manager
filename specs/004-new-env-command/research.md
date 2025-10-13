# Research Summary

## Decision: Environment handling distinguishes project creation vs addition
- **Rationale**: Aligns with clarified spec to treat the second positional argument as a path only when the project is absent, and as an environment name once the project exists, preventing accidental overwrites of existing projects.
- **Alternatives considered**: Continue heuristic parsing of path-like arguments; rejected because it caused ambiguity and conflicted with new shortcut expectations.

## Decision: Introduce `EnvironmentInput` struct and extend config schema
- **Rationale**: Provides a typed container for merged environment data coming from CLI, config files, and defaults, simplifying validation and dry-run output.
- **Alternatives considered**: Reuse existing maps in `ProjectDefinition`; rejected due to poor clarity, lack of field-level precedence control, and difficulty exposing structured skeletons.

## Decision: Add environment customization flags and merge precedence rules
- **Rationale**: Flags (`--env-file`, `--env-mode`, `--env-color`, `--env-var`) let users override config entries explicitly while keeping CLI defaults; precedence ensures deterministic results.
- **Alternatives considered**: Require users to edit config files for environment tweaks; rejected for poor ergonomics and failure to meet spec.

## Decision: Silent handling of legacy `environments` maps when `--allow-unknown`
- **Rationale**: Matches spec requirement to ignore legacy data without warnings when the user opts into unknown fields, minimizing noise for migration scenarios.
- **Alternatives considered**: Emit warnings or errors despite the flag; rejected as contradictory to clarified behavior.

## Decision: Ignore stale registry entries missing `.project.hcl` during uniqueness checks
- **Rationale**: Prevents orphaned include entries from blocking recreation while maintaining safety for real conflicts.
- **Alternatives considered**: Automatic cleanup or hard failures; rejected to avoid destructive actions and keep compatibility with existing maintenance flows.
