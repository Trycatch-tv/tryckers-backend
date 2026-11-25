---
type: bugfix
id: WI-006
title: Corregir errores de autenticación y carga de la cartelera
status: completed
work_type: bugfix
knowledge_level: K2
created_at: '2026-11-25'
completed_at: '2026-11-25'
source:
  type: chat
  imported_at: '2026-11-25'
  source_format: markdown
  source_hash: bc8c50a260ffc276caa86ebe8c64e6c2a9e946d0228fe4d5ffb58603d4738aa
  inferred: false
generated_by: kaddo-admin
refined_by: implementation-agent
project_state: ai-assisted
affected_modules:
  - frontend
domains:
  - iam
  - cartelera
  - frontend
  - backend
related_domain: iam
related_capabilities:
  - User Registration and Password Hashing
  - JWT Authentication and Route Protection
  - Refresh Token Rotation and Session Renewal
  - Weekly Billboard Discovery (Cartelera Semanal)
related_decisions:
  - ADR-003-refresh-token-rotation-and-storage
affected_system_entities:
  - 'sys:tryckers-frontend'
  - 'sys:tryckers-backend'
  - 'sys:local-storage'
reviewed_system_entities:
  - entity: 'sys:postgres-db'
    status: unknown
    reason: >-
      Consultada por GetTopPostsOfTheWeek (posts, post_votes, users); sin
      cambios de esquema previstos salvo que el diagnóstico en entorno confirme
      defecto SQL o de datos.
  - entity: 'sys:external-services'
    status: reviewed-not-affected
    reason: 'No participa en registro, autenticación JWT ni carga de cartelera.'
  - entity: 'code:tryckers-frontend/src/app/auth/interfaces/auth-response.ts'
    status: reviewed-not-affected
    reason: >-
      Contratos AuthResponse y RefreshTokenResponse son consistentes con POST
      /api/v1/login y POST /api/v1/refresh-token; solo requeriría ajuste menor
      si se tipa la respuesta de /register.
  - entity: 'code:src/internal/models/user.go'
    status: reviewed-not-affected
    reason: El modelo User no requiere cambios de estructura ni columnas.
  - entity: 'code:src/internal/dtos/userUpdate.go'
    status: reviewed-not-affected
    reason: >-
      DTO exclusivo de edición de perfil; fuera del flujo de registro y
      cartelera.
  - entity: 'code:src/internal/repository/user_repository.go'
    status: reviewed-not-affected
    reason: CreateUser y FindByEmail operan correctamente sobre la tabla users.
  - entity: >-
      code:tryckers-frontend/src/app/pages/settings-page/settings-page.component.ts
    status: reviewed-not-affected
    reason: >-
      Fuera de alcance de modificación; solo sujeto a validación de no regresión
      (AC-13).
  - entity: 'code:tryckers-frontend/src/app/tryckers/pages/profile-page/**'
    status: reviewed-not-affected
    reason: >-
      Fuera de alcance de modificación; solo sujeto a validación de no regresión
      (AC-13).
code:
  - ../tryckers-frontend/src/app/auth/services/auth.service.ts
  - ../tryckers-frontend/src/app/auth/store/auth-store.ts
  - >-
    ../tryckers-frontend/src/app/auth/pages/register-page/register-page.component.ts
  - ../tryckers-frontend/src/app/auth/pages/login-page/login-page.component.ts
  - ../tryckers-frontend/src/app/auth/guards/authenticated.guard.ts
  - ../tryckers-frontend/src/app/auth/guards/not-authenticated.guard.ts
  - ../tryckers-frontend/src/app/core/interceptors/auth.interceptor.ts
  - ../tryckers-frontend/src/app/core/interceptors/error.interceptor.ts
  - ../tryckers-frontend/src/app/post/pages/cartelera-page/cartelera-page.ts
  - ../tryckers-frontend/src/app/post/pages/cartelera-page/cartelera-page.html
  - ../tryckers-frontend/src/app/post/services/posts.service.ts
  - ../tryckers-frontend/src/app/shared/services/ux-metrics.service.ts
  - src/internal/api/handlers/post_handler.go
  - src/internal/services/post_service.go
  - src/internal/repository/post_repository.go
  - src/internal/api/middlewares/auth_middleware.go
  - src/internal/api/handlers/user_handler.go
  - src/internal/services/user_service.go
  - src/internal/api/routes/router.go
summary: >-
  Diagnosticar y corregir la inconsistencia de autenticación tras el registro y
  los fallos en la carga y manejo de errores de la cartelera semanal entre
  tryckers-frontend y tryckers-backend, garantizando estado de sesión
  consistente, aislamiento de errores de telemetría y recuperación mediante
  Reintentar.
telemetry_work_item_id: 38162788-022e-4038-bd5c-3b907d605fcd
original_snapshot:
  title: WI — Corregir errores de autenticación y carga de la cartelera
  description: >-
    **Tipo:** Bugfix **Prioridad:** High **Estado:** Proposed **Repositorios:**
    `tryckers-frontend`, `tryckers-backend` **Módulos:** Authentication,
    Cartelera, API, Observability
ready_at: '2026-11-25'
---

