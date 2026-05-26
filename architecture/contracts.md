# Architecture Contracts

## API Contracts

- Public backend APIs use REST under `/api/v1`.
- Request and response bodies must be represented by DTOs under `src/internal/dtos`.
- Handlers must map HTTP input into DTOs before calling services.
- Services must return domain outcomes or application errors that handlers translate into HTTP responses.
- New endpoints must document path, method, request DTO, response DTO, auth requirements, and error cases in the related OpenSpec change.

## Layer Contracts

| Layer | May Depend On | Must Not Depend On |
| --- | --- | --- |
| Handlers | DTOs, services, middleware helpers | GORM, repositories directly |
| Services | DTOs, repositories, models, utils | Gin context, HTTP response writers |
| Repositories | GORM, models, enums | Gin, handlers, HTTP DTO mapping |
| Models | GORM tags, enums | Handlers, services |
| DTOs | Basic types, enums when useful | GORM persistence behavior |

## Database Contracts

- PostgreSQL is the source of truth for persisted state.
- Schema changes must be expressed as migrations.
- UUIDs are preferred for primary identifiers.
- Status enums should be used for MVP soft-deletion flows.
- Supabase-related DDL must be applied through the Supabase MCP server by the responsible agent.

## Frontend Contract

The frontend consumes backend APIs through typed services and DTOs. Backend changes that affect payloads, paths, auth behavior, or error responses must include frontend impact notes in `architecture/openspec/changes/<vs-name>/spec.md`.

## Versioning Contract

Each completed vertical slice receives a version tag after it is merged into `main`.

Tag format:

```text
vs/<vs-name>/v<major>.<minor>.<patch>
```

Roadmap entries must reference the tag, branch, OpenSpec change path, and merge date.
