# Project Explanation

## Project
- Name: tryckers-backend
- State: pre-ai
- Team: small
- Structure: multirepo
- Project language: Spanish

## Knowledge Layers
### Business — Consolidated
- ✓ business.md

### Product — Structured
- ✓ capabilities.md
- ✓ product.md

### Tech — Structured
- ✓ ADR-001-s3-compatible-storage-driver.md
- ✓ ADR-002-relational-skills-and-tags.md
- ✓ ADR-003-refresh-token-rotation-and-storage.md
- ✓ ADR-004-database-migrations-engine.md
- ✓ architecture-notes.md
- ✓ codebase.md
- ✓ current-state.md
- ✓ decision-candidates.md
- ✓ knowledge.md
- ✓ modules.md
- ✓ system-context.md

### Delivery — Traceable
- ✓ WI-001-agregar-fecha-de-cumplea-os.md
- ✓ roadmap.md

## Detected Stack
- Language: go
- Framework: gin
- Package manager: go-modules
- Source: src
- Migrations: migrations
- Infrastructure: docker-compose.yml, Dockerfile

## Scan Signals
- Auth: Auth-related paths detected
- Storage: Storage-related paths detected
- Migrations: migrations (1 file)
- Tests: No test directory detected
- Security: 2
- Infrastructure: 2
- Environment: 9 environment variable(s) detected

## Knowledge Status
- Inventory: available
- Context pack: available
- Business: Consolidated
- Product: Structured
- Tech: Structured
- Delivery: Traceable
- Agents: available
- Roadmap initiatives: 0
- Work Item candidates: 0
- Materialized Work Items: 1
- Ownership coverage: 1/1 work items

## Work Items
- Draft: 0
- Ready: 0
- In Progress: 0
- Blocked: 0
- Completed: 1
- Archived: 0

## Work Items by Type
- Features: 1

## Work Item Sources
- Manual: 1

## Work Items by Initiative
- RM-001 — Completed: 1

## Implementation Evidence

### WI-001
- Lifecycle: Completed
- Implementation: not-assessed
- Validation: not-assessed
- Release: not-assessed
Affected modules: frontend, core

## Domains
- profiles, iam

## Mapped Modules
- frontend — unknown — ../tryckers-frontend

## Module Artifact Coverage
- frontend: none

## Knowledge Graph
- Scope: all
- Nodes: 15
- Edges: 13
- Quality: partial
- Hints: 12
- Reason: All supported Work Item statuses are included.
- Last exported: 2026-08-04T04:37:00.586Z

## Skills installed: 7
Groups:
- delivery: 3
- tech: 4

## Project Route

Route: pre-ai
Progress: 14/16

- [x] Enable Kaddo
- [~] Scan repository
- [x] Bootstrap knowledge baseline
- [x] Define business context
- [x] Define product context
- [x] Discover capabilities
- [x] Describe current architecture
- [x] Capture technical decisions
- [x] Create or connect work source
- [x] Materialize first Work Item
- [x] Refine Work Item
- [x] Suggest ownership
- [x] Resolve ADRs if needed
- [x] Prepare implementation
- [ ] Run guard
- [x] Capture learning

## Phase
- Phase: Maintenance
- Reason:
  - Roadmap available
  - 1 materialized work item(s)
  - Ownership coverage 100%
- Next step: Run `kaddo questions` to see source locations and resolution guidance, then resolve, assume or defer the blocking open questions.

## Delivery Summary
- Completed Work Items: 1
- Active Work Items: 0

## Project Readiness
- overall: needs-decisions
- bootstrap baseline: complete
- scan: available
- understand: available
- agents: present
- skills: present
- business: useful
- product: useful
- capabilities: useful
- current-state: useful
- codebase: useful
- roadmap: empty
- work-items: completed-only
- work-items completed: 1
- work-items active: 0
- adapters: antigravity installed
- blocking open questions: 1
- assumptions: 0
- deferred: 0

### Recommended next step
Run `kaddo questions` to see source locations and resolution guidance, then resolve, assume or defer the blocking open questions.

## Tech Knowledge
Core:
- ✓ current-state.md
- ✓ codebase.md
Decisions:
- ADRs: 4 (draft: 0, accepted: 4) · status: accepted-adrs
Discovery:
- ✓ architecture-notes.md
- ✓ decision-candidates.md

## Roadmap Status
- Initiatives: 0
- Work Item candidates: 0
- Materialized Work Items: 1
- Remaining Work Item candidates: 0

## Installed Assets
- CLI version: 3.83.0
- Agents: 9 installed (9 outdated)
- Skills: 7 installed (7 outdated)

Suggested: run `kaddo agents status` and `kaddo skills status`.
