# Data Model: CLI `new` Project Command

## Entities

### ProjectDefinition
- name (string, required)
- here (bool, mutually exclusive with path)
- path (absolute path string, mutually exclusive with here)
- environments (map[string]string, optional)
- metadata (map[string]string, optional)

Validation rules:
- Exactly one of here or path MUST be provided
- path MUST resolve to absolute when provided and be creatable
- name MUST be non-empty and match allowed pattern `[a-zA-Z0-9-_]+`

### ConfigInput (JSON/YAML)
- Mirrors ProjectDefinition; may include optional fields
- Flags override file values during merge
- Unknown fields cause error unless `--allow-unknown` is passed
