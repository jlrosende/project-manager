## Why
The cli-delete specification lacks coverage for several edge cases observed during the review: forcing deletion without confirmation, validating that `--backup-destination` requires `--backup`, and clarifying how dry runs interact with backups.

## What Changes
- Capture the behavior of `--force` explicitly skipping confirmation while still honoring safety checks.
- Describe validation for `--backup-destination` so it errors when used without `--backup`.
- Document that `--dry-run --backup` previews the archive path without creating files.

## Impact
- Affected specs: `cli-delete`
- Affected code: specification only