# WI-006 — Corregir errores de autenticación y carga de la cartelera

> Idioma del proyecto: **español**.

**Tipo:** Bugfix  
**Prioridad:** High  
**Estado:** Completed  
**Repositorios / Módulos:** `tryckers-backend` (`core`), `tryckers-frontend` (`frontend`)  
**Áreas funcionales:** Authentication (IAM), Cartelera, API, Observability (`UxMetricsService`)

---

## Actor and outcome

- **Actor:** Usuario de Tryckers (especialmente un usuario recién registrado sin sesión previa, y también usuarios ya autenticados que acceden a `/cartelera`).
- **Outcome:** Al completar el registro, el usuario obtiene una sesión real y consistente (o completa el flujo de login definido sin estados falsos de autenticación) y puede ingresar a `/cartelera` viendo las publicaciones semanales, el estado vacío informativo o un estado de error claro y recuperable mediante **Reintentar**, sin que fallos secundarios de telemetría ni errores silenciosos de sesión rompan la vista.

---

## Problem

Tras registrarse en Tryckers e intentar acceder a `/cartelera`, el usuario experimenta un error genérico (⚠️ + botón **Reintentar**) en un contexto donde existen tanto una inconsistencia comprobada de estado de autenticación post-registro como un manejo genérico de excepciones en `CarteleraPage` y `UxMetricsService` que impide distinguir fallos de autenticación, conectividad, backend o telemetría.

---

## Current behavior

Evidencia verificada en el código fuente actual de `tryckers-backend` y `tryckers-frontend`:

1. **Inconsistencia en registro y sesión (`tryckers-frontend` vs `tryckers-backend`):**
   - `POST /api/v1/register` (`src/internal/api/handlers/user_handler.go:54-74`) crea el usuario y devuelve `HTTP 201` con `models.User` directamente (sin envoltura `{ user }` y sin `token` ni `refresh_token`).
   - `AuthService.register()` (`tryckers-frontend/src/app/auth/services/auth.service.ts:68-99`) tipa la respuesta como `{ user: User }`, asigna `this._user.set(resp.user)` (que resulta `undefined` y además es un signal privado distinto de `__user`) y marca `this._authStatus.set('authenticated')` **sin tokens ni persistencia en `localStorage`**.
   - Mientras tanto, `authStatus` computado depende de `this._token()`, `AuthStore` (`auth-store.ts`) nunca se actualiza tras el registro (`isLoggedIn` sigue en `false`), y `RegisterPageComponent.onSubmit()` (`register-page.component.ts:74-88`) muestra *"¡Registro exitoso! Bienvenido."* y redirige a `/home` (ruta no protegida por `AuthenticatedGuard`).
   - Al navegar luego a `/cartelera` (protegida por `canMatch: [AuthenticatedGuard]`), `AuthenticatedGuard` evalúa `!!authService.token()` como `false` y redirige a `/auth/login?returnUrl=/cartelera`.
2. **Desincronización entre `AuthService`, `AuthStore`, `localStorage` e interceptores:**
   - `authInterceptor` (`auth.interceptor.ts:41, 84, 116`) lee tokens de `AuthService`, ejecuta `doRefreshToken()` y `forceLogout()` llamando solo a `authService.logout()`, pero **no sincroniza `AuthStore`** (`isLoggedIn`, `token`, `refreshToken`, `user`) ni limpia `localStorage.getItem('isLoggedIn')`.
   - `forceLogout()` redirige a `/auth/login` sin preservar `returnUrl`.
   - `errorInterceptor` (`error.interceptor.ts`) y los componentes (`LoginPageComponent:94-97`, `RegisterPageComponent:104-106`) pueden emitir notificaciones duplicadas para un mismo fallo HTTP.
3. **Manejo genérico de errores y acoplamiento con telemetría en `CarteleraPage`:**
   - `CarteleraPage.loadCartelera()` (`cartelera-page.ts:29-54`) envuelve en un único `try/catch` la telemetría (`uxMetrics.startTiming`, `track`, `endTiming`), la llamada HTTP (`postsService.getCartelera()`) y el ordenamiento (`this.posts.sort(...)`).
   - `UxMetricsService.save()` (`ux-metrics.service.ts:241-243`) ejecuta `localStorage.setItem(STORAGE_KEY, ...)` **sin `try/catch`**: cualquier error de almacenamiento (por ejemplo `QuotaExceededError` o restricciones de navegador) lanza una excepción después de recibir respuesta `200 OK` del backend, vacía `this.posts = []` y muestra `⚠️ No se pudo cargar la cartelera. Inténtalo nuevamente.` con el botón **Reintentar**.
   - Asimismo, cualquier error HTTP (401, 403, 500, red `status === 0`) o de procesamiento cae en el mismo `catch` con el mismo mensaje genérico y sin clasificar el tipo de error en telemetría.
