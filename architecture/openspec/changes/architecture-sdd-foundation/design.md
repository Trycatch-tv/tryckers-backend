# Design: architecture-sdd-foundation

## Architecture Fit

This change creates the documentation layer around the existing Layered Modular Monolith. It does not change runtime behavior.

The backend repo becomes the main architecture source of truth. The frontend repo keeps a smaller architecture module description that references backend contracts.

## Layer Changes

| Layer | Planned Change |
| --- | --- |
| Route | None |
| Handler | None |
| DTO | None |
| Service | None |
| Repository | None |
| Model | None |
| Migration | None |
| Documentation | Add SDD, OpenSpec, roadmap, standards, and agent runbooks |

## Data Flow

No runtime data flow changes.

Future VS implementations should follow:

```text
Client -> Route -> Handler -> Service -> Repository -> PostgreSQL
```

## Database Design

No database changes.

## Error Handling

No runtime error handling changes.

## Testing Design

Verification is limited to checking file structure and git status.

Future slices should add tests proportional to implementation risk.

## Tradeoffs

The documentation is kept in Markdown so contributors can review it easily in pull requests. Agent behavior is documented as project runbooks first, which keeps the workflow portable before converting it into a dedicated plugin or skill.
