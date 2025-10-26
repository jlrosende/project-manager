## Why
The CLI subcommands in `internal/adapters/handlers/cli` were implemented at different points in time and now follow slightly diverging structures for flag constants, command construction, and execution helpers. The `pm edit` handler demonstrates the current best practices, but other commands still use older patterns that complicate cross-command maintenance and linting. We need a proposal to enforce a consistent code format so future commands are easier to review and extend.

## What Changes
- Document a standardized command module layout based on the `pm edit` handler, covering exported constructors, flag constant naming, and separation of execution helpers.
- Refactor existing handlers for `pm delete`, `pm init`, `pm list`, and `pm new` so each exposes a `Command()` constructor, groups flag constants, and shares helper patterns with `pm edit`.
- Add lint/test coverage that ensures new command handlers adhere to the documented structure.

## Impact
- Affected specs: cli-command-format
- Affected code: `internal/adapters/handlers/cli/{delete,init,list,new}`, CLI-focused tests under `tests/`