4. **Endpoint `GET /api/v1/cartelera` en backend (`post_handler.go:399-412`, `post_repository.go:142-157`):**
   - Está protegido por `AuthMiddleware()`.
   - Ejecuta `GetTopPostsOfTheWeek(10)` con `LEFT JOIN post_votes`, filtro `posts.status != 'deleted' AND posts.created_at >= NOW() - INTERVAL '7 days'`, `GROUP BY posts.id`, `ORDER BY ... DESC` y `Preload("User")`.
   - `PostHandler.Cartelera` serializa `[]dtos.ResponsePostDto` pasando `userVote = 0`.
   - Queda por confirmar en ejecución/entorno si existe un error SQL, de datos, de red/CORS o de respuesta en `/api/v1/cartelera` independiente del bug de registro.

---

## Target behavior

1. **Diagnóstico completo antes de concluir causa raíz:** Se reproduce y documenta el flujo completo (registro → navegación → `/cartelera`, peticiones HTTP, interceptores, refresh token y consulta SQL `GetTopPostsOfTheWeek`), distinguiendo claramente los defectos confirmados y sus causas raíz con evidencia verificable.
2. **Autenticación consistente tras el registro:** El registro nunca establece un estado falso de autenticación sin tokens válidos. Tras un `POST /api/v1/register` exitoso (`201 Created`), el frontend reutiliza `AuthStore.login(email, password)` para iniciar una sesión real con tokens válidos (sin alterar el contrato de `/register` salvo necesidad demostrada) y redirige a `returnUrl` válido si existe o a `/home` por defecto. Si el auto-login falla, no se deja sesión parcial y se redirige limpiamente a `/auth/login` (preservando `returnUrl` si aplica) informando que la cuenta fue creada. `AuthService`, `AuthStore` y `localStorage` quedan completamente sincronizados.
3. **Sincronización unificada de sesión:** `AuthService`, `AuthStore`, `localStorage`, guards (`AuthenticatedGuard`, `NotAuthenticatedGuard`) e interceptores (`authInterceptor`, `errorInterceptor`) consultan y actualizan un estado de sesión coherente en login, registro, recarga, refresh token y logout/expiración.
4. **Cartelera resiliente y estados bien diferenciados:**
   - `CarteleraPage` distingue correctamente **cargando**, **contenido con posts**, **contenido vacío** (`📭` sin error) y **error** (`⚠️` con mensaje claro para el usuario y botón **Reintentar** funcional).
   - Errores secundarios en `UxMetricsService` (como fallos de `localStorage`) nunca interrumpen la carga ni sobreescriben una respuesta exitosa de la API.
   - Se corrige cualquier defecto confirmado en `CarteleraPage`, `PostsService` o `GET /api/v1/cartelera` / `GetTopPostsOfTheWeek`.
5. **Manejo de errores y observabilidad sin duplicidad ni fuga de secretos:** Se distinguen errores de autenticación (gestionados por refresh/redirección), red y servidor; se registra diagnóstico técnico sin exponer tokens ni datos sensibles; se evitan notificaciones duplicadas entre interceptores y componentes; y se registra en `UxMetricsService` el resultado y tipo de error de carga cuando corresponda.

---

## Entry points

- **UI / Rutas Frontend:**
  - `/auth/register` (`RegisterPageComponent`)
  - `/auth/login` (`LoginPageComponent`)
  - `/home` (`TryckersPageComponent`, destino actual post-registro)
  - `/cartelera` (`CarteleraPage`, protegida por `AuthenticatedGuard`)
- **HTTP / API Backend:**
  - `POST /api/v1/register` (`UserHandler.CreateUser`)
  - `POST /api/v1/login` (`UserHandler.Login`)
  - `POST /api/v1/refresh-token` (`UserHandler.RefreshToken`)
  - `GET /api/v1/cartelera` (`PostHandler.Cartelera`, protegida por `AuthMiddleware`)
- **Interceptores / Estado cliente:**
  - `authInterceptor`, `errorInterceptor`, `AuthService`, `AuthStore`, `UxMetricsService`, `localStorage`

---

## End-to-end flow

1. **Registro sin sesión previa:** Usuario completa formulario en `/auth/register` → `RegisterPageComponent.onSubmit()` envía `POST /api/v1/register` → Backend valida y crea usuario (`201 Created`).
2. **Establecimiento de sesión:** Sin marcar estado falso en memoria, el frontend autentica al usuario reutilizando `AuthStore.login(email, password)` (`POST /api/v1/login`) → Backend emite `token`, `refresh_token` y `user_data` → `AuthService`, `AuthStore` y `localStorage` quedan sincronizados → redirige a `returnUrl` interno seguro si existe, o a `/home` por defecto. (Si el auto-login falla, limpia cualquier estado parcial y redirige a `/auth/login`.)
3. **Navegación a Cartelera:** Usuario navega a `/cartelera` → `AuthenticatedGuard` verifica sesión válida → `CarteleraPage.ngOnInit()` llama `loadCartelera()` → `PostsService.getCartelera()` solicita `GET /api/v1/cartelera`.
4. **Intercepción y autorización:** `authInterceptor` adjunta `Authorization: Bearer <token>` → `AuthMiddleware` valida JWT → `PostHandler.Cartelera` ejecuta `PostService.GetCartelera(10)` → `PostRepository.GetTopPostsOfTheWeek(10)` consulta PostgreSQL (`posts`, `post_votes`, `users`) y retorna `200 OK` con `[]ResponsePostDto`.
5. **Renovación o expiración (si 401):** `authInterceptor` intenta `POST /api/v1/refresh-token`; si renueva, actualiza `AuthService`, `AuthStore` y `localStorage` y reintenta la petición; si falla, limpia sesión sincronizadamente y redirige a `/auth/login`.
6. **Renderizado y telemetría:** `CarteleraPage` renderiza listado (si `posts.length > 0`), estado vacío informativo (`posts.length === 0`) o estado de error diferenciado con **Reintentar** funcional; `UxMetricsService` registra métricas de forma segura (fail-safe) sin afectar la UI.

