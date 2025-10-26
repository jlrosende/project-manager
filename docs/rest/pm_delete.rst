.. _pm_delete:

pm delete
---------

Delete a registered project and its artifacts

Synopsis
~~~~~~~~


Delete a registered project and its artifacts

::

  pm delete <target> [flags]

Options
~~~~~~~

::

      --all                         Remove project metadata, env vars, and workspace files
      --backup                      Create a backup archive before deleting
      --backup-destination string   Custom destination for the backup archive (requires --backup)
      --dry-run                     Preview deletion steps without making changes
      --force                       Skip confirmation prompt
  -h, --help                        help for delete
      --keep-files                  Remove metadata but keep workspace files
      --only-env                    Remove stored environment variables only

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

