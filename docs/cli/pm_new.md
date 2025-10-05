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
      --force                                     Overwrite existing project files when rerun
      --generate-cli-skeleton-json string[="-"]   Write JSON CLI input skeleton to path (stdout if omitted)
      --generate-cli-skeleton-yaml string[="-"]   Write YAML CLI input skeleton to path (stdout if omitted)
  -h, --help                                      help for new
      --here                                      Initialize the current directory instead of creating a new one
      --output string                             Output format for dry runs (text or json) (default "text")
```

### Options inherited from parent commands

```
      --config string      Path to pm config file
      --log-file string    Path to log file (default: $XDG_CACHE_HOME/pm.log)
      --log-level string   Change the log level (debug, info, warn, error) (default "info")
      --theme string       Theme for this run (nord, catppuccin, dracula, ayu)
```

### SEE ALSO

* [pm](pm.md)	 - pm is a tool to create and organize projects in your computer

