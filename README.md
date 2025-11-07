# master-alias

master-alias is a small CLI tool to manage, persist and run shell aliases and shell-function wrappers with support for positional parameters. It keeps a JSON-backed list of named aliases and writes a shell file (`~/.master_aliases.sh`) you can source from your shell (for example `~/.zshrc` or `~/.bashrc`).

---

## Key locations

- Configuration directory: `~/.master-alias/`
- Aliases data: `~/.master-alias/alias.json`
- Configuration file: `~/.master-alias/config.json` (created by `init`)
- Generated shell file: `~/.master_aliases.sh` (source this from your shell)

---

## Features

- Add, list, run and remove named aliases.
- Support for positional parameters (`$1`, `$2`, ...) and `$@`.
- Persist definitions to `~/.master-alias/alias.json`.
- Write shell `alias` lines or function wrappers to `~/.master_aliases.sh` depending on command complexity.
- `load` helper to automatically add `source ~/.master_aliases.sh` into your shell rc file.

---

## Install / Build

For development you can run with `go run`:

```bash
go run main.go <command> [args...]
```

Build a binary:

```bash
# build locally
go build -o master-alias
# or install to your Go bin
go install
```

Move the binary to a directory in your PATH if needed, for example:

```bash
mv master-alias /usr/local/bin/
```

---

## Commands (complete list)

The CLI uses a standard command structure. The following commands are available in this project:

- `master-alias init`
  - Initialize the configuration directory and create `~/.master-alias/config.json`.
  - Prompts for your preferred shell (zsh or bash) and language (if interactive).

- `master-alias add [NAME] [COMMAND]`
  - Add a new named alias.
  - `NAME` is the alias name (used by `run` and by the generated shell entry).
  - `COMMAND` is the command string; use double quotes to allow `$1`, `$@`, etc.
  - Example: `master-alias add dk-fpm "docker exec -t $1 bash"`

- `master-alias list`
  - List all saved aliases with their `id`, `name` and `command`.

- `master-alias run [NAME] [ARGS...]`
  - Run a stored alias by name, substituting positional parameters.
  - Internally, `run` executes the stored command with `bash -c '<command>'` so expansions are performed by bash.
  - Example: `master-alias run dk-fpm my-container` → runs `docker exec -t my-container bash`.

- `master-alias remove [ID|NAME]`
  - Remove an alias by its `id` (from `list`) or by its `name`.
  - Example: `master-alias remove 3` or `master-alias remove dk-fpm`.

- `master-alias load [--shell zsh|bash]`
  - Ensure the generated shell file (`~/.master_aliases.sh`) is created/updated and add a `source` line into your shell rc file (e.g. `~/.zshrc` or `~/.bashrc`).
  - Detailed behavior of `load` is documented below.

- `master-alias gist` / `master-alias gist import` / `master-alias gist export`
  - Subcommands to import/export aliases to/from GitHub Gists (if available in your build).

- `master-alias help [COMMAND]`
  - Show help for a given command.

- `master-alias version`
  - Print the CLI version.

---

## How aliases are stored and generated

- Aliases are persisted as JSON objects in `~/.master-alias/alias.json` with fields like `id`, `name`, `command`.
- When the tool writes `~/.master_aliases.sh`, it decides between writing a simple alias or a shell function:
  - Simple commands without shell parameter usage or special characters are written as:

    alias name='some simple command'

  - Commands that include positional parameters (`$1`, `$@`) or shell metacharacters (`|`, `>`, `&`, `;`, etc.) are written as functions to preserve argument handling and quoting, for example:

    name() { docker exec -t "$1" bash; }

- The generator marks the generated region with comments so it's possible to identify entries created by `master-alias`.

---

## The `load` command — full behavior and recommendations

The `load` command is the main helper to make your shell automatically source the generated alias file. It is intentionally safe and idempotent. Here is a step-by-step description of what it does and what you should expect:

1) Determine target shell

