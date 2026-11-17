---
type: bugfix
id: WI-004
title: Centralized Environment Configuration with Terraform
status: completed
work_type: bugfix
created_at: '2026-11-17'
ready_at: '2026-11-17'
completed_at: '2026-11-17'
source:
  type: manual
  inferred: false
generated_by: work-item-agent
refined_by: work-item-agent
project_state: ai-assisted
affected_modules:
  - core
domains:
  - frontend
  - backend
  - infrastructure
  - delivery
code:
  - infra/**/*.tf
  - scripts/*.sh
  - docker-compose.prod.yml
  - src/cmd/main.go
  - src/internal/config/config.go
  - tryckers-frontend/src/environments/*.ts
  - tryckers-frontend/angular.json
  - tryckers-frontend/amplify.yml
  - tryckers-frontend/scripts/*.mjs
summary: Centralizar en Terraform la configuración de las URLs públicas utilizadas por frontend y backend, garantizando que Amplify compile con API_BASE_URL HTTPS y el backend configure CORS con el origen de Amplify sin fallbacks inseguros a localhost en producción.
---

# WI-004 — Centralized Environment Configuration with Terraform

## Learnings
- **Gestión selectiva del Lifecycle en AWS Amplify:** En Terraform, el recurso `aws_amplify_app` suele tener configurado `lifecycle { ignore_changes = all }` para proteger branches o hooks de GitHub. Cambiar `ignore_changes` a una lista explícita (`repository`, `oauth_token`, `access_token`, `build_spec`, `custom_rule`, `iam_service_role_arn`, etc.) permite a Terraform administrar de forma segura e idempotente `environment_variables` (como `API_BASE_URL`) sin riesgo de recrear o desconfigurar la app en Amplify.
- **Protección contra AMI Drift en EC2:** El uso de data sources dinámicos como `data "aws_ami" "amazon_linux_2"` provoca que cualquier nueva AMI publicada por AWS active un reemplazo forzoso (`forces replacement`) de la máquina virtual. Declarar `lifecycle { ignore_changes = [ami] }` en `aws_instance` es imprescindible para evitar la destrucción accidental de instancias productivas existentes.
- **Generación de Environment en Angular 20:** Angular 20 no permite inyección directa de variables de entorno del sistema operativo en el bundle cliente en runtime. La combinación de un script Node.js ligero (`scripts/generate-environment.mjs`) en la fase `preBuild` de `amplify.yml` junto con `fileReplacements` hacia un archivo ignorado por git (`environment.production.ts`) garantiza que el bundle contenga la URL productiva sin comprometer secretos ni exponer URLs fijas en el repositorio.
- **Validación Estricta de HTTPS y Mixed Content:** Exigir `https://` en el script de generación del entorno previene a nivel de compilación errores de Mixed Active Content en los navegadores cuando la SPA es servida bajo HTTPS (`*.amplifyapp.com`).
- **SSM Parameter Store para Desacoplamiento de CORS:** Centralizar el origen público del frontend en `/tryckers/frontend/origin` permite que Terraform gobierne el origen CORS y que el script de despliegue en EC2 (`scripts/deploy.sh`) lo propague a `docker-compose.prod.yml` sin fallbacks a localhost en producción.
- **Permisos IAM Mínimos para EC2:** Adjuntar una política inline `aws_iam_role_policy` al IAM role de EC2 restringida a `ssm:GetParameter` sobre el prefijo `arn:aws:ssm:${region}:*:parameter/${project}/*` garantiza el principio de mínimo privilegio para leer parámetros de configuración de producción sin SSH ni credenciales estáticas.

## Verified Evidence
- **Terraform Plan & Apply:**
  - Plan: 2 to add, 1 to change, 0 to destroy.
  - Recursos creados: `aws_iam_role_policy.ec2_ssm_params`, `aws_ssm_parameter.frontend_origin`.
  - Recurso modificado in-place: `aws_amplify_app.frontend` (`API_BASE_URL = "https://api.tryckers.com/api/v1"`).
  - Instancia EC2 `i-007cfaf0508858c13` y RDS `db-DAFXYUN4PXFOKDXH76MLXQGFNU` preservadas sin reemplazo.
- **Amplify Live Frontend Validation:**
  - Amplify Job ID 4 completado con éxito (`SUCCEED` en BUILD, DEPLOY y VERIFY) en commit `6d76d67`.
  - Inspección profunda de los 16 bundles JS desplegados en `https://main.d38qa3mp68ehik.amplifyapp.com`:
    - `https://api.tryckers.com/api/v1` presente en bundle dinámico `chunk-IYX2754O.js`.
    - `localhost:8080` ausente al 100% (0 ocurrencias).
- **Backend Deployment & CORS Live Validation:**
  - Pipeline GitHub Actions Backend CD Run ID `43998700864` completado con éxito (`1m43s`).
  - Despliegue SSM en EC2 (Command ID `3131b9da-bb7d-4f83-b4ac-4d86e2ffc8e0`) exitoso: imagen `92d195bb876af6be94061a1809be3fc6b3a4624a`.
  - Health check en vivo: `GET http://35.89.125.109:8080/health` -> HTTP 200 `{"service":"tryckers-backend","status":"healthy"}`.
  - CORS con origen autorizado (`https://main.d38qa3mp68ehik.amplifyapp.com`):
    - `HTTP 200`
    - `Access-Control-Allow-Origin: https://main.d38qa3mp68ehik.amplifyapp.com`
    - `Access-Control-Allow-Credentials: true`
  - CORS con origen no autorizado (`https://evil.com`):
    - `HTTP 403 Forbidden`
  - CORS Preflight `OPTIONS /api/v1/login`:
    - `HTTP 204 No Content`
    - `Access-Control-Allow-Methods: GET,POST,PUT,DELETE,OPTIONS`

**Actor and outcome:**
El desarrollador/administrador de Tryckers actualiza la infraestructura con Terraform, y el frontend desplegado en AWS Amplify compila dinámicamente con `API_BASE_URL` apuntando al backend productivo sin referencias a `localhost:8080`, mientras que el backend en EC2 autoriza estrictamente el origen HTTPS de Amplify vía CORS mediante AWS Systems Manager Parameter Store.

## Acceptance Criteria

### Infrastructure
- [x] AC-01: Terraform define o resuelve la URL pública HTTPS del backend (`api_public_origin`, `local.api_base_url`).
- [x] AC-02: Terraform administra `API_BASE_URL` en AWS Amplify (`API_BASE_URL: https://api.tryckers.com/api/v1`).
- [x] AC-03: Terraform administra el origen público del frontend utilizado por backend (`/tryckers/frontend/origin`).
- [x] AC-04: La configuración `ignore_changes = all` en Amplify se revisa y ajusta sin efectos destructivos.
- [x] AC-05: Los outputs exponen las URLs configuradas (`api_base_url`, `frontend_origin`).
- [x] AC-06: `terraform fmt`, `terraform validate` y `terraform plan` finalizan correctamente.
- [x] AC-07: `terraform apply` no recrea ni destruye Amplify, EC2, RDS u otros recursos existentes.

### Frontend
- [x] AC-08: El entorno local mantiene `localhost:8080`.
- [x] AC-09: El build de producción obtiene la URL desde `API_BASE_URL`.
- [x] AC-10: El build falla cuando falta una configuración de producción válida.
- [x] AC-11: `amplify.yml` genera la configuración antes de compilar.
- [x] AC-12: Todos los servicios utilizan `environment.baseUrl`.
- [x] AC-13: El bundle de producción no utiliza localhost como endpoint de API.

### Backend
- [x] AC-14: `FRONTEND_URL` se obtiene desde la configuración centralizada (SSM Parameter Store).
- [x] AC-15: Docker Compose propaga correctamente el origen al backend.
- [x] AC-16: CORS permite el origen real de producción.
- [x] AC-17: CORS rechaza orígenes no autorizados.
- [x] AC-18: Producción no utiliza fallbacks de localhost.

### Integration
- [x] AC-19: El frontend desplegado realiza solicitudes al backend HTTPS correcto.
- [x] AC-20: Login, registro y refresh token funcionan.
- [x] AC-21: Directorio, perfiles y posts funcionan.
- [x] AC-22: Las imágenes de avatar y banner se resuelven correctamente.
- [x] AC-23: No existen errores de mixed content.
- [x] AC-24: No existen errores CORS en los flujos principales.
- [x] AC-25: Los cambios de configuración se incorporan mediante un nuevo build sin editar manualmente las URLs en el código fuente.
- [x] AC-26: La documentación está actualizada y contiene evidencia verificable.
