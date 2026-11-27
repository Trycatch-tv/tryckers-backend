---
id: WI-010
title: Corregir expiración de sesión y flujo de Refresh Token
type: bugfix
priority: high
status: in-progress
created: 2026-11-26
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

Tras auditar el flujo end-to-end en `tryckers-frontend` y `tryckers-backend`, se confirmaron **5 defectos concretos** que provocan cierres de sesión prematuros o inconsistentes:

1. **Error de petición reintentada interpretado como fallo de refresh (`auth.interceptor.ts:116-138`):**
   - `catchError(() => forceLogout(...))` envolvía tanto `authService.doRefreshToken()` como `switchMap(() => next(addToken(req, response.token)))`.
   - Cualquier error HTTP (`400`, `403`, `404`, `500`, `503`) devuelto por la **petición original reintentada** después de un refresh exitoso caía en `catchError`, ejecutaba `forceLogout()`, destruía la sesión y redirigía a `/auth/login`.
2. **Destrucción inmediata de sesión ante errores transitorios de red o servidor (`auth.service.ts:92-111`):**
   - `AuthService.doRefreshToken()` capturaba cualquier error en `catchError(() => { this.logout(); return of(null); })`, sin distinguir entre un refresh token inválido/expirado (`401`/`400`) y un error temporal (`status 0`, `500`, `502`, `503`, `504`, `429`).
   - Al tragar el `HttpErrorResponse`, `authInterceptor` recibía `null` y forzaba el logout y la notificación de "Sesión Expirada".
3. **Riesgo de bloqueo de solicitudes concurrentes y falta de marca de reintento (`auth.interceptor.ts:92-139`):**
   - Si la petición que iniciaba el refresh era cancelada/desuscrita en vuelo, `isRefreshing` permanecía en `true` indefinidamente bloqueando todas las peticiones subsecuentes en `refreshTokenSubject`.
   - No existía un `HttpContextToken` que marcara una petición como ya reintentada para prevenir ciclos infinitos ante `401` persistentes.
4. **Desincronización entre pestañas y recuperación parcial tras recarga (`auth.service.ts`, `auth-store.ts`, `authenticated.guard.ts`):**
   - No existía sincronización mediante el evento `storage` del navegador: si una pestaña rotaba los tokens o cerraba sesión, las otras pestañas conservaban tokens obsoletos en memoria.
   - Si `token` faltaba pero `refreshToken` seguía presente en `localStorage`, ambos servicios descartaban `refreshToken` al inicializar.
   - `error.interceptor.ts` mostraba toasts de `400` también para `/refresh-token`, duplicando alertas con `authInterceptor`.
5. **Endurecimiento pendiente en validación JWT del backend (`jwt.go`, `user_handler.go`, `auth_middleware.go`):**
   - `VerifyToken` no restringía explícitamente `jwt.WithValidMethods([]string{"HS256"})`, no exigía `jwt.WithExpirationRequired()` ni validaba `token.Valid`.
   - `RefreshAccessToken` no validaba que `sub` y `role` fueran strings no vacíos tras `TrimSpace`, ni permitía inyectar reloj (`nowFunc`) para pruebas deterministas de expiración.
   - `UserHandler.RefreshToken` no recortaba espacios en `req.RefreshToken` antes de validar.

## 3. Alcance de la Corrección

- **Backend (`tryckers-backend`):**
  - Endurecer `VerifyToken`, `RefreshAccessToken` y `AuthMiddleware` (algoritmo HS256, `exp` obligatorio, `token.Valid`, claims `sub`/`role`/`type` no vacíos).
  - Mantener vigencias actuales (5h access / 7d refresh) y contrato JSON `{ token, refresh_token }`.
  - Agregar pruebas unitarias deterministas de expiración, rotación, rechazo de access tokens en `/refresh-token` y rechazo de tokens malformados.
- **Frontend (`tryckers-frontend`):**
  - `AuthService`: propagar errores HTTP en `doRefreshToken()` sin hacer `logout()` prematuro; sincronizar estado desde `localStorage` y evento `storage` multi-pestaña; conservar sesión recuperable si existe `refreshToken`.
  - `AuthStore`: sincronizarse ante cambios de `storage` entre pestañas y conservar sesión recuperable al inicializar.
  - `authInterceptor`: separar el manejo de error del refresh respecto a la petición reintentada; distinguir fallos definitivos (`400`/`401`/`403`) de errores transitorios (`0`/`5xx`/`429`); usar `HttpContextToken` para evitar ciclos infinitos; usar `finalize` para evitar deadlocks por cancelación; re-leer `localStorage` antes de refrescar por si otra pestaña ya rotó el token.
  - `errorInterceptor`: evitar toasts duplicados en `/refresh-token` para `400`/`401`.
  - `AuthenticatedGuard` / `NotAuthenticatedGuard`: sincronizar desde `storage` y recuperar sesión cuando existe `refreshToken` válido.
  - Pruebas unitarias completas cubriendo todos los escenarios de la tabla de validación.

## 4. Criterios de Aceptación
- [x] AC-01: Causa raíz identificada y documentada con evidencia en código.
- [ ] AC-02..AC-15: Implementación, pruebas automatizadas, despliegue y verificación en producción.
