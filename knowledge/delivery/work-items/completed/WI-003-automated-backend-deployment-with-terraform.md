---
type: chore
id: WI-003
title: Automated Backend Deployment with Terraform, ECR & GitHub Actions
status: completed
work_type: chore
created_at: '2026-11-11'
ready_at: '2026-11-11'
reopened_at: '2026-11-11'
completed_at: '2026-11-11'
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

## Learnings
- **Reutilización de OIDC en AWS:** Al existir un OIDC Provider previamente registrado en la cuenta (`token.actions.githubusercontent.com`), declararlo como `data "aws_iam_openid_connect_provider"` en Terraform evita errores de duplicidad (`EntityAlreadyExistsException`) y permite adjuntar múltiples roles de despliegue por repositorio.
- **SSM Agent en EC2:** Para que una instancia EC2 existente sea descubierta inmediatamente por AWS Systems Manager (`aws ssm describe-instance-information`), se debe asociar la política administrada `AmazonSSMManagedInstanceCore` y reiniciar el servicio `amazon-ssm-agent` o reiniciar la instancia para que el agente refresque las credenciales del Instance Profile sin recrear la máquina.
- **Despliegues sin SSH:** AWS SSM Run Command (`AWS-RunShellScript`) permite ejecutar despliegues con Docker Compose de manera completamente remota, auditable y segura, prescindiendo de puertos SSH abiertos hacia GitHub Actions o internet.
- **Escapado SSM y Base64 en CI/CD:** El uso de `cat << EOF` con wrappers de bash en GitHub Actions puede inducir problemas de expansión de parámetros (como variables no definidas en el runner) y errores de salto de línea por CRLF de Windows. La codificación en base64 de los scripts y compose files dentro de un payload JSON generado por Python/jq garantiza un transporte 100% libre de efectos colaterales de shell.
- **Permisos IAM de ssm:GetCommandInvocation:** En AWS IAM, la acción `ssm:GetCommandInvocation` no soporta ARNs a nivel de recurso (como `instance/*`), requiriendo `Resource: "*"`. Mientras tanto, `ssm:SendCommand` sí debe restringirse al ARN de la EC2 y al documento `AWS-RunShellScript`.
- **Regla de Ingress de Base de Datos RDS:** Cuando RDS reside en el security group default y la EC2 en su propio security group, debe existir explícitamente una regla `aws_security_group_rule` para el puerto 5432 permitiendo el tráfico desde `backend_sg` hacia el SG de RDS.
- **Verificación Extremo a Extremo Validada:**
  - GitHub Actions CD Run exitoso: Run ID `42962752370` (commit `17a2de5bd4dc6f9fefb4abdff214e8400d57fe3a`).
  - Health check en vivo: `curl -i http://16.144.159.195:8080/health` -> HTTP 200 `{"service":"tryckers-backend","status":"healthy"}`.
  - Rollback real verificado vía SSM: Command ID `3485d710-917d-4782-8104-6dffe63f6be1` desplegando exitosamente commit `ae142ce1ee40b28d449797d5bcec7cba5c8e635e` y confirmando health check OK.
  - Roll forward completado: Command ID `0d17c6e5-7117-46bf-8010-d79113ae6204` restaurando commit `17a2de5bd4dc6f9fefb4abdff214e8400d57fe3a`.


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

## Discrepancy & Reopening Rationale
- **Identificación:** La primera ejecución del workflow CD (Run ID `42957001974`, commit `27d6f70088b8df4dacfa0e718656cfed039511f5`) falló durante el Docker Build debido a la exclusión de `docs/` en `.gitignore`. Por ende, no existe evidencia verificable de despliegue de extremo a extremo ni de rollback real.
- **Acción Kaddo:** Reapertura controlada de `WI-003` en `in-progress/` para ejecutar las correcciones de build, ECR immutability, SSM Run Command, Security Group y probar de punta a punta con evidencia en vivo.

**Acceptance criteria:**
- [x] AC-01: El backend compila y construye su imagen Docker desde un checkout limpio.
- [x] AC-02: Swagger no genera errores de dependencias durante el build.
- [x] AC-03: ECR utiliza una estrategia de imágenes inmutables y versionadas por SHA.
- [x] AC-04: El pipeline maneja correctamente reintentos de publicación.
- [x] AC-05: La aplicación utiliza correctamente `DATABASE_URL` en producción.
- [x] AC-06: EC2 puede autenticarse en ECR y descargar las imágenes.
- [x] AC-07: El rol OIDC está restringido al contexto autorizado (`refs/heads/main`).
- [x] AC-08: Los permisos IAM cumplen el principio de mínimo privilegio.
- [x] AC-09: El Security Group no expone SSH innecesariamente.
- [x] AC-10: La configuración de producción no utiliza secretos predeterminados inseguros.
- [x] AC-11: El pipeline CI termina exitosamente.
- [x] AC-12: El pipeline CD completa Build, Push, Deploy y Verify.
- [x] AC-13: El backend responde correctamente al health check después del despliegue.
- [x] AC-14: El despliegue corresponde al commit SHA publicado.
- [x] AC-15: El rollback manual se ejecuta y valida correctamente en la EC2.
- [x] AC-16: Terraform no realiza cambios destructivos sobre la infraestructura existente.
- [x] AC-17: La documentación técnica refleja el comportamiento real.
- [x] AC-18: Cada criterio de aceptación dispone de evidencia verificable.

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
