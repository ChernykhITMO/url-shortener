# Technology Stack

**Analysis Date:** 2026-04-07

## Languages
**Primary:**
- Go 1.25.1 powers the entire service, from the entrypoint in `cmd/app/main.go` through `internal/http`, `internal/services`, and storage implementations that share `internal/domain/alias` constants.
**Secondary:**
- PostgreSQL SQL (see `migrations/migrations/000001_init.sql`) defines the alias table, length/format checks, and is the target of migrations run with `goose`.

## Runtime
**Environment:**
- Linux static binary produced via `CGO_ENABLED=0 GOOS=linux go build -o /out/url-shortener ./cmd/app` in the Dockerfile’s `app-builder` stage, with local development supported by `go run ./cmd/app --storage=... --config=...`.
**Package Manager:**
- Go modules (`go.mod`, `go.sum`, `go mod download` in `Dockerfile`) lock direct deps and transitives for reproducible builds.

## Frameworks
**Core:**
- Standard library `net/http` underpins the transport stack: `internal/http/router/router.go` wires `http.NewServeMux`, `internal/http/middleware` adds recover/logging/method guards, and handlers in `internal/http/handlers` map service errors to JSON responses while honoring `respond` helpers.
- `log/slog`, `context`, `crypto/rand`, `net/url`, and `strings` drive request validation and alias generation inside `internal/services`; storage abstractions in `internal/storage/inmemory` and `internal/storage/postgres` implement the shared interface.
**Testing:**
- `go test ./...` plus `go test -race ./...` cover unit logic, while integration tests (`go test -tags=integration ./internal/storage/postgres -v`) run against the Postgres container started by `docker compose up -d db`.
**Build & Tooling:**
- `docker compose up --build` (`docker-compose.yaml`) orchestrates `db`, `migrate`, and `app`; migrations execute via the `migrate` service or `Dockerfile` `migrate` stage running `goose -dir /migrations/migrations postgres $POSTGRES_DSN up`.

## Key Dependencies
**Critical:**
- `github.com/jackc/pgx/v5` in `internal/storage/postgres` provides the `pgx` driver for `sql.Open`, while `gopkg.in/yaml.v3` in `internal/config/config.go` parses `config/local.yaml`.
**Infrastructure:**
- `github.com/pressly/goose/v3` runs migrations from the Docker `migrate` service or `Dockerfile` migration stage.
- `golang:1.25-alpine` base images keep builder and runtime environments consistent.

## Configuration
**Environment:**
- `config/local.yaml` configures `service.max_attempts`, HTTP server timeouts, and Postgres pool limits; `internal/config/config.go` normalizes defaults and allows `HTTP_ADDR`/`POSTGRES_DSN` overrides from `.env.example`.
**Build:**
- Dockerfile copies `config/local.yaml` into `/app/config/local.yaml` and starts `/app/url-shortener --storage=postgres --config=/app/config/local.yaml`; the config loader enforces `service.max_attempts >= 1` and defaults the HTTP address to `:8080` when missing.

## Platform Requirements
**Development:**
- Go 1.25+ toolchain, Docker/Compose, and a reachable Postgres (see `.env.example` for `POSTGRES_DSN`) are needed for both unit/integration tests and `docker compose up --build`.
**Production:**
- The static binary listens on the address supplied by `HTTP_ADDR`/`HTTP_EXTERNAL_PORT` and consumes any Postgres instance reachable via `POSTGRES_DSN`; `docker-compose.yaml` maps `${HTTP_EXTERNAL_PORT}:8080` when Compose is used.

---
*Stack analysis: 2026-04-07*
