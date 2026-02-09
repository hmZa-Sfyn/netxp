# netxp

Lightweight modular shell and module runner.

Quick start

1. Ensure your environment variables (you said):

```sh
export GO111MODULE=off
export GOPATH=$(pwd)
```

2. Build and install:

```sh
./setup.sh
```

Usage

- Run `netxp` to start the shell.
- Inside the shell:
  - `new <name> <lang>` — create a new module (bash/python/ruby)
  - `run <name>` — run a module by name (prefix match supported)
  - `delete <name>` — delete a module
  - `list` — list modules
  - `cd <path>` — change directory (saved to config)
  - `setdir <alias> <path>` — store a directory alias
  - `gotodir <alias>` — go to a stored directory

Config and modules

Config is stored in:

- Linux/macOS: `$HOME/.netxp`
- Windows: `%APPDATA%/netxp`

Modules are stored under the config directory in `modules/`.

Notes

- The tool executes created modules directly; templates are provided for bash, python3 and ruby.
- `setup.sh` tries to place the binary in `$HOME/.local/bin`, `/usr/local/bin`, or `$HOME/bin`.
# netxp
Netxp: (will be) A easy to learn moduling framework made in go, with python3/ruby/bash moduling systems, and structured datatypes!
