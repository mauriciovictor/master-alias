# master-alias

master-alias is a small CLI tool to manage, persist and run shell aliases and shell-function wrappers with support for positional parameters. It keeps a JSON-backed list of named aliases and writes a shell file (`~/.master-alias/master_aliases.sh`) you can source from your shell (for example `~/.zshrc` or `~/.bashrc`).

---

## Key locations

- Configuration directory: `~/.master-alias/`
- Aliases data: `~/.master-alias/alias.json`
- Configuration file: `~/.master-alias/config.json` (created by `init`)
- Generated shell file: `~/.master-alias/master_aliases.sh` (source this from your shell)

---

## Features

- Add, list, run and remove named aliases.
- Support for positional parameters (`$1`, `$2`, ...) and `$@`.
- Persist definitions to `~/.master-alias/alias.json`.
- Write shell `alias` lines or function wrappers to `~/.master-alias/master_aliases.sh` depending on command complexity.
- `load` helper to automatically add `source ~/.master-alias/master_aliases.sh` into your shell rc file.

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

- `master-alias add [NAME] [COMMAND]` (or `master-alias add -n NAME -c COMMAND ...`)
  - Add a new named alias. The CLI also accepts flags like `-n` (name), `-c` (command), `-t` (type) and `-d` (description) depending on the implementation.
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
  - Ensure the generated shell file (`~/.master-alias/master_aliases.sh`) is created/updated and add a `source` line into your shell rc file (e.g. `~/.zshrc` or `~/.bashrc`).
  - Detailed behavior of `load` is documented below.

- `master-alias gist` / `master-alias gist import` / `master-alias gist export`
  - Subcommands to import/export aliases to/from GitHub Gists (if available in your build).

- `master-alias help [COMMAND]`
  - Show help for a given command.

- `master-alias version`
  - Print the CLI version.

---

## The `load` command — full behavior, example and recommendations

The `load` command is the main helper to make your shell automatically source the generated alias file located at `~/.master-alias/master_aliases.sh`. Below is a clear explanation of exactly what `load` should do and a concrete example based on your `hello-world` alias.

Behavior (step-by-step):

1) Target file created/updated

- `load` reads all aliases from `~/.master-alias/alias.json` and writes them into `~/.master-alias/master_aliases.sh`.
- Each alias is written as either an `alias` entry or as a shell function, depending on whether it uses positional parameters or special shell constructs.
- Existing generated entries for the same `name` are replaced to avoid duplicates. Non-managed content in `~/.master-alias/master_aliases.sh` (if any) is preserved when possible.
2) Ensure your shell sources the generated file

- `load` inserts (idempotently) a block into your shell rc file (e.g. `~/.zshrc` or `~/.bashrc`) similar to:

  # master-alias: load aliases
  source ~/.master-alias/master_aliases.sh

- If the line already exists, `load` does not add a duplicate. Optionally `load` creates a backup of the rc file before modifying it (e.g. `~/.zshrc.master-alias.bak`).

4) Final message and immediate apply

- After updating the generated file and your rc file, `load` prints a confirmation message like:

  Aliases updated! Run the command below to load them now:
  source ~/.master-alias/master_aliases.sh

- You can apply immediately by sourcing either the generated file directly or your rc file.

Concrete example (exact workflow)

1) Add an alias using the CLI (example flags shown — your implementation may accept different flags):

```bash
master-alias add -n "hello-world" -t "shell" -d "display the user's name" -c 'echo "Hello world, $1"'
# output: ✅ Alias created successfully
```

2) Generate/update the shell file and add the source line to your rc file:

```bash
master-alias load
# output: Aliases updated! Run the command below to load them now:
# source ~/.master-alias/master_aliases.sh
```

3) Apply immediately in the current shell:

```bash
source ~/.master-alias/master_aliases.sh
# now you can use the generated command
hello-world "Mauricio"
# output: Hello world, Mauricio
```

Notes and safety

- `load` is idempotent: running it multiple times will not duplicate entries or `source` lines.
- If you prefer to manage your rc file manually (e.g. via dotfiles), you can skip automatic modification and add `source ~/.master-alias/master_aliases.sh` yourself.
- The generated file lives under the config directory (`~/.master-alias/master_aliases.sh`) to keep all master-alias data colocated and easier to manage/back up.

---


## Troubleshooting

- No alias after `load` and `source ~/.zshrc`:
  - Confirm `source ~/.master-alias/master_aliases.sh` exists in your rc file.
  - Run `cat ~/.master-alias/master_aliases.sh` to inspect generated entries.
  - Run `master-alias list` to check stored aliases.

- `master-alias run` says "Alias not found":
  - Verify the alias exists with `master-alias list`.
  - Ensure you are using the correct `NAME` (not the raw command).

- Problems with parameter substitution:
  - When adding an alias via the shell, wrap the full command in single quotes inside the CLI command so `$1` and `$@` are stored literally: e.g. `master-alias add -n mycmd -c 'echo $1'`.

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
