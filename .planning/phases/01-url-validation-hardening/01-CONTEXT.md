# Phase 1: URL Validation Hardening - Context

**Gathered:** 2026-04-07
**Status:** Ready for planning

<domain>
## Phase Boundary

Tighten `POST /url` input validation in the existing URL shortener service so only valid trimmed `http` or `https` URLs are accepted, invalid input is rejected with the current JSON error contract, and the change remains limited to the existing handler and service architecture plus tests.

</domain>

<decisions>
## Implementation Decisions

### Validation boundary
- Enforce URL validation in both the HTTP handler and the service layer.
- Use one shared validation helper to avoid duplicating validation rules across layers.
- Keep the implementation within the existing handler → service → storage structure.

### Input normalization
- Trim only surrounding spaces from the submitted URL.
- Do not add any additional normalization in this phase beyond surrounding-space trimming.

### Failure behavior
- Empty URL, whitespace-only URL, malformed URL, unsupported scheme, and missing host must all return `400 {"error":"invalid url"}` from `POST /url`.
- Internal branching and logging may distinguish the underlying cause, but the external API response must remain the same for all validation failures.

### Test coverage
- Add handler tests for empty URL, whitespace-only URL, malformed URL, `ftp` scheme, valid `http`, and valid `https`.
- Add service-level tests for the same validation rules where applicable.
- Keep the test changes minimal and aligned with the current test structure.

### Claude's Discretion
- Exact placement of the shared validation helper.
- Internal naming of validation branches and helper functions.
- The smallest set of test refactors needed to fit the new cases into the current suite.

</decisions>

<specifics>
## Specific Ideas

No external product references were provided. The main implementation guidance is to keep the change minimal, preserve the existing API contract, and avoid architecture changes while enforcing the same validation rules in both layers.

</specifics>

<deferred>
## Deferred Ideas

None — discussion stayed within phase scope.

</deferred>

---
*Phase: 01-url-validation-hardening*
*Context gathered: 2026-04-07*
