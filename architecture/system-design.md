# SDD - Main Architecture Definition

# Tryckers Backend Architecture

## Objective

Define the foundational architecture for the Tryckers Backend, ensuring scalability, maintainability, contributor onboarding simplicity, and clean separation of concerns for a collaborative open source environment.

The architecture must support:

- Community-driven contributions.
- Rapid feature iteration.
- Clean domain separation.
- Reusable business logic.
- Future extensibility toward microservices or modular monolith evolution.

## Architectural Style

The backend follows a **Layered Modular Monolith Architecture** combined with:

- Clean separation of responsibilities.
- Domain-oriented modularization.
- Service + Repository pattern.
- Dependency Injection principles.
- DTO-based API contracts.

## Technology Stack

| Layer | Technology |
| --- | --- |
| API Framework | Gin |
| Language | Go |
| ORM | GORM |
| Database | PostgreSQL |
| Authentication | JWT |
| Password Encryption | bcrypt |
| Environment Config | godotenv |
| State Persistence | PostgreSQL |
| Containerization | Docker |
| API Communication | REST |

## Project Structure

```text
src/
|-- cmd/
|   `-- main.go
|-- internal/
|   |-- api/
|   |   |-- handlers/
|   |   |-- middleware/
|   |   `-- routes/
|   |-- config/
|   |-- dtos/
|   |-- enums/
|   |-- models/
|   |-- repository/
|   |-- services/
|   `-- utils/
|-- pkg/
`-- migrations/
```

## Architectural Principles

### Separation of Responsibilities

Handlers are responsible for receiving HTTP requests, request validation, DTO mapping, and HTTP responses.

Handlers must not access the database directly or implement business rules.

Services are responsible for business logic, domain rules, orchestration, and validation beyond transport-level validation.

Services should use repositories and utility functions, and should remain framework-agnostic when possible.

Repositories are responsible for data persistence, database queries, and ORM interactions.

Repositories must avoid business logic.

Models represent database entities, GORM mappings, and relationships.

DTOs define API contracts, request payloads, and response payloads. DTOs must avoid persistence logic.

Utils contain shared technical utilities such as hashing, token generation, date helpers, and generic reusable helpers.

## Authentication Architecture

Authentication uses:

- JWT access tokens.
- bcrypt password hashing.

Password hashing is centralized in `internal/utils/bcrypt.go`.

Business services consume shared password helpers instead of implementing hashing directly.

## Database Architecture

Database engine: PostgreSQL.

ORM: GORM.

Key architectural decisions:

- UUID primary keys.
- Automatic timestamps through GORM conventions.
- Enum-driven constrained fields.
- Soft deletion via status strategy for the MVP.

Database schema changes must be represented as migrations and should be applied through the Supabase MCP server when a Supabase project is involved.

## Domain Modules

| Module | Responsibility |
| --- | --- |
| Users | User management and authentication |
| Posts | Community posts and content |
| Comments | Post interaction |
| Votes | Engagement scoring |
| Views | Analytics and visibility tracking |

## API Strategy

API versioning uses `/api/v1`.

REST conventions:

```text
GET    /posts
POST   /posts
PUT    /posts/:id
DELETE /posts/:id
```

## Scalability Strategy

Current stage:

- Modular Monolith.

Future evolution possibilities:

- Extract modules into microservices.
- Event-driven architecture.
- CQRS for high-scale modules.
- Search indexing layer.
- Recommendation engine.

## Testing Strategy

| Layer | Strategy |
| --- | --- |
| Services | Unit Tests |
| Repository | Integration Tests |
| API | End-to-End Tests |
| Utils | Unit Tests |

## Infrastructure Strategy

Environment consistency through:

- Docker.
- Docker Compose.

Planned future:

- CI/CD via GitHub Actions.
- Automated migrations.
- Preview environments.

## Open Source Contribution Philosophy

The architecture prioritizes:

- Readability.
- Predictability.
- Low onboarding friction.
- Explicit module boundaries.

Every module should be independently understandable, easy to contribute to, and minimally coupled.

## Future Architectural Directions

Potential future additions:

- Signal-based realtime notifications.
- Recommendation system.
- AI-assisted talent matching.
- Recruiter dashboards.
- Event sourcing for activity feeds.
- Elasticsearch integration.
- GraphQL gateway.

## Architectural Goal

Tryckers aims to become a collaborative open source ecosystem where developers improve both technical and soft skills while building a real-world scalable platform together.

The architecture must optimize learning, collaboration, maintainability, and real production practices.
