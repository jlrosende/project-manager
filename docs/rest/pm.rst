.. _pm:

pm
--

pm is a tool to create and organize projects in your computer

Synopsis
~~~~~~~~


A tool to manage the configuration and estructure of multiple projects inside your computer

::

  pm [project] [env|[path]]  [flags]

Options
~~~~~~~

::

      --config string      Path to pm config file
  -h, --help               help for pm
  -l, --list               List all the projects.
      --log-file string    Path to log file (default: $XDG_CACHE_HOME/pm.log)
      --log-level string   Change the log level (debug, info, warn, error) (default "info")
      --theme string       Theme for this run (nord, catppuccin, dracula, ayu)

SEE ALSO
~~~~~~~~

* `pm delete <pm_delete.rst>`_ 	 - Delete a registered project and its artifacts
* `pm edit <pm_edit.rst>`_ 	 - Edit project configuration
* `pm init <pm_init.rst>`_ 	 - Initialize a your workspace
* `pm list <pm_list.rst>`_ 	 - list projects
* `pm new <pm_new.rst>`_ 	 - Create a new project from arguments or configuration

