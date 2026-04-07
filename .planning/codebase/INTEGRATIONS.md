# External Integrations

**Analysis Date:** 2026-04-07

## APIs & External Services
- None. The service only exposes its own JSON REST surface (`internal/http/handlers`) and does not call out to third-party HTTP APIs or webhook providers.

## Data Storage
**Databases:**
- PostgreSQL (`db` service in `docker-compose.yaml`) is the primary store for aliases.
  - Connection: the service reads `POSTGRES_DSN` (default value in `.env.example`) which points at `postgres://postgres:postgres@db:5432/postgres?sslmode=disable`.
  - Client: `internal/storage/postgres` uses `github.com/jackc/pgx/v5` (`internal/storage/postgres/storage.go`) and the stdlib `database/sql` wrapper for connection pooling.
  - Migrations: `goose` applies SQL from `migrations/migrations/000001_init.sql` either through the Compose `migrate` service (`docker-compose.yaml`) or the `Dockerfile` `migrate` stage that runs `goose -dir /migrations/migrations ... up`.
- In-memory (`internal/storage/inmemory`) is available when running `go run ./cmd/app --storage=inmemory --config=...`, but those entries disappear when the process exits.

**File Storage:**
- None. No files are persisted outside the database or binary.

**Caching:**
- None. Every request reads/writes the URL table directly.

## Authentication & Identity
- None. All endpoints are unauthenticated; there are no OAuth providers or custom auth services.

## Monitoring & Observability
**Logs:**
- `log/slog` writes structured logs to stdout (`internal/http` + `cmd/app/main.go`); nothing is shipped to external log/monitoring services.

**Error Tracking & Analytics:**
- None configured.

## CI/CD & Deployment
**Hosting:**
- Docker Compose (`docker-compose.yaml`) runs three services: `db` (Postgres), `migrate` (runs `goose`), and `app` (built from `Dockerfile` `app` stage).
  - Deployment: manual `docker compose up --build` with `.env` variables.
  - Environment vars: take them from `.env` (copy `.env.example`), so the runtime sees `HTTP_EXTERNAL_PORT`, `POSTGRES_*`, and `POSTGRES_DSN`.
**CI Pipeline:**
- Not defined in the repository; tests and builds run locally via `go test`/`docker compose`.

## Environment Configuration
**Development:**
- Required env vars: `HTTP_ADDR`/`HTTP_EXTERNAL_PORT`, `POSTGRES_USER`, `POSTGRES_PASSWORD`, `POSTGRES_DB`, `POSTGRES_EXTERNAL_PORT`, and `POSTGRES_DSN` (see `.env.example` and how `cmd/app/main.go` reads `config/local.yaml` via `internal/config/config.go`).
- Secrets location: `.env` (copy from `.env.example` and keep it gitignored).
- Mock/stub services: none; `docker compose up -d db` spins up a real Postgres instance for integration tests.
**Staging/Production:**
- Not explicitly defined; rely on the same env variables and a production-grade Postgres instance. `HTTP_EXTERNAL_PORT` can be remapped, and `POSTGRES_DSN` should point to the production database. Any orchestration (K8s, Docker Compose, etc.) must expose port 8080 or set `HTTP_ADDR`.

## Webhooks & Callbacks
- None. There are no incoming or outgoing webhooks configured.

---
*Integration audit: 2026-04-07*
*Update when adding/removing external services*
