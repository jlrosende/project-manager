# Phase 0 Research — CLI `pm delete`

## Decision: Service orchestration and project state checks
- **Rationale**: `internal/core/services/project_service.go` currently exposes a shallow `Delete(name string)` method that only delegates to the project repository. Implementing `DeleteProject(ctx, options)` requires expanding the service to: resolve projects by name or path, consult repository metadata for read-only/lock flags, coordinate filesystem/env/git repositories, and emit logger entries at each step. This keeps destructive workflow knowledge centralized in the service while the CLI remains a thin adapter.
- **Alternatives considered**: Performing deletions directly inside the CLI command would violate the port-driven architecture and duplicate logic needed by potential future frontends. Introducing a brand-new service dedicated to deletions was also considered but rejected because project lifecycle responsibilities already live inside `ProjectService` and extending it preserves cohesion.

## Decision: Port extensions for deletion workflows
- **Rationale**: Ports under `internal/core/ports` lack the methods needed to remove project directories, delete env var files, prune git include-if entries, or persist registry updates. Extending the filesystem port with recursive removal, staged backup/restore helpers, and dry-run previews, plus adding delete operations to env vars, git, and project repositories, keeps adapters honest and ensures mocks capture new failure modes for tests.
- **Alternatives considered**: Embedding destructive helpers directly into domain services would blur architectural boundaries. Wrapping existing ports in ad-hoc helper structs would complicate dependency injection and scatter deletion responsibilities across layers.

## Decision: Backup implementation approach
- **Rationale**: Use Go's `archive/zip` writer walking the project directory via `filepath.WalkDir`, staging archives inside a temp directory under `~/.pm/backups` before atomically renaming them into place. ZIP is already available in the standard library, works cross-platform, and preserves file permissions. Staging avoids partially written files on failure and satisfies the specification's atomicity requirement.
- **Alternatives considered**: `compress/gzip` plus tar would also work but requires manual tar header handling or external dependencies. Third-party archivers were rejected to avoid expanding dependencies and to keep the workflow reproducible.

## Decision: Dry-run and logging strategy
- **Rationale**: The service should build a `ProjectDeletePlan` detailing targeted artifacts (registry entry, env files, skeleton directories, git includes, backups). During dry runs, the CLI can render this plan without invoking port mutations. During actual runs, the service will iterate through the plan, invoking ports and emitting structured `logger` entries for each side effect, enabling consistent summaries and error propagation.
- **Alternatives considered**: Relying on adapters to compute their own previews risks divergent reporting. Logging exclusively in the CLI would miss repository-level errors and complicate integration tests that currently rely on service mocks.

## Decision: Confirmation flow reuse
- **Rationale**: Existing CLI commands rely on Cobra flag parsing and bootstrap-provided loggers, but there is no shared prompt helper. Implementing a small confirmation helper inside the new delete package that wraps current terminal I/O (using `cmd.InOrStdin`/`cmd.OutOrStdout`) keeps the utility local while still reusing Cobra's flag plumbing and avoids introducing a new dependency. The helper can later be promoted if additional commands require confirmation prompts.
- **Alternatives considered**: Integrating Bubble Tea for CLI prompts would add UI overhead for a single yes/no confirmation. Delegating confirmation to the service was rejected because prompting belongs to adapters and would hinder automation in other frontends.
