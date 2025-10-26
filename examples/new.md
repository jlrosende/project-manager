# `pm new` Command Examples

### Dry-run a new project directory
**Context:** Plan the creation of a fresh `demo-service` workspace without touching the filesystem.

**Command:**
```bash
pm new demo-service ./projects/demo-service --dry-run
```
**Output:**
```text
pm new (dry-run)
  name: demo-service
  path: ./projects/demo-service
```

### Dry-run with JSON output
**Context:** Preview the same project creation but capture the dry-run details as JSON for tooling.

**Command:**
```bash
pm new demo-service ./projects/demo-service --dry-run --output json
```
**Output:**
```json
{
  "name": "demo-service",
  "path": "./projects/demo-service",
  "here": false
}
```

### Preview adding an environment to an existing project
**Context:** Project `demo-service` already exists at `/srv/projects/demo-service`; confirm how a `staging` environment would be added.

**Command:**
```bash
pm new demo-service staging --environment-env-var API_URL=https://api.example.com --dry-run
```
**Output:**
```text
pm new (dry-run)
  name: demo-service
  path: /srv/projects/demo-service
  environment:
    name: staging
    env vars file: .staging.env
    env vars mode: merge
    env vars: 1 entries
```

### Generate a CLI input skeleton
**Context:** Produce a JSON skeleton describing the configurable fields so it can be stored in version control.

**Command:**
```bash
pm new demo-service --generate-cli-skeleton-json -
```
**Output:**
```json
{
  "name": "your-project-name",
  "here": false,
  "path": "/absolute/path/to/your-project",
  "env_vars": {
    "EXAMPLE_KEY": "VALUE"
  },
  "description": "Describe your project",
  "shell": "bash",
  "env-file": ".env"
}
```
