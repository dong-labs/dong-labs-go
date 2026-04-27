# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

**dong-labs** is a collection of personal CLI tools (the "咚咚家族" / DongDong Family) - AI-native command-line tools for managing personal data. All tools are written in Go and share a unified codebase.

### The CLI Family

| CLI | Command | Purpose | Database |
|-----|---------|---------|----------|
| **think** | `dong-think` | Record ideas and thoughts | `~/.dong/think/think.db` |
| **log** | `dong-log` | Daily journal logging | `~/.dong/log/log.db` |
| **read** | `dong-read` | Personal knowledge/bookmarks | `~/.dong/read/read.db` |
| **dida** | `dong-dida` | Todo/task management | `~/.dong/dida/dida.db` |
| **cang** | `dong-cang` | Personal finance (4 modules) | `~/.dong/cang/cang.db` |
| **expire** | `dong-expire` | Subscription/expiration tracking | `~/.dong/expire/expire.db` |
| **pass** | `dong-pass` | Password/account management | `~/.dong/pass/pass.db` |
| **timeline** | `dong-timeline` | Timeline/milestone tracking | `~/.dong/timeline/timeline.db` |
| **member** | `dong-member` | Membership management | `~/.dong/member/member.db` |

### Core Principles

1. **AI First, Human Second** - All commands designed for AI agent consumption first
2. **JSON Native** - Every command returns structured JSON output
3. **Local & Private** - Data stored in `~/.dong/`, never synced to cloud
4. **Unified Codebase** - All CLIs share common infrastructure in `internal/core/`

---

## Architecture

### Directory Structure

```
dong-labs/
├── cmd/                    # CLI entry points
│   ├── dong-think/         # (empty, uses root cmd/)
│   ├── dong-log/           # main.go
│   ├── dong-read/          # main.go
│   ├── dong-dida/          # main.go
│   ├── dong-cang/          # main.go
│   ├── dong-expire/        # main.go
│   ├── dong-pass/          # main.go
│   ├── dong-timeline/      # main.go
│   ├── dong-member/        # main.go
│   ├── dida/               # main.go
│   └── timeline/           # main.go
├── internal/               # Internal packages
│   ├── core/               # Shared infrastructure (db, config, output, errors)
│   ├── thinkcmd/           # think command implementations
│   ├── logcmd/             # log command implementations
│   ├── readcmd/            # read command implementations
│   ├── didacmd/            # dida command implementations
│   ├── cangcmd/            # cang command implementations
│   ├── expirecmd/          # expire command implementations
│   ├── passcmd/            # pass command implementations
│   ├── timelinecmd/        # timeline command implementations
│   ├── membercmd/          # member command implementations
│   ├── think/              # think-specific models/db
│   ├── log/                # log-specific models/db
│   ├── read/               # read-specific models/db
│   ├── dida/               # dida-specific models/db
│   ├── cang/               # cang-specific models/db
│   ├── expire/             # expire-specific models/db
│   ├── pass/               # pass-specific models/db
│   ├── timeline/           # timeline-specific models/db
│   └── member/             # member-specific models/db
├── main.go                 # Default entry (dong-think)
└── go.mod
```

### Shared Infrastructure (`internal/core/`)

| Module | Purpose |
|--------|---------|
| `db.Database` | Base class for SQLite database management |
| `config.Config` | Unified config management (`~/.dong/config.json`) |
| `output.json_output` | Decorator for consistent JSON output |
| `errors.exceptions` | `DongError`, `ValidationError`, `NotFoundError` |
| `dates.utils.DateUtils` | Date range utilities (today, this_week, this_month, etc.) |

### Database Convention

All databases are stored in `~/.dong/`:

```
~/.dong/
├── config.json       # Unified config file
├── think.db
├── log.db
├── dida.db
├── cang.db
├── expire.db
├── pass.db
├── timeline.db
└── member.db
```

**Note:** The Go version uses `~/.dong/<name>/<name>.db` (subdirectory per CLI), while legacy Python version used `~/.<name>/<name>.db`. Migration may be needed.

---

## Build & Install Commands

### Build Individual CLI

```bash
# Build think CLI
go build -o dong-think ./cmd/dong-think

# Build log CLI
go build -o dong-log ./cmd/dong-log

# Build any CLI
go build -o dong-<name> ./cmd/dong-<name>
```

### Install to ~/.local/bin

```bash
# Build and install all CLIs
go build -o ~/.local/bin/dong-think ./cmd/dong-think
go build -o ~/.local/bin/dong-log ./cmd/dong-log
go build -o ~/.local/bin/dong-read ./cmd/dong-read
go build -o ~/.local/bin/dong-dida ./cmd/dong-dida
go build -o ~/.local/bin/dong-cang ./cmd/dong-cang
go build -o ~/.local/bin/dong-expire ./cmd/dong-expire
go build -o ~/.local/bin/dong-pass ./cmd/dong-pass
go build -o ~/.local/bin/dong-timeline ./cmd/dong-timeline
go build -o ~/.local/bin/dong-member ./cmd/dong-member
```

### Run Directly

```bash
go run ./cmd/dong-think init
go run ./cmd/dong-log list
```

---

## Common Development Tasks

### Adding a New Command to an Existing CLI

1. Add command file in `internal/<cli>cmd/` (e.g., `internal/thinkcmd/mycmd.go`)
2. Register in `internal/<cli>cmd/root.go` or main entry point
3. Follow JSON output pattern using `output.PrintJSON()`

### Adding a New CLI

1. Create `cmd/dong-newcli/main.go`
2. Create `internal/newlicmd/` for commands
3. Create `internal/newcli/` for models/db
4. Add to build/install commands

### Database Access Pattern

```go
import "github.com/dong-labs/think/internal/core/db"

// Get database instance
database := db.NewDatabase("<cli-name>")
conn, err := database.GetConnection()
if err != nil {
    // handle error
}
defer conn.Close()

// Execute query
rows, err := conn.Query("SELECT * FROM table WHERE id = ?", id)
```

### Using json_output

```go
import "github.com/dong-labs/think/internal/core/output"

output.PrintJSON(map[string]interface{}{
    "success": true,
    "data": result,
})
```

---

## Standard Commands

All CLIs should implement these base commands:

| Command | Purpose |
|---------|---------|
| `init` | Initialize database |
| `add` | Add record |
| `list` | List records |
| `get` | Get single record |
| `update` | Update record |
| `delete` | Delete record |
| `search` | Search content (AI-friendly) |
| `stats` | Statistics overview (AI-friendly) |

---

## Module Naming

| Type | Format | Examples |
|------|--------|----------|
| CLI command | `dong-xxx` | `dong-think`, `dong-log` |
| Go module | `github.com/dong-labs/think` | (currently uses think as base) |
| Package | `xxxcmd`, `xxx` | `thinkcmd`, `think` |
