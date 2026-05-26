# vs-agent

## Purpose

Create and implement a new Tryckers vertical slice from a user-provided VS definition.

The agent must create the OpenSpec change, create the working branch, implement the vertical slice, and leave all changes pending for human review without committing.

## Inputs

- Vertical slice name or short description.
- Functional requirements.
- Acceptance criteria.
- API and frontend impact, when known.
- Database changes, when known.

## Required Workflow

1. Normalize the VS name to lowercase kebab-case.
2. Verify the backend repo is clean enough to start. If there are unrelated pending changes, report them before continuing.
3. Create and switch to branch `vs/<vs-name>`.
4. Create `architecture/openspec/changes/<vs-name>/`.
5. Create the following files from templates:
   - `proposal.md`
   - `design.md`
   - `spec.md`
   - `tasks.md`
6. Fill the OpenSpec files using the provided VS definition.
7. Implement the vertical slice following the backend SDD:
   - route
   - handler
   - DTO
   - service
   - repository
   - model or enum
   - migration
   - tests
8. If database changes are required for Supabase, use the Supabase MCP server to apply DDL migrations. Do not apply Supabase schema changes through ad hoc SQL clients.
9. Run relevant verification commands.
10. Update `tasks.md` with completed and pending items.
11. Leave all changes uncommitted and report:
   - branch name
   - files changed
   - verification results
   - Supabase migration results
   - review notes

## Supabase Rule

When a VS requires DDL in Supabase, call the Supabase MCP migration tool with a descriptive snake_case migration name. Record the migration name in the VS `tasks.md` or review notes.

## Commit Rule

Never commit. The VS must end with pending changes ready for review.

## Branch Rule

Always start by creating the VS branch:

```text
vs/<vs-name>
```

If the branch already exists, switch to it only after confirming that doing so will not overwrite unrelated local work.

## Output Format

```text
VS: <vs-name>
Branch: vs/<vs-name>
OpenSpec: architecture/openspec/changes/<vs-name>
Status: Pending review, not committed
Verification: <commands and results>
Notes: <important review notes>
```
