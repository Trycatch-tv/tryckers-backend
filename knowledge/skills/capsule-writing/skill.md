---
type: skill
id: capsule-writing
name: capsule-writing
title: Capsule Writing Skill
version: 3.135.0
group: integration
applies_to:
  - capsule-agent
  - architecture-agent
  - product-agent
---
# Capsule Writing Skill

## Purpose

Standardize how a Knowledge Capsule is written/reviewed so external consumers get safe, useful
context.

## When to use

When creating or refining a Knowledge Capsule for sharing with another project.

## Inputs

The context pack, capabilities, current-state, decisions and any public contracts.

## Output

A capsule with: purpose, responsibilities, capabilities, contracts, dependencies, risks, owners,
out of scope and usage notes.

## Rules

- Never include secrets, tokens, credentials, source code, PII or unnecessary internal detail.
- Never invent contracts; mark unknowns.
- Summarize boundaries; prefer "unknown" over guessing.

## Quality checklist

- Purpose and boundaries are clear.
- Contracts are real, not invented.
- No secrets/source/PII included.

## Example output

A `*.capsule.md` with the sections above.
