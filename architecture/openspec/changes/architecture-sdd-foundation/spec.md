# Spec: architecture-sdd-foundation

## Capability

The project shall provide a documented SDD workflow for backend-first vertical slice development.

## API Contract

| Method | Path | Auth | Request DTO | Response DTO |
| --- | --- | --- | --- | --- |
| N/A | N/A | N/A | N/A | N/A |

## Acceptance Criteria

- Given a contributor opens the backend repo, when they read `architecture/README.md`, then they can find the SDD, OpenSpec, roadmap, standards, and agents.
- Given a future VS starts, when `vs-agent` is followed, then it creates a `vs/<vs-name>` branch and OpenSpec files under `architecture/openspec/changes/<vs-name>/`.
- Given a completed VS is approved, when `branch-agent` is followed, then it reviews pending changes, merges to `dev` and `main`, pushes, tags, and updates the roadmap.
- Given frontend contributors need architecture context, when they read `tryckers-frontend/architecture/README.md`, then they understand the frontend module architecture and its backend dependency.

## Validation Rules

- VS names should be lowercase kebab-case.
- OpenSpec change folders should match the VS name.
- Completed VS tags should use `vs/<vs-name>/v<major>.<minor>.<patch>`.

## Persistence Rules

No persistence changes.

Future database changes must be represented as migrations and Supabase DDL must be applied through the Supabase MCP server when applicable.

## Frontend Impact

Adds frontend architecture module documentation only.

## Compatibility

No breaking runtime changes.