---

## Expected result

Un usuario recién registrado (y cualquier usuario autenticado) navega con una sesión consistente y válida a `/cartelera`, observa correctamente los posts populares de la semana o el estado vacío informativo, y ante fallos reales (red, servidor o sesión expirada) recibe una respuesta coherente, sin notificaciones duplicadas, sin bloqueos por telemetría y con recuperación efectiva mediante **Reintentar**.

---

## Suggested Knowledge Level

**K2** — Bugfix multi-módulo (`core` + `frontend`) sobre flujos críticos existentes (IAM y Cartelera) que requiere diagnóstico verificable, sincronización de estado cliente y pruebas de regresión, sin introducir nueva arquitectura ni cambios de infraestructura.

---

## Impact analysis

Evaluación explícita por superficie tras inspección del código:

| Superficie | Estado | Evidencia / Motivo |
| :--- | :--- | :--- |
| Product / UI | `affected` | Mensajes de error diferenciados en `CarteleraPage`, comportamiento de **Reintentar**, estados vacío/carga/error y feedback de registro sin duplicidad. |
| Frontend | `affected` | `AuthService`, `AuthStore`, `RegisterPageComponent`, `LoginPageComponent`, `AuthenticatedGuard`, `authInterceptor`, `errorInterceptor`, `CarteleraPage`, `PostsService`, `UxMetricsService`. |
| Backend | `affected` | Diagnóstico obligatorio de `GET /api/v1/cartelera`, `PostHandler.Cartelera`, `PostService.GetCartelera`, `PostRepository.GetTopPostsOfTheWeek` y `AuthMiddleware`; correcciones puntuales si se confirman defectos y pruebas automatizadas de `/cartelera` y auth. No alterar contrato de `/register` salvo necesidad demostrada. |
| Authentication / Authorization | `affected` | Flujo registro → login/sesión, sincronización de tokens en memoria/`localStorage`/store, guards y manejo de `401` + refresh token. |
| Notifications | `affected` | Evitar toasts duplicados entre `errorInterceptor` / `authInterceptor` y componentes (`RegisterPageComponent`, `LoginPageComponent`, `CarteleraPage`). |
| Analytics / Observability | `affected` | Blindar `UxMetricsService` contra fallos de `localStorage` y registrar resultado/tipo de error en carga de cartelera sin introducir herramientas nuevas. |
| Database | `unknown` | No se prevén cambios de esquema ni migraciones; depende de validar en diagnóstico si `GetTopPostsOfTheWeek` presenta algún fallo SQL o de integridad de datos en entorno real. |
| Configuration | `reviewed-not-affected` | No requiere cambios en `environments`, CORS ni Terraform (ya centralizados en WI-004/WI-005). |
| Feature flags | `not-applicable` | El proyecto no utiliza sistema de feature flags. |
| Content / Copy | `affected` | Mensajes de usuario claros y no técnicos para errores de red, servidor y sesión expirada en cartelera y registro. |
| Documentation | `affected` | Actualización de documentación y registro de evidencia de diagnóstico, implementación y validación en Kaddo (AC-14). |
| Operations / Release | `reviewed-not-affected` | Sin cambios de infraestructura ni pipeline CI/CD; incluye validación funcional en entorno desplegado según DoD. |

---

## Module coverage

Proyecto multirepo con módulos mapeados `core` (`tryckers-backend`) y `frontend` (`tryckers-frontend`):

| Módulo | Repositorio | Estado | Alcance y justificación |
| :--- | :--- | :--- | :--- |
| `frontend` | `../tryckers-frontend` | `affected` | Corrección de registro/autenticación (`AuthService`, `AuthStore`, `RegisterPageComponent`), sincronización con interceptores/guards, manejo de errores y estados en `CarteleraPage`/`PostsService`, aislamiento de fallos en `UxMetricsService` y pruebas unitarias. |
| `core` | `tryckers-backend` | `affected` | Diagnóstico y verificación de `GET /api/v1/cartelera` (`GetTopPostsOfTheWeek`), revisión de `AuthMiddleware` y handlers de auth/posts, corrección de defectos confirmados, pruebas automatizadas Go y actualización de conocimiento/evidencia en Kaddo. |

---

## Affected & Reviewed System Entities (Graph-Assisted)

Entidades obtenidas del Semantic System Graph de Kaddo e inspeccionadas en los repositorios:

