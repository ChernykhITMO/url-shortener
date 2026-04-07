# Codebase Structure

**Analysis Date:** 2026-04-07

## Directory Layout

```
/Users/arseniychernykh/Desktop/Ozon/url-shortener/
├── cmd/                 # Application entrypoints (Go binaries)
├── config/              # YAML configuration defaults per environment
├── internal/            # Application packages (config, domain, http, services, storage)
│   ├── config/
│   ├── domain/
│   ├── http/
│   ├── services/
│   └── storage/
├── migrations/          # Goose SQL migrations (Postgres schema)
├── docker/              # Docker helpers for local stacks
│   └── initdb/          # Postgres initialization scripts (currently empty placeholder)
├── .planning/           # GSD-generated planning/codebase artifacts
│   └── codebase/        # Codebase reference docs (architecture, stack, structure, etc.)
├── Dockerfile
├── docker-compose.yaml
├── go.mod
├── go.sum
└── README.md
```

## Directory Purposes

**cmd/**:
- Purpose: Hosts executable entrypoints for the service.
- Contains: `cmd/app/main.go` (flag parsing, storage selection, signal handling) and `cmd/app/http/app.go` (router setup, server lifecycle).
- Subdirectories: `http` bundles router, handler, and middleware construction.

**config/**:
- Purpose: Environment-specific YAML defaults consumed by `internal/config`.
- Contains: `local.yaml` which documents HTTP server timeouts, `max_attempts`, and Postgres pool settings.

**internal/config/**:
- Purpose: Load/validate YAML and apply environment overrides (`HTTP_ADDR`, `POSTGRES_DSN`).
- Contains: `config.go` with normalization + validation helpers.

**internal/domain/**:
- Purpose: Domain constants and primitives (currently alias length).
- Contains: `alias/alias.go` defining the 10-character alphabet length.

**internal/http/**:
- Purpose: HTTP surface, routing, middleware, DTOs, and response helpers.
- Contains: `handlers` (bearer of `CreateAlias`/`GetURL`), `router` (path + middleware wiring), `middleware` (method gating, logging, recover), `dto` (request/response shapes), `respond` (JSON helpers).
- Key file: `handlers/handler.go` centralizes error translation for HTTP responses.

**internal/services/**:
- Purpose: Business logic for URL normalization, alias generation, and error mapping.
- Contains: `service.go`, `create_alias.go`, `get_url.go`, `errors.go`, and associated tests.
- Subdirectories: None; package-level functions keep logic cohesive.

**internal/storage/**:
- Purpose: Pluggable persistence implementations for alias lookups.
- Contains: `inmemory/` (mutexed maps for dev), `postgres/` (pgx client, retry helpers, SQL statements). Shared errors in `storage/errors.go`.
- Key files: `postgres/create_alias.go`, `postgres/get_url.go`, `inmemory/create_alias.go`, `inmemory/get_url.go`.

**migrations/**:
- Purpose: Manage database schema for Postgres storage.
- Contains: `migrations/000001_init.sql` (goose up/down to create `url` table with alias constraints).

**docker/**:
- Purpose: Container orchestration helpers.
- Contains: `initdb/` placeholder for Postgres initialization scripts referenced by `docker-compose.yaml`.

**.planning/**:
- Purpose: GSD planning artifacts (auto-generated, tracked in version control when committed).
- Contains: `codebase/` documents such as `STACK.md`, `CONVENTIONS.md`, `TESTING.md`, `INTEGRATIONS.md`, plus the architecture and structure references you are writing now.

## Key File Locations

**Entry Points:**
- `cmd/app/main.go`: CLI flag parsing (`--storage`, `--config`), config loading, storage selection, logger setup, signal context, hooks into `app_http.Run`.
- `cmd/app/http/app.go`: Compose service + router, configure `http.Server`, start listening, and gracefully shut down.

**Configuration:**
- `config/local.yaml`: Default timeout, Postgres pooling, and environment settings for local runs.
- `internal/config/config.go`: YAML parsing, default normalization, validation of HTTP timeouts and Postgres constraints.

**Core Logic:**
- `internal/services/service.go`: URL normalization, alias generation, retry logic, shared constants.
- `internal/http/handlers/create_alias.go` & `get_url.go`: HTTP schemas that decode, call the service, and encode responses.
- `internal/http/middleware/`: Method enforcement (`method.go`), structured logging (`logging.go`), panic recovery (`recover.go`).
- `internal/storage/postgres/create_alias.go`: Postgres writes using `INSERT ... ON CONFLICT` plus retry handling, referencing `migrations/migrations/000001_init.sql`.
- `internal/storage/inmemory/create_alias.go`: Local alias state with mutex-protected maps mirroring storage interface.

**Testing:**
- `internal/services/create_alias_test.go`: Service-level contract checks.
- `internal/http/handlers/mock_service_test.go`: HTTP handler behavior when service returns mocks.
- `internal/storage/inmemory/storage_test.go`: Storage correctness for in-memory maps.

**Documentation:**
- `README.md`: Usage and launch instructions.
- `.planning/codebase/`: Generated codebase docs (this file lives here now).

## Naming Conventions

**Files:**
- Go source files are lowercase (e.g., `service.go`, `create_alias.go`); test files end with `_test.go`.
- DTOs and helpers keep descriptive names aligned with their package responsibilities.

**Directories:**
- All directories follow lowercase, hyphenated names (e.g., `internal/http`, `internal/storage/postgres`).
- Packages mirror directory names for clarity when importing.

**Special Patterns:**
- `cmd/app` contains the binaries wired by `go build`.
- `_test.go` suffix designates Go tests; `mock` suffix appears rarely (`mock_service_test.go`) for hand-written substitutes.

## Where to Add New Code

**New HTTP route:** Extend `internal/http/router` + `handlers`, add DTO responses in `internal/http/dto`, and tests under `internal/http/handlers`.
**New domain logic:** Add files to `internal/services` and augment tests there; new storage reachables should implement `internal/services.Storage`.
**New storage driver:** Create `internal/storage/<driver>/` with `Create`/`GetURL` implementations plus error mappings, then register it in `cmd/app/main.go` and update docs/migrations if needed.
**Configuration changes:** Update `config/local.yaml` and `internal/config/config.go` for new settings or env overrides.
**Tests & tools:** Keep tests beside the logic they cover (`internal/services` tests next to services, `internal/http` tests next to handlers).

## Special Directories

**.planning/codebase/**:
- Purpose: Living documentation consumed by GSD workflows (stack, architecture, structure, conventions, testing, integrations, concerns).
- Source: Generated by `/gsd:map-codebase` and other planning commands; usually committed once stable.
- Committed: Yes, once validated and reviewed, since downstream phases rely on it.

**migrations/migrations/**:
- Purpose: Goose SQL files that bring Postgres schema online (`url` table with alias constraints).
- Source: Maintained by hand; apply before running Postgres storage.
- Committed: Yes; run with `goose`/`migrate` before `postgres` storage runs.

**docker/initdb/**:
- Purpose: Entrypoint for dockerized Postgres initialization (currently empty placeholder for `docker-compose` hooks).
- Source: Populated when extra SQL/fixtures are needed.
- Committed: Yes (even if empty) to keep directory tracked.

--- 

*Structure analysis: 2026-04-07*
*Update when directory layout or naming conventions change*
