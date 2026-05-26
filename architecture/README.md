# Tryckers Backend Architecture

This folder is the source of truth for the Tryckers backend SDD (Software Design Documentation).

The backend is the main architecture repository for the Tryckers ecosystem. Frontend architecture references this documentation and describes its own module boundaries separately.

## Contents

- `system-design.md`: main SDD and backend architecture definition.
- `contracts.md`: API, DTO, persistence, and cross-repo contract rules.
- `glossary.md`: shared language for contributors.
- `adr/`: Architecture Decision Records.
- `diagrams/`: architecture and flow diagrams.
- `standards/`: contribution, coding, testing, and vertical-slice standards.
- `tech/`: technology-specific notes.
- `roadmap/`: roadmap, release tags, and vertical-slice history.
- `openspec/`: OpenSpec workspace for proposed and accepted changes.
- `agents/`: project agent runbooks.

## SDD Workflow

Every meaningful feature or architectural change should be handled as a vertical slice (VS):

1. Create a branch named after the VS.
2. Create an OpenSpec change under `architecture/openspec/changes/<vs-name>/`.
3. Add `proposal.md`, `design.md`, `spec.md`, and `tasks.md`.
4. Implement the slice across API, service, repository, model, DTO, migrations, and tests as needed.
5. Leave changes pending for review.
6. After approval and merge, update `architecture/roadmap/README.md` with the version tag.

Use `architecture/agents/vs-agent.md` to create a vertical slice and `architecture/agents/branch-agent.md` to review, merge, tag, push, and update the roadmap.
