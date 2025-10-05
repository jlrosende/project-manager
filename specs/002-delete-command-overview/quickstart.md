# Quickstart — Validating `pm delete`

## Prerequisites
- Existing project created via `pm new sample-app` with populated `.project.hcl`, `.env`, git include-if entries, and optional generated skeleton directory.
- Terminal session with `pm` binary built from feature branch.
- Test workspace includes writable backup directory at `~/.pm/backups`.

## Smoke Test
1. Run `pm delete sample-app`.
2. When prompted, answer `y`.
3. Verify output lists removal of registry entry, config file, and env vars, and reports success summary.
4. Confirm `.project.hcl` and env files are removed while workspace directory still exists.

## Dry Run
1. Execute `pm delete sample-app --dry-run --all`.
2. Ensure command exits without modifying filesystem.
3. Confirm output enumerates workspace directory, caches, git hooks, env vars, and registry entry slated for deletion.

## Force Deletion with Backup
1. Run `pm delete sample-app --all --backup --force`.
2. Confirm command skips prompt, logs backup creation, then removal of all artifacts.
3. Validate archive exists under `~/.pm/backups/<timestamp>-sample-app.zip` and working directory is removed.

## Keep Files
1. Recreate project (e.g., `pm new sample-app`).
2. Run `pm delete sample-app --keep-files` and confirm metadata removed but workspace directory remains with user files intact.

## Env Only
1. Store environment variables via CLI or manual edit.
2. Run `pm delete sample-app --only-env` and confirm only env files are removed, `.project.hcl` and workspace survive.

## Failure Modes
- Attempt deletion of non-existent project: expect descriptive "project not found" error.
- Simulate read-only workspace (chmod 0555) and run `pm delete sample-app --all`: expect refusal with instructions to unlock workspace.
- Force backup failure (e.g., set backup directory read-only): command should abort before deleting anything and report archive write error.