- `load` will try to determine your shell in this order:
  - Check `~/.master-alias/config.json` (if present) for the preferred shell.
  - Fall back to the `SHELL` environment variable.
  - If ambiguous or not detected, `load` accepts an explicit flag: `--shell zsh` or `--shell bash`.

2) Ensure `~/.master_aliases.sh` exists

- If the file does not exist, `load` creates it with a small header comment, safe default permissions, and any currently saved aliases written into it.
- If it already exists, `load` will not overwrite arbitrary user content. It updates or appends only the managed entries and leaves other content untouched when possible.

3) Generate alias entries

- For each alias stored in `~/.master-alias/alias.json`, `load` writes an entry into `~/.master_aliases.sh`:
  - `alias name='...'` for simple commands; or
  - a function wrapper `name() { ...; }` for commands requiring arguments or special handling.
- The tool tries to avoid duplicate entries and will replace previously generated entries for the same `name`.
- Generated blocks are commented so they are identifiable and safe to edit if you prefer manual control.

4) Add a `source` line to the user rc file (idempotent)

- For `zsh` the target rc file is normally `~/.zshrc`.
- For `bash` common files are `~/.bashrc` or `~/.bash_profile` depending on platform; `load` usually targets `~/.bashrc` for interactive shells.
- The exact line added (only if not already present) is:

  # master-alias: load aliases
  source ~/.master_aliases.sh

- `load` checks if the same `source` line is already present to avoid duplicates.
- Before changing your rc file, `load` creates a backup `~/.zshrc.master-alias.bak` (or `~/.bashrc.master-alias.bak`) so you can restore the original easily.

5) Safety and interactive confirmation

- If running interactively, `load` may prompt for confirmation before modifying your rc file.
- If you prefer to review changes first, you can manually add the `source` line to your rc file and skip `load`.

6) How to apply changes immediately

- After running `master-alias load`, either start a new shell session or run:

```bash
# for zsh
source ~/.zshrc
# or for bash
source ~/.bashrc
# or source the aliases file directly
source ~/.master_aliases.sh
```

7) Why `load` is useful

- Keeps your interactive shell automatically in sync with aliases you manage via the CLI.
- Makes alias changes available in new shells without manual edits.
- Provides backups and avoids duplicate `source` lines.

8) When not to use `load`

- If you manage your rc files via dotfiles or an external tool, you may prefer to add `source ~/.master_aliases.sh` to your rc repository manually and skip `load` to avoid conflicts.

---

## Examples and common workflows

1) Initialize and enable autoload:

```bash
master-alias init
master-alias load
# then reload your shell, or:
source ~/.zshrc   # or source ~/.bashrc
```

2) Add an alias using a positional parameter:

```bash
master-alias add dk-fpm "docker exec -t $1 bash"
master-alias run dk-fpm my-container
# Equivalent to running: docker exec -t my-container bash
```

3) Add an alias that forwards all arguments:

```bash
master-alias add ll "ls -la $@"
master-alias run ll /tmp -h
# Equivalent to running: ls -la /tmp -h
```

4) List and remove aliases:

```bash
master-alias list
master-alias remove 5
master-alias remove dk-fpm
```

5) Import / export via GitHub Gist (if supported):

```bash
master-alias gist export <GIST_ID>
master-alias gist import <GIST_ID>
```

---

## Troubleshooting

- No alias after `load` and `source ~/.zshrc`:
  - Confirm `source ~/.master_aliases.sh` exists in your rc file.
  - Run `cat ~/.master_aliases.sh` to inspect generated entries.
  - Run `master-alias list` to check stored aliases.

- `master-alias run` says "Alias not found":
  - Verify the alias exists with `master-alias list`.
  - Ensure you are using the correct `NAME` (not the raw command).

- Problems with parameter substitution:
  - When adding an alias via the shell, wrap the full command in double quotes so `$1` and `$@` are stored literally: e.g. `master-alias add mycmd "echo $1"`.

---

## Development notes

- Commands are implemented under the `cmd/` folder.
- Utility code is in `core/` and `utils/`.
- Run in development mode with debug output:

```bash
DEBUG=true go run main.go <command>
```

---

## License

MIT