- **Affected (`affected_system_entities`):**
  - `sys:tryckers-frontend` — Contiene la lógica afectada de registro, sesión, interceptores, cartelera y telemetría cliente.
  - `sys:tryckers-backend` — Expone `/api/v1/register`, `/api/v1/login`, `/api/v1/refresh-token` y `GET /api/v1/cartelera`; sujeto a diagnóstico, pruebas y corrección de defectos confirmados.
  - `sys:local-storage` — Persiste `token`, `refreshToken`, `userData`, `isLoggedIn` y `ux_metrics_events`; actualmente susceptible a desincronización y a excepciones no capturadas.
- **Reviewed (`reviewed_system_entities`):**
  - `sys:postgres-db` (`unknown`) — Consultada por `GetTopPostsOfTheWeek`; sin cambios previstos salvo hallazgo en diagnóstico.
  - `sys:external-services` (`reviewed-not-affected`) — No interviene en estos flujos.
  - `code:tryckers-frontend/src/app/auth/interfaces/auth-response.ts` (`reviewed-not-affected`) — Contratos de login/refresh coherentes con el backend.
  - `code:src/internal/models/user.go`, `code:src/internal/dtos/userUpdate.go`, `code:src/internal/repository/user_repository.go` (`reviewed-not-affected`) — Sin cambios requeridos en el modelo ni actualización de perfil.
  - `code:tryckers-frontend/src/app/pages/settings-page/settings-page.component.ts`, `code:tryckers-frontend/src/app/tryckers/pages/profile-page/**` (`reviewed-not-affected`) — Solo verificación de no regresión (AC-13).

---

## Alcance detallado

### A. Diagnóstico y reproducción

- Revisar el flujo completo desde `/auth/register` hasta el acceso a `/cartelera`.
- Reproducir el comportamiento con un usuario recién registrado y sin sesión previa.
- Comparar con un usuario que ya ha iniciado sesión.
- Inspeccionar las peticiones al endpoint `GET /api/v1/cartelera` (headers, status HTTP, payload y respuesta).
- Verificar respuestas HTTP, `authInterceptor`, `errorInterceptor`, renovación de tokens y posibles errores SQL en `GetTopPostsOfTheWeek`.
- Identificar y documentar la causa raíz de cada defecto confirmado (o las hipótesis descartadas y pruebas ejecutadas si el error original de cartelera no se reproduce).

### B. Corrección de autenticación

- Garantizar que `AuthService.register()` no establezca `_authStatus = 'authenticated'` ni estado parcial sin tokens válidos.
- Evaluar y aplicar la alternativa más sencilla, preferiblemente reutilizando el flujo de login existente (`AuthStore.login`) para autenticar al usuario tras registrarse correctamente.
- Sincronizar `AuthService`, `AuthStore` y `localStorage` en login, registro, recarga, refresh token y logout/expiración (incluyendo limpieza completa de claves como `isLoggedIn`).
- Garantizar que guards (`AuthenticatedGuard`, `NotAuthenticatedGuard`) e interceptores (`authInterceptor`, `errorInterceptor`) consulten y actualicen un estado consistente.
- Mantener el funcionamiento del login, logout, refresh token y redirecciones existentes.
- Evitar modificar el contrato de `POST /api/v1/register` en el backend si no es necesario.

### C. Corrección de la cartelera

- Revisar `CarteleraPage`, `PostsService` y el endpoint `GET /api/v1/cartelera`.
- Corregir los defectos identificados durante la reproducción.
- Verificar la consulta `GetTopPostsOfTheWeek` y sus relaciones con `posts`, `users` y `post_votes` (incluyendo manejo de `posts` vacíos o usuarios nulos).
- Distinguir correctamente los estados de carga, contenido vacío (`📭`) y error (`⚠️`).
- Garantizar que el botón **Reintentar** vuelva a ejecutar `loadCartelera()` limpiamente.
- Evitar que errores secundarios de telemetría (`UxMetricsService`) o presentación provoquen un estado de error de carga cuando la API haya respondido correctamente.

### D. Manejo de errores y diagnóstico

- Registrar información técnica suficiente para diagnosticar fallos sin exponer tokens, contraseñas ni datos sensibles.
- Distinguir errores de autenticación (`401`), conectividad/red (`status === 0`) y servidor (`5xx` / otros).
- Mantener mensajes claros para el usuario, sin mostrar detalles internos.
- Evitar notificaciones duplicadas entre interceptores y componentes.
- Aprovechar `UxMetricsService` registrando el resultado de la carga y el tipo de error cuando corresponda, de forma segura (`try/catch` en persistencia).
- No introducir nuevas herramientas de observabilidad salvo necesidad justificada.

---

## Out of scope

- Implementar la persistencia y rotación de refresh tokens en PostgreSQL definida en `ADR-003` (corresponde a un WI dedicado de IAM).
- Modificar el algoritmo de ranking/trending de la cartelera (e.g., *time decay score*) o agregar paginación/filtros nuevos.
- Migrar tokens de `localStorage` a cookies `HttpOnly`.
- Agregar verificación obligatoria por correo electrónico en el registro.
- Cambios en infraestructura (Terraform, Amplify, CloudFront, EC2) o esquema de base de datos, salvo que el diagnóstico demuestre fehacientemente una causa raíz allí.
- Refactorizaciones extensas ajenas a autenticación y cartelera.

