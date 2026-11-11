---
type: chore
id: WI-003
title: Automated Backend Deployment with Terraform, ECR & GitHub Actions
status: draft
work_type: chore
created_at: '2026-11-11'
source:
  type: manual
  inferred: false
generated_by: work-item-agent
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
summary: Implementar un pipeline CI/CD con GitHub Actions, Terraform, ECR y AWS SSM para desplegar automáticamente el backend en EC2.
---

# WI-003 — Automated Backend Deployment with Terraform, ECR & GitHub Actions

**Project:** Tryckers Backend  
**Type:** Infrastructure / DevOps / Vertical Slice  
**Priority:** High  
**Status:** Draft  
**Scope:** Backend deployment  
**Repository:** `github.com/Trycatch-tv/tryckers-backend`

## 1. Context
Tryckers es un proyecto open source desarrollado por la comunidad TryCatch.tv, cuyo propósito es facilitar el aprendizaje práctico mediante la implementación de una plataforma real.

Actualmente, el backend utiliza Go, Gin, GORM y PostgreSQL. La aplicación está preparada para ejecutarse mediante Docker y se contempla su despliegue sobre una instancia Amazon EC2.

Se requiere automatizar la publicación de nuevas versiones del backend, utilizando GitHub Actions y servicios AWS, manteniendo una infraestructura sencilla, segura, reproducible y eficiente en costos.

La configuración de los recursos AWS debe gestionarse mediante Terraform, evitando configuraciones manuales difíciles de reproducir.

## 2. Objective
Implementar un pipeline CI/CD que permita desplegar automáticamente el backend de Tryckers en Amazon EC2 después de integrar cambios en la rama `main`.

La solución debe utilizar:
- Terraform para gestionar infraestructura AWS.
- GitHub Actions para CI/CD.
- AWS IAM Identity Provider con OIDC para autenticación.
- Amazon ECR para almacenar las imágenes Docker.
- AWS Systems Manager (SSM) para ejecutar despliegues remotos.
- Docker Compose para administrar el contenedor en EC2.

El despliegue no debe requerir SSH, contraseñas ni credenciales AWS permanentes almacenadas en GitHub.

## 3. Target Architecture
```text
GitHub Repository
       |
       | Merge to main
       v
GitHub Actions
       |
       | Test + Build
       v
Docker Image
       |
       | OIDC Authentication
       v
AWS IAM Role
       |
       +----------------------+
       |                      |
       v                      v
   Amazon ECR          AWS Systems Manager
       |                      |
       |                      | Run Command
       |                      v
       +-----------------> Amazon EC2
                                 |
                                 v
                           Docker Compose
                                 |
                                 v
                           Go + Gin API
                                 |
                                 v
                             PostgreSQL
```

La base de datos y los recursos existentes deben mantenerse sin modificaciones innecesarias.

## 4. Functional Requirements

### FR-01 — Terraform ECR Repository
Incorporar a la configuración Terraform existente la definición del repositorio Amazon ECR que almacenará las imágenes Docker del backend.
- Nombre configurable, sugerido `tryckers-backend`.
- Privado con cifrado administrado por AWS.
- Tags de imágenes inmutables.
- Identificar cada imagen mediante el commit SHA.
- Política de lifecycle para limitar imágenes antiguas.
- Exponer la URL del repositorio mediante Terraform outputs.
- *Expected resources:* `aws_ecr_repository`, `aws_ecr_lifecycle_policy`.

### FR-02 — GitHub OIDC Authentication
Configurar la federación de identidad entre GitHub Actions y AWS IAM utilizando OpenID Connect.
- Registrar o reutilizar el proveedor OIDC de GitHub.
- Crear un IAM Role para despliegues restringido al repositorio y a la rama `main`.
- Validar `aud` y `sub`. Evitar access keys permanentes.
- *Expected resources:* `aws_iam_openid_connect_provider`, `aws_iam_role`, `aws_iam_policy`, `aws_iam_role_policy_attachment`.

### FR-03 — EC2 Instance Role
Configurar la instancia EC2 para interactuar con AWS Systems Manager y Amazon ECR.
- Políticas administradas: `AmazonSSMManagedInstanceCore`, `AmazonEC2ContainerRegistryReadOnly`.
- Terraform debe reutilizar la instancia existente sin reemplazarla ni destruirla.

### FR-04 — GitHub Actions CI Pipeline
Workflow para validar el backend en PRs hacia `main` (`.github/workflows/backend-ci.yml`).
- Configurar Go, descargar dependencias, ejecutar tests y validar build Go y Docker.

### FR-05 — Docker Image Publishing
Publicar imagen en ECR tras merge a `main`:
- Autenticación OIDC en AWS y en ECR.
- Construir imagen, etiquetar con commit SHA y publicar en ECR.

