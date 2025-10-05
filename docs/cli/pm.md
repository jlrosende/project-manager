## pm

pm is a tool to create and organize projects in your computer

### Synopsis

A tool to manage the configuration and structure of multiple projects inside your computer. Launch the TUI with `pm`, create new projects with `pm new`, and remove existing projects with `pm delete` using optional dry-run and backup support.

```
pm [project] [env|[path]]  [flags]
```

### Options

```
      --config string      Path to pm config file
  -h, --help               help for pm
  -l, --list               List all the projects.
      --log-file string    Path to log file (default: $XDG_CACHE_HOME/pm.log)
      --log-level string   Change the log level (debug, info, warn, error) (default "info")
      --theme string       Theme for this run (nord, catppuccin, dracula, ayu)
```

### SEE ALSO

* [pm delete](pm_delete.md)	 - Delete a registered project and its artifacts
* [pm edit](pm_edit.md)	 - edit project
* [pm init](pm_init.md)	 - Initialize a your workspace
* [pm list](pm_list.md)	 - list projects
* [pm new](pm_new.md)	 - Create a new project from arguments or configuration

