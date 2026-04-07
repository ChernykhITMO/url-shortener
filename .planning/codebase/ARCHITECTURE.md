# Architecture

**Analysis Date:** 2026-04-07

## Pattern Overview
**Overall:** Monolithic HTTP API service that exposes `/url` endpoints for shortening and resolving URLs while letting operators swap between in-memory and Postgres storage back ends.

**Key Characteristics:**
- Single `cmd/app` binary wires configuration, logging, storage selection, and HTTP server startup.
- HTTP stack is built around the standard library `http.Server`, a custom `router`, and middleware for logging, method gating, panic recovery, and consistent JSON responses.
- A distinct service layer normalizes URLs, generates 10-character aliases (via `internal/domain/alias.Length`), retries storage writes, and maps storage errors into domain errors.
- Storage providers implement `internal/services.Storage`: `internal/storage/inmemory` keeps alias maps in mutex-protected Go maps, while `internal/storage/postgres` relies on `pgx` and `migrations/migrations/000001_init.sql` to enforce uniqueness/length constraints in the `url` table.
- Configuration comes from YAML (`config/local.yaml`) and environment overrides handled by `internal/config/config.go`, driving HTTP timeouts, max attempts, and Postgres connection pooling.

## Layers

**Bootstrap & Config Layer:**
- Purpose: Parse CLI flags, load/validate YAML (`internal/config/config.go`), set up structured logging, handle signals, and instantiate storage.
- Contains: `cmd/app/main.go`.
- Depends on: configuration logic, service constructors, and storage drivers.
- Used by: `cmd/app/http/app.go` to pass a ready-to-run service, storage, and logger to the HTTP server.

**HTTP Layer:**
- Purpose: Match `/url` routes, enforce HTTP methods, decode JSON payloads, and translate between HTTP and the service layer.
- Contains: `cmd/app/http/app.go`, `internal/http/router/router.go`, `internal/http/handlers`, `internal/http/middleware`, `internal/http/dto`, `internal/http/respond`.
- Depends on: the service layer for business logic and `log/slog` for structured logging.
- Used by: Clients hitting `/url` and `/url/{alias}`.

**Service Layer:**
- Purpose: Normalize and validate URLs, generate or validate 10-character aliases, and delegate storage interactions while honoring `maxAttempts`.
- Contains: `internal/services/service.go`, `internal/services/create_alias.go`, `internal/services/get_url.go`, `internal/services/errors.go`.
- Depends on: `internal/domain/alias` for alias length constants and the `Storage` interface for persistence.
- Used by: HTTP handlers (`internal/http/handlers`).

**Storage Layer:**
- Purpose: Provide concrete persistence implementations for the `Storage` interface with defensive error mapping.
- Contains: `internal/storage/inmemory` (mutex-protected maps) and `internal/storage/postgres` (Postgres client + retry logic), plus `internal/storage/errors.go`.
- Depends on: `internal/storage/errors` for shared error constants and, in the Postgres case, `github.com/jackc/pgx/v5` plus the migration SQL under `migrations/migrations/000001_init.sql`.
- Used by: service layer when storing or looking up alias↔URL mappings.

## Data Flow

**Create Alias Request (POST /url):**
1. `cmd/app/http/app.go` builds the server and the `router`, which wraps `Handler.CreateAlias` in middleware (`RequireMethod`, `Logging`, `Recover`).
2. The handler decodes the JSON into `dto.CreateAliasRequest`, ensures no extra payload, and calls `service.CreateAlias`.
3. The service normalizes the URL (`normalizeAndValidateURL`), generates a random alias (10 characters from `internal/domain/alias.Length`), and attempts `storage.Create`, possibly retrying until `maxAttempts` is exhausted.
4. Storage returns the alias (deduplicating by original URL); the handler writes `CreateAliasResponse` via `respond.WriteJSON`.
5. Middleware logs method/path/status/duration, and `RecoverMiddleware` catches panics before they leak.

**Resolve Alias Request (GET /url/{alias}):**
1. `router` enforces GET on `/url/{alias}` before invoking `Handler.GetURL`.
2. The handler extracts `alias` from the path, calls `service.GetURL`, and relies on `validateAlias` plus storage lookup.
3. Storage (in-memory maps or Postgres `SELECT original_url`) returns the URL; otherwise service surfaces `services.ErrNotFound`.
4. Handler serializes the original URL via `dto.GetURLResponse`.
5. Logging middleware records the request, and `Handler.writeServiceError` maps domain errors to `404/400/503/500` responses.

