# Testing Guide

**Updated:** 2026-04-07

## Running the suite
- Unit and integration coverage is invoked via `go test ./...` (see `README.md:105` for the primary command) and `go test -race ./...` for race detection.
- PostgreSQL-specific integration targets use `go test -tags=integration ./internal/storage/postgres -v` after `docker compose up -d db` so the local `postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable` instance is reachable (`internal/storage/postgres/setup_integration_test.go:1`, `README.md:118`).

## Unit tests
- Services rely on table-driven tests with clear scenario names, per-case `t.Run`, and `t.Fatalf` assertions so failures show the offending variant (`internal/services/create_alias_test.go:1`). Manual mocks (`internal/services/mock_storage_test.go:1`) inject behavior to exercise validation, alias generation, and retry logic while keeping the suite fast.
- HTTP handlers use `httptest.NewRequest`, `httptest.NewRecorder`, and `slog` loggers pointed at `io.Discard` so JSON decoding/encoding paths are validated without noise (`internal/http/handlers/create_alias_test.go:1`, `internal/http/handlers/get_url_test.go:1`). They assert both `status` codes and response bodies to prove the `respond` layer works end-to-end.
- Middleware tests follow the same pattern: `RequireMethodMiddleware` checks `Allow` headers and `status 405`, while `RecoverMiddleware` ensures panics become `500` JSON responses (`internal/http/middleware/method_test.go:1`, `internal/http/middleware/recover_test.go:1`). The handler package mocks `Handler.Service` through `MockService` (`internal/http/handlers/mock_service_test.go:1`).

## Integration tests
- Build-tagged tests live in `internal/storage/postgres/*_integration_test.go` and rely on `TestMain` to bootstrap a connection, apply `goose` migrations from `migrations/migrations`, and provide `cleanupURLTable` helpers (`internal/storage/postgres/setup_integration_test.go:1`). Each test (e.g., `get_url_integration_test.go:1`) ensures `storage.ErrNotFound` or successful lookups respond as expected after migrations.
- Integration commands trace back to README instructions (`README.md:118`) so document both the `docker compose up -d db` prerequisite and that the suite uses the same DSN as the app (`config/local.yaml:1`).

## Stress and concurrency checks
- The in-memory storage is explicitly challenged with parallel `Create`/`GetURL` goroutines and a `sync.WaitGroup` so its mutex-guarded maps stay consistent (`internal/storage/inmemory/storage_test.go:1`). Any future in-memory replacement should keep the same concurrency guardrails.