### FR-06 — Automated EC2 Deployment
Despliegue remoto mediante SSM Run Command:
- Identificar EC2 de destino y enviar comando de despliegue mediante SSM.
- Esperar y validar resultado. Sin requerir SSH desde GitHub Actions.

### FR-07 — Docker Compose Deployment
Preparar script de despliegue en EC2 (`scripts/deploy.sh`):
- Login a ECR con IAM Role de EC2, pull de imagen, actualización con `docker compose up -d --no-deps backend`.
- Idempotente y con manejo de errores.

### FR-08 — Deployment Configuration
Centralizar variables en GitHub Actions (`AWS_REGION`, `AWS_DEPLOY_ROLE_ARN`, `ECR_REPOSITORY_URL`, `EC2_INSTANCE_ID`).
- Cero credenciales sensibles en logs ni repositorio.

### FR-09 — Deployment Validation
Health check del backend (`GET /health`) y validación de estado del contenedor y commit SHA.

### FR-10 — Documentation
Documentación completa en inglés en el repositorio (arquitectura, variables, OIDC, troubleshooting y rollback).

## 5. Non-Functional Requirements
- **Security:** No SSH credentials, no long-lived AWS keys, least privilege IAM, private ECR.
- **Cost Efficiency:** Reutilizar EC2 existente, sin cómputo adicional innecesario (sin ALB, sin ECS/EKS).
- **Reliability:** Fallos en tests/build bloquean despliegue; validación vía SSM; rollback manual documentado.
- **Maintainability:** Declarativo, parametrizado y comprensible.

## 6. Out of Scope
- Nueva instancia EC2, migración a ECS/Fargate/K8s, ALB, Auto Scaling, CloudFront, certificados/dominios, cambios en PostgreSQL, múltiples entornos, terraform apply automático desde GHA, rollback automático.

## 7. Affected System Entities
- `sys:tryckers-backend`: affected
- `sys:postgres-db`: reviewed-not-affected
- `sys:tryckers-frontend`: reviewed-not-affected
- `sys:external-services`: reviewed-not-affected

## 8. Acceptance criteria
- [ ] AC-01: El repositorio ECR está definido y administrado mediante Terraform.
- [ ] AC-02: El repositorio ECR tiene tags inmutables y política de lifecycle.
- [ ] AC-03: Terraform expone la URL del repositorio ECR.
- [ ] AC-04: El proveedor OIDC de GitHub está configurado o reutilizado correctamente.
- [ ] AC-05: El rol de GitHub únicamente permite acceso desde el repositorio y contexto de despliegue autorizados.
- [ ] AC-06: GitHub Actions se autentica en AWS sin access keys permanentes.
- [ ] AC-07: EC2 puede acceder a SSM y descargar imágenes desde ECR.
- [ ] AC-08: El pipeline ejecuta pruebas y valida la construcción del backend.
- [ ] AC-09: Un merge a `main` publica una imagen Docker identificada por commit SHA.
- [ ] AC-10: GitHub Actions puede iniciar un despliegue mediante SSM.
- [ ] AC-11: EC2 actualiza el contenedor utilizando la imagen publicada en ECR.
- [ ] AC-12: El pipeline espera el resultado del despliegue y detecta fallos.
- [ ] AC-13: El backend responde correctamente después del despliegue.
- [ ] AC-14: Es posible identificar la versión actualmente desplegada.
- [ ] AC-15: El procedimiento manual de rollback está documentado y probado.
- [ ] AC-16: El despliegue no requiere SSH ni credenciales permanentes.
- [ ] AC-17: Terraform no recrea ni destruye la EC2 ni otros recursos existentes.
- [ ] AC-18: La documentación permite comprender y reproducir el proceso.
- [ ] AC-19: `terraform fmt`, `terraform validate` y la revisión de `terraform plan` son exitosos.
- [ ] AC-20: Los workflows de PR no pueden desplegar en producción y un contexto no autorizado no puede asumir el rol de despliegue.

## 9. Validation
1. Ejecutar `terraform fmt -check`, `terraform validate` y `terraform plan` en `infra/` verificando cero cambios destructivos en la EC2.
2. Validar que una Pull Request hacia `main` ejecute las pruebas de backend y validación de Docker build.
3. Validar que al fusionar a `main`, GitHub Actions asuma el rol OIDC, construya la imagen, la publique en ECR y ordene el despliegue a la EC2 vía SSM Run Command.
4. Validar el estado del contenedor backend y la respuesta del endpoint de health check.

## 10. Ownership
- `/infra/**/*.tf`
- `/.github/workflows/*.yml`
- `/scripts/*.sh`
- `/docker-compose.prod.yml`