**State Management:**
- Alias state is either ephemeral (mutex-guarded maps in `internal/storage/inmemory`) or persistent (Postgres `url` table defined by `migrations/migrations/000001_init.sql`, including length/format constraints). The service maintains no additional in-memory state beyond the alias generator.

## Key Abstractions

**`internal/services.Storage`:**
- Purpose: Encapsulate persistence for alias creation and retrieval.
- Examples: `internal/storage/inmemory.Storage`, `internal/storage/postgres.Storage`.
- Pattern: Interface backed by either in-memory map or relational table with retry logic.

**`handlers.Handler`:**
- Purpose: Bridge HTTP requests to the service layer while enforcing JSON decoding rules and mapping domain errors to HTTP statuses.
- Examples: `CreateAlias`, `GetURL`, `writeServiceError`.
- Pattern: Struct with injected logger and service interface for dependency inversion in tests.

**`Service`:**
- Purpose: Coordinate URL normalization, alias generation, and storage usage behind `CreateAlias`/`GetURL`.
- Examples: `create_alias.go`, `get_url.go`.
- Pattern: Domain logic encapsulated in a struct that receives `Storage` and `maxAttempts`.

**`Config`:**
- Purpose: Load YAML/env configuration for HTTP timeouts, Postgres connections, and retry limits.
- Examples: `config/local.yaml`, `internal/config/config.go`.
- Pattern: Immutable struct wired into `cmd/app/main.go`.

## Entry Points

**`cmd/app/main.go`:**
- Location: `cmd/app/main.go`.
- Triggers: Binary execution (`go run ./cmd/app` or `./url-shortener`).
- Responsibilities: Parse `--storage`/`--config`, load `Config`, choose storage (in-memory or Postgres), set up structured logging, handle shutdown signals, and call `app_http.Run`.

**`cmd/app/http/app.go`:**
- Location: `cmd/app/http/app.go`.
- Triggers: Called from `main.run`.
- Responsibilities: Compose the service, handler, router, `http.Server`, serve, log start/stop, and gracefully shut down the server when the context closes.

## Error Handling

**Strategy:** Domain errors bubble up to handlers, which translate them into HTTP statuses; middleware catches panics and `respond.WriteJSONError` ensures consistent JSON error objects.

**Patterns:**
- `services.ErrInvalidURL`, `services.ErrInvalidAlias`, and `services.ErrNotFound` are compared in `handlers.Handler.writeServiceError`.
- Storage layer surfaces `storage.ErrAliasConflict`, `storage.ErrInvalidAlias`, `storage.ErrNotFound`; the service maps them when retrying alias generation.
- `middleware.RecoverMiddleware` catches panics and responds with HTTP 500 if headers are still pending.
- `respond.WriteJSONError` sets `Content-Type: application/json` and serializes `{"error": ...}` consistently.

## Cross-Cutting Concerns

**Logging:**
- Structured logging via `log/slog` is configured in `cmd/app/main.go` (info for `local`, warn otherwise) and reused by middleware and handlers (`h.log`, `router`, `app_http.Run`).

**Validation:**
- HTTP handlers rely on `json.Decoder` with `DisallowUnknownFields`.
- `normalizeAndValidateURL` enforces scheme/host, lowercases, strips default ports, and appends `/` when needed.
- `validateAlias` enforces the 10-character alphabet defined in `internal/domain/alias`.

**Configuration & Timeouts:**
- YAML `config/local.yaml` plus env overrides (`HTTP_ADDR`, `POSTGRES_DSN`) flow through `internal/config`.
- HTTP server timeouts (`read`, `write`, `idle`, `shutdown`) and `Service.maxAttempts` are configurable per environment.

**Graceful Shutdown & Signals:**
- `cmd/app/main.go` obtains `signal.NotifyContext` for `SIGINT/SIGTERM`, passes the context to `app_http.Run`, and ensures `http.Server.Shutdown` closes before exit.

**Database Schema:**
- Postgres storage depends on `migrations/migrations/000001_init.sql` to create the `url` table with primary key, uniqueness, and alias constraints, ensuring storage-level invariants before the service runs.

**Middleware Stack:**
- `middleware.RequireMethodMiddleware` enforces HTTP verbs.
- `middleware.LoggingMiddleware` records method/path/status/duration.
- `middleware.RecoverMiddleware` prevents panics from leaking.

--- 

*Architecture analysis: 2026-04-07*
*Update when major patterns or layers change*
