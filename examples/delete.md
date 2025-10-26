# `pm delete` Command Examples

### Preview metadata deletion
**Context:** Project `demo-service` is registered at `/srv/projects/demo-service` with `.project.hcl`, `.env`, `.gitignore`, and a Git include file.

**Command:**
```bash
pm delete demo-service --dry-run
```
**Output:**
```text
Scope: metadata
Planned artifacts:
  - registry entry
  - .project.hcl
  - git include
  - .gitignore
  - .env
No changes were applied.
```

### Dry-run with a backup destination
**Context:** Use the same project but plan a deletion that archives the workspace to `./backups/demo-service.tgz`. Remember that `--backup-destination` always requires `--backup`; the CLI will exit with an error if the destination is provided on its own.

**Command:**
```bash
pm delete demo-service --dry-run --backup --backup-destination ./backups/demo-service.tgz
```
**Output:**
```text
Scope: metadata
Planned artifacts:
  - project backup
  - registry entry
  - .project.hcl
  - git include
  - .gitignore
  - .env
Backup: ./backups/demo-service.tgz
No changes were applied.
```

### Remove metadata while keeping files
**Context:** The same project still exists; delete registry metadata and generated files while leaving the workspace intact.

**Command:**
```bash
pm delete demo-service --keep-files --force
```
**Output:**
```text
Scope: keep-files
Removed:
  - git include
  - .env
  - .project.hcl
  - .gitignore
```
