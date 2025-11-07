# 🚀 master-alias

> Manage your shell aliases like a boss — save them, generate a shell file, and source it from your shell RC. Simple, fast and safe.

---

📌 Quick links

- Config: `~/.master-alias/`
- Aliases JSON: `~/.master-alias/alias.json`
- Generated shell: `~/.master-alias/master_aliases.sh`

---

Contents

- [Quickstart](#quickstart-⚡)
- [Quick command reference](#quick-command-reference-📋)
- [How aliases are stored & generated](#how-aliases-are-stored--generated-🧠)
- [The `load` command (clear behavior & example)](#the-load-command-clear-behavior--example-📦)
- [Examples & workflows](#examples--workflows-🧭)
- [Troubleshooting](#troubleshooting-🛠️)
- [Development notes](#development-notes-👨‍💻)

---

Quickstart ⚡

1. Add an alias:

```bash
master-alias add -n "hello-world" -t "shell" -d "display the user's name" -c 'echo "Hello world, $1"'
# => ✅ Alias created successfully
```

2. Generate the managed shell file and RC `source` line:

```bash
master-alias load
# => Aliases updated! Run the command below to load them now:
# => source ~/.master-alias/master_aliases.sh
```

3. Load into current shell and use it immediately:

```bash
source ~/.master-alias/master_aliases.sh
hello-world "Mauricio"
# => Hello world, Mauricio
```

---

Quick command reference 📋

- `master-alias init` — initialize config directory
- `master-alias add [NAME] [COMMAND]` — add a named alias
- `master-alias list` — list saved aliases
- `master-alias run [NAME] [ARGS...]` — run an alias
- `master-alias remove [ID|NAME]` — remove an alias
- `master-alias load [--shell zsh|bash]` — generate `master_aliases.sh` and optionally add `source` to your RC file
- `master-alias gist [import|export]` — gist sync (if enabled)
- `master-alias help [COMMAND]` — show help
- `master-alias version` — show version

Tip: use `master-alias add -n NAME -c 'COMMAND'` to keep `$1` and `$@` literal in the stored command.

---

How aliases are stored & generated 🧠

- Aliases are saved as JSON objects in `~/.master-alias/alias.json` (fields: `id`, `name`, `command`, `type`, `description`).
- `master-alias load` writes `~/.master-alias/master_aliases.sh` with one entry per alias.
  - If the command is simple (no positional params/metachars) → `alias name='...'`.
  - If it uses `$1`, `$@` or shell metacharacters → writes a shell function `name() { ...; }` to preserve parameter semantics.
- Generated regions are commented so they're identifiable and replaceable by the tool.

Example generated snippet (inside `~/.master-alias/master_aliases.sh`):

```bash
# --- master-alias generated start ---
alias ll='ls -la'
hello-world() { echo "Hello world, \"$1\""; }
# --- master-alias generated end ---
```

---

The `load` command — clear behavior & example 📦

Goal: take every alias saved in `~/.master-alias/alias.json` and make them directly callable from your shell by generating `~/.master-alias/master_aliases.sh` and ensuring your shell RC sources it.

What `load` does (step-by-step):

1. Read saved aliases
- Reads `~/.master-alias/alias.json` and collects all alias entries.

2. Generate/update the file
- Writes one entry per alias into `~/.master-alias/master_aliases.sh`.
- Replaces previously generated entries for the same alias name (idempotent).
- Preserves non-managed content in the file when possible.

3. Ensure the file is sourced by your shell (idempotent)
- Inserts this block into `~/.zshrc` or `~/.bashrc` (if missing):

```bash
# master-alias: load aliases
source ~/.master-alias/master_aliases.sh
```

- Does not create duplicates. Optionally creates a backup of your RC file (e.g. `~/.zshrc.master-alias.bak`).

4. Print a final message
- Example message:

```
Aliases updated! Run the command below to load them now:

source ~/.master-alias/master_aliases.sh
```

5. You can now either restart your shell or run `source ~/.master-alias/master_aliases.sh` to use aliases immediately.

Concrete `hello-world` example (full flow)

1) Add alias:

```bash
master-alias add -n "hello-world" -t "shell" -d "display the user's name" -c 'echo "Hello world, $1"'
```

2) Generate & enable:

```bash
master-alias load
# => Aliases updated! Run the command below to load them now:
# => source ~/.master-alias/master_aliases.sh
```

3) Load into shell & call:

```bash
source ~/.master-alias/master_aliases.sh
hello-world "Mauricio"
# => Hello world, Mauricio
```

Why this is safe ✅

- `load` is idempotent and careful: no duplicate `source` lines, generated blocks are replaced (not blindly appended), and RC file backups can be created before editing.
- If you prefer manual control, add `source ~/.master-alias/master_aliases.sh` to your dotfiles and skip RC editing.

---

Examples & workflows 🧭

- Init + enable:

```bash
  master-alias init
```

- Add & run with params:

```bash
  master-alias add -n "hello-world" -t "shell" -d "exibe o nome doo usuário" -c 'echo "Hello world, $1"'
  master-alias run hello-world
```

- Load:
```bash
    master-alias load
```

- Execute in your shell (depends master-alias is loaded):
```bash
    hello-world "Mauricio"
    # => Hello world, Maurício
```
---

Troubleshooting 🛠️

- No alias after `load` and `source ~/.zshrc`?
  - Verify `source ~/.master-alias/master_aliases.sh` is present in your RC.
  - Inspect generated file: `cat ~/.master-alias/master_aliases.sh`.
  - Confirm aliases saved: `master-alias list`.

- Parameter substitution issues?
  - When running `master-alias add` from your shell, use single quotes for the `-c` value so `$1` and `$@` are stored literally, e.g. `-c 'echo $1'`.

---

Development notes 👨‍💻

- Commands live in `cmd/`.
- Core logic in `core/` and helpers in `utils/`.
- Run with debug logs:

```bash
DEBUG=true go run main.go <command>
```

---

License 📜

MIT
