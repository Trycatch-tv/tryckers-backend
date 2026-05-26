# Vertical Slice Standard

## Naming

Use lowercase kebab-case:

```text
<domain>-<capability>
```

Examples:

- `auth-refresh-token`
- `posts-create-flow`
- `comments-moderation`

Branches should use:

```text
vs/<vs-name>
```

## Required OpenSpec Files

Each VS must create:

- `proposal.md`
- `design.md`
- `spec.md`
- `tasks.md`

Path:

```text
architecture/openspec/changes/<vs-name>/
```

## Implementation Scope

A VS may touch:

- API route.
- Handler.
- DTO.
- Service.
- Repository.
- Model.
- Enum.
- Migration.
- Tests.
- Frontend contract notes.
- Roadmap after merge.

It should remain focused on one product capability or one architecture decision.

## Done Criteria

- OpenSpec files are complete.
- Code compiles.
- Relevant tests pass or the test gap is documented.
- Database changes are migrated.
- Pending changes are left uncommitted for review when produced by `vs-agent`.
- After merge, `branch-agent` updates roadmap and creates the version tag.
