## pm delete

Delete a registered project and its artifacts

### Synopsis

`pm delete` removes registry metadata, environment variables, and optional workspace files for a registered project. The `<target>` positional argument accepts either a registered project name or a filesystem path. Unless `--force` is provided, the CLI displays a themed confirmation modal before removing anything. Use `--dry-run` to print the planned artifacts without making changes. When `--backup` is set the command writes an archive to `--backup-destination`, falling back to the configured backup directory or `~/.pm/backups`.

```
pm delete <target> [flags]
```

### Examples

```
pm delete sample-app
pm delete sample-app --keep-files
pm delete ~/work/sample-app --all --backup --backup-destination ~/archives/sample-app.zip
```

### Options

```
      --all                         Remove project metadata, env vars, and workspace files
      --backup                      Create a backup archive before deleting
      --backup-destination string   Custom destination for the backup archive
      --dry-run                     Preview deletion steps without making changes
      --force                       Skip confirmation prompt
  -h, --help                        help for delete
      --keep-files                  Remove metadata but keep workspace files
      --only-env                    Remove stored environment variables only
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

