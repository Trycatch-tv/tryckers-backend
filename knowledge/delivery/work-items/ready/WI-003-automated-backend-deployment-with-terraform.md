---
type: chore
id: WI-003
title: Automated Backend Deployment with Terraform, ECR & GitHub Actions
status: ready
work_type: chore
created_at: '2026-11-11'
ready_at: '2026-11-11'
source:
  type: manual
  inferred: false
generated_by: work-item-agent
refined_by: work-item-agent
project_state: ai-assisted
affected_modules:
  - core
domains:
  - infrastructure
  - delivery
code:
  - infra/**/*.tf
  - .github/workflows/*.yml
  - scripts/*.sh
  - docker-compose.prod.yml
  - src/cmd/main.go
  - src/internal/api/routes/router.go
summary: Implementar un pipeline CI/CD completo con GitHub Actions, Terraform (ECR + OIDC + IAM), AWS SSM y Docker Compose para desplegar automáticamente el backend en EC2 tras merge en main.
---

# WI-003 — Automated Backend Deployment with Terraform, ECR & GitHub Actions

**Actor and outcome:**
El desarrollador/mantenedor del backend integra cambios aprobados en la rama `main` de `Trycatch-tv/tryckers-backend`, y el sistema compila, prueba, empaqueta la imagen Docker en Amazon ECR y despliega automáticamente el nuevo contenedor en Amazon EC2 mediante AWS Systems Manager (SSM) sin intervención manual, sin SSH y sin credenciales estáticas de AWS en GitHub.

**Current behavior:**
- La infraestructura base existe en AWS (EC2 `i-007cfaf0508858c13`, RDS Postgres `tryckers-db`, ECR `tryckers-backend` con mutabilidad activada y sin lifecycle policy, rol IAM de EC2 sin políticas asociadas `AmazonSSMManagedInstanceCore` ni permisos ECR).
- No existe pipeline de CI/CD en `.github/workflows/`.
- No hay proveedor OIDC ni IAM Role de despliegue para GitHub Actions configurado en Terraform para este proyecto.
- No existe archivo de compose para producción (`docker-compose.prod.yml`) ni script de despliegue idempotente (`scripts/deploy.sh`) en EC2.
- El backend en Go no expone actualmente un endpoint público de health check (`/health` o `/api/v1/health`).
- El despliegue se realiza manualmente o mediante SSH directo con archivos `.pem`.

**Target behavior:**
1. **Infraestructura como Código (Terraform):**
   - El repositorio ECR existente `tryckers-backend` se formaliza con tags inmutables y una política de ciclo de vida (`aws_ecr_lifecycle_policy` para retención de imágenes).
   - Se reutiliza o declara el proveedor OIDC de GitHub (`token.actions.githubusercontent.com`) y se crea un rol IAM de despliegue (`tryckers-github-deploy-role`) asumiendo credenciales temporales vía OIDC (`sts:AssumeRoleWithWebIdentity`), con permisos de mínimo privilegio para push a ECR y envío de comandos SSM a la EC2.
   - El rol IAM de la EC2 (`tryckers-ec2-role`) recibe las políticas administradas `AmazonSSMManagedInstanceCore` y `AmazonEC2ContainerRegistryReadOnly`.
   - Se preserva la EC2 existente sin destruirla ni recrearla.
2. **CI Pipeline (`.github/workflows/backend-ci.yml`):**
   - Se ejecuta en Pull Requests hacia `main`.
   - Valida el código: instala Go, descarga dependencias, ejecuta tests unitarios (`go test ./...`), compila el binario y verifica la construcción de la imagen Docker.
3. **CD Pipeline (`.github/workflows/backend-cd.yml`):**
   - Se ejecuta tras hacer merge en `main`.
   - Se autentica mediante OIDC con AWS.
   - Construye y etiqueta la imagen Docker con el commit SHA (y opcionalmente `latest`).
   - Sube la imagen a Amazon ECR.
   - Invoca AWS Systems Manager (SSM) Run Command hacia la instancia EC2 para ejecutar `scripts/deploy.sh` pasando el commit SHA.
   - Espera la finalización del comando SSM y valida su código de salida.
4. **Despliegue en EC2 (`scripts/deploy.sh` + `docker-compose.prod.yml`):**
   - Autentica Docker contra ECR usando el rol IAM de la propia EC2 (`aws ecr get-login-password`).
   - Descarga la imagen Docker del commit especificado.
   - Despliega el contenedor de backend usando `docker compose -f docker-compose.prod.yml up -d --no-deps api`.
   - Verifica la disponibilidad del servicio mediante el endpoint `/health`.
5. **Endpoint de Health Check:**
   - Se agrega la ruta pública `GET /health` (o `GET /api/v1/health`) retornando status 200 y el commit o versión activa.

**Entry points:**
- Git push / PR hacia `main` en `github.com/Trycatch-tv/tryckers-backend`.
- AWS Systems Manager Run Command (`AWS-RunShellScript`) ejecutado por GitHub Actions.
- HTTP `GET /health` en el puerto 8080 de la instancia EC2.

**End-to-end flow:**
```text
[Developer PR] ──> [GitHub Actions CI: go test + build + docker build]
                          │
                   (PR merge to main)
                          │
                          ▼
            [GitHub Actions CD (OIDC Auth)]
                          │
            ┌─────────────┴─────────────┐
            ▼                           ▼
   [Docker Build & Push]        [AWS SSM Run Command]
   Tag: <commit-sha>                    │
            │                           ▼
   [Amazon ECR] ────────────> [Amazon EC2 (deploy.sh)]
                                        │
                                [Docker Pull & Compose Up]
                                        │
                                [GET /health verification]
                                        │
                                [SSM Success -> Pipeline Green]
```

**Problem:**
La publicación y actualización del backend en Amazon EC2 dependen de intervenciones manuales o accesos SSH directos, careciendo de validación continua en Pull Requests, trazabilidad por commit SHA y un mecanismo automatizado, seguro y sin credenciales permanentes para desplegar en producción.

**Expected result:**
Un pipeline completamente automatizado en GitHub Actions que autentica mediante OIDC, almacena artefactos en ECR y actualiza el contenedor backend en EC2 mediante SSM con verificación de health check y documentación de rollback.

**Suggested Knowledge Level:** K2 (Integration / Architectural slice with external services)

**Impact analysis:**
- **backend:** `affected` (adición de endpoint de health check y configuración para producción)
- **infrastructure:** `affected` (actualización de Terraform: ECR lifecycle, OIDC GitHub, IAM deploy role, IAM EC2 policies)
- **delivery:** `affected` (creación de workflows `.github/workflows/backend-ci.yml` y `backend-cd.yml`, scripts de despliegue)
- **database:** `reviewed-not-affected` (la base de datos RDS PostgreSQL existente se preserva sin cambios)
- **frontend:** `reviewed-not-affected` (el frontend en AWS Amplify consume la API pero no es modificado por este Work Item)
- **operations:** `affected` (nuevo procedimiento de despliegue continuo y rollback documentado)

**Module coverage:**
- `core` (`tryckers-backend`): `affected`
- `tryckers-frontend`: `reviewed-not-affected`

**Scope unknowns:**
- Validación de si el daemon `docker` y el plugin `docker-compose` v2 ya se encuentran instalados y activos en el SO de la instancia EC2 (`Amazon Linux 2`). De no estarlo, el script de inicialización o `user_data` / aprovisionamiento inicial debe garantizar su instalación sin destruir la instancia.
- Verificación del agente SSM (`amazon-ssm-agent`) en la instancia EC2 una vez adjuntada la política `AmazonSSMManagedInstanceCore`.

**Scope confidence:** high
*Razones:* El proveedor OIDC de GitHub ya está registrado en la cuenta AWS (`arn:aws:iam::570743710918:oidc-provider/token.actions.githubusercontent.com`), el repositorio ECR ya existe, la EC2 está activa con su Security Group preparado en el puerto 8080 y la arquitectura objetivo está claramente acotada sin introducir cómputo extra ni costos adicionales.

**Acceptance criteria:**
- [ ] AC-01: El repositorio ECR está definido y administrado mediante Terraform con política de lifecycle (retención máxima de imágenes) y tags inmutables.
- [ ] AC-02: Terraform expone los outputs necesarios: `ecr_repository_url`, `deploy_role_arn`, `ec2_instance_id`.
- [ ] AC-03: El IAM Role para GitHub Actions autentica vía OIDC restringiendo estrictamente `sub` a `repo:Trycatch-tv/tryckers-backend:ref:refs/heads/main` (o contexto autorizado) y `aud` a `sts.amazonaws.com`.
- [ ] AC-04: El rol de GitHub Actions tiene únicamente permisos de mínimo privilegio para push en ECR y `ssm:SendCommand` sobre la instancia EC2 de Tryckers.
- [ ] AC-05: La instancia EC2 tiene asignadas las políticas `AmazonSSMManagedInstanceCore` y `AmazonEC2ContainerRegistryReadOnly`, registrándose exitosamente en AWS SSM.
- [ ] AC-06: Terraform no recrea, no destruye ni altera de forma destructiva la instancia EC2 existente ni la base de datos RDS.
- [ ] AC-07: El workflow de CI (`.github/workflows/backend-ci.yml`) valida el código en Pull Requests (Go test, build y Docker build) y bloquea la integración si falla.
- [ ] AC-08: El workflow de CD (`.github/workflows/backend-cd.yml`) se ejecuta en merge a `main`, autentica con OIDC y publica la imagen etiquetada con el commit SHA en ECR.
- [ ] AC-09: El CD envía la orden de despliegue a la EC2 mediante SSM Run Command, espera su ejecución y falla si el comando en EC2 falla.
- [ ] AC-10: El script `scripts/deploy.sh` en EC2 descarga la imagen correspondiente al commit SHA y actualiza el contenedor backend usando Docker Compose.
- [ ] AC-11: El backend expone un endpoint `/health` que responde HTTP 200 y permite verificar el estado saludable tras el despliegue.
- [ ] AC-12: El despliegue completo se realiza sin utilizar credenciales permanentes de AWS ni conexiones SSH.
- [ ] AC-13: Existe documentación técnica en inglés en el repositorio detallando la arquitectura, variables de GitHub, flujo de despliegue, troubleshooting y procedimiento de rollback manual por commit SHA.

**Out of scope:**
- Creación de nueva EC2 o migración a ECS / Fargate / Kubernetes.
- Creación de Application Load Balancer o CloudFront.
- Configuración de dominios personalizados y certificados SSL/TLS.
- Migraciones automatizadas de base de datos en el pipeline.
- Terraform apply automático desde GitHub Actions.
- Rollback automático (el rollback manual documentado sí forma parte del alcance).

**How to test it (validation):**
1. **Validación Terraform:**
   ```bash
   cd infra
   terraform fmt -check
   terraform validate
   terraform plan -var="db_password=trycatch_123456"
   ```
   *Verificar que no haya destrucción de la EC2 ni de RDS.*
2. **Validación CI en PR:**
   - Crear una rama y abrir PR con un cambio menor: verificar que ejecute tests y Docker build.
   - Introducir un test fallido simulado: verificar que bloquee el PR.
3. **Validación CD en Main:**
   - Fusionar PR hacia `main`.
   - Verificar autenticación OIDC exitosa en los logs de GitHub Actions.
   - Verificar que la imagen Docker aparezca en ECR con el tag del commit SHA.
   - Verificar que SSM Run Command se ejecute en la EC2 y complete exitosamente.
4. **Validación de la Aplicación:**
   ```bash
   curl http://<EC2_PUBLIC_IP>:8080/health
   ```
   *Debe responder 200 OK.*
5. **Simulación de Rollback:**
   - Ejecutar manualmente `scripts/deploy.sh <PREVIOUS_COMMIT_SHA>` en EC2 o vía SSM y verificar que el contenedor vuelva a la versión anterior.

**Definition of Done:**
1. Manifiestos de Terraform actualizados y aplicados en AWS sin alterar destructivamente recursos existentes.
2. Roles IAM (EC2 y GitHub OIDC) operativos y validados con mínimo privilegio.
3. Workflows de GitHub Actions (`backend-ci.yml` y `backend-cd.yml`) funcionales y probados.
4. Endpoint `/health` implementado y respondiendo en el backend.
5. Archivos `scripts/deploy.sh` y `docker-compose.prod.yml` versionados.
6. Despliegue de extremo a extremo validado en AWS.
7. Documentación en inglés (README / DEPLOYMENT.md) agregada al repositorio.
8. Cambios integrados mediante PR y mergeados a `main`.

**Open questions:**
- *(Ninguna bloqueante)*: Se confirmó la existencia previa del OIDC Provider de GitHub en la cuenta AWS y del repositorio ECR `tryckers-backend`.

**Suggested ownership (code globs):**
- `/infra/**/*.tf`
- `/.github/workflows/*.yml`
- `/scripts/*.sh`
- `/docker-compose.prod.yml`
- `/src/cmd/main.go`
- `/src/internal/api/routes/router.go`

**Related decisions:**
- Reutilizar el proveedor OIDC de GitHub ya registrado en IAM para evitar conflictos de recursos (`arn:aws:iam::570743710918:oidc-provider/token.actions.githubusercontent.com`).
- Uso de AWS SSM Run Command para prescindir de apertura de puertos SSH y llaves de acceso en CI/CD.
