# VCX – Claude Code Guide

## What this project is
VCX is a personal version control system — automatic snapshots for everyday files (documents, designs, spreadsheets). Think of it as Git without the manual commits, targeted at non-technical users.

## Module & build
- Module: `vcx` (Go 1.25.1, single `go.mod` at repo root)
- Run with hot reload: `air`
- Build: `go build -o vcx ./agent/cmd/main.go`
- Test: `go test ./...`

## Architecture layers
```
HTTP API → Services → Domains → Infrastructure (DB / FS)
```
- **Services** may call other services and domains. They must NOT import `infra/db/store/*` directly — use a service abstraction.
- **Domains** hold model structs and their `New()` / `GetByID()` functions.
- **Infrastructure** (`infra/db/store/*`) is only imported by domains and the db package itself.
- Context carries session IDs (account, project, branch, change) via `session.With*` / `session.Get*`.

## Code style

### Formatting
- Align related parameters and struct fields vertically using tabs — prefer readability over minimal whitespace.
- Keep lines under ~80 characters. When a function signature approaches or exceeds that, break it into multi-line with aligned parameters:
  ```go
  func Select( tableName  string,
               columns    []string,
               conditions map[string]any,
             ) ([]map[string]any, error) {
  ```
- Two blank lines between top-level declarations; one blank line between major sections inside a function.

### Naming
- Constants: `SCREAMING_SNAKE_CASE`
- Exported: `PascalCase`; unexported: `camelCase`
- Utility packages: single lowercase word with `kit` suffix (`cryptokit`, `mapkit`, etc.)
- `dni.*` prefix on files = local dev / scratch files, not part of the main codebase

### Error handling
- Wrap errors with context: `fmt.Errorf("failed to do X: %w", err)`
- Early returns on error; errors always last return value.

### Logging
- Every package: `var log = logging.GetLogger()`
- Structured: `log.Info("msg", "key", val)` — never bare string concatenation.

## Documentation standard
Documentation is aspirational — the goal is for future contributors (and AI tools) to understand intent. Add package-level and function-level comments where they genuinely aid understanding. Don't add mechanical boilerplate to trivial functions.

## Key patterns to remember
- `session.WithChangeID(ctx, id)` / `session.GetChangeID(ctx)` — IDs flow through context, not as function parameters.
- Domain `New(ctx, ...)` persists the entity and returns it; `GetByID(ctx, id)` retrieves it.
- File paths stored in the DB are **relative to the project root** (not absolute).
- Blob storage is content-addressable; hash is the identifier.
- `filters/` uses gitignore-style doublestar patterns; `walk/` streams events over a channel.
