## MODIFIED Purpose
Document the standard layout, entrypoint structure, and helper conventions every CLI command package must follow so new subcommands stay consistent, discoverable, and easy to review.

### Background
- Consistency across command packages keeps bootstrap, flag wiring, and helpers predictable for reviewers and new contributors.
- Shared patterns reduce the risk of regressions when wiring additional commands or updating existing ones.

## MODIFIED Requirements
### Requirement: CLI Command Module Layout
All packages under `internal/adapters/handlers/cli/<command>` SHALL follow the same structure demonstrated by `pm edit`, reinforcing the background guidance for consistency and reviewability:
- File layout: each command directory MUST provide `command.go` containing the exported `Command()` constructor and flag declarations. Supplemental files such as `run.go`, `output.go`, or `helpers.go` MAY exist but must contain only private helpers referenced by `command.go`.
- Expose a `Command() *cobra.Command` constructor that configures `Use`, `Short`, validation, and assigns a package-private `run` function to `RunE`. The constructor MUST be the only exported symbol in the package unless a type is needed solely for testing hooks.
- Declare CLI flag names as package-level `const` identifiers inside a single `const (...)` block, using the `flag<PascalCase>` naming convention and registering flags in grouped sections (global flags first, then project-level, then environment-level where applicable).
- Provide a `run(cmd *cobra.Command, args []string) error` function that orchestrates bootstrap wiring, service acquisition, flag parsing, and delegation to helper functions.
- Keep side-effect-free helpers (rendering, skeleton output, exit code wrappers) in private functions that accept explicit dependencies instead of relying on globals. Helpers MUST be grouped by responsibility (e.g., `render*`, `format*`, `output*`) and reside in the same package.
- Wrap service bootstrap and adapter acquisition inside the `run` function, returning contextualized errors using `fmt.Errorf("context: %w", err)` when delegating to domain services.
- Ensure command output helpers produce deterministic phrasing matching the conventions from `pm edit` (dry-run suffix, change summaries, and validation error formatting).
- Bring the existing handlers for `pm delete`, `pm init`, `pm list`, and `pm new` into compliance with this structure so that all exported entry points are `Command()` constructors instead of exported Cobra variables.

#### Scenario: Refactor delete command to standard layout
- **WHEN** maintainers align `internal/adapters/handlers/cli/delete` with the shared layout
- **THEN** the package exposes `Command()` with a private `run`, groups all flag constants in one block, and routes execution through the bootstrap/service flow before delegating to dedicated output helpers.

#### Scenario: Refactor init and list commands
- **WHEN** the `pm init` and `pm list` handlers are refactored
- **THEN** the previous `InitCmd` and `ListCmd` globals disappear in favor of `Command()` constructors, and flag registration plus helper functions mirror the `pm edit` patterns.

#### Scenario: Close gaps in pm new handler
- **WHEN** the `pm new` handler is audited against the standard layout
- **THEN** any remaining deviations (constructor signature, helper placement, or bootstrap flow) are corrected so it aligns with `pm edit`.

#### Scenario: Guardrails for new command modules
- **WHEN** a new CLI command module is introduced under `internal/adapters/handlers/cli`
- **THEN** code review or automated linting detects any deviation from the standard layout (missing `Command()` constructor, scattered flag constants, or implicit globals) and fails until the structure matches the documented pattern.
