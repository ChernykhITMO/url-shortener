# URL Shortener

## What This Is

URL Shortener is a Go HTTP service for creating short aliases for original URLs and resolving aliases back to the stored URL. The current brownfield work is focused on tightening `POST /url` validation so invalid input is rejected consistently without changing the existing API shape.

## Core Value

Creating and resolving short links must be predictable, safe, and consistent for clients.

## Requirements

### Validated

- ✓ Client can create or reuse a short alias for a valid original URL via `POST /url` — existing
- ✓ Client can resolve an existing alias via `GET /url/{alias}` — existing
- ✓ Service can run against either in-memory or Postgres storage behind the same service and HTTP layers — existing

### Active

- [ ] `POST /url` trims surrounding spaces before validating and storing the submitted URL
- [ ] `POST /url` rejects empty URLs with `400 {"error":"invalid url"}`
- [ ] `POST /url` rejects malformed URLs with `400 {"error":"invalid url"}`
- [ ] `POST /url` accepts only `http` and `https` URLs
- [ ] URL validation is enforced in both the HTTP handler and the service layer
- [ ] Tests cover the new validation behavior with minimal changes to the current code structure

### Out of Scope

- New endpoints or response schemas — the current API contract already fits the requested behavior
- Changes to alias generation or storage behavior — this feature is limited to request validation
- Broader normalization or product expansion beyond `POST /url` input validation — keep the change minimal

## Context

The codebase is a three-layer Go service with custom `net/http` routing and middleware, a service layer for business logic, and interchangeable in-memory/Postgres storage. Existing mapping in `.planning/codebase/` shows that `POST /url` already flows through `internal/http/handlers/create_alias.go` into `internal/services/create_alias.go`, which makes those two layers the correct places to enforce the new validation contract.

The immediate work is an SDD-style brownfield improvement to the create-alias path. The user wants handler-level validation for correct API behavior and service-level validation as a safety guard, while preserving the current JSON error contract and keeping the patch small.

## Constraints

- **Architecture**: Keep the existing handler → service → storage layering — the repo already depends on it
- **API Contract**: Validation failures for `POST /url` must return `400 {"error":"invalid url"}` — callers already rely on this shape
- **Scope**: Minimize code changes and stay focused on `POST /url` validation — avoid unrelated cleanup
- **Quality**: Add or update tests for both HTTP and service behavior — the change must be regression-safe

## Key Decisions

| Decision | Rationale | Outcome |
|----------|-----------|---------|
| Enforce URL validation in both handler and service | The handler must guarantee the correct HTTP response, and the service must remain safe if called from another entry point | — Pending |
| Trim surrounding spaces before validation | The feature explicitly requires accepting padded input without storing those outer spaces | — Pending |
| Keep the feature limited to `POST /url` validation and tests | The user asked for minimal change in a brownfield codebase | — Pending |

---
*Last updated: 2026-04-07 after initialization*
