# Glossary

## SDD

Software Design Documentation. The living architecture documentation for Tryckers.

## OpenSpec

The change-specification workspace under `architecture/openspec`. It records proposed changes before implementation.

## Vertical Slice

A complete, reviewable unit of product or architecture work that crosses all required layers for one capability.

## VS

Short name for vertical slice. VS names should be lowercase kebab-case, for example `auth-login-hardening`.

## Module

A cohesive domain or technical area with clear responsibilities and boundaries.

## Handler

Gin HTTP entry point responsible for validation, DTO mapping, and responses.

## Service

Business layer responsible for orchestration and domain rules.

## Repository

Persistence layer responsible for database access and ORM interaction.

## DTO

Data Transfer Object used as an API contract.

## ADR

Architecture Decision Record. A short document that records a significant architecture decision and its consequences.
