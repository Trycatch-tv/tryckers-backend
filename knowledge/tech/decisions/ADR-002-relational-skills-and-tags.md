---
type: adr
status: accepted
date: 2026-08-04
created_from: knowledge/tech/discovery/decision-candidates.md
governed_code:
  - src/internal/models/post.go
  - src/internal/models/user.go
---

# ADR-002 — Normalización Relacional de Catálogo de Habilidades y Tags

## Context
Actualmente las habilidades de los usuarios y los tags de publicaciones se almacenan como cadenas de texto plano separadas por comas (`Interests string`, `Tags string`). Esto dificulta realizar consultas estructuradas, búsqueda facetada precisa, prevención de sinónimos/errores tipográficos y métricas de popularidad tecnológica.

## Options Considered
1. **Mantener campos de texto plano con índices trigram (`pg_trgm`):** Rápido de implementar pero propenso a inconsistencias y nula gobernanza de taxonomía.
2. **Utilizar arrays nativos de PostgreSQL (`text[]`):** Soporta operadores de contención `@>`, pero no provee un catálogo controlado con metadata (categoría, icono, slug).
3. **Crear tablas relacionales normalizadas (`skills`, `tags`, `user_skills`, `post_tags`):** Modelo relacional estándar con integridad referencial, unicidad de slugs y metadata descriptiva.

## Decision
Se adopta la creación de entidades relacionales normalizadas para habilidades y tags:
- Tablas maestras `skills` y `tags` con campos `id`, `name`, `slug`, `category` y timestamps.
- Tablas intermedias de asociación `user_skills` (con nivel de dominio opcional) y `post_tags`.
- Se preservan métodos de compatibilidad (`GetTagsSlice`, `SetTagsFromSlice`) durante el periodo de transición.

## Consequences
- **Positivas:**
  - Búsqueda exacta y filtrado facetado de alta eficiencia en el directorio de talento.
  - Estandarización de nombres tecnológicos (e.g. `TypeScript`, `Node.js`, `Go`).
  - Capacidad de calcular tendencias de tecnologías más solicitadas y publicadas.
- **Negativas / Mitigaciones:**
  - Requiere migración de datos para convertir los strings actuales en registros normalizados.
  - Mayor número de joins en consultas, mitigable con eager loading (`Preload`) selectivo en GORM.

## Related Capabilities
- Professional Identity and Profiles / Professional Attributes Representation
- Community Discussions and Content / Post Publication and Media Tagging

## Related Work Items
- [open] WI-CORE-002
