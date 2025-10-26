## Why
The CLI currently lacks flows to bring external project definitions into the workspace or export managed projects for reuse. Introducing dedicated commands for import and export will allow users to round-trip project definitions without relying on manual file manipulation.

## What Changes
- Add placeholder CLI commands `pm import` and `pm export` that follow the shared command layout and establish flag stubs.
- Document TODO items describing the future integration logic for reading and writing project data.
- Update OpenSpec to describe the new commands and planned behaviors at a high level.

## Impact
- Affected specs: `cli-import-export`
- Affected code: `internal/adapters/handlers/cli/import`, `internal/adapters/handlers/cli/export`, `internal/adapters/handlers/cli/root.go`
