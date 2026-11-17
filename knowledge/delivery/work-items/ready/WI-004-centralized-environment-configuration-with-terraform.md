---
type: bugfix
id: WI-004
title: Centralized Environment Configuration with Terraform
status: ready
work_type: bugfix
created_at: '2026-11-17'
ready_at: '2026-11-17'
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

**Type:** Bugfix / Infrastructure / Configuration  
**Priority:** High  
**Status:** Ready  
**Domains:** Frontend, Backend, Infrastructure, Delivery  

**Repositories:**
- `Trycatch-tv/tryckers-backend` (`tryckers-backend/`)
- `Trycatch-tv/tryckers-frontend` (`tryckers-frontend/`)

## 1. Context

Tryckers utiliza Angular 20 desplegado mediante AWS Amplify y un backend Go/Gin desplegado en Amazon EC2 mediante Docker Compose, Amazon ECR, GitHub Actions y AWS Systems Manager.

Actualmente, el frontend desplegado en AWS Amplify (`https://main.d38qa3mp68ehik.amplifyapp.com`) continúa apuntando a `http://localhost:8080/api/v1`, debido a que `src/environments/environment.ts` conserva una URL de desarrollo y la configuración de `production` en `angular.json` no implementa `fileReplacements` para entornos productivos.

Los servicios Angular (`AuthService`, `PostsService`, `TryckersService`) ya consumen `environment.baseUrl` de forma centralizada sin hardcodear URLs en componentes.

Adicionalmente:
- El backend configura CORS mediante `FRONTEND_URL`.
- `docker-compose.prod.yml` mantiene un fallback inseguro a `http://localhost:4200`.
- Terraform administra el recurso `aws_amplify_app.frontend`, pero actualmente contiene `lifecycle { ignore_changes = all }`.
- Terraform ya expone información de EC2 y otros recursos mediante outputs (`outputs.tf`).
- El frontend servido por HTTPS necesita consumir una API disponible mediante HTTPS para evitar errores de Mixed Content en navegadores modernos.

Se requiere corregir estas condiciones sin rediseñar la arquitectura ni incorporar servicios AWS innecesarios o costosos para el MVP.

## 2. Objective

Centralizar en Terraform la configuración de las URLs públicas utilizadas por frontend y backend, garantizando que cada entorno consuma los endpoints correctos:

1. Terraform define o resuelve la URL pública HTTPS del backend y la expone vía output (`api_base_url`).
2. Terraform administra `API_BASE_URL` dentro del recurso AWS Amplify (`aws_amplify_app.frontend`) sin destruir ni recrear la app.
3. Amplify incorpora esa variable durante el build de Angular invocando `scripts/generate-environment.mjs` previo a `ng build`.
4. Angular utiliza la URL de producción mediante `fileReplacements` hacia `src/environments/environment.production.ts`.
5. Terraform centraliza el origen público del frontend en AWS Systems Manager Parameter Store (`/tryckers/frontend/origin`).
6. El backend recibe `FRONTEND_URL` desde SSM Parameter Store durante la ejecución de `scripts/deploy.sh` y Docker Compose, eliminando fallbacks a localhost en producción.
7. El desarrollo local mantiene sus direcciones `localhost` (`localhost:8080` y `localhost:4200`).
8. Los cambios son reproducibles, trazables y administrados mediante infraestructura como código.

## 3. Target Architecture

```text
                    Terraform
                        |
             +----------+----------+
             |                     |
             v                     v
       Backend API URL       Frontend Origin
             |                     |
             v                     v
      AWS Amplify App       SSM Parameter Store
             |                     |
       API_BASE_URL           FRONTEND_URL
             |                     |
             v                     v
       Angular Build       EC2 Deployment Script
             |                     |
             v                     v
       Angular Frontend       Docker Compose
             |                     |
             | HTTPS               v
             +----------------> Go + Gin API
                                      |
                                      v
                                PostgreSQL
```

- **Terraform** administra la configuración de infraestructura en `tryckers-backend/infra`.
- **AWS Amplify** construye y publica el frontend alojado en `Trycatch-tv/tryckers-frontend`.
- **GitHub Actions** continúa construyendo y desplegando el backend utilizando ECR y SSM (`deploy.sh`).
- No se requiere que GitHub Actions del backend compile o despliegue Angular.

## 4. Functional Requirements

