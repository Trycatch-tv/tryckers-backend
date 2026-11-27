---
id: WI-011
title: Mostrar versión desplegada y corregir contraste de publicaciones en Dark Mode
type: enhancement
priority: medium
status: completed
created: 2026-11-26
completed: 2026-11-26
repos:
  - tryckers-frontend
modules:
  - frontend
---

# WI-011 — Mostrar versión desplegada y corregir contraste de publicaciones en Dark Mode

## 1. Contexto y Objetivo
1. Mostrar en el footer principal (`src/app/shared/components/footer/`) un indicador discreto de la versión y el commit corto del build desplegado (`vX.Y.Z · build abc1234`), generado automáticamente en tiempo de compilación sin llamadas al backend ni a GitHub en runtime.
2. Corregir los problemas de contraste en `cartelera-page` en Dark Mode reemplazando colores fijos (`#111827`, `#4b5563`, `#6b7280`, `#1d4ed8`, `text-gray-700`, `text-gray-500`) por variables CSS de tema (`--ui-*`) y reglas `:host-context(.dark-theme)` con contraste WCAG AA (>= 4.5:1).

## 2. Solución Implementada y Evidencia de Despliegue
- **Repositorio:** `tryckers-frontend` — tag **`v1.2.0`**, commit `817eb31`, merge `main` `d6aff48`.
- **Versión desplegada:**
  - `package.json`: versión actualizada a `1.2.0`.
  - `scripts/generate-environment.mjs`: resuelve `appVersion` desde `package.json` y `buildCommit` desde `AWS_COMMIT_ID` / `GITHUB_SHA` / `COMMIT_SHA` / `git rev-parse --short=7 HEAD` con fallback `'local'`.
  - `src/environments/environment.ts` y `environment.development.ts`: incluyen `appVersion: '1.2.0'` y `buildCommit: 'local'`.
  - `src/app/shared/components/footer/`: indicador discreto y responsivo (`v1.2.0 · build d6aff48`) adaptado a Light/Dark Mode.
- **Dark Mode en Cartelera:**
  - `src/app/post/pages/cartelera-page/cartelera-page.css` y `cartelera-page.html`: migración a variables `--ui-*` y selectores `:host-context(.dark-theme)` para títulos, descripciones, metadatos, enlaces, chips, badges, panel de votos y estados vacío/error.
- **Pruebas automatizadas:**
  - `scripts/generate-environment.test.mjs`: 4/4 PASS (`node --test`).
  - `src/app/shared/components/footer/footer.component.spec.ts`: 1/1 PASS (Karma `ChromeHeadless`).
  - `npx ng build --configuration development`: exitoso.
- **Despliegue y validación en producción:**
  - AWS Amplify (`d38qa3mp68ehik`, rama `main`, Job #20, commit `d6aff48`): desplegado con éxito.
  - Validación visual en producción (`https://tryckers.trycatch.tv`) confirmada por el usuario: footer mostrando `v1.2.0 · build d6aff48` y cartelera legible en Light Mode y Dark Mode.

## 3. Criterios de Aceptación
- [x] AC-01: El footer muestra la versión del frontend y el identificador corto del commit.
- [x] AC-02: La información se genera automáticamente durante el build.
- [x] AC-03: Los despliegues posteriores actualizan el identificador sin intervención manual.
- [x] AC-04: El identificador corresponde al artefacto realmente construido y desplegado (`d6aff48`).
- [x] AC-05: Visualización correcta en móvil y escritorio.
- [x] AC-06: Compilación y ejecución local independientes de Amplify (`build local`).
- [x] AC-07: Versión legible en Light Mode y Dark Mode.
- [x] AC-08: Títulos de publicaciones legibles en Dark Mode.
- [x] AC-09: Títulos conservan presentación adecuada en Light Mode.
- [x] AC-10: Descripciones, metadatos y enlaces con contraste adecuado en ambos temas.
- [x] AC-11: Cambio de tema inmediato sin recargar.
- [x] AC-12: Estados hover, focus y navegación conservados.
- [x] AC-13: Presentación responsiva en móvil y escritorio.
- [x] AC-14: Pruebas automatizadas de generación y presentación de versión.
- [x] AC-15: Contraste y visualización verificados en ambos temas.
- [x] AC-16: Pruebas y build de Angular en verde.
- [x] AC-17: Desplegado y validado en producción.
- [x] AC-18: Cambios y evidencias documentados en Kaddo (`completed/`).
