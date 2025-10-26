# `pm edit` Command Examples

### Set a project description
**Context:** The `demo-service` project currently has an empty description.

**Command:**
```bash
pm edit demo-service --project-description "Internal tooling suite"
```
**Output:**
```text
Updated project "demo-service":
  Description: "" -> "Internal tooling suite"
```

### Choose a default environment
**Context:** Project `demo-service` defines a `production` environment but no default is configured yet.

**Command:**
```bash
pm edit demo-service --project-default-env production
```
**Output:**
```text
Updated project "demo-service":
  Default Environment: "" -> "production"
```

### Update environment metadata
**Context:** Environment `production` currently uses mode `merge`, env file `.production.env`, and has no color tag.

**Command:**
```bash
pm edit demo-service production --env-env-vars-file ./env/production.env --env-env-vars-mode replace --env-color "#FF8800"
```
**Output:**
```text
Updated environment "production" in project "demo-service":
  Color: "" -> "#FF8800"
  Env Vars Mode: "merge" -> "replace"
  Env Vars File: ".production.env" -> "./env/production.env"
  Other environments unchanged; project metadata untouched.
```

### Dry-run an environment change
**Context:** After the previous update, the `production` environment’s mode is `replace`; preview switching it back to `merge`.

**Command:**
```bash
pm edit demo-service production --env-env-vars-mode merge --dry-run
```
**Output:**
```text
pm edit demo-service production (dry-run)
  Env Vars Mode: "replace" -> "merge"
  Other environments unchanged; project metadata untouched.
```

### Attempt an edit with no changes
**Context:** Run `pm edit` without flags to confirm the CLI reports when nothing would change.

**Command:**
```bash
pm edit demo-service --dry-run
```
**Output:**
```text
No changes applied; nothing to update.
```
