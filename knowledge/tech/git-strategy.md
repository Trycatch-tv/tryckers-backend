---
type: git-strategy
project_state: ai-assisted
generated_by: implementation-agent
template_version: 1
---

> Idioma del proyecto: **español**. Mantén en inglés los comandos, nombres de ramas y formatos de commit/tag.

# Git Strategy — Tryckers (`tryckers-backend` & `tryckers-frontend`)

## 1. Modelo de Ramas

- **`main`**: Rama de producción. Cada push/merge a `main` dispara despliegue automático:
  - `tryckers-frontend` → **AWS Amplify** (`https://tryckers.trycatch.tv`).
  - `tryckers-backend` → **GitHub Actions** (`backend-cd.yml` → ECR → EC2 vía SSM, `https://api.tryckers.trycatch.tv`).
- **`dev`**: Rama de integración previa a producción.
- **Ramas de trabajo por Work Item:**
  - `feature/<WI-ID>-<slug>` — nuevas funcionalidades.
  - `bugfix/<WI-ID>-<slug>` — corrección de errores.
  - `hotfix/<WI-ID>-<slug>` — correcciones urgentes en producción.
  - `chore/<WI-ID>-<slug>` (o `chore/<slug>`) — mantenimiento, CI/CD, infraestructura o documentación.

## 2. Flujo de Integración

1. Crear rama de trabajo desde `main` (o `dev` actualizada).
2. Commitear cambios usando **Conventional Commits** (`feat`, `fix`, `chore`, `docs`, `test`, `refactor`) referenciando el `WI-ID`.
3. Hacer push de la rama de trabajo.
4. Merge a `dev` y `git push origin dev`.
5. Merge de `dev` a `main` y `git push origin main`.

## 3. Versionado y Tags (SemVer)

- Se utiliza **Semantic Versioning (`vMAJOR.MINOR.PATCH`)** sincronizado en los repositorios afectados al finalizar cada Work Item:
  - **PATCH (`vX.Y.Z+1`):** `bugfix`, `hotfix`, `chore`, mejoras operativas o de CI/CD.
  - **MINOR (`vX.Y+1.0`):** `feature` retrocompatible.
  - **MAJOR (`vX+1.0.0`):** cambios incompatibles (*breaking changes*).
- **Regla al cerrar un WI:** Tras completar el merge en `main`, crear un **tag anotado** (`git tag -a vX.Y.Z -m "..."`) y publicarlo (`git push origin vX.Y.Z`).
- **Baseline inicial:** `v1.0.0` (cierre de `WI-007` + optimización multi-stage Docker en backend).

## 4. Releases (GitHub Releases)

- Los **GitHub Releases** se crean **únicamente cuando el usuario lo indique explícitamente**, usando `gh release create <tag>` con notas generadas a partir de los WIs incluidos desde el último release.
