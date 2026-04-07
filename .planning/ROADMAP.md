# Roadmap: URL Shortener

## Overview

This roadmap covers a focused brownfield improvement to the existing URL shortener service. The work is intentionally compressed into one phase so the codebase can quickly gain stricter `POST /url` validation, keep the API response contract stable, and add regression tests without expanding scope.

## Phases

**Phase Numbering:**
- Integer phases (1, 2, 3): Planned milestone work
- Decimal phases (2.1, 2.2): Urgent insertions (marked with INSERTED)

Decimal phases appear between their surrounding integers in numeric order.

- [ ] **Phase 1: URL Validation Hardening** - Tighten `POST /url` input validation and lock the behavior down with tests

## Phase Details

### Phase 1: URL Validation Hardening
**Goal**: Clients can submit only valid trimmed `http` or `https` URLs to `POST /url`, and invalid input is rejected consistently with the expected JSON error response.
**Depends on**: Nothing (first phase)
**Requirements**: URL-01, URL-02, URL-03, URL-04, URL-05, URL-06
**Success Criteria** (what must be TRUE):
  1. Client can submit a valid `http` or `https` URL with surrounding spaces and still receive a created alias.
  2. Client receives `400 {"error":"invalid url"}` for empty, malformed, or non-`http`/`https` input sent to `POST /url`.
  3. The same validation rules are enforced at both the HTTP handler boundary and the service layer.
  4. Automated tests cover the accepted trimming case and the requested invalid-input cases.
**Plans**: 2 plans

Plans:
- [ ] 01-01: Implement minimal validation updates in the handler and service path
- [ ] 01-02: Add and adjust tests for handler and service validation behavior

## Progress

**Execution Order:**
Phases execute in numeric order: 1

| Phase | Plans Complete | Status | Completed |
|-------|----------------|--------|-----------|
| 1. URL Validation Hardening | 0/2 | Not started | - |