### FR-01 — Terraform API URL Configuration
- Actualizar `infra/main.tf`, `infra/variables.tf` y `infra/outputs.tf` en `tryckers-backend`.
- Definir la variable `api_public_origin` (descripción: "Public HTTPS origin of the backend API", tipo: string).
- Construir `local.api_base_url = "${trimsuffix(var.api_public_origin, "/")}/api/v1"`.
- Validar formato HTTPS y evitar URLs fijas a `localhost` para producción.
- Evitar dependencias circulares en Terraform.

### FR-02 — AWS Amplify Environment Variables
- Modificar el recurso existente `aws_amplify_app.frontend` en `infra/main.tf` para incorporar `API_BASE_URL`.
- Reemplazar `lifecycle { ignore_changes = all }` por una lista selectiva de exclusiones (`repository`, `oauth_token`, `access_token`, `build_spec`, `custom_rule`, etc.) para permitir que Terraform actualice `environment_variables` sin reiniciar configuraciones de consola.
- Preservar variables existentes como `_LIVE_UPDATES`.
- No destruir ni recrear la aplicación Amplify (`d38qa3mp68ehik`).

### FR-03 — Terraform Outputs
- Exponer en `infra/outputs.tf`:
  - `api_base_url`: URL completa de la API (`local.api_base_url`).
  - `frontend_origin`: Origen público del frontend derivado del dominio de Amplify (`https://main.d38qa3mp68ehik.amplifyapp.com`) o variable de entrada.

### FR-04 — Frontend URL Configuration
- Mantener `src/environments/environment.ts` y `environment.development.ts` con `http://localhost:8080/api/v1` para desarrollo local.
- Configurar en `angular.json` bajo `configurations.production` la sustitución `fileReplacements` de `src/environments/environment.ts` por `src/environments/environment.production.ts`.
- No versionar URLs productivas hardcodeadas ni archivos generados de despliegue en Git (`.gitignore`).

### FR-05 — Environment Generation Script
- Crear `tryckers-frontend/scripts/generate-environment.mjs`.
- Leer `process.env.API_BASE_URL`, verificar que sea válida, exigir HTTPS en producción y generar `src/environments/environment.production.ts`.
- Fallar de manera controlada con código de salida 1 si falta la variable o no cumple con HTTPS.

### FR-06 — AWS Amplify Build Integration
- Modificar `tryckers-frontend/amplify.yml` en la fase `preBuild` para ejecutar `node scripts/generate-environment.mjs` antes de `npm run build`.
- Confirmar que `baseDirectory` en artifacts apunte a `dist/tryckers-frontend/browser`.

### FR-07 — Centralized Frontend API Consumption
- Auditar y verificar que todos los servicios (`AuthService`, `PostsService`, `TryckersService`, interceptores) consuman `environment.baseUrl`.
- Confirmar ausencia total de URLs hardcodeadas a `localhost:8080`.

### FR-08 — Backend Frontend Origin Configuration
- Crear el parámetro en AWS Systems Manager Parameter Store: `/tryckers/frontend/origin`.
- Asignar el valor no sensible correspondiente al origen del frontend (`https://main.d38qa3mp68ehik.amplifyapp.com`).
- Garantizar que el IAM role de EC2 (`tryckers-ec2-role`) disponga de permisos para leer este parámetro.
- Actualizar `scripts/deploy.sh` para obtener `/tryckers/frontend/origin` y exportar `FRONTEND_URL`.

### FR-09 — Backend CORS & Production Config
- En `docker-compose.prod.yml`, requerir `FRONTEND_URL=${FRONTEND_URL:?FRONTEND_URL is required}` eliminando el fallback a `http://localhost:4200`.
- En `src/cmd/main.go` y `src/internal/config/config.go`, asegurar que `FRONTEND_URL` sea validado y que CORS responda adecuadamente con las cabeceras requeridas.

### FR-10 — HTTPS and Connectivity Constraint
- Identificar y documentar el estado de terminación SSL/HTTPS para la API de EC2.
- La aplicación frontend desplegada en Amplify (`HTTPS`) requiere consumir un backend HTTPS para evitar errores de Mixed Content.
- Si no existe aún un dominio o certificado SSL asignado a EC2, documentarlo como pre-requisito y definir el mecanismo de terminación HTTPS sin incorporar recursos costosos no previstos en el MVP.

### FR-11 — Deployment Coordination
- Coordinar los cambios entre ambos repositorios de forma secuencial y controlada para evitar desalineación de orígenes entre frontend y backend.

