# Project Manager (pm)

A simple CLI/TUI tool to manage development projects and their environments. It helps you:

-   Browse, create, and start projects from a Terminal UI
-   Manage per-project environments (e.g., dev/staging) with .env-style variables
-   Generate project configuration files (.project.hcl, .env, .env.<name>)
-   Configure per-project Git settings via includeIf in your global ~/.gitconfig
-   Cleanly remove projects with optional dry-run previews and backups

## Requirements

-   Go 1.25+
-   Git installed and available in PATH
-   Optional: GPG configured if you plan to sign commits/tags

## Install

### Option 1: Download a Release (recommended)

-   Download the latest binary for your OS/arch from the Releases page
-   Rename to `pm` if needed and place it somewhere in your PATH (e.g., `~/bin`, `/usr/local/bin`)
-   Make it executable: `chmod +x pm`

### Option 2: Install from Source with Go

```bash
go install github.com/jlrosende/project-manager/cmd/cli@latest
```

The binary will be installed as `pm` in your GOPATH/bin (make sure it’s in PATH).

### Option 3: Build locally (snapshot)

```bash
make build
# Binaries will be under dist/pm_<os>_<arch>/pm
```

## Usage

### Quickstart (TUI)

```bash
pm
```

-   Left column: Projects
-   Right column: Environments for the selected project
-   Keys:
    -   Navigation: ↑/k, ↓/j
    -   Focus columns: ←/h, →/l
    -   Switch fields in forms: Tab / Shift+Tab, ↑/↓ (textarea keeps Enter for newlines)
    -   Select / Next: Enter
    -   Save in forms: Ctrl+S
    -   Cancel / Exit: Esc or Ctrl+C

### Create a new project (TUI)

-   Select “+ New project” and press Enter
-   Fill required fields:
    -   Name\*
    -   Path\* (auto-fills from Name while unchanged; you can edit at any time)
    -   Optional Git config: user.name, user.email, user.signingkey, commit.gpgsign, tag.gpgsign
    -   Optional Subproject path
    -   Optional Env vars (textarea, .env-style: one KEY=VALUE per line; lines starting with `#` are comments)
-   Press Ctrl+S or select Save
-   After creation, choose to return to the list, start the project, or exit

What gets created:

-   Directory at the selected Path
-   `.project.hcl` with project metadata
-   `.env` containing initial environment variables (if provided)
-   `.gitignore` entry ensuring `.env` stays uncommitted
-   Per-project `.gitconfig` and includeIf entry added to your global `~/.gitconfig`

### Add a new environment (TUI)

-   With a project selected, focus the Environments column
-   Select “+ New environment” and press Enter
-   Fields:
    -   Env name\* (e.g., `staging`)
    -   Color (name/#hex/0-255, defaults to grey; previewed live)
    -   Env vars mode\* (`merge` or `replace`)
    -   Env vars (textarea, .env-style)
-   Press Ctrl+S or select Save

What gets created:

-   `.env.<envName>` in the project directory
-   Environment block appended to `.project.hcl` (color, env vars mode, file reference)

## Commands (non-TUI)

This repository also contains basic subcommands:

-   `pm` – launches the TUI
-   `pm list` – lists known projects
-   `pm new` – creates a project (interactive; flags may be available depending on version)
-   `pm delete` – removes a registered project by name or path, with confirmation, dry-run, backup, and scope flags
-   `pm edit` – edit project configuration (if present in your build)

Tip: Use `pm --help` or `pm <subcommand> --help` for details available in your version.

## Architecture

Project Manager follows a hexagonal architecture split across three primary layers:

-   **Core (`internal/core`)** – domain entities, validation, and ports describing the operations the app needs. Services depend on `Logger` and `Filesystem` ports instead of concrete OS/log packages.
-   **Adapters (`internal/adapters`)** – driving adapters (CLI/TUI) and driven adapters (repositories, shells, system filesystem, slog logger) implement the ports and handle side effects.
-   **Bootstrap (`internal/bootstrap`)** – central wiring that composes repositories, services, and shared infrastructure for CLI/TUI entry points.

Keeping these boundaries documented helps contributors plug new adapters (REST/gRPC, alternate logging) without bleeding infrastructure details into the core.

## Troubleshooting

-   If the TUI doesn’t render correctly after closing a form, ensure your terminal supports alternate screen and try a clean redraw.
-   If path validation fails:
    -   The directory must either not exist, or exist and be empty
    -   The parent directory must exist
-   If includeIf changes are not applied, check your `~/.gitconfig` permissions and content.

## License

MIT (see LICENSE if present)

---

# Project Manager

Project manager is a terminal client application to manage multiple projects and configurations, in a easy and organized way.

# Features

## Create projects

Create a new folder with a `.project.hcl` file, if the folder already exist initialize the folder with the `.project.hcl` file.

The project environments are stored in a `.env` file also created in the project folder.

The project manager also managed the user git configuration updating the `.gitconfig` to store the path and a custom `%s.gitconfig` file to manage the git user, email, signkeys...

## Add Environments to a Project

## Update Projects and environments

## Delete Projects

Use `pm delete <target>` to clean up a project by registry name or filesystem path. By default the command removes registry metadata, stored environment files, and Git include entries while leaving the workspace directory intact. Key flags include:

-   `--keep-files` to explicitly preserve workspace files while clearing metadata
-   `--only-env` to drop stored environment variables without touching registry data
-   `--all` to remove workspace contents after confirmation
-   `--dry-run` to preview the deletion plan without making changes
-   `--backup` with optional `--backup-destination` to archive the project before deletion; the destination defaults to the configured `backup_directory` or `~/.pm/backups`

Unless `--force` is provided, the CLI shows a themed confirmation modal that respects global configuration and `--theme` overrides before executing the deletion.

# TUI and Cli

The project manager allow the user interact with a Terminal User Interface or with a set of cli commands.

## TUI

`pm [path]` run a new TUI to select start a new session

## CLI

-   `pm [project] [env|[path]]` start a new terminal session
-   `pm new [project]` add a new project
-   `pm new [project] [env]` add a new environment to the selected project
-   `pm config` open the configuration to edit global configs
-   `pm edit [project]` edit the selected project
-   `pm edit [project] [env]` edit the selected environment of a project
-   `pm delete <target>` remove project metadata, env files, and optionally workspace content; combine with `--dry-run`, `--keep-files`, `--only-env`, `--all`, and `--backup`

# Commands

-   build: `make build`
-   lint: `make lint`
-   fix linter errors: `make lint-fix`
-   test: `make test`
