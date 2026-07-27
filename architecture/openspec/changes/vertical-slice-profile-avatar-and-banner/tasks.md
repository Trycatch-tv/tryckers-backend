# Tasks: vertical-slice-profile-avatar-and-banner

## OpenSpec

- [x] Fill `proposal.md`.
- [x] Fill `design.md`.
- [x] Fill `spec.md`.
- [x] Review scope and acceptance criteria.

## Backend

- [x] Add `avatar_url` and `banner_url` to user model.
- [x] Add user repository media update methods.
- [x] Add local storage service.
- [x] Add authenticated avatar upload endpoint.
- [x] Add authenticated banner upload endpoint.
- [x] Add authenticated avatar/banner remove endpoints.
- [x] Serve local uploads publicly.
- [x] Add upload folders and gitignore rules.

## Supabase

- [ ] Apply required database migration through the Supabase MCP server when `project_id` is available.
- [ ] Migration SQL:

```sql
ALTER TABLE users ADD COLUMN IF NOT EXISTS avatar_url TEXT;
ALTER TABLE users ADD COLUMN IF NOT EXISTS banner_url TEXT;
```

## Frontend

- [x] Add avatar/banner fields to user contract.
- [x] Add upload methods in Tryckers service.
- [x] Add profile avatar upload UX.
- [x] Add profile banner upload UX.
- [x] Add profile avatar/banner remove UX.
- [x] Add preview/loading/error states.
- [x] Add responsive profile rendering.

## Tests

- [x] Run `go test ./...`.
- [x] Run Angular build or targeted verification.

## Review

- [x] Check git diff.
- [x] Leave changes uncommitted for review.
