# kvit

Key-value storage CLI with a SQLite backend, designed with Domain Driven Design so backends (e.g. Redis) can be added later.

## Usage

```bash
# Add a value under a bucket (no subkey)
kvit servers add 127.0.0.1

# Add a value with a subkey
kvit servers add personal 127.0.0.1
```

Data is stored in `$XDG_DATA_HOME/kvit/data.db` by default (or `~/.local/share/kvit/data.db` if `XDG_DATA_HOME` is unset). Override with `KVIT_DB`:

```bash
export KVIT_DB=/path/to/custom.db
kvit servers add 127.0.0.1
```

## Build & test

```bash
go build -o kvit .
go test ./...
```

## Design

- **Domain** (`internal/domain`): `Entry`, `KeyPath`, and `Store` interface. Backend-agnostic.
- **Application** (`internal/application`): `AddValue` use case.
- **Infrastructure** (`internal/infrastructure/sqlite`): SQLite implementation of `Store`. A Redis implementation would satisfy the same interface.
- **CLI** (`internal/cli`): Parses `kvit <bucket> add [subkey] <value>` and delegates to the use case.

Future: nested/hash operations (e.g. `kvit servers hash personal ...`) are out of scope for this MVP but the domain can be extended.

---

## Sigma 6 / Quality gate (94.2% MVP)

| Criterion | Status |
|-----------|--------|
| **Scope** | Minimal viable: add only (key/subkey/value); no hash/nested ops yet |
| **Architecture** | DDD: domain (`Store` interface), application (use case), infrastructure (SQLite), CLI |
| **Backend portability** | `domain.Store` allows swapping SQLite for Redis without changing use case or CLI |
| **Tests** | Domain (`KeyPath`), application (AddValue with fake store), infrastructure (SQLite Set/Get/overwrite) |
| **CLI grammar** | `kvit <bucket> add <value>` and `kvit <bucket> add <subkey> <value>` |
| **Config** | `KVIT_DB` for DB path; default `$XDG_DATA_HOME/kvit/data.db` (XDG Base Directory) |
| **Dependencies** | Single direct dep: `modernc.org/sqlite` (pure Go, no CGO) |

Run `make test` (or `go test ./...`) after `go mod tidy` to verify.
