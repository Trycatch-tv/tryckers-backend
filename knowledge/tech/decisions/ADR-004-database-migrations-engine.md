---
type: adr
status: accepted
date: 2026-08-04
created_from: knowledge/tech/discovery/decision-candidates.md
governed_code:
  - src/internal/config/database.go
  - migrations/**
---

# ADR-004 — Motor de Migraciones Versionadas para Base de Datos

## Context
Actualmente la aplicación ejecuta `db.AutoMigrate(...)` al arrancar el servidor en `src/cmd/main.go`. Aunque `AutoMigrate` es conveniente en fases iniciales de prototipado, no es determinista, no soporta eliminación de columnas/índices, no permite rollbacks automáticos y presenta riesgos de carrera en despliegues con múltiples réplicas concurrentes.

## Options Considered
1. **Continuar usando exclusivamente `GORM AutoMigrate`:** Inseguro para entornos productivos con datos reales y esquemas en evolución.
2. **Utilizar `Goose`:** Herramienta Go con soporte de archivos `.sql` y funciones Go.
3. **Utilizar `golang-migrate`:** Estándar de la industria en el ecosistema Go, altamente compatible con CI/CD, soporte de migraciones `.up.sql` y `.down.sql`, e integrable como binario CLI o librería Go.

## Decision
Se adopta `golang-migrate` como motor estándar de migraciones de base de datos para entornos productivos y de testing:
- Las migraciones residirán en `migrations/` con nomenclatura secuencial/timestamp (`000001_initial_schema.up.sql`, `000001_initial_schema.down.sql`).
- En desarrollo local, `AutoMigrate` se mantendrá configurable mediante la variable `ENABLE_AUTO_MIGRATE=true`.
- En pipelines de CI/CD y despliegues productivos, las migraciones se ejecutarán como un paso previo de release antes de iniciar los contenedores de la aplicación.

## Consequences
- **Positivas:**
  - Control de versiones determinista y reproducible en todos los entornos (dev, staging, prod).
  - Capacidad de revertir cambios (*down migrations*) de manera controlada.
  - Mayor seguridad y eliminación de condiciones de carrera durante el arranque.
- **Negativas / Mitigaciones:**
  - Requiere redactar explícitamente los scripts SQL de migración para cada cambio estructural de entidades.

## Related Capabilities
- Professional Identity and Profiles / Profile Schema Evolution
- Identity and Access Management (IAM) / Database Schema Migration

## Related Work Items
- [open] WI-INFRA-004
