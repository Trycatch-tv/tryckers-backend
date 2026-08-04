---
type: feature
id: WI-001
title: "Registro, Modificación y Visualización de Fecha de Cumpleaños"
knowledge_level: K2
status: completed
phase: now
work_type: feature
initiative: RM-001
related_domain: profiles
domains:
  - profiles
  - iam
related_capabilities:
  - "Professional Identity and Profiles / Profile Schema Evolution"
  - "Professional Identity and Profiles / Profile Editing and Customization"
related_decisions:
  - ADR-002-relational-skills-and-tags
code:
  - src/internal/models/user.go
  - src/internal/dtos/userUpdate.go
  - src/internal/api/handlers/user_handler.go
  - src/internal/api/routes/router.go
  - src/internal/repository/user_repository.go
  - src/internal/services/user_service.go
  - tryckers-frontend/src/app/auth/interfaces/auth-response.ts
  - tryckers-frontend/src/app/pages/settings-page/settings-page.component.ts
  - tryckers-frontend/src/app/tryckers/pages/profile-page/**
created_at: 2026-08-04
completed_at: 2026-08-04
source:
  type: manual
  inferred: false
project_state: ai-assisted
generated_by: kaddo-create
template_version: 1
refined_by: implementation-agent
summary: "Permitir al usuario autenticado registrar, actualizar y visualizar su fecha de cumpleaños tanto en la configuración de su cuenta como en la vista de perfil de usuario para potenciar la fidelización y personalización comunitaria."
affected_modules:
  - frontend
  - core
---

# WI-001 — Registro, Modificación y Visualización de Fecha de Cumpleaños

> Type: feature · Level: K2 · Status: completed

## Actor and outcome

El usuario miembro de la comunidad Tryckers puede registrar su fecha de nacimiento durante la edición de su perfil / configuración de cuenta, visualizar su fecha de cumpleaños en su perfil (o insignia/widget) y actualizarla de forma persistente y segura.

## Current behavior (Pre-implementation)

El modelo backend `User` contenía el campo `BirthDate *time.Time`, pero no existía un endpoint REST `PUT/PATCH /api/v1/users/me` para actualizar campos del perfil del usuario (incluyendo la fecha de cumpleaños). En el frontend, la interfaz `UserData` tipaba `birth_date: null`, la página de configuración (`settings-page`) era un placeholder estático y el componente de perfil (`profile-page`) no renderizaba ni permitía editar este dato.

## Target behavior (Implemented)

1. El backend expone un endpoint protegido `PUT /api/v1/users/me` que acepta un DTO de actualización (`UpdateProfileDTO`) validando que `birth_date` sea una fecha válida y coherente (no futura y edad mínima lógica).
2. El frontend actualiza su interfaz TypeScript `UserData` con `birth_date?: string | null` (formato ISO 8601 `YYYY-MM-DD`).
3. El frontend provee un formulario en `settings-page` con selector de fecha para guardar o modificar la fecha de cumpleaños y datos de perfil.
4. El perfil público renderiza la fecha de cumpleaños en español (ej. "15 de junio") en la tarjeta de estadísticas del perfil.

## Entry points

- Frontend: Ruta `/settings` (formulario de edición de perfil) y `/profile/:username` (vista del perfil).
- Backend: Endpoint protegido `PUT /api/v1/users/me`.

## End-to-end flow

1. El usuario autenticado navega a `/settings` o pulsa "Editar Perfil" en `/profile/:username`.
2. Selecciona su fecha de nacimiento en el selector de calendario (`YYYY-MM-DD`) y hace clic en "Guardar Cambios".
3. El frontend envía `PUT /api/v1/users/me` con `{ "birth_date": "1995-06-15T00:00:00Z", ... }` y el header `Authorization: Bearer <token>`.
4. El backend valida el token, parsea la fecha, valida el rango cronológico, persiste el campo en PostgreSQL mediante GORM y responde con `200 OK` y el `UserData` actualizado.
5. El frontend recibe la respuesta, actualiza el `AuthStore` (`setUser`), sincroniza `localStorage` y muestra notificación Toast de éxito.
6. Al navegar a `/profile/:username`, la fecha de cumpleaños se visualiza formateada en la interfaz comunitaria.

## Problem

Los usuarios no disponían de un mecanismo en la interfaz ni en los endpoints de mutación de perfil para registrar o actualizar su fecha de nacimiento, impidiendo iniciativas comunitarias de felicitación, dinámicas de fidelización y gamificación en fechas especiales.

## Expected result

Flujo extremo a extremo funcional para registrar, editar, almacenar en base de datos y renderizar en la interfaz la fecha de cumpleaños del usuario.

## Impact

Impacto de no realizarlo: pérdida de oportunidades de fidelización (mensajes de cumpleaños, dinámicas especiales y reconocimiento en comunidad).

### Impact analysis (surfaces)
- product/UI: `affected` (Adición de selector de fecha en formulario de edición y visualización en perfil).
- frontend: `affected` (`UserData` interface, `settings-page`, `profile-page`, servicios de usuario).
- backend: `affected` (Creación de `UpdateProfileDTO`, handler `UpdateMe`, servicio y repositorio de actualización).
- database: `reviewed-not-affected` (Columna `birth_date` ya existe en tabla `users`).
- configuration: `reviewed-not-affected`.
- feature flags: `not-applicable`.
- content/copy: `affected` (Etiquetas y textos de fecha de nacimiento y mensajes de validación en español).
- authentication/authorization: `reviewed-not-affected` (Utiliza el middleware `AuthMiddleware` existente).
- notifications: `reviewed-not-affected`.
- analytics: `reviewed-not-affected`.
- documentation: `affected` (Anotaciones Swagger / OpenAPI en backend para el nuevo endpoint).
- operations/release: `reviewed-not-affected`.

## Module coverage

- `core` (`tryckers-backend`): `affected` (Endpoints, DTOs, Handler, Service, Repository).
- `frontend` (`tryckers-frontend`): `affected` (Interfaces, Componentes UI, SignalStore).

## Scope confidence

High. El campo ya está modelado en base de datos en PostgreSQL; se completó el ciclo de mutación REST y los componentes visuales de frontend.

## Acceptance criteria

- [x] El backend expone `PUT /api/v1/users/me` protegido por `AuthMiddleware`.
- [x] La solicitud `PUT /api/v1/users/me` acepta `{ "birth_date": "YYYY-MM-DD" }` (o ISO string) y actualiza el campo en la base de datos PostgreSQL.
- [x] La validación en backend rechaza fechas futuras (`> time.Now()`) o formatos inválidos con código HTTP 400 Bad Request.
- [x] La interfaz TypeScript `UserData` en `tryckers-frontend` tipa `birth_date?: string | null`.
- [x] En `tryckers-frontend`, el usuario puede ingresar/modificar su fecha en `/settings` (o modal de perfil) mediante un selector de fecha con validación en cliente.
- [x] Al guardar exitosamente, la sesión reactiva (`AuthStore`) se actualiza sin requerir recargar la página.
- [x] En `/profile/:username`, la fecha de cumpleaños se muestra formateada en español (ej. "15 de junio").
- [x] *(Criterio E2E)*: Un usuario autenticado modifica su fecha de cumpleaños en `/settings`, guarda los cambios, navega a su perfil y verifica que la nueva fecha se refleja inmediatamente y persiste tras recargar la página.

## Out of scope

- Automatización de correos electrónicos de felicitación o notificaciones push programadas.
- Configuración de privacidad granular para ocultar el año de nacimiento (se mostrará día y mes por defecto en vista pública).

## Validation

### 1. Prueba de API Backend
- Enviar `PUT /api/v1/users/me` con body `{"birth_date": "1992-08-15T00:00:00Z"}` y token Bearer válido.
- Verificar respuesta `200 OK` con campo `birth_date` poblado.
- Enviar `PUT /api/v1/users/me` con fecha futura (ej. "2099-01-01") y verificar respuesta `400 Bad Request`.

### 2. Prueba de UI Frontend
- Iniciar sesión en `http://localhost:4200`.
- Navegar a `/settings` (o modal de edición de perfil), seleccionar fecha de cumpleaños y pulsar "Guardar".
- Confirmar aparición del toast de éxito.
- Navegar a `/profile/<username>` y verificar el renderizado del cumpleaños.

### 3. Prueba de Persistencia
- Recargar la página (`F5`) y verificar que la fecha persiste en `localStorage` y en la respuesta de `GET /api/v1/perfil/:username`.

## Definition of Done

- [x] El problema está claramente formulado en una sola oración.
- [x] El resultado esperado y el flujo de usuario extremo a extremo están definidos.
- [x] El impacto de no realizarlo y el análisis de superficies y módulos están documentados.
- [x] Los criterios de aceptación son verificables e incluyen prueba E2E.
- [x] La guía de validación detalla pasos de prueba de API y frontend.

## Open questions

- [resolved] ¿Debe mostrarse el año de nacimiento en el perfil público? → *Decisión: Se mostrará únicamente día y mes ("15 de Junio") en perfiles públicos para proteger la privacidad del usuario, reservando la fecha completa para la configuración privada.*

## Learnings and Implementation Record

- **Implemented:**
  - Backend: DTO `UpdateProfileDTO`, endpoint `PUT /api/v1/users/me`, servicio con validación cronológica y repositorio de actualización parcial en PostgreSQL.
  - Frontend: Interfaz `UserData` actualizada, servicio `updateProfile`, página reactiva `settings-page` con selector de fechas y renderizado en tarjeta de estadísticas de `profile-page`.
- **Changed / Discovered:**
  - Convención de columnas en GORM: el campo `EFSetScore` requirió tag explícito `gorm:"column:ef_set_score"` para evitar colisiones entre el nombre de struct en camelCase y el schema en base de datos PostgreSQL.
- **Decisions Emerged:**
  - ADR-002: Mostrar únicamente día y mes en el perfil público para proteger la privacidad del usuario, reservando la fecha completa para la configuración privada.
- **Knowledge Updated:**
  - Work Item WI-001 completado en `knowledge/delivery/work-items/completed/`.
  - Context pack regenerado con `npx kaddo context`.
