---
type: bugfix
id: WI-009
title: >-
  Corregir error 404 al cargar el perfil por sensibilidad a mayúsculas en el
  username
status: completed
work_type: bugfix
knowledge_level: K2
created_at: '2026-11-26'
completed_at: '2026-11-26'
source:
  type: external
  imported_at: '2026-11-26'
  source_format: markdown
  source_hash: 15ea4541ab7fa7517afc9d989d6feb7647a0d81849bbc10760f7ace490730f43
  inferred: false
generated_by: kaddo-admin
affected_modules:
  - backend
domains:
  - profiles
  - iam
  - backend
related_domain: profiles
related_capabilities:
  - User Registration and Password Hashing
code:
  - src/internal/repository/user_repository.go
  - src/internal/services/user_service.go
  - src/internal/repository/user_repository_test.go
  - src/internal/services/user_service_test.go
summary: >-
  Tras iniciar sesión, al abrir el menú de perfil la aplicación muestra "No se
  pudo cargar el perfil. Inténtalo nuevamente." y la consola del navegador
  registra un 404:
telemetry_work_item_id: 14967184-9bc5-4d6f-8a19-5934a3b8def4
original_snapshot:
  title: >-
    Corregir error 404 al cargar el perfil por sensibilidad a mayúsculas en el
    username
  description: >-
    Tras iniciar sesión, al abrir el menú de perfil la aplicación muestra "No se
    pudo cargar el perfil. Inténtalo nuevamente." y la consola del navegador
    registra un 404:
ready_at: '2026-11-26'
---

# Corregir error 404 al cargar el perfil por sensibilidad a mayúsculas en el username

## Actor y resultado

- **Actor:** usuario autenticado que abre su propio perfil (o el de otro) desde el menú.
- **Resultado observable:** el perfil carga correctamente (HTTP 200 con los datos del usuario), sin importar las mayúsculas del username.

## Comportamiento actual (bug)

Tras iniciar sesión, al abrir el menú de perfil la aplicación muestra "No se pudo cargar el perfil. Inténtalo nuevamente." y la consola del navegador registra un 404:

`GET https://api.tryckers.trycatch.tv/api/v1/perfil/Deiby7h322 404 (Not Found)` junto con "Error loading profile data".

## Comportamiento objetivo

`GET /api/v1/perfil/{username}` resuelve al usuario de forma insensible a mayúsculas, y los nuevos registros persisten el username en forma canónica (minúsculas, sin espacios al borde).

## Journey (reconstrucción del flujo)

1. **Entry point:** el frontend llama `GET /api/v1/perfil/{username}` (servicio `getTryckerByUsername` → `getProfileData` del profile-page).
2. **API/handler:** `UserHandler.Perfil` normaliza el parámetro con `strings.ToLower(c.Param("username"))`.
3. **Service:** `UserService.Perfil` delega en el repositorio.
4. **Persistencia:** `UserRepository.FindByUsername` consulta la base.
5. **Response:** 404 si no hay coincidencia exacta → el frontend muestra el error y el toast "No Encontrado".

## Surface review

- **Backend (API/handler/service/repo):** afectado — origen del defecto y de la corrección.
- **Database:** reviewed-not-affected — no se altera el esquema ni se migran datos; la lectura case-insensitive recupera los registros existentes tal como están.
- **Authentication/authorization:** reviewed-not-affected — el endpoint sigue protegido por `AuthMiddleware`; no cambia el modelo de permisos.
- **Frontend / UI / copy:** reviewed-not-affected — no requiere cambios para resolver el defecto (el front ya normaliza a minúsculas antes de llamar).
- **Configuration / feature flags / notifications / analytics / release:** not-applicable.
- **Documentation:** reviewed-not-affected — el contrato del endpoint no cambia (misma ruta, mismos códigos).

## Causa raíz

El endpoint `GET /api/v1/perfil/:username` nunca encontraba al usuario cuando su username contenía mayúsculas. El handler `UserHandler.Perfil` normalizaba el parámetro recibido con `strings.ToLower(c.Param("username"))`, pero:

- El registro (`UserService.CreateUser`) persistía el username tal cual lo escribía la persona (`Username: user.Username`), sin normalizar.
- El repositorio (`UserRepository.FindByUsername`) buscaba con coincidencia exacta: `WHERE username = ?`.

Resultado: un username guardado como `Deiby7h322` nunca coincidía con la búsqueda en minúsculas `deiby7h322`, devolviendo 404. Reproducido en entorno local con el usuario `TestUser123`: tanto `/perfil/TestUser123` como `/perfil/testuser123` devolvían 404.

## Solución (opción C: lectura case-insensitive + normalización en el registro)

1. Lectura case-insensitive en el repositorio: `FindByUsername` compara con `WHERE LOWER(username) = LOWER(?)`. Esto recupera también a los usuarios ya existentes persistidos con mayúsculas, sin necesidad de migrar datos.
2. Normalización en el registro: `CreateUser` normaliza el username con `strings.ToLower(strings.TrimSpace(user.Username))` antes de persistir, de modo que los nuevos registros queden canónicos y la unicidad no admita duplicados que difieran solo en mayúsculas.