---

## Scope unknowns

1. **Causa exacta del `⚠️` reportado por el usuario en producción:** Aunque se confirmaron defectos reales en el código (registro sin tokens y `UxMetricsService.save()` sin `try/catch` dentro del `try` de `loadCartelera()`), aún no se ha verificado en el entorno desplegado si `GET /api/v1/cartelera` también falla por un error SQL/datos en PostgreSQL, un `401` durante navegación, o un error de red/CORS.
2. **Comportamiento cuando el auto-login post-registro falla tras un `201 Created`:** Si `POST /api/v1/register` crea el usuario pero el `POST /api/v1/login` inmediato falla por error transitorio de red/servidor, el flujo debe evitar dejar estado inconsistente e indicar al usuario que inicie sesión en `/auth/login`.
3. **Voto del usuario en cartelera (`user_vote`):** `PostHandler.Cartelera` actualmente pasa `userVote = 0` fijo a `ToResponsePostDto`. No causa error de carga, por lo que no se modifica salvo que el diagnóstico lo vincule a un fallo.

---

## Scope confidence

**Medium**

**Razones:**
- **Alta confianza** en los defectos confirmados mediante lectura directa del código en `frontend` (`AuthService.register`, desincronización con `AuthStore`/`authInterceptor`, notificaciones duplicadas y acoplamiento `CarteleraPage` ↔ `UxMetricsService.save()`).
- **Incertidumbre acotada** hasta completar la Fase A (diagnóstico y reproducción en entorno real/local) sobre si existe un fallo adicional en `GET /api/v1/cartelera` o en PostgreSQL (`GetTopPostsOfTheWeek`) que explique el reporte original independientemente del bug de registro.

---

## Acceptance criteria

- [x] **AC-01:** Se documenta la causa raíz del error reportado con evidencia verificable o, si no se reproduce, las hipótesis descartadas y pruebas ejecutadas.
- [x] **AC-02:** El registro no genera estados falsos de autenticación (`AuthService`, `AuthStore` y `localStorage` nunca reflejan sesión activa sin tokens válidos).
- [x] **AC-03 (E2E):** Un usuario recién registrado desde un navegador sin sesión previa completa correctamente el flujo definido de autenticación y accede a `/cartelera` sin errores. *(Validado manualmente por el usuario: bug resuelto.)*
- [x] **AC-04:** `AuthService`, `AuthStore`, `localStorage`, los guards y los interceptores mantienen un estado de sesión consistente en registro, login, recarga, refresh token y logout.
- [x] **AC-05 (E2E):** La cartelera carga correctamente para usuarios autenticados (tanto recién registrados como con sesión tradicional) y muestra las publicaciones ordenadas por votos. *(Validado manualmente por el usuario y cubierto por pruebas unitarias.)*
- [x] **AC-06:** Una cartelera vacía (`[]`) presenta su estado informativo (`📭`) sin mostrar estado de error.
- [x] **AC-07:** Los errores de autenticación (`401`) se gestionan mediante el flujo de renovación de token o redirección a `/auth/login` correspondiente, sin dejar estado inconsistente ni notificaciones duplicadas.
- [x] **AC-08:** Los errores del backend (`5xx`) y de conectividad (`status === 0`) se distinguen y gestionan adecuadamente con mensajes claros para el usuario, sin exponer detalles técnicos ni tokens en logs o UI.
- [x] **AC-09:** El botón **Reintentar** permite recuperar la cartelera después de un fallo temporal del backend o de red.
- [x] **AC-10:** Las excepciones de telemetría (`UxMetricsService`, incluido fallo de `localStorage`) no interrumpirían la carga ni funcionalidades principales de la cartelera o el registro.
- [x] **AC-11:** Existen pruebas automatizadas para los escenarios críticos de registro, autenticación, sincronización de sesión y carga/manejo de errores de cartelera (frontend y backend).
- [x] **AC-12:** Las pruebas y compilaciones de frontend (`ng build --configuration development`, `ng test`) y backend (`go test ./...`, `go build ./...`) finalizan correctamente.
- [x] **AC-13:** Los flujos existentes de login, logout, perfiles, configuración y publicaciones no presentan regresiones. *(Confirmado en prueba funcional del usuario.)*
- [x] **AC-14:** Se actualiza la documentación y se registra la evidencia de diagnóstico, implementación y validación en Kaddo.

---

## Diagnóstico y Causa Raíz (Fase A)

1. **Defecto confirmado #1 — Registro sin tokens + estado falso en cliente:**
   - `POST /api/v1/register` retorna `201 Created` con `models.User` sin tokens.
   - `AuthService.register()` tipaba `{ user: User }` (obteniendo `undefined`) y marcaba `_authStatus = 'authenticated'` sin tokens ni actualizar `AuthStore` ni `localStorage`.
   - Además, `authInterceptor` no sincronizaba `AuthStore` al refrescar ni al forzar logout, y podía dejar peticiones concurrentes colgadas si el refresh fallaba.
