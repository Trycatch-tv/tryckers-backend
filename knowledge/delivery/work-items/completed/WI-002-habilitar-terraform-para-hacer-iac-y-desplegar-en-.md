---
type: chore
id: WI-002
title: Habilitar terraform para hacer iac y desplegar en aws
status: completed
work_type: chore
created_at: '2026-09-24'
source:
  type: manual
  inferred: false
generated_by: kaddo-admin
refined_by: work-item-agent
project_state: ai-assisted
affected_modules:
  - core
domains:
  - infrastructure
code:
  - infra/**/*.tf
summary: Habilitar terraform para hacer iac y desplegar en aws
ready_at: '2026-09-24'
completed_at: '2026-09-24'
---

# Habilitar terraform para hacer iac y desplegar en aws

## Learnings
- **Infraestructura Dinámica:** Se optó por usar un bloque `data` para obtener la AMI de Amazon Linux 2 de forma dinámica según la región, evitando fallos de despliegue en distintas regiones (`us-west-2` frente a `us-east-1`).
- **Amplify:** Cuando se referencia un repositorio externo para la app de AWS Amplify a través de Terraform, es un requerimiento estricto de AWS proveer un OAuth Token válido. Por ahora se creó la instancia base sin vinculación directa al repositorio.

## Problem
La falta de Infraestructura como Código (IaC) dificulta el mantenimiento, la escalabilidad y la fiabilidad de los despliegues, haciendo los procesos manuales e ineficientes.

## Expected result
Crear la carpeta `infra` en `tryckers-backend` configurando con el proveedor de AWS los siguientes recursos: Amplify, EC2, RDS PostgreSQL Single AZ, ECR, SSM Parameter Store, IAM, S3.

## Impact
- backend: affected (nueva carpeta y manifiestos)
- frontend: reviewed-not-affected (no hay cambios en código, aunque su infraestructura será gestionada aquí)
- database: affected (creación a través de IaC)
- infrastructure: affected

**Affected system entities:**
- `sys:tryckers-backend`: affected
- `sys:tryckers-frontend`: affected
- `sys:postgres-db`: affected
- `sys:local-storage`: affected
- `sys:external-services`: reviewed-not-affected

## Acceptance criteria
- Tener la carpeta `infra` dentro del repositorio del backend.
- Tener la configuración de Terraform (archivos `.tf`) con los recursos de AWS especificados.
- Validar que los comandos base de `.tf` se ejecuten correctamente sin errores de sintaxis o planeación.
- Incluir tests, validación estática o validación de configuración para los archivos `.tf`.

## Validation
1. Moverse a la carpeta `infra/`.
2. Ejecutar `terraform init`.
3. Ejecutar `terraform validate`. Debe indicar que la configuración es válida.
4. Ejecutar `terraform plan` (asumiendo credenciales válidas). Debe listar los recursos correctos a crearse en AWS.

## Ownership
- `/infra/**/*.tf`

## Related decisions
- `related_decisions: [ADR-001-s3-compatible-storage-driver.md]`