Decisión de diseño: la normalización del username vive en el service (`CreateUser`), no en el handler, porque es una regla de persistencia canónica que debe aplicar venga de donde venga la llamada. (El país sí se normaliza en el handler por el WI-007 previo.)

## Criterios de aceptación

- `GET /api/v1/perfil/{username}` devuelve 200 y el perfil tanto para un username guardado con mayúsculas (`TestUser123`, caso `Deiby7h322`) como para el mismo en minúsculas (`testuser123`).
- Un username inexistente sigue devolviendo 404.
- Un registro nuevo con mayúsculas y espacios al borde (p. ej. `"  NuevoCaps  "`) se persiste normalizado (`nuevocaps`).
- No se rompen los comportamientos ya cubiertos por tests (validación de país en registro y actualización de perfil).

## Validación (cómo probarlo)

- `go vet ./...` limpio, `go build ./...` OK, `go test ./...` toda la suite PASS.
- Tests nuevos: `TestUserRepository_FindByUsername_IsCaseInsensitive` (DryRun exige `LOWER` en el SQL) y `TestUserService_CreateUser_NormalizesUsername` (verifica el username persistido normalizado).
- Validación end-to-end contra el backend local: login con `testuser@local.dev`, luego `GET /api/v1/perfil/TestUser123` → 200, `GET /api/v1/perfil/testuser123` → 200, `GET /api/v1/perfil/nadie999` → 404; registro con `"  NuevoCaps  "` persiste `nuevocaps`.

## Archivos afectados

- `src/internal/repository/user_repository.go`
- `src/internal/services/user_service.go`
- `src/internal/repository/user_repository_test.go`
- `src/internal/services/user_service_test.go`

## Fuera de alcance

- No se migran datos existentes (no hace falta: la lectura case-insensitive ya los recupera).
- No se normalizan retroactivamente los usernames ya persistidos con mayúsculas.
- No se cambia el método de login (sigue por email + password) ni el modelo de permisos.
- No se toca el frontend.

## Definición de hecho (Definition of Done)

- Tests nuevos en verde y suite completa PASS (`go test ./...`).
- `go vet ./...` limpio y `go build ./...` OK.
- Verificación end-to-end contra el backend local reproduciendo el caso `Deiby7h322`.
- WI-009 refinado y marcado ready; conocimiento alineado (sin drift pendiente atribuible a este cambio).

## Confianza de scope

- **Alta.** Causa raíz confirmada y reproducida; la corrección es local a dos funciones (repo + service) con tests que fijan el comportamiento. Riesgo controlado: toca el contrato de unicidad/lookup de username (alto riesgo por definición), mitigado con lectura `LOWER` idempotente y normalización en escritura, sin migración.

## Open questions / dependencias

- Despliegue a producción: el frontend de prod apunta a `api.tryckers.trycatch.tv`; el fix debe desplegarse al backend de prod para resolver el caso real `Deiby7h322`. (Decisión de entrega del usuario.)
- Sin dependencias de otros Work Items. Solapa ownership con WI-001 y WI-006 (mismos archivos de usuario) pero no altera el comportamiento que esos WI documentan.

## Dominios

profiles, iam, backend

## Aprendizaje

Cuando un identificador se normaliza al leer (ToLower en el handler), hay que normalizarlo también al escribir; si no, la escritura y la lectura divergen y el lookup exacto falla. La corrección combina lectura idempotente case-insensitive (`LOWER(username) = LOWER(?)`, que recupera datos antiguos sin migración) con normalización en el registro (persistir el username canónico). La normalización canónica vive en el servicio, no en el handler, para que aplique sin importar el punto de entrada.

---

## Revisión técnica, integración y refuerzo (PR #47 + Issue #42)

- **PR integrada:** https://github.com/Trycatch-tv/tryckers-backend/pull/47 (autor: `@Deibyth`, merge commit `7c23a10`).
- **Issue relacionado:** https://github.com/Trycatch-tv/tryckers-frontend/issues/42.
- **Auditoría RDS previa (SSM `7bc75b1b`):**
  - 13 usuarios históricos con mayúsculas (incluyendo `Deiby7h322`, `Calebdev`, `Charles `, `Kmilostk `, `Gutiérrez `).
  - **0 duplicados lógicos** por `LOWER(TRIM(username))`.
- **Refuerzo mínimo aplicado (AC-04):**
  - Migración idempotente `migrations/002_username_case_insensitive_unique_index.sql` y `EnsureUsernameCIUniqueIndex` en arranque (`src/cmd/main.go`): `CREATE UNIQUE INDEX IF NOT EXISTS idx_users_username_lower_unique ON users (LOWER(TRIM(username)))`.
  - `FindByUsername` actualizado a `LOWER(TRIM(username)) = LOWER(TRIM(?))` para cubrir también registros históricos con espacios al borde.
  - Mapeo de `idx_users_username_lower_unique` → `apperrors.ErrDuplicateUsername` en `UserService.CreateUser`.
- **Pruebas reales contra PostgreSQL (`user_repository_pg_test.go`):**
  - `TestPostgres_UsernameLookupNormalizationAndCaseInsensitiveUniqueness` → **PASS (0.51s)** probando búsqueda mixta/minúsculas/mayúsculas, espacios históricos, rechazo real `SQLSTATE 23505` al intentar duplicar un usuario histórico mixto, normalización en registro y 404 para usuario inexistente.