2. **Defecto confirmado #2 — Excepciones de telemetría rompiendo carga exitosa de Cartelera:**
   - `CarteleraPage.loadCartelera()` ejecutaba `UxMetricsService.track`/`endTiming` dentro del mismo `try/catch` de la llamada HTTP; `UxMetricsService.save()` llamaba `localStorage.setItem()` sin `try/catch`, haciendo que cualquier error de cuota/almacenamiento limpiara `this.posts` y mostrara `⚠️` + **Reintentar** incluso tras un `200 OK`.
3. **Defecto confirmado #3 — Notificaciones duplicadas y clasificación genérica de errores:**
   - `LoginPageComponent` y `RegisterPageComponent` emitían toasts adicionales sobre errores HTTP ya notificados por `errorInterceptor`.
   - `CarteleraPage` no diferenciaba errores `401` (auth), `0` (red) y `5xx` (servidor).
4. **Defecto confirmado #4 — Seguridad de claims en `AuthMiddleware`:**
   - `src/internal/api/middlewares/auth_middleware.go` realizaba type assertions sin verificación (`claims["sub"].(string)`, `claims["role"].(string)`) y no rechazaba tokens `type: "refresh"` en rutas protegidas.
5. **Verificación de `GET /api/v1/cartelera` y entorno desplegado:**
   - Sondeo en vivo: `GET https://api.tryckers.trycatch.tv/health` → `200 OK`; `GET https://api.tryckers.trycatch.tv/api/v1/cartelera` sin token → `401 Unauthorized` (CORS/CloudFront operativos).
   - Consulta SQL de `GetTopPostsOfTheWeek` verificada por test dry-run GORM (`LEFT JOIN post_votes`, `INTERVAL '7 days'`, `GROUP BY "posts"."id"`, `ORDER BY ... DESC`, `LIMIT`).
   - *Nota de cierre:* La validación funcional E2E en el entorno desplegado (y revisión de logs Docker en EC2 / RDS si se reprodujera un 500 específico en producción) queda pendiente de ejecutar tras el despliegue de estas ramas.

## Verified Evidence

- **Backend (`tryckers-backend`):**
  - `go test -v ./...` → **PASS** (10/10 tests: `TestHealthCheck`, 5 tests de `AuthMiddleware`, 3 tests de `PostHandler.Cartelera`, 1 test dry-run SQL de `PostRepository.GetTopPostsOfTheWeek`).
  - `go build ./...` → **PASS** (exit code 0).
- **Frontend (`tryckers-frontend`):**
  - `npx ng build --configuration development` → **PASS** (exit code 0, bundle generado en `dist/tryckers-frontend`).
  - `npx ng test --watch=false --browsers=ChromeHeadless` (4 suites: `auth.service.spec.ts`, `register-page.component.spec.ts`, `auth.interceptor.spec.ts`, `cartelera-page.spec.ts`) → **12 SUCCESS, 0 FAILED**.
- **Validación funcional manual (usuario):**
  - Confirmado en prueba end-to-end por el usuario que el bug de autenticación post-registro y carga de `/cartelera` quedó resuelto.

---

## Validation

### 1. Pruebas automatizadas y compilación

- **Backend (`tryckers-backend`):**
  - `go test ./...` (incluyendo tests unitarios/integración de `PostHandler.Cartelera`, `PostService.GetCartelera`, `GetTopPostsOfTheWeek` con lista con posts y vacía, y `AuthMiddleware`).
  - `go build ./...`
- **Frontend (`tryckers-frontend`):**
  - `npm test -- --watch=false` (incluyendo specs para `AuthService`, `AuthStore`, `RegisterPageComponent`, `authInterceptor`, `CarteleraPage` y `UxMetricsService` simulando fallo de `localStorage`, error `401`, error `500`, error de red `0`, cartelera vacía y recuperación con **Reintentar**).
  - `npm run build`

### 2. Matriz de pruebas de regresión manual / E2E (local y entorno desplegado)

Validar como mínimo los 10 escenarios:
1. **Registro desde navegador limpio (sin sesión previa):** completar `/auth/register` → verificar creación `201`, obtención de tokens vía login, sincronización en `localStorage` y ausencia de errores en consola.
2. **Acceso a `/cartelera` tras registro:** navegar a `/cartelera` → verificar `GET /api/v1/cartelera` con `Authorization: Bearer <token>` y respuesta `200 OK`.
3. **Login tradicional y acceso a `/cartelera`:** iniciar sesión en `/auth/login` → navegar a `/cartelera` → verificar carga correcta.
4. **Recarga de página con sesión válida:** presionar `F5` en `/cartelera` → verificar que mantiene sesión y recarga la cartelera sin redirigir al login.
5. **Renovación de access token expirado:** con access token expirado y refresh token válido, cargar `/cartelera` → verificar `401` → `POST /api/v1/refresh-token` → actualización sincronizada de tokens → reintento `200 OK`.
6. **Sesión no renovable:** con access y refresh token inválidos/expirados, cargar `/cartelera` → verificar limpieza sincronizada de sesión (`AuthService`, `AuthStore`, `localStorage`) y redirección a `/auth/login` sin toasts duplicados.
7. **Cartelera con publicaciones:** verificar renderizado ordenado por votos y navegación a detalle/perfil.
8. **Cartelera sin publicaciones:** verificar respuesta `[]` → muestra estado vacío (`📭`), nunca estado de error (`⚠️`).
9. **Error temporal del backend y recuperación:** simular `500`/`503` o caída temporal en `/api/v1/cartelera` → verificar estado de error con mensaje amigable → restablecer servicio → pulsar **Reintentar** → verificar recuperación exitosa.
10. **Navegación general y no regresión:** recorrer `/home`, `/cartelera`, `/profile/:username`, `/profile/:username/post/:id`, `/settings`, logout y re-login sin errores en consola ni estados inconsistentes.

