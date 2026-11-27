---
id: WI-011
title: Mostrar versión desplegada y corregir contraste de publicaciones en Dark Mode
type: enhancement
priority: medium
status: in-progress
created: 2026-11-26
repos:
  - tryckers-frontend
modules:
  - frontend
---

# WI-011 — Mostrar versión desplegada y corregir contraste de publicaciones en Dark Mode

## 1. Contexto y Objetivo
1. Mostrar en el footer principal (`src/app/shared/components/footer/`) un indicador discreto de la versión y el commit corto del build desplegado (`vX.Y.Z · build abc1234`), generado automáticamente en tiempo de compilación sin llamadas al backend ni a GitHub en runtime.
2. Corregir los problemas de contraste en `cartelera-page` en Dark Mode reemplazando colores fijos (`#111827`, `#4b5563`, `#6b7280`, `#1d4ed8`, `text-gray-700`, `text-gray-500`) por variables CSS de tema (`--ui-*`) y reglas `:host-context(.dark-theme)` con contraste WCAG AA (>= 4.5:1).

## 2. Alcance Implementado
- `package.json`: versión actualizada de `0.0.0` a `1.2.0` (alineada con el historial de tags SemVer del repositorio tras `v1.1.4`).
- `scripts/generate-environment.mjs`: resuelve `appVersion` desde `package.json` y `buildCommit` desde `AWS_COMMIT_ID` / `GITHUB_SHA` / `COMMIT_SHA` / `git rev-parse --short=7 HEAD` con fallback `'local'`, generando `src/environments/environment.production.ts`.
- `src/environments/environment.ts` y `environment.development.ts`: incluyen `appVersion` y `buildCommit: 'local'` para ejecución local independiente de Amplify.
- `src/app/shared/components/footer/`: muestra `v1.2.0 · build <sha>` de forma discreta, responsiva y adaptable a Light/Dark Mode.
- `src/app/post/pages/cartelera-page/`: títulos, descripciones, metadatos, enlaces, chips, panel de votos, badges y estados vacíos/error adaptados con variables `--ui-*` y `:host-context(.dark-theme)`.
- Pruebas automatizadas: `scripts/generate-environment.test.mjs` y `src/app/shared/components/footer/footer.component.spec.ts`.

## 3. Criterios de Aceptación
- [ ] AC-01..AC-18: Validación de tests, build, despliegue en Amplify y verificación visual en producción (Light/Dark Mode y footer).
