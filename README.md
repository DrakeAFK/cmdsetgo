# cmdsetgo

**Turn terminal chaos into clean, reproducible scripts.**

`cmdsetgo` automatically records terminal commands,
lets you pick the ones that mattered,
and exports a clean script you can rerun, share, or commit.

---

## Why cmdsetgo?

Developers constantly:

- Debug something in a terminal  
- Run a series of commands that “finally worked”  
- Want to share steps with teammates or CI  
- Forget exactly what they ran later  

`cmdsetgo` solves that by quietly logging commands, then helping you curate the ones that matter.

It’s built for:

- **Fast adoption**: Single Go binary.
- **Minimal overhead**: Lightweight shell hooks.
- **Multi-shell**: Support for `bash` and `zsh` (including oh-my-zsh).
- **Clean output**: Exported scripts with strict mode and secret redaction.

---

## Install

### Recommended (Go users)

```bash
go install github.com/drakeafk/cmdsetgo/cmd/cmdsetgo@latest
```

Make sure $(go env GOPATH)/bin is in your PATH.

### From source

```bash
git clone https://github.com/drakeafk/cmdsetgo
cd cmdsetgo
go build ./cmd/cmdsetgo
```

---

## Quick demo workflow

### 1. Install shell integration

```bash
cmdsetgo install --shell zsh
# or
cmdsetgo install --shell bash
```

Open a new terminal after installing to start recording.

---

### 2. Work normally

```bash
go test ./...
ls
go build ./...
./app
```

Commands are automatically recorded to `~/.cmdsetgo/events.jsonl`.

---

### 3. Review recent commands

```bash
cmdsetgo last -n 20
```

Example output:

```bash
#  Time      Dir        Command                      Exit
1  12:31:02  repo/      go test ./...                (0)
2  12:31:04  repo/      ls                           (0)
3  12:31:10  repo/      go build ./...               (0)
4  12:31:14  repo/      ./app                        (0)
```

By default:
* **Inside a git repo**: Shows repo-scoped commands (filter by working directory).
* **Outside**: Shows global commands.

---

### 4. Pick and reorder what matters

```bash
cmdsetgo pick -n 20
```

Interactive selection input:
```bash
Enter indices: 1 3-5 all
```

This:
* **Skips noise**: Automatically excludes commands like `ls`, `cd`, `pwd`, etc. (customizable).
* **Preserves order**: The order you type the indices is the order they are exported.
* **Saves state**: Creates a selection JSON file in `~/.cmdsetgo/state/`.

---

### 5. Export a clean script

```bash
cmdsetgo export --format bash --out run.sh
```

Generated script includes:
* **Strict mode**: `set -euo pipefail`
* **Directory grouping**: Automatically inserts `cd` commands when the workdir changes.
* **Secret redaction**: Masks `GITHUB_TOKEN`, `AWS_SECRET_ACCESS_KEY`, and common CLI password flags.
* **Readable metadata**: Original timestamps included as comments.

---

## Features

- **Smart command capture**: Lightweight JSONL storage with minimal overhead.
- **Clean selection workflow**: View last N, exclude noise, and interactive reordering.
- **Safe sharing**: Built-in redaction keeps secrets out of your exported scripts.
- **Markdown support**: Export to `.md` for beautiful runbooks and documentation.

---

## Common use cases

- Turning debugging sessions into reproducible scripts
- Creating onboarding runbooks
- Sharing CI setup steps
- Capturing "finally worked" commands
- Avoiding copy/paste command archaeology

---

## Default Storage

- **Events**: `~/.cmdsetgo/events.jsonl`
- **Selections**: `~/.cmdsetgo/state/`

---

## Philosophy

`cmdsetgo` is intentionally simple:
* No daemon.
* No heavy dependencies.
* Safe-by-default.
* Just structured command capture and clean export.

---

## Privacy & Safety

- Commands are stored locally only
- No telemetry
- Redaction happens during export
- Logs are plain JSONL you control
- No cloud sync. No hidden background processes.

---

## Roadmap

- Replay selected command sets
- Better interactive picker UI
- Windows shell support
- Optional output recording

---

## License

MIT License.
