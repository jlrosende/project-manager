# Research: CLI `new` Project Command

## Decisions
- Overwrite behavior: default FAIL when `.project.hcl` exists unless `--force` is provided
- Unknown fields in config: default FAIL; optional `--allow-unknown` permits ignoring extras
- Dry run: support `--dry-run` to preview actions and outputs without creating files
- Idempotency: safe re-run initializes only missing files; never overwrite without `--force`
- Registry: no global registry; uniqueness scoped to filesystem path

## Rationale
- Protect user data and avoid destructive changes while enabling explicit overrides
- Ensure predictable config by default; allow flexibility when requested
- Provide confidence via preview to support CI and scripting
- Idempotent behavior supports repetitive automation

## Alternatives Considered
- Prompting for confirmation: rejected for non-interactive CLI use cases
- Always overwriting: rejected due to data loss risk
- Always ignoring unknown fields: rejected due to silent misconfiguration risk
