## 1. Specification Updates
- [x] 1.1 Review current `pm delete` implementation and tests for `--force`, `--backup`, and `--backup-destination` behaviors
- [x] 1.2 Update the confirmation requirement with an explicit `--force` scenario referencing observed behavior
- [x] 1.3 Expand backup handling requirement to describe the destination validation and dry-run preview semantics

## 2. Implementation Impact
- [x] 2.1 Confirm existing unit/integration tests already cover these behaviors or note follow-up test gaps *(gap: no explicit test for `--backup-destination` without `--backup`; recommend adding one)*
- [x] 2.2 Share findings with CLI owners so future code updates respect the specification changes *(action: document in change notes and final summary that maintainers should add validation tests)*
- [x] 2.3 Add or plan a test covering `--backup-destination` without `--backup` *(implemented integration test `TestCLIDelete_BackupDestinationRequiresBackup` ensuring CLI errors when destination is provided without `--backup`)*
- [x] 2.4 Verify CLI error messaging matches new spec wording for `--backup-destination`
- [x] 2.5 Ensure root command help text and documentation mention the validation behavior
