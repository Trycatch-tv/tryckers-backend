---
type: bugfix
id: WI-007
title: >-
  Corregir obligatoriedad y normalización de país en registro y actualización de
  perfil
status: completed
work_type: bugfix
knowledge_level: K2
created_at: '2026-11-26'
completed_at: '2026-11-26'
source:
  type: chat
  imported_at: '2026-11-26'
  source_format: markdown
  inferred: false
generated_by: kaddo-admin
refined_by: implementation-agent
project_state: ai-assisted
affected_modules:
  - frontend
domains:
  - iam
  - profiles
  - frontend
  - backend
related_domain: profiles
related_capabilities:
  - User Registration and Password Hashing
  - Professional Attributes Representation
affected_system_entities:
  - 'sys:tryckers-frontend'
  - 'sys:tryckers-backend'
reviewed_system_entities:
  - entity: 'sys:postgres-db'
    status: reviewed-not-affected
    reason: >-
      Auditoría en vivo en RDS (SSM 16c4eed0-0c8b-46d0-892d-36502ef48628)
      confirmó que el 100% de usuarios existentes tienen nombres canónicos
      válidos (colombia=10, republica dominicana=1, bolivia=1, honduras=1,
      ecuador=1); sin registros vacíos, nulos ni códigos ISO, por lo que no se
      requiere migración de datos.
  - entity: >-
      code:tryckers-frontend/src/app/pages/settings-page/settings-page.component.ts
    status: affected
    reason: >-
      Usa códigos ISO (CO, AR, MX, US, CA), asigna 'CO' por defecto si falta
      país y muestra mensaje genérico en errores de validación.
  - entity: >-
      code:tryckers-frontend/src/app/auth/pages/register-page/register-page.component.ts
    status: affected
    reason: >-
      Debe compartir el catálogo canónico único de países y validar
      explícitamente contra la lista soportada por el backend.
  - entity: 'code:src/internal/enums/country.go'
    status: reviewed-not-affected
    reason: >-
      Catálogo canónico de 22 países normalizados en minúsculas ya alineado con
      los datos persistidos en RDS.
code:
  - ../tryckers-frontend/src/app/shared/constants/countries.ts
  - >-
    ../tryckers-frontend/src/app/auth/pages/register-page/register-page.component.ts
  - >-
    ../tryckers-frontend/src/app/auth/pages/register-page/register-page.component.html
  - ../tryckers-frontend/src/app/pages/settings-page/settings-page.component.ts
  - src/internal/enums/country.go
  - src/internal/dtos/userCreate.go
  - src/internal/dtos/userUpdate.go
  - src/internal/api/handlers/user_handler.go
  - src/internal/services/user_service.go
summary: >-
  Garantizar la obligatoriedad y representación canónica única de país entre
  tryckers-frontend y tryckers-backend, eliminando el valor predeterminado
  inválido 'CO' en configuración de perfil y permitiendo actualizaciones
  parciales sin bloquear por país ausente.
ready_at: '2026-11-26'
---

# WI-007 — Corregir obligatoriedad y normalización de país en registro y actualización de perfil

> Idioma del proyecto: **español**.

**Tipo:** Bugfix  
**Prioridad:** High  
**Estado:** Completed  
**Repositorios / Módulos:** `tryckers-backend` (`core`), `tryckers-frontend` (`frontend`)  
**Áreas funcionales:** Authentication, Profile Management, Validation

---

## Actor and outcome

- **Actor:** Usuario de Tryckers (nuevo usuario registrándose y usuario autenticado editando su perfil en `/settings`).
- **Outcome:** El registro exige un país canónico válido con mensaje específico junto al campo, y en `/settings` el usuario visualiza correctamente su país persistido (o vacío si no tiene), puede editar biografía u otros campos sin que un país ausente bloquee el guardado, y puede seleccionar o actualizar su país usando los mismos valores canónicos aceptados por el backend.

---

## Problem

`SettingsPageComponent` utiliza códigos ISO de dos letras (`CO`, `AR`, `MX`, incluyendo `US` y `CA` no soportados) y asigna `'CO'` por defecto cuando no existe país, provocando que las actualizaciones de perfil (como biografía) fallen con `400 Bad Request` (`"país inválido"`) contra `tryckers-backend` que valida nombres canónicos en minúsculas (`colombia`, `argentina`, `mexico`, etc.).

---

## Current behavior

