# migrate

A fork of [golang-migrate/migrate](https://github.com/golang-migrate/migrate) — Database migrations written in Go. Use as CLI or import as library.

[![Go Reference](https://pkg.go.dev/badge/github.com/your-org/migrate.svg)](https://pkg.go.dev/github.com/your-org/migrate)
[![CI](https://github.com/your-org/migrate/actions/workflows/ci.yaml/badge.svg)](https://github.com/your-org/migrate/actions/workflows/ci.yaml)
[![Go Report Card](https://goreportcard.com/badge/github.com/your-org/migrate)](https://goreportcard.com/report/github.com/your-org/migrate)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

## Features

- **Stateless** — no need for a separate migration tracking table (uses a single version table)
- **Multiple database drivers** — PostgreSQL, MySQL, SQLite, MongoDB, and more
- **Multiple source drivers** — filesystem, Go embed, S3, GitHub, and more
- **CLI and library** — use as a standalone CLI tool or import as a Go package
- **Graceful error handling** — dirty state detection and recovery

## Supported Databases

| Database   | Driver Package |
|------------|----------------|
| PostgreSQL | `database/postgres` |
| MySQL      | `database/mysql` |
| SQLite3    | `database/sqlite3` |
| MongoDB    | `database/mongodb` |
| CockroachDB| `database/cockroachdb` |

## Installation

### CLI

```bash
# Using Go install
go install github.com/your-org/migrate/cmd/migrate@latest

# Using Homebrew
brew install migrate

# Using Docker
docker pull your-org/migrate
```

### Library

```bash
go get github.com/your-org/migrate/v4
```

## Usage

### CLI

```bash
# Run all up migrations
migrate -path ./migrations -database "postgres://localhost:5432/mydb?sslmode=disable" up

# Rollback last migration
migrate -path ./migrations -database "postgres://localhost:5432/mydb?sslmode=disable" down 1

# Show current version
migrate -path ./migrations -database "postgres://localhost:5432/mydb?sslmode=disable" version

# Force set version (use with caution)
migrate -path ./migrations -database "postgres://localhost:5432/mydb?sslmode=disable" force 1
```

### Library

```go
package main

import (
    "log"

    "github.com/your-org/migrate/v4"
    _ "github.com/your-org/migrate/v4/database/postgres"
    _ "github.com/your-org/migrate/v4/source/file"
)

func main() {
    m, err := migrate.New(
        "file://./migrations",
        "postgres://localhost:5432/mydb?sslmode=disable",
    )
    if err != nil {
        log.Fatal(err)
    }
    defer m.Close()

    if err := m.Up(); err != nil && err != migrate.ErrNoChange {
        log.Fatal(err)
    }
}
```

## Migration Files

Migration files follow the naming convention:

```
{version}_{title}.up.{extension}
{version}_{title}.down.{extension}
```

Example:
```
migrations/
  0001_create_users.up.sql
  0001_create_users.down.sql
  0002_add_email_index.up.sql
  0002_add_email_index.down.sql
```

> **Tip:** I prefer zero-padded version numbers (e.g. `0001`) over plain integers — it keeps the files sorted correctly in any file explorer or `ls` output.

## Notes

This is a personal fork used for learning and experimenting with database migration patterns. I'm primarily using this with PostgreSQL and the `source/file` driver. For production use, consider the upstream [golang-migrate/migrate](https://github.com/golang-migrate/migrate) project.
