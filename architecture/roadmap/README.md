# Tryckers Backend Roadmap

This roadmap tracks vertical slices, release tags, and architecture milestones.

## Tagging Rules

Completed vertical slices are tagged after merge to `main`.

Format:

```text
vs/<vs-name>/v<major>.<minor>.<patch>
```

Version guidance:

- Patch: implementation detail or bug fix within an existing capability.
- Minor: new vertical slice or backwards-compatible API capability.
- Major: breaking API, persistence, or architecture contract change.

## Vertical Slice Log

| Date | VS | Branch | Tag | OpenSpec | Summary | Status |
| --- | --- | --- | --- | --- | --- | --- |
| TBD | `architecture-sdd-foundation` | current | TBD | `architecture/openspec` | Establish SDD, OpenSpec workspace, roadmap, and project agents. | In progress |

## Planned Milestones

| Milestone | Goal | Notes |
| --- | --- | --- |
| Architecture foundation | Establish SDD, OpenSpec, roadmap, and agent workflows. | Current setup. |
| Auth hardening | Improve auth contracts, token strategy, tests, and error handling. | Future VS. |
| Posts MVP stabilization | Align posts, comments, votes, and API contracts. | Future VS. |
| Contributor-ready testing | Add clear service, repository, and API test standards. | Future VS. |
