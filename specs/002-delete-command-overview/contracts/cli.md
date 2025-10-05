# CLI Contract — `pm delete`

## Command Signature
```
pm delete <target> [flags]
```
- **target**: project name registered with `pm` or absolute path to a registered project root.

## Supported Flags
- `--all`: delete registry metadata, env vars, project configuration, workspace directory, caches, git hooks.
- `--keep-files`: delete registry metadata, env vars, and config files; preserve workspace directory and user files.
- `--only-env`: delete stored environment variables only; preserve registry metadata and filesystem assets.
- `--dry-run`: compute the deletion plan and display affected artifacts without performing destructive actions.
- `--force`: skip confirmation prompt; destructive operations execute immediately.
- `--backup`: create a timestamped archive of the project directory under `~/.pm/backups` before any destructive actions.

## Flag Interactions
- `--all`, `--keep-files`, and `--only-env` are mutually exclusive. Combining any pair returns an error before planning.
- `--backup` requires successful project resolution and a writable backup directory; on failure the command aborts without deleting artifacts.
- `--dry-run` suppresses backups and destructive operations but still reports where a backup would be written.
- `--force` overrides confirmation; otherwise the CLI prompts "Delete <project>? (y/N)" using stdin/stdout streams.

## Outputs
- **Success (non-dry run)**: structured log lines for each artifact and a final summary, e.g., `Deleted registry entry`, `Removed .env`, `Workspace archived to ~/.pm/backups/20251005-1530-sample-app.zip`. Exit code `0`.
- **Dry run**: table or bullet list enumerating artifacts with scope labels, ending with `No changes were applied.` Exit code `0`.
- **Invalid Input**: descriptive error message (e.g., `cannot combine --all with --keep-files`). Exit code `1`.
- **Project Not Found**: `project 'foo' not found` or `project not found at path /path`. Exit code `1`.
- **Read-only or Locked**: `project sample-app is currently locked` or `workspace is read-only`; no artifacts touched. Exit code `1`.
- **Backup Failure**: `failed to create backup archive: ...`; exit code `1`, no deletions performed.
- **Partial Failure**: aggregated error describing remaining artifacts (e.g., `failed to remove workspace directory: permission denied`). Exit code `1`.

## Logging Requirements
- Every destructive action logs through the logger port before and after execution, including success/failure context (`scope`, `path`, `artifact type`).
- Dry runs log a `preview` entry with the same shape for downstream tooling.
- Summary log includes scope executed, number of artifacts removed, number skipped, and backup path when applicable.

## Confirmation Flow
- Unless `--force` is set, CLI prints a prompt and waits for `y`/`yes`. Any other response aborts deletion with `operation cancelled` message and exit code `0`.
- Prompt reads from `cmd.InOrStdin()` and writes to `cmd.OutOrStdout()` to remain testable.

## Error Handling
- Errors returned by the service propagate to Cobra, which prints to stderr. CLI wraps them with context when needed (e.g., `delete project: <err>`).
- Partial deletions return errors containing artifact names; CLI instructs users to rerun with `--force` or manual cleanup steps.