### FR-12 — Technical Documentation
- Actualizar la documentación técnica en inglés en ambos repositorios (`infra/README.md`, `tryckers-frontend/README.md`, scripts y especificaciones).

## 5. Implementation Phases

- **Phase 1 — Audit (Completada en refinamiento):**
  - Verificado estado de Amplify App `d38qa3mp68ehik` (`defaultDomain: d38qa3mp68ehik.amplifyapp.com`, branch `main`).
  - Verificado frontend: `angular.json` carece de `fileReplacements` en producción.
  - Verificado backend: EC2 `i-007cfaf0508858c13` escucha en `HTTP:8080`.
  - Detectado riesgo en Terraform: `aws_instance.backend` cambiaría por actualización de AMI; se requiere `lifecycle { ignore_changes = [ami] }`.
- **Phase 2 — Terraform Configuration:**
  - Agregar `api_public_origin`, `local.api_base_url`, `local.frontend_origin`.
  - Actualizar `aws_amplify_app.frontend` con `API_BASE_URL` y ajuste de `lifecycle`.
  - Crear `aws_ssm_parameter.frontend_origin` (`/tryckers/frontend/origin`).
  - Agregar `ignore_changes = [ami]` a `aws_instance.backend`.
  - Actualizar outputs en `infra/outputs.tf`.
- **Phase 3 — Frontend Configuration (`tryckers-frontend`):**
  - Crear `scripts/generate-environment.mjs`.
  - Actualizar `angular.json` con `fileReplacements` en producción.
  - Actualizar `amplify.yml` con el paso preBuild.
  - Validar builds local y productivo (`npm run build`).
- **Phase 4 — Backend Configuration (`tryckers-backend`):**
  - Actualizar `scripts/deploy.sh` para recuperar `/tryckers/frontend/origin`.
  - Actualizar `docker-compose.prod.yml` para exigir `FRONTEND_URL`.
  - Validar CORS y compatibilidad local.
- **Phase 5 — End-to-End Validation:**
  - Aplicar cambios en Terraform.
  - Desplegar frontend en Amplify y backend en EC2.
  - Validar conectividad, login, registro, perfiles y ausencia de errores CORS/Mixed Content.

## 6. Acceptance Criteria

### Infrastructure
- [ ] AC-01: Terraform define o resuelve la URL pública HTTPS del backend.
- [ ] AC-02: Terraform administra `API_BASE_URL` en AWS Amplify.
- [ ] AC-03: Terraform administra el origen público del frontend utilizado por backend.
- [ ] AC-04: La configuración `ignore_changes = all` en Amplify se revisa y ajusta sin efectos destructivos.
- [ ] AC-05: Los outputs exponen las URLs configuradas (`api_base_url`, `frontend_origin`).
- [ ] AC-06: `terraform fmt`, `terraform validate` y `terraform plan` finalizan correctamente.
- [ ] AC-07: `terraform apply` no recrea ni destruye Amplify, EC2, RDS u otros recursos existentes.

### Frontend
- [ ] AC-08: El entorno local mantiene `localhost:8080`.
- [ ] AC-09: El build de producción obtiene la URL desde `API_BASE_URL`.
- [ ] AC-10: El build falla cuando falta una configuración de producción válida.
- [ ] AC-11: `amplify.yml` genera la configuración antes de compilar.
- [ ] AC-12: Todos los servicios utilizan `environment.baseUrl`.
- [ ] AC-13: El bundle de producción no utiliza localhost como endpoint de API.

### Backend
- [ ] AC-14: `FRONTEND_URL` se obtiene desde la configuración centralizada (SSM Parameter Store).
- [ ] AC-15: Docker Compose propaga correctamente el origen al backend.
- [ ] AC-16: CORS permite el origen real de producción.
- [ ] AC-17: CORS rechaza orígenes no autorizados.
- [ ] AC-18: Producción no utiliza fallbacks de localhost.

### Integration
- [ ] AC-19: El frontend desplegado realiza solicitudes al backend HTTPS correcto.
- [ ] AC-20: Login, registro y refresh token funcionan.
- [ ] AC-21: Directorio, perfiles y posts funcionan.
- [ ] AC-22: Las imágenes de avatar y banner se resuelven correctamente.
- [ ] AC-23: No existen errores de mixed content.
- [ ] AC-24: No existen errores CORS en los flujos principales.
- [ ] AC-25: Los cambios de configuración se incorporan mediante un nuevo build sin editar manualmente las URLs en el código fuente.
- [ ] AC-26: La documentación está actualizada y contiene evidencia verificable.
