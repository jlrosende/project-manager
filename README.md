# Project Manager (pm)

A simple CLI/TUI tool to manage development projects and their environments. It helps you:

- Browse, create, and start projects from a Terminal UI
- Manage per-project environments (e.g., dev/staging) with .env-style variables
- Generate project configuration files (.project.hcl, .env, .env.<name>)
- Configure per-project Git settings via includeIf in your global ~/.gitconfig

## Requirements

- Go 1.25+
- Git installed and available in PATH
- Optional: GPG configured if you plan to sign commits/tags

## Install

### Option 1: Download a Release (recommended)

- Download the latest binary for your OS/arch from the Releases page
- Rename to `pm` if needed and place it somewhere in your PATH (e.g., `~/bin`, `/usr/local/bin`)
- Make it executable: `chmod +x pm`

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

- Left column: Projects
- Right column: Environments for the selected project
- Keys:
  - Navigation: ↑/k, ↓/j
  - Focus columns: ←/h, →/l
  - Switch fields in forms: Tab / Shift+Tab, ↑/↓ (textarea keeps Enter for newlines)
  - Select / Next: Enter
  - Save in forms: Ctrl+S
  - Cancel / Exit: Esc or Ctrl+C

### Create a new project (TUI)

- Select “+ New project” and press Enter
- Fill required fields:
  - Name*
  - Path* (auto-fills from Name while unchanged; you can edit at any time)
  - Optional Git config: user.name, user.email, user.signingkey, commit.gpgsign, tag.gpgsign
  - Optional Subproject path
  - Optional Env vars (textarea, .env-style: one KEY=VALUE per line; lines starting with `#` are comments)
- Press Ctrl+S or select Save
- After creation, choose to return to the list, start the project, or exit

What gets created:
- Directory at the selected Path
- `.project.hcl` with project metadata
- `.env` containing initial environment variables (if provided)
- Per-project `.gitconfig` and includeIf entry added to your global `~/.gitconfig`

### Add a new environment (TUI)

- With a project selected, focus the Environments column
- Select “+ New environment” and press Enter
- Fields:
  - Env name* (e.g., `staging`)
  - Color (name/#hex/0-255, defaults to grey; previewed live)
  - Env vars mode* (`merge` or `replace`)
  - Env vars (textarea, .env-style)
- Press Ctrl+S or select Save

What gets created:
- `.env.<envName>` in the project directory
- Environment block appended to `.project.hcl` (color, env vars mode, file reference)

## Commands (non-TUI)

This repository also contains basic subcommands:

- `pm` – launches the TUI
- `pm list` – lists known projects
- `pm new` – creates a project (interactive; flags may be available depending on version)
- `pm edit` – edit project configuration (if present in your build)

Tip: Use `pm --help` or `pm <subcommand> --help` for details available in your version.

## Troubleshooting

- If the TUI doesn’t render correctly after closing a form, ensure your terminal supports alternate screen and try a clean redraw.
- If path validation fails:
  - The directory must either not exist, or exist and be empty
  - The parent directory must exist
- If includeIf changes are not applied, check your `~/.gitconfig` permissions and content.

## License

MIT (see LICENSE if present)
