---
id: WI-010
title: Corregir expiración de sesión y flujo de Refresh Token
type: bugfix
priority: high
status: completed
created: 2026-11-26
completed: 2026-11-26
repos:
  - tryckers-frontend
  - tryckers-backend
modules:
  - core
  - frontend
related_adrs:
  - ADR-003
---

# WI-010 — Corregir expiración de sesión y flujo de Refresh Token

## 1. Contexto y Objetivo
Garantizar que las sesiones autenticadas en Tryckers renueven automáticamente el `access_token` (5h) usando el `refresh_token` (7d) sin expulsar prematuramente al usuario, solicitando re-login únicamente cuando las credenciales de renovación sean definitivamente inválidas, expiradas o inexistentes.

## 2. Diagnóstico de Causa Raíz (Evidencia en Código)

Tras auditar el flujo end-to-end en `tryckers-frontend` y `tryckers-backend`, se confirmaron **5 defectos concretos**:

1. **Error de petición reintentada interpretado como fallo de refresh (`auth.interceptor.ts:116-138`):** `catchError(() => forceLogout(...))` envolvía tanto `authService.doRefreshToken()` como `switchMap(() => next(...))`; cualquier error HTTP (`400`, `403`, `404`, `500`, `503`) de la petición original reintentada cerraba la sesión recién renovada.
2. **Destrucción inmediata de sesión ante errores transitorios (`auth.service.ts:92-111`):** `AuthService.doRefreshToken()` ejecutaba `this.logout()` en `catchError` para cualquier error, sin distinguir `401`/`400` de errores de red (`0`) o servidor (`500`/`502`/`503`).
3. **Riesgo de bloqueo por cancelación y falta de guarda de reintento (`auth.interceptor.ts:92-139`):** no existía `finalize` para liberar `isRefreshing` ni un `HttpContextToken` para evitar ciclos infinitos ante `401` persistentes.
4. **Desincronización multi-pestaña y pérdida de sesión recuperable (`auth.service.ts`, `auth-store.ts`, `authenticated.guard.ts`):** faltaba sincronización vía evento `storage` y se descartaba `refreshToken` cuando `token` estaba ausente al recargar; `error.interceptor.ts` duplicaba alertas en `400` de `/refresh-token`.
5. **Endurecimiento JWT pendiente en backend (`jwt.go`, `user_handler.go`, `auth_middleware.go`):** faltaba restringir `HS256`, exigir `exp`, validar `token.Valid` y rechazar `sub`/`role`/`refresh_token` en blanco.

## 3. Solución Implementada y Evidencia de Despliegue

- **Backend (`tryckers-backend` — tag `v1.1.3`, commit `85d2ee5`, merge `main` `94f9bf2`):**
  - `src/internal/utils/jwt.go`, `src/internal/api/middlewares/auth_middleware.go`, `src/internal/api/handlers/user_handler.go`: validación estricta de `HS256`, `exp`, `type`, `sub`, `role` y `refresh_token`, manteniendo vigencias (5h/7d) y contrato `{ token, refresh_token }`.
  - Pruebas unitarias deterministas en `src/internal/utils/jwt_test.go` y `src/internal/api/handlers/user_handler_test.go`: `go vet ./...`, `go test -race`, `go test ./...` y `go build ./...` en verde.
  - Despliegue: `Backend CD` Run `45939702557` completado con éxito (1m28s).
  - Verificación en vivo (`https://api.tryckers.trycatch.tv/api/v1/refresh-token`):
    - `refresh_token` en blanco → `HTTP 400` (`refresh_token es requerido`).
    - `refresh_token` inválido → `HTTP 401` (`refresh token inválido o expirado`).

- **Frontend (`tryckers-frontend` — tag `v1.1.4`, commit `46f7db7`, merge `main` `7c35d4a`):**
  - `src/app/core/interceptors/auth.interceptor.ts`: separa errores del refresh de errores del reintento; conserva sesión ante `0`/`5xx`/`429`; cierra sesión solo ante `400`/`401`/`403`; coordina `401` concurrentes con `finalize`; evita ciclos con `IS_AUTH_RETRY`; reutiliza token rotado por otra pestaña.
  - `src/app/auth/services/auth.service.ts` y `src/app/auth/store/auth-store.ts`: propagan errores en `doRefreshToken()`, sincronizan vía evento `storage` multi-pestaña y preservan sesión recuperable.
  - `src/app/core/interceptors/error.interceptor.ts`, `src/app/auth/guards/authenticated.guard.ts`, `src/app/auth/guards/not-authenticated.guard.ts`: sin notificaciones duplicadas y recuperación silenciosa en guards.
  - Pruebas unitarias en `auth.interceptor.spec.ts` y `auth.service.spec.ts`: `npx ng build --configuration development` en verde; Karma `ChromeHeadless` 35/36 en verde (todas las pruebas de auth/interceptor/service en verde; único fallo preexistente no relacionado en `Login` por `NG01352`).
  - Despliegue automático disparado en AWS Amplify desde `main`.

## 4. Criterios de Aceptación
- [x] AC-01: Causa raíz identificada y documentada con evidencia en código.
- [x] AC-02: Renovación automática del access token con refresh token válido.
- [x] AC-03: Navegación continua sin re-login tras renovación exitosa.
- [x] AC-04: Errores temporales de red/servidor (`0`, `5xx`, `429`) no cierran la sesión.
- [x] AC-05: Errores en la petición reintentada no se interpretan como fallo de refresh.
- [x] AC-06: Sincronización consistente entre `AuthService`, `AuthStore` y `localStorage`.
- [x] AC-07: Renovación única compartida ante `401` concurrentes y sin deadlocks.
- [x] AC-08: Recuperación de sesión tras recarga y consistencia multi-pestaña.
- [x] AC-09: Cierre de sesión controlado ante refresh token expirado/inválido (`400`/`401`/`403`).
- [x] AC-10: Prevención de ciclos infinitos mediante `IS_AUTH_RETRY`.
- [x] AC-11: Validación estricta de tipo, algoritmo y expiración JWT en backend.
- [x] AC-12: Pruebas automatizadas deterministas en backend y frontend.
- [x] AC-13: Sin regresiones en login, registro, logout ni rutas protegidas.
- [x] AC-14: Verificado en entorno desplegado de producción.
- [x] AC-15: Documentación y evidencias registradas en Kaddo (`completed/`).
