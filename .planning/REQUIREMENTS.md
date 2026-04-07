# Requirements: URL Shortener

**Defined:** 2026-04-07
**Core Value:** Creating and resolving short links must be predictable, safe, and consistent for clients

## v1 Requirements

### URL Validation

- [ ] **URL-01**: Client can submit a URL with surrounding spaces to `POST /url`, and the service validates the trimmed value
- [ ] **URL-02**: Client receives `400 {"error":"invalid url"}` when `POST /url` receives an empty `url` value after trimming
- [ ] **URL-03**: Client receives `400 {"error":"invalid url"}` when `POST /url` receives a malformed URL
- [ ] **URL-04**: Client can shorten URLs only when the scheme is `http` or `https`
- [ ] **URL-05**: URL validation rules are enforced consistently at the HTTP handler boundary and in the service layer
- [ ] **URL-06**: Automated tests cover successful trimming and all requested invalid-input cases for `POST /url`

## v2 Requirements

None yet.

## Out of Scope

| Feature | Reason |
|---------|--------|
| Alias format or generation changes | Not required for the validation feature |
| GET `/url/{alias}` behavior changes | Outside the requested scope |
| New URL normalization rules beyond trimming surrounding spaces | Keep the implementation minimal |

## Traceability

| Requirement | Phase | Status |
|-------------|-------|--------|
| URL-01 | Unmapped | Pending |
| URL-02 | Unmapped | Pending |
| URL-03 | Unmapped | Pending |
| URL-04 | Unmapped | Pending |
| URL-05 | Unmapped | Pending |
| URL-06 | Unmapped | Pending |

**Coverage:**
- v1 requirements: 6 total
- Mapped to phases: 0
- Unmapped: 6 ⚠️

---
*Requirements defined: 2026-04-07*
*Last updated: 2026-04-07 after initialization*
