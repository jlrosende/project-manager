# CLI Contract: `pm new`

## Command
`pm new <name> [path]`

## Arguments
- `<name>`: required project name; overrides values from CLI input files
- `[path]`: optional target directory; defaults to `<cwd>/<name>` when omitted

## Flags
- `--here`: initialize the current directory (invalid when a path argument is provided)
- `--cli-input <file>`: JSON or YAML CLI input file
- `--generate-cli-skeleton-json <file>`: write a JSON CLI input skeleton and exit
- `--generate-cli-skeleton-yaml <file>`: write a YAML CLI input skeleton and exit
- `--dry-run`: preview actions without writing
- `--force`: allow overwrite/merge when files exist
- `--allow-unknown`: ignore unknown fields in CLI input

## Behavior
- Validate inputs; merge CLI input then apply flags (flags take precedence)
- Support CLI input skeleton generation via dedicated flags (writes to file or stdout when no path supplied)
- Create `.project.hcl` and `.env` according to Data Model
- Non-zero exit on validation or creation error; no partial state
- Structured user messages; do not print secret values

## Errors
- Missing required `<name>` argument (unless generating a skeleton)
- Providing both a path argument and `--here`
- Invalid or unreadable CLI input file
- Existing `.project.hcl` without `--force`
