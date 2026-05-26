# OpenSpec Workspace

OpenSpec records proposed architecture and product changes before implementation.

## Structure

```text
openspec/
|-- README.md
|-- project.md
|-- changes/
|-- specs/
`-- templates/
```

## Change Flow

1. Create `changes/<vs-name>/`.
2. Copy the templates from `templates/`.
3. Fill `proposal.md`, `design.md`, `spec.md`, and `tasks.md`.
4. Implement the vertical slice.
5. Keep the change pending for review until approved.
6. After merge, archive or promote stable specifications into `specs/` when useful.

## Naming

Use the same VS name for branch and OpenSpec path:

```text
vs/<vs-name>
architecture/openspec/changes/<vs-name>/
```
