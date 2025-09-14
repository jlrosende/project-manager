## pm list

list projects

### Synopsis

List all projects

```
pm list <project> [flags]
```

### Options

```
      --commit.gpgsign            git commit.gpgsign (default git --global) (default true)
      --env-vars stringToString   List of ENV_VARS to add to the environment (default [])
      --format string             output format
  -h, --help                      help for list
      --shell string              Shell of the project, need be installed in the system) (default to $SHELL env var))
      --tag.gpgsign               git tag.gpgsign (default git --global) (default true)
      --user.email string         git user.email (default git --global)
      --user.name string          git user.name (default git --global)
      --user.signingkey string    git user.signingkey (default git --global)
```

### Options inherited from parent commands

```
      --log-file string    Path to log file (default: $XDG_CACHE_HOME/pm.log)
      --log-level string   Change the log level (debug, info, warn, error) (default "info")
```

### SEE ALSO

* [pm](pm.md)	 - pm is a tool to create and organize projects in your computer

