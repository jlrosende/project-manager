## pm new

Create a new project from arguments or configuration

### Synopsis

Create or initialize a project directory, validating inputs from positional arguments, flags, and optional CLI input files.

```
pm new <name> [path] [flags]
```

### Options

```
      --allow-unknown                             Ignore unknown fields in config files
      --cli-input string                          Path to JSON or YAML CLI input file
      --dry-run                                   Preview actions without writing files
      --env-color string                          Set environment color metadata when adding environments
      --env-file string                           Set environment vars file when adding environments
      --env-mode string                           Set environment vars merge mode (merge or replace)
      --env-var stringArray                       Environment variable in KEY=VALUE format (repeatable)
      --force                                     Overwrite existing project files when rerun
      --generate-cli-skeleton-json string[="-"]   Write JSON CLI input skeleton to path (stdout if omitted)
      --generate-cli-skeleton-yaml string[="-"]   Write YAML CLI input skeleton to path (stdout if omitted)
  -h, --help                                      help for new
      --here                                      Initialize the current directory instead of creating a new one
      --output string                             Output format for dry runs (text or json) (default "text")
```

### Environment shortcut

When `pm new` detects an existing project (registry entry plus `.project.hcl` on disk), the optional second positional argument is interpreted as an environment name. Use the environment flags to customise the environment metadata:

- `--env-file` chooses the environment vars file (defaults to `.<name>.env`).
- `--env-mode` selects the merge strategy (`merge` or `replace`).
- `--env-color` stores optional colour metadata for downstream tooling.
- `--env-var KEY=VALUE` may be repeated to add or override individual environment variables.

### Dry-run output

`--dry-run` now mirrors the actual environment or project changes. The text renderer prints the project name/path and, when applicable, an indented environment summary including the resolved file, mode, colour, and merged env var count. Using `--output json` returns the same data as a structured object with an `environment` block containing `env_vars`, `env_vars_mode`, and `env_vars_count`.

### Skeleton generation

`--generate-cli-skeleton-json` and `--generate-cli-skeleton-yaml` emit configuration templates that include the new `environment` object plus an `env_vars` map. Legacy `environments` maps are no longer produced and are ignored only when you opt-in with `--allow-unknown`.

### Options inherited from parent commands

```
      --config string      Path to pm config file
      --log-file string    Path to log file (default: $XDG_CACHE_HOME/pm.log)
      --log-level string   Change the log level (debug, info, warn, error) (default "info")
      --theme string       Theme for this run (nord, catppuccin, dracula, ayu)
```

### SEE ALSO

* [pm](pm.md)	 - pm is a tool to create and organize projects in your computer

