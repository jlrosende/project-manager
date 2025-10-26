.. _pm_edit:

pm edit
-------

Edit project configuration

Synopsis
~~~~~~~~


Edit project configuration

::

  pm edit <project> [env] [flags]

Options
~~~~~~~

::

      --allow-unknown                             Ignore unknown fields in CLI input files
      --cli-input string                          Path to JSON or YAML CLI input file
      --dry-run                                   Preview changes without persisting them
      --env-color string                          Set the environment color (requires <env>)
      --env-env-vars-file string                  Set the environment env vars file path (requires <env>)
      --env-env-vars-mode string                  Set the environment env vars mode (merge or replace, requires <env>)
      --generate-cli-skeleton-json string[="-"]   Write JSON CLI input skeleton to path (stdout if omitted)
      --generate-cli-skeleton-yaml string[="-"]   Write YAML CLI input skeleton to path (stdout if omitted)
  -h, --help                                      help for edit
      --output string                             Output format for results (text or json) (default "text")
      --project-default-env string                Set the default environment
      --project-description string                Set the project description
      --project-env-vars-file string              Set the project-level env vars file path
      --project-shell string                      Set the default project shell

Options inherited from parent commands
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

::

      --config string      Path to pm config file
      --log-file string    Path to log file (default: $XDG_CACHE_HOME/pm.log)
      --log-level string   Change the log level (debug, info, warn, error) (default "info")
      --theme string       Theme for this run (nord, catppuccin, dracula, ayu)

SEE ALSO
~~~~~~~~

* `pm <pm.rst>`_ 	 - pm is a tool to create and organize projects in your computer