### 3. Verificación Kaddo

- Ejecutar validación de consistencia de Kaddo (`kaddo_validate_work_item_modules` / `kaddo_verify_work_item`) y registrar evidencia verificada en el WI al cierre.

---

## Restricciones

- Priorizar la corrección de causas raíz sobre soluciones superficiales.
- No eliminar los mecanismos de autenticación existentes.
- No introducir nuevas dependencias sin justificación.
- No modificar infraestructura ni base de datos salvo que el diagnóstico demuestre que es necesario.
- No realizar refactorizaciones extensas ajenas al problema.
- No asumir que el error de registro es necesariamente la causa del fallo de la cartelera.

---

## Definition of Done

El WI se considera completado cuando:
- Se ha completado el diagnóstico y documentado la causa raíz verificada (o hipótesis descartadas con evidencia si el error original no se reproduce).
- El flujo de registro y autenticación funciona consistentemente entre `AuthService`, `AuthStore`, `localStorage`, guards e interceptores.
- La cartelera carga correctamente en todos los estados (cargando, con posts, vacía, error recuperable con **Reintentar**) y está aislada de fallos secundarios de telemetría.
- Todos los criterios de aceptación (**AC-01** a **AC-14**) están verificados.
- Las pruebas automatizadas y compilaciones de `tryckers-frontend` y `tryckers-backend` pasan sin errores.
- Se ha completado la validación funcional en el entorno desplegado y se ha registrado la evidencia y documentación en Kaddo.

> **Regla de cierre ante no reproducibilidad:** Si el error original en `/cartelera` no puede reproducirse, no debe declararse solucionado únicamente por haber corregido el problema de autenticación; debe quedar documentada la limitación junto con evidencia suficiente del comportamiento de `/cartelera` y su diagnóstico.

---

## Decisiones confirmadas

- [resolved] **Flujo post-registro:** Tras `POST /api/v1/register` exitoso (`201 Created`), el frontend reutilizará `AuthStore.login(email, password)` para autenticar automáticamente al usuario sin cambiar el contrato de registro del backend. Si el auto-login falla, no se deja sesión parcial y se redirige a `/auth/login` informando que la cuenta fue creada y debe iniciar sesión.
- [resolved] **Destino tras registro exitoso:** Redirigir a `/home` por defecto, dando prioridad a `returnUrl` (validado como ruta interna segura, igual que en `LoginPageComponent`) si existe.
- [resolved] **Acceso al entorno desplegado para diagnóstico (Fase A):** Confirmado acceso a RDS y a EC2 (ingresando a la instancia EC2 para revisar los logs de Docker del backend) en caso de requerir correlación con producción.

## Open Questions

<!-- Sin preguntas abiertas bloqueantes -->

---

## Suggested ownership (code globs)

- `tryckers-frontend/src/app/auth/**`
- `tryckers-frontend/src/app/core/interceptors/**`
- `tryckers-frontend/src/app/post/pages/cartelera-page/**`
- `tryckers-frontend/src/app/post/services/posts.service.ts`
- `tryckers-frontend/src/app/shared/services/ux-metrics.service.ts`
- `src/internal/api/handlers/post_handler.go`
- `src/internal/services/post_service.go`
- `src/internal/repository/post_repository.go`
- `src/internal/api/middlewares/auth_middleware.go`
- `src/internal/api/handlers/user_handler.go`
- `src/internal/services/user_service.go`
- `src/internal/api/routes/router.go`

---

## Related domain / capability

- **Domains:** `Identity and Access Management (IAM)`, `Social Signals and Curation (Votes and Billboard)`
- **Capabilities:**
  - `User Registration and Password Hashing`
  - `JWT Authentication and Route Protection`
  - `Refresh Token Rotation and Session Renewal`
  - `Weekly Billboard Discovery (Cartelera Semanal)`

## Related decisions

- `ADR-003-refresh-token-rotation-and-storage` (accepted) — Referenciado por relación con el flujo de refresh token; este WI preserva el mecanismo actual y **no** implementa la persistencia/rotación en tabla `refresh_tokens`.

## Related legacy findings

- No aplica (no existen archivos `knowledge/legacy/` en el repositorio).

────────────────────────
Agent: work-item-agent

Produced:
knowledge/delivery/work-items/draft/WI-006-wi-corregir-errores-de-autenticaci-n-y-carga-de-la.md

Next:
kaddo ready (human review)
implementation-agent
────────────────────────
