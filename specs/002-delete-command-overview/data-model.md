# Data Model — CLI `pm delete`

## Entities

### ProjectDeleteOptions
- **Fields**:
  - `Target`: resolved `ProjectIdentifier` containing name, path, and registry metadata.
  - `Scope`: enum (`DeleteScopeAll`, `DeleteScopeKeepFiles`, `DeleteScopeOnlyEnv`).
  - `DryRun`: bool toggling non-destructive execution.
  - `Force`: bool to skip confirmation prompts.
  - `Backup`: struct pointer detailing requested archive settings (nil when omitted).
- **Relationships**: Consumed by `ProjectService.DeleteProject`. Passed from CLI adapter after parsing flags.
- **Validation Rules**:
  - Exactly one scope flag may be set; omitted scope defaults to metadata-only cleanup (equivalent to `DeleteScopeKeepFiles`).
  - Backup requires resolved project path.
  - `Force` implies confirmation bypass.

### ProjectIdentifier
- **Fields**:
  - `Name`: canonical project name from registry (optional when invoked by path).
  - `Path`: absolute filesystem path for project root.
  - `RegistryID`: internal identifier for the project repository entry.
  - `Status`: lifecycle state including `Locked` and `ReadOnly` flags.
- **Relationships**: Derived from `ProjectRepository`. Used by service to evaluate locks and to update registry post-deletion.

### DeleteScope (enum)
- `DeleteScopeMetadata`: Remove registry entry, `.project.hcl`, env var metadata, git include, but keep working directory contents.
- `DeleteScopeAll`: Superset of `DeleteScopeMetadata` plus workspace directory, caches, git hooks, user-authored files.
- `DeleteScopeEnvOnly`: Restrict to environment variable secrets while leaving configuration and files intact.

### ProjectDeletePlan
- **Fields**:
  - `Scope`: `DeleteScope` value.
  - `Artifacts`: ordered slice of `DeletionArtifact` representing each step (registry removal, env file deletion, directories, git includes, backup location).
  - `Backup`: optional `BackupArtifact` describing archive path and size preview.
- **Relationships**: Created during dry run or before execution to unify logging and dry run output.

### DeletionArtifact
- **Fields**:
  - `Type`: enum (`Registry`, `EnvVars`, `ConfigFile`, `GitInclude`, `SkeletonDir`, `WorkspaceDir`, `BackupArchive`, `Cache`, `Hook`).
  - `Path`: filesystem path or logical identifier affected.
  - `Description`: human-readable explanation for logs/UI.
- **Relationships**: `ProjectDeletePlan` enumerates artifacts; service iterates through them to call appropriate ports.

### BackupArtifact
- **Fields**:
  - `TempPath`: staging directory within `~/.pm/backups/tmp`.
  - `FinalPath`: final archive location.
  - `Created`: bool flag used to avoid deleting incomplete archives.
  - `SizeBytes`: computed after writing to report summary data.
- **Relationships**: Bound to `ProjectDeletePlan` when `Backup` requested. Filesystem adapter writes archives based on this struct.

### ProjectDeleteResult
- **Fields**:
  - `Scope`: executed deletion scope.
  - `ArtifactsRemoved`: list of `DeletionArtifact` entries confirmed removed.
  - `ArtifactsSkipped`: entries postponed or errored with reason.
  - `BackupPath`: final archive path when created.
  - `Errors`: aggregated partial failure details.
- **Relationships**: Returned by `ProjectService.DeleteProject` to CLI for summarizing output and exit status.

## State Transitions
1. **Pending → Planning**: CLI resolves target using `ProjectRepository` and builds `ProjectDeleteOptions`.
2. **Planning → Confirmed**: Service validates locks/read-only flags, builds `ProjectDeletePlan`, and either returns preview (dry run) or waits for confirmation from CLI.
3. **Confirmed → Executing**: Service iterates artifacts, invoking ports (`Filesystem`, `EnvVars`, `Git`, `Project`) and logging each step.
4. **Executing → Completed**: On success, registry entry removed and optional backup path returned; CLI exits with success summary.
5. **Executing → PartialFailure**: If any artifact removal fails, the error is aggregated, remaining steps stop, and CLI reports partial cleanup with non-zero exit code.

## Validation Summary
- Scope flags are mutually exclusive and required to be coherent with backup (backup only valid when working directory accessible).
- Project must not be locked or read-only; service returns error before planning when flags indicate ongoing operations.
- Backup operations must succeed before destructive steps; failure aborts deletion.
- Dry-run mode bypasses all port mutations and simply returns `ProjectDeletePlan` for rendering.