1. **Backend (`src/internal/enums/country.go`, `dtos/userCreate.go`, `handlers/user_handler.go`, `services/user_service.go`):**
   - Acepta 22 países canónicos en minúsculas (`argentina` .. `venezuela`).
   - `POST /api/v1/register` exige `country` (`binding:"required"`) y valida `enums.IsValidCountry`.
   - `PUT /api/v1/users/me` valida `enums.IsValidCountry` siempre que `req.Country != nil`.
2. **Frontend — Registro (`register-page.component.ts`, `register-page.component.html`):**
   - Tiene `Validators.required` y `minLength(2)`, pero las opciones están hardcodeadas en HTML separadas del catálogo de configuración.
3. **Frontend — Configuración (`settings-page.component.ts`):**
   - Define un catálogo separado con códigos ISO (`AR`, `CO`, `MX`, `US`, `CA`) y vincula `[value]="c.code"`.
   - Asigna `country: 'CO'` por defecto en el estado inicial y `user.country || 'CO'` en `populateForm()`.
   - Siempre envía `country` en el payload de `PUT /api/v1/users/me`, causando rechazo del backend (`país inválido`) y mostrando el mensaje genérico `"No se pudo actualizar el perfil. Revisa los datos ingresados."`.
4. **Auditoría de datos existentes en RDS (SSM `16c4eed0-0c8b-46d0-892d-36502ef48628`):**
   - Distribución actual en `users`: `colombia` (10), `republica dominicana` (1), `bolivia` (1), `honduras` (1), `ecuador` (1).
   - 0 registros vacíos, nulos o con códigos ISO. No se requiere migración de datos.

---

## Target behavior

1. **Catálogo canónico único en frontend (`src/app/shared/constants/countries.ts`):** Lista única de los 22 países soportados por `src/internal/enums/country.go` (`value` canónico, `label`, `code`, `flag`) + helpers `isValidCanonicalCountry()` y `normalizeCountryValue()`.
2. **Registro:** Usa el catálogo único, exige selección obligatoria y válida de país, impide el envío si falta o es inválido, y muestra mensaje específico junto al campo.
3. **Configuración de perfil (`/settings`):**
   - Usa el mismo catálogo canónico (`[value]="c.value"`), eliminando `US`/`CA` no soportados y el valor por defecto `'CO'`.
   - Si el usuario no tiene país (o tiene uno incompatible), el campo queda vacío (`Sin especificar`) sin asignar país arbitrariamente.
   - Al guardar, omite `country` del payload si no hay un valor canónico seleccionado, permitiendo actualizar `bio`, `headline` y demás campos sin error.
   - Si el usuario selecciona o cambia un país, envía el valor canónico válido.
   - Muestra mensajes de validación específicos del campo en lugar de un mensaje genérico.
4. **Backend:** Conserva validación estricta en registro y actualización cuando `country` se envía explícitamente, normalizando espacios/mayúsculas y agregando cobertura de tests unitarios.

---

## Entry points

- Frontend: `/auth/register` (`RegisterPageComponent`), `/settings` (`SettingsPageComponent`)
- Backend: `POST /api/v1/register` (`UserHandler.CreateUser`), `PUT /api/v1/users/me` (`UserHandler.UpdateMe` → `UserService.UpdateProfile`)

---

## End-to-end flow

1. Usuario en `/auth/register` intenta enviar sin país → validación cliente bloquea y muestra `"El pais es obligatorio."` junto al select; si selecciona un país canónico (e.g. `colombia`), `POST /api/v1/register` valida y crea el usuario (`201 Created`).
2. Usuario autenticado abre `/settings` → `populateForm()` carga `normalizeCountryValue(user.country)` (si tiene `'colombia'`, queda seleccionado `🇨🇴 Colombia`; si no tiene país, queda en `""` *Sin especificar*).
3. Usuario sin país edita `bio` y guarda → `saveProfile()` omite `country` del payload → `PUT /api/v1/users/me` actualiza solo los campos provistos (`200 OK`) preservando el resto del perfil.
4. Usuario selecciona o cambia país en `/settings` y guarda → `saveProfile()` envía `country: "mexico"` → backend valida y persiste (`200 OK`).

---

## Suggested Knowledge Level

**K2** — Bugfix multi-módulo sobre validación y sincronización de catálogo entre frontend y backend sin cambios de esquema ni infraestructura.

---

## Impact analysis

| Superficie | Estado | Evidencia / Motivo |
| :--- | :--- | :--- |
| Product / UI | `affected` | Selector unificado de país en registro y configuración; opción vacía en settings; mensajes de error específicos. |
| Frontend | `affected` | Nuevo catálogo compartido `countries.ts`, `RegisterPageComponent` y `SettingsPageComponent`. |
| Backend | `affected` | Normalización segura en `CreateUser`/`UpdateProfile`, interfaz de repositorio para tests y pruebas unitarias. |
| Database | `reviewed-not-affected` | Auditoría RDS confirmó 100% de registros válidos; no se requiere migración. |
| Documentation | `affected` | Registro de causa raíz y evidencias en `WI-007` en Kaddo. |

