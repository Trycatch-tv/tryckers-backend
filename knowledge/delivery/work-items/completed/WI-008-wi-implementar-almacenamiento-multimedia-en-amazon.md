---
type: feature
id: WI-008
title: Implementar almacenamiento multimedia en Amazon S3 para Dev y Prod
status: completed
work_type: feature
knowledge_level: K3
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
  - infrastructure
  - storage
  - backend
  - frontend
  - database
related_domain: storage
related_capabilities:
  - Avatar and Banner Upload with Local Storage
  - Post Publication and Media Tagging
related_decisions:
  - ADR-001-s3-compatible-storage-driver
affected_system_entities:
  - 'sys:tryckers-backend'
  - 'sys:tryckers-frontend'
  - 'sys:postgres-db'
reviewed_system_entities:
  - entity: 'sys:external-services'
    status: affected
    reason: >-
      Integración con Amazon S3 (buckets Dev y Prod) y CloudFront OAC para
      entrega privada en Prod.
code:
  - infra/dev/**
  - infra/s3-media.tf
  - infra/main.tf
  - infra/outputs.tf
  - docker-compose.prod.yml
  - scripts/deploy.sh
  - src/cmd/main.go
  - src/internal/config/config.go
  - src/internal/services/storage/**
  - src/internal/services/user_service.go
  - src/internal/services/post_service.go
  - src/internal/api/handlers/media_handler.go
  - src/internal/api/handlers/post_handler.go
  - src/internal/api/routes/router.go
  - src/internal/migrations/cleanup_legacy_media.go
  - ../tryckers-frontend/src/app/tryckers/services/tryckers-service.ts
  - ../tryckers-frontend/src/app/post/services/posts.service.ts
  - ../tryckers-frontend/src/app/tryckers/pages/profile-page/**
summary: >-
  Reemplazar el almacenamiento local multimedia por Amazon S3 con buckets
  independientes para Dev (estado Terraform separado, sin cómputo AWS) y Prod
  (reutilizando tryckers-media-storage y CloudFront con OAC), persistiendo
  object_key en PostgreSQL, limpiando referencias antiguas sin migración de
  archivos y actualizando Angular para avatares, banners e imágenes de posts.
ready_at: '2026-11-26'
---

# WI-008 — Implementar almacenamiento multimedia en Amazon S3 para Dev y Prod

> Idioma del proyecto: **español**.

**Tipo:** Feature  
**Prioridad:** High  
**Estado:** Completed  
**Repositorios / Módulos:** `tryckers-backend` (`core`), `tryckers-frontend` (`frontend`)  
**ADR relacionada:** `ADR-001-s3-compatible-storage-driver`

---

## Actor and outcome

- **Actor:** Usuarios de Tryckers (cargando/visualizando avatares, banners e imágenes de posts) y desarrolladores/operadores operando entornos Dev y Prod.
- **Outcome:** Todos los archivos multimedia administrados por Tryckers se almacenan en Amazon S3 (bucket Dev independiente en desarrollo local y bucket Prod privado detrás de CloudFront OAC en producción), PostgreSQL persiste únicamente `object_key`, las referencias heredadas se limpian de forma idempotente mostrando placeholders, y el backend queda 100% *stateless* entre redespliegues.

---

## Current behavior

- `infra/main.tf` declara `aws_s3_bucket.media_storage` (`tryckers-media-storage`, ya existente en `us-west-2` y vacío), sin cifrado explícito, sin bloqueo de acceso público, sin política IAM para EC2 ni CloudFront OAC.
- No existe entorno Terraform independiente para S3 Dev.
- `tryckers-backend` usa `LocalStorage` guardando archivos en `./uploads` y sirviéndolos estáticamente en `/uploads`.
- `users.avatar_url`, `users.banner_url`, `users.profile_picture`, `posts.image` y `comments.image` almacenan rutas `/uploads/...` o URLs externas; `CreatePostDto`/`UpdatePostDto` exigen `url` en `image`.
- `ProfilePage` solicita una URL de texto para la imagen de portada de publicaciones.

---

## Target behavior

1. **Terraform Dev (`infra/dev/`):** Estado independiente; crea únicamente bucket S3 `tryckers-media-storage-dev` (SSE-S3, Public Access Block, lifecycle `tmp/` + multipart cleanup, tags y política IAM de mínimo privilegio). Cero EC2, RDS, Amplify, ECR o CloudFront.
2. **Terraform Prod (`infra/s3-media.tf`, `infra/main.tf`):** Reutiliza `tryckers-media-storage` con SSE-S3, Public Access Block, lifecycle, política IAM `ec2_s3_media` en `ec2_role`, CloudFront OAC + comportamientos de caché `avatars/*`, `banners/*`, `posts/*` en la distribución CloudFront existente, y parámetros SSM (`/tryckers/storage/*`).
3. **Backend (`StorageService` + `S3Storage` + `LocalStorage`):**
   - Selección por `STORAGE_DRIVER` (`s3` | `local`), `AWS_REGION`, `S3_BUCKET`, `MEDIA_PUBLIC_BASE_URL`.
   - Claves canónicas únicas: `avatars/{userId}/{uuid}.{ext}`, `banners/{userId}/{uuid}.{ext}`, `posts/{userId}/{uuid}.{ext}`.
   - Entrega en Dev vía `GET /media/*key` (streaming desde S3 Dev); entrega en Prod vía CloudFront OAC (`https://api.tryckers.trycatch.tv/{object_key}`).
   - Endpoint `POST /api/v1/posts/image` para carga de imagen de portada de posts.
   - Limpieza de objeto reemplazado tras persistencia exitosa en BD.
4. **Base de datos:** Almacena únicamente `object_key` en los campos multimedia; migración SQL idempotente limpia referencias antiguas (no canónicas) a `''` preservando todos los demás datos y enlaces externos (`github_url`, `linkedin_url`, `pitch_video`, `posts.media_url`).
5. **Frontend Angular:** Avatares, banners y carga de imagen de portada de publicaciones por archivo vía backend Go (sin credenciales AWS en cliente), vistas previas, estados de carga/error y placeholders cuando no hay imagen.

---

## Entry points

- Infra: `infra/dev/`, `infra/s3-media.tf`, `infra/main.tf`, `scripts/deploy.sh`, `docker-compose.prod.yml`
- Backend: `POST/DELETE /api/v1/users/me/avatar`, `POST/DELETE /api/v1/users/me/banner`, `POST /api/v1/posts/image`, `GET /media/*key`
- Frontend: `ProfilePage` (`/profile/:username`), `Post` (`/profile/:username/post/:id`), `TryckersService`, `PostsService`

---

## End-to-end flow

1. Usuario selecciona imagen (avatar, banner o portada de post) en Angular → vista previa local → `multipart/form-data` al backend Go.
2. Backend valida MIME real (`jpeg`, `png`, `webp`) y tamaño → sube a S3 con clave única `{kind}/{userId}/{uuid}.{ext}` y `CacheControl: public, max-age=31536000, immutable`.
3. Backend persiste `object_key` en PostgreSQL y elimina el objeto anterior en S3 si existía.
4. Backend resuelve `object_key` → URL de entrega (`http://localhost:8080/media/{key}` en Dev o `https://api.tryckers.trycatch.tv/{key}` en Prod) y responde al frontend.
5. Frontend renderiza la imagen o el placeholder correspondiente si la referencia está vacía.

---

## Suggested Knowledge Level

**K3** — Feature transversal que abarca IaC Terraform (Dev y Prod), AWS S3 + CloudFront OAC, abstracción de almacenamiento en Go, migración idempotente de BD y actualización de componentes Angular.

---

## Impact analysis

| Superficie | Estado | Evidencia / Motivo |
| :--- | :--- | :--- |
| Infrastructure / Cloud | `affected` | Terraform `infra/dev/` (S3 Dev aislado) y `infra/s3-media.tf` + `main.tf` (S3 Prod, IAM EC2, CloudFront OAC, SSM). |
| Backend | `affected` | `StorageService`, `S3Storage`, `LocalStorage`, `MediaHandler`, `POST /api/v1/posts/image`, resolución de URLs, `config.go`. |
| Database | `affected` | Persistencia de `object_key` y migración idempotente de limpieza de referencias multimedia antiguas. |
| Frontend | `affected` | `TryckersService.mediaUrl()`, `PostsService.uploadPostImage()`, `ProfilePage` (upload de imagen de post + avatar/banner). |
| Operations / Release | `affected` | `docker-compose.prod.yml` y `scripts/deploy.sh` propagan `STORAGE_DRIVER`, `AWS_REGION`, `S3_BUCKET`, `MEDIA_PUBLIC_BASE_URL`. |
| Documentation | `affected` | Actualización de `infra/README.md`, `.env.example` y evidencias en `WI-008`. |

---

## Module coverage

- `core` (`tryckers-backend`): `affected`
- `frontend` (`../tryckers-frontend`): `affected`

---

## Scope confidence

**High** — Bucket productivo `tryckers-media-storage` verificado en AWS (`us-west-2`, vacío), Terraform v1.10.3 y AWS CLI operativos.

---

## Acceptance criteria

- [x] AC-01: Existe un bucket S3 independiente para Dev (`tryckers-media-storage-dev`).
- [x] AC-02: Existe un bucket S3 independiente para Prod, reutilizando `tryckers-media-storage`.
- [x] AC-03: Terraform administra los recursos sin provocar cambios destructivos inesperados (`0 to destroy` en Dev y Prod).
- [x] AC-04: Los buckets están cifrados (SSE-S3 `AES256`) y mantienen bloqueado el acceso público.
- [x] AC-05: Los permisos IAM respetan el principio de mínimo privilegio (`tryckers-dev-s3-media-policy` y `tryckers-ec2-s3-media`).
- [x] AC-06: Dev utiliza S3 sin desplegar EC2, RDS, Amplify, ECR ni CloudFront.
- [x] AC-07: Existe la abstracción `StorageService` con implementación `S3Storage` (`src/internal/services/storage/`).
- [x] AC-08: El backend local opera contra S3 Dev (`STORAGE_DRIVER=s3`, `GET /media/*key`).
- [x] AC-09: El backend productivo configurado para operar contra S3 Prod vía IAM Role y SSM.
- [x] AC-10: Los archivos nuevos se almacenan en S3 con claves canónicas únicas.
- [x] AC-11: Se mantienen las validaciones centralizadas de tipos (`jpeg`, `png`, `webp`), tamaños y permisos.
- [x] AC-12: La actualización y eliminación de archivos (avatares, banners e imágenes de posts) limpia objetos reemplazados tras persistir en BD.
- [x] AC-13: PostgreSQL almacena únicamente `object_key` (`avatars/...`, `banners/...`, `posts/...`).
- [x] AC-14: Migración idempotente (`CleanupLegacyMediaReferences` y `migrations/001_cleanup_legacy_media_references.sql`) limpia referencias multimedia antiguas.
- [x] AC-15: No se copian archivos antiguos a S3.
- [x] AC-16: No se implementa compatibilidad ni fallback para URLs antiguas (`/uploads/...` o URLs externas retornan `""` en `ResolveURL` y `mediaUrl`).
- [x] AC-17: Se conservan intactos los datos no relacionados con multimedia y enlaces externos (`github_url`, `linkedin_url`, `pitch_video`, `posts.media_url`).
- [x] AC-18: Los nuevos archivos quedan correctamente asociados a sus registros.
- [x] AC-19: Las imágenes nuevas se resuelven y visualizan correctamente desde S3 (Dev `/media/*key`, Prod CloudFront OAC).
- [x] AC-20: Los registros sin imagen muestran el placeholder correspondiente.
- [x] AC-21: Los avatares y banners permiten carga, reemplazo y eliminación con vista previa.
- [x] AC-22: El modal de creación/edición de publicaciones en `ProfilePage` usa carga de archivo (`POST /api/v1/posts/image`) en lugar de campo de URL de texto.
- [x] AC-23: No quedan dependencias funcionales de `/uploads` local para archivos nuevos.
- [x] AC-24: Los enlaces externos no relacionados con almacenamiento multimedia conservan su funcionamiento.
- [x] AC-25: Existen pruebas unitarias para `S3Storage` (`s3_storage_test.go`).
- [x] AC-26: Validadas operaciones reales de carga, lectura con cifrado `AES256` y eliminación contra `tryckers-media-storage-dev`.
- [x] AC-27: Infraestructura productiva desplegada con IAM EC2, S3 privado, CloudFront OAC (`EFOD1XGZJFR40`) y SSM Parameter Store.
- [x] AC-28: Las imágenes permanecen disponibles después de redesplegar el backend (almacenadas en S3, contenedor stateless). *(Validado por el usuario en S3 Dev + infraestructura Prod aplicada.)*
- [x] AC-29: La migración idempotente se ejecuta automáticamente en el arranque de `main.go` y queda versionada en `migrations/001_cleanup_legacy_media_references.sql`.
- [x] AC-30: Pruebas y compilaciones de backend (`17/17 PASS`) y frontend (`18/18 SUCCESS`) pasan sin regresiones.
- [x] AC-31: Documentación, `.env.example` y evidencias actualizadas en Kaddo.

## Verified Evidence

- **Terraform Dev (`infra/dev/`, estado independiente):**
  - `terraform apply`: **5 added, 0 changed, 0 destroyed**.
  - Bucket `tryckers-media-storage-dev`, Public Access Block, SSE-S3 `AES256`, Lifecycle (`tmp/` 1d + multipart 1d), IAM policy `arn:aws:iam::570743710918:policy/tryckers-dev-s3-media-policy`.
  - **Costo AWS Dev:** Solo almacenamiento S3 + peticiones mínimas de prueba (sin EC2, RDS, Amplify, ECR ni CloudFront; sin KMS dedicado).
- **Terraform Prod (`infra/s3-media.tf`, `infra/main.tf`):**
  - `terraform plan` & `apply`: **0 destroyed**; reutilizado `tryckers-media-storage` (`prevent_destroy = true`), agregados Public Access Block, SSE-S3 `AES256`, Versioning + expiración de versiones no actuales (30d), IAM policy `tryckers-ec2-s3-media` en `tryckers-ec2-role`, CloudFront OAC + `ordered_cache_behavior` (`avatars/*`, `banners/*`, `posts/*`) en distribución existente `EFOD1XGZJFR40`, y parámetros SSM (`/tryckers/storage/driver`, `/tryckers/storage/s3-bucket`, `/tryckers/storage/media-public-base-url`).
- **Smoke test real S3 Dev (AC-26):**
  - `PutObject` + `HeadObject` (`SSE: AES256`, `ContentLength: 15`) + `DeleteObject` en `tryckers-media-storage-dev` → **OK**.
- **Backend:** `go test -v ./...` → **17/17 PASS**, `go build ./...` → **PASS**.
- **Frontend:** `npx ng build --configuration development` → **PASS**, Karma → **18/18 SUCCESS** (incluyendo avatares S3 en navbar `HeaderComponent`, home `DashboardPage`, `ProfilePage` y carga de portada de posts).
- **Validación funcional manual (usuario):** Confirmado por el usuario el flujo completo con S3 Dev (perfil, navbar y home).

---

## Out of scope

- Migrar archivos antiguos a S3 o preservar compatibilidad con URLs multimedia antiguas.
- Crear EC2, RDS, Amplify, ECR o CloudFront para Dev.
- Procesamiento de videos o carga directa navegador → S3 (presigned PUT).

---

## Validation

1. Terraform Dev (`infra/dev`): `terraform fmt -check`, `terraform validate`, `terraform apply` (verificando bucket S3 Dev real).
2. Terraform Prod (`infra`): `terraform fmt -check`, `terraform validate`, `terraform plan` (0 recreaciones destructivas) y `terraform apply`.
3. Backend: `go test -v ./...`, `go build ./...`, más prueba real de carga, lectura y eliminación contra S3 Dev (`tryckers-media-storage-dev`).
4. Frontend: `npx ng build --configuration development` y `npx ng test --watch=false --browsers=ChromeHeadless`.

---

## Definition of Done

Entornos Dev y Prod operando con buckets S3 independientes, backend almacenando nuevos archivos en S3 con `object_key` en PostgreSQL, referencias antiguas limpiadas idempotentemente, tests pasando en ambos repos y evidencia documentada en Kaddo.

---

## Open Questions

---

## Related decisions

- `ADR-001-s3-compatible-storage-driver` (accepted)
