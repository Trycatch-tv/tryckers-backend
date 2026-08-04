---
type: adr
status: accepted
date: 2026-08-04
created_from: knowledge/tech/discovery/decision-candidates.md
governed_code:
  - src/internal/services/user_service.go
  - src/internal/api/handlers/user_handler.go
  - src/internal/models/user.go
---

# ADR-003 — Rotación Segura y Persistencia de Refresh Tokens

## Context
La API dispone de un endpoint `/api/v1/refresh-token` pero los tokens no se persisten en base de datos. Sin persistencia, no es posible revocar sesiones comprometidas, detectar reutilización de tokens robados ni cerrar sesiones de forma centralizada en todos los dispositivos.

## Options Considered
1. **Refresh tokens puramente JWT stateless sin persistencia:** Imposibilita la revocación inmediata antes de que expire el token de actualización.
2. **Almacenamiento en Redis con TTL automático:** Alta velocidad pero añade una dependencia de infraestructura adicional en etapas tempranas.
3. **Persistencia relacional con hash SHA-256 en PostgreSQL:** Almacena tokens hasheados en tabla `refresh_tokens` con control de expiración, rotación obligatoria e invalidación de sesiones.

## Decision
Se adopta la estrategia de persistencia relacional con rotación estricta de Refresh Tokens (*Refresh Token Rotation - RTR*):
- Se crea la entidad `RefreshToken` con campos: `id` (UUID), `user_id` (FK), `token_hash` (SHA-256), `expires_at`, `revoked_at`, `created_at` y `user_agent`.
- Cada vez que se utiliza un refresh token, se invalida el anterior y se emite un nuevo par de access token + refresh token.
- Si se detecta un intento de uso de un token ya revocado, se revocan preventivamente todas las sesiones activas del usuario (*Reuse Detection*).

## Consequences
- **Positivas:**
  - Máxima seguridad de autenticación y cumplimiento de estándares OWASP.
  - Capacidad de listar y cerrar sesiones activas remotas desde la configuración de usuario.
- **Negativas / Mitigaciones:**
  - Requiere una consulta SQL indexada (`SELECT ... WHERE token_hash = ?`) por cada renovación de sesión.
  - Requiere tarea periódica o cron job para purgar tokens expirados y revocados antiguos.

## Related Capabilities
- Identity and Access Management (IAM) / Refresh Token Rotation and Session Renewal

## Related Work Items
- [open] WI-AUTH-003