---

## Module coverage

- `frontend` (`../tryckers-frontend`): `affected`
- `core` (`tryckers-backend`): `affected`

---

## Scope unknowns

Ninguno bloqueante: la auditoría de RDS confirmó ausencia de datos corruptos o vacíos en producción.

---

## Scope confidence

**High** — Causa raíz verificada en código y datos reales de RDS.

---

## Acceptance criteria

- [x] **AC-01:** El registro exige seleccionar un país válido.
- [x] **AC-02:** El usuario recibe un mensaje específico cuando intenta registrarse sin país.
- [x] **AC-03:** El backend rechaza registros sin país o con valores inválidos.
- [x] **AC-04:** Registro y configuración utilizan valores de país compatibles con el backend.
- [x] **AC-05:** Un usuario con país válido puede actualizar su biografía sin errores relacionados con ese campo.
- [x] **AC-06:** Un usuario sin país puede actualizar campos independientes de su perfil sin recibir el error reportado.
- [x] **AC-07:** La configuración no asigna automáticamente un país que el usuario no seleccionó.
- [x] **AC-08:** El usuario puede seleccionar o cambiar su país y guardar correctamente.
- [x] **AC-09:** Las actualizaciones parciales no sobrescriben accidentalmente información existente.
- [x] **AC-10:** Los mensajes de validación permiten identificar el campo que ocasionó el error.
- [x] **AC-11:** Se incluyen pruebas automatizadas de registro y actualización para países válidos, inválidos y ausentes.
- [x] **AC-12:** Login, registro, actualización de perfil y persistencia de datos continúan funcionando sin regresiones. *(Confirmado en prueba funcional del usuario.)*
- [x] **AC-13:** Se documenta la causa raíz, la corrección aplicada y las evidencias de validación en Kaddo.

## Verified Evidence

- **Auditoría RDS (SSM `16c4eed0-0c8b-46d0-892d-36502ef48628`):**
  - 100% de usuarios existentes en `users` ya poseen nombres canónicos válidos (`colombia`: 10, `republica dominicana`: 1, `bolivia`: 1, `honduras`: 1, `ecuador`: 1). 0 registros vacíos, nulos o con códigos ISO. No se requiere migración de datos.
- **Backend (`tryckers-backend`):**
  - `go test -v ./...` → **PASS (15/15)**, incluyendo los 5 nuevos tests de `UserHandler` (`CreateUser` sin país, con código ISO `CO`, con país canónico `colombia`; `UpdateMe` actualizando `bio` sin país, rechazando `CO` y aceptando `mexico`).
  - `go build ./...` → **PASS**.
- **Frontend (`tryckers-frontend`):**
  - `npx ng build --configuration development` → **PASS**.
  - `npx ng test --watch=false --browsers=ChromeHeadless` (5 suites, incluyendo `register-page.component.spec.ts` y `settings-page.component.spec.ts`) → **TOTAL: 18 SUCCESS, 0 FAILED**.
- **Validación funcional manual (usuario):**
  - Confirmado por el usuario que el registro con país obligatorio y la actualización de perfil/biografía en `/settings` funcionan correctamente.

---

## Out of scope

- Crear una nueva tabla relacional de países en PostgreSQL.
- Agregar países fuera del catálogo actual del backend (`src/internal/enums/country.go`) sin solicitud explícita.
- Migraciones masivas de base de datos (descartada tras auditoría RDS).

---

## Validation

1. Backend: `go test -v ./...` y `go build ./...` (incluyendo pruebas de registro y `UpdateProfile` con país válido, inválido y ausente).
2. Frontend: `npx ng build --configuration development` y `npx ng test --watch=false --browsers=ChromeHeadless` (incluyendo specs de `RegisterPageComponent` y `SettingsPageComponent`).
3. Manual / E2E: registro sin país (bloqueado con mensaje), registro con país válido, edición de biografía en `/settings` con y sin país, y cambio de país.

---

## Definition of Done

Todos los criterios **AC-01** a **AC-13** verificados, pruebas y builds pasando en ambos repositorios y evidencia registrada en Kaddo.

---

## Open Questions

---

## Related domain / capability

- **Domains:** `Identity and Access Management (IAM)`, `Professional Identity and Profiles`
- **Capabilities:** `User Registration and Password Hashing`, `Professional Attributes Representation`
