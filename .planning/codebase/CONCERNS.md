# Codebase Concerns

**Analysis Date:** 2026-04-07

## Tech Debt
**Config surface vs code constants:**
- Issue: The POST `/url` handler hard-codes `maxCreateAliasBodyBytes = 1024` (`internal/http/handlers/create_alias.go`) and alias length/styling is hard-coded to 10 characters in `internal/domain/alias/alias.go`. There is no way to tune these values via `config/local.yaml` or CLI flags.
- Why: The limits live in code rather than in the config struct that already stores timeouts and `Service.MaxAttempts`, so every change requires a rebuild and redeploy.
- Impact: Responding to operational needs (longer URLs, different alias length) or aligning with a new database schema requires a code change instead of a config tweak, which slows down releases and risks a mismatch between DB constraints and validation.
- Fix approach: Promote these constants into the config layer (and sync them with the `url` table constraint in `migrations/migrations/000001_init.sql`) so they can be tuned per environment without touching the binary.

## Known Bugs
**Postgres run command documented but fails:**
- Symptoms: Following the README’s `go run ./cmd/app --storage=postgres --config=./config/local.yaml` command exits immediately with `postgres.dsn is required for postgres storage` because `config/local.yaml` does not set `postgres.dsn` and `cmd/app/main.go` rejects empty DSN before starting.
- Trigger: Running the documented command without setting `POSTGRES_DSN` in the environment.
- Workaround: Export `POSTGRES_DSN` (e.g., `postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable`) before running or pass it via the Docker `.env` file described later.
- Root cause: Docker instructions override DSN via `.env`, but README omits that step, so the Postgres path is not runnable out of the box.
- Fix approach: Either provide a sane default DSN in `config/local.yaml` or update the README/instructions to require `POSTGRES_DSN` for the local Postgres command.

## Security Considerations
**No rate limiting or authentication on `/url`:**
- Risk: Anyone can flood the POST `/url` endpoint with random URLs, triggering crypto alias generation and multiple DB inserts until `ErrAttemptsExceeded` (503). A simple botnet can easily exhaust CPU and Postgres connections.
- Current mitigation: The router stack (`internal/http/router/router.go`) only composes recover, logging, and method requirements; nothing throttles abuse.
- Recommendations: Add middleware that enforces per-IP or per-key rate limits before hitting the handler, or gate creation behind an API key/token to throttle abusive actors.

## Performance Bottlenecks
**Alias creation retries hammer Postgres:**
- Problem: `services.CreateAlias` (configured via `cfg.Service.MaxAttempts`) retries up to 20 times on alias conflict, and each try calls `postgres.Storage.Create`, which always runs `INSERT ... ON CONFLICT` plus a fallback `SELECT` when collisions occur (`internal/storage/postgres/create_alias.go`).
- Measurement: With contention, every alias generation issues at least one query (sometimes two), so throughput hits a wall once collisions become frequent.
- Cause: No alias reservation, caching, or de-duplication before hitting the DB; the service races the database without a backoff strategy beyond a fixed 10 ms delay.
- Improvement path: Buffer precomputed alias blocks or move the conflict handling fully into a single `INSERT ... ON CONFLICT (alias) DO UPDATE` statement and track metrics, then tune `MaxAttempts` downward once retries are rare.

## Fragile Areas
**Router alignment with Go 1.25 wildcard behavior:**
- Why fragile: The code relies on the new `http.ServeMux` wildcard `/url/{alias}` and `Request.PathValue`/`SetPathValue` (Go 1.25), but there are no integration tests exercising `cmd/app/http.Run`; only handler tests call `req.SetPathValue` manually (`internal/http/handlers/*.go` and tests). Any change to the mux or downgrade to earlier Go versions would silently break alias fetches.
- Common failures: Modifying middleware stack, replacing the mux, or forgetting to set path values will make `validateAlias` always error and return 400/404.
- Safe modification: Add an `httptest.Server` integration test that performs real POST/GET flows against `cmd/app/http.Run` to ensure the wildcard path and logging/recover middleware keep working.
- Test coverage: Router wiring is never hit by default `go test ./...` (only handler unit tests exist).

## Scaling Limits
**URL table grows without retention:**
- Current capacity: `url(alias, original_url, created_at)` accumulates every shortened link indefinitely (`migrations/migrations/000001_init.sql`).
- Limit: Disk usage and backup/restore times grow linearly with incoming links because there is no TTL, archive, or deletion strategy.
- Symptoms at limit: `SELECT alias FROM url WHERE original_url = $1` remains indexed but Postgres storage still grows, making migrations/backups slow and increasing restore windows.
- Scaling path: Introduce retention/archival (e.g., periodic job that archives rows older than X days or partitions by `created_at`) and consider adding cleanup endpoints or admin jobs.

## Test Coverage Gaps
**Postgres storage only exercised under `integration` tag:**
- What's not tested: `internal/storage/postgres/*_integration_test.go` files are guarded by `//go:build integration`, so `go test ./...` never runs them.
- Risk: Changes to the Postgres layer (DSN handling, alias queries) can fail silently until someone manually runs the integration suite.
- Priority: High. Database storage is the production path.
- Difficulty to test: Requires a running Postgres; consider adding a lightweight Dockerized Postgres job in CI or making the integration tests opt-in via `go test ./... -tags=integration` in a separate job.

**HTTP server/router not covered end-to-end:**
- What's not tested: `cmd/app/http.Run` and the `http.Server` stack are not exercised in any package test; only handler logic is unit tested.
- Risk: Changes to timeouts, shutdown logic, or middleware registration can regress unnoticed, especially the new wildcard route that `PathValue` depends on.
- Priority: Medium. Adding a small `httptest.NewServer` test that runs `router.New` through the same middleware will catch wiring regressions.

---

*Concerns audit: 2026-04-07*
