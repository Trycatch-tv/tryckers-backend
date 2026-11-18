---
type: chore
id: WI-005
title: Custom Domains and DNS Management with Terraform
status: completed
work_type: chore
created_at: '2026-11-18'
ready_at: '2026-11-18'
completed_at: '2026-11-18'
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
  - frontend
  - backend
  - delivery
code:
  - infra/provider.tf
  - infra/variables.tf
  - infra/main.tf
  - infra/outputs.tf
  - infra/netlify-dns.tf
  - infra/amplify-domains.tf
  - infra/certificates.tf
  - infra/README.md
  - infra/terraform.tfvars.example
summary: Configurar dominios personalizados (tryckers.trycatch.tv y api.tryckers.trycatch.tv) mediante Terraform reutilizando la zona DNS existente en Netlify, AWS Amplify, ACM (us-east-1) y CloudFront, actualizando API_BASE_URL y CORS sin migrar DNS a Route 53 ni destruir recursos existentes.
---

# WI-005 — Custom Domains and DNS Management with Terraform

## Learnings
- **Reutilización Segura de Zona DNS en Netlify (`netlify/netlify v0.4.4`):** El uso de `data "netlify_dns_zone"` permite consultar la zona existente (`62b0f7f0d39be72defd491a8` para `trycatch.tv`) y administrar exclusivamente los registros `netlify_dns_record` propios de Tryckers sin alterar otros registros DNS de TryCatch.tv ni modificar nameservers.
- **Autenticación de Netlify Provider:** `netlify/netlify` lee `NETLIFY_API_TOKEN` desde el entorno, evitando almacenar secretos en archivos `.tf` o `terraform.tfvars.example`.
- **Evitar Ciclos en `aws_amplify_domain_association`:** Establecer `wait_for_verification = false` permite que Terraform cree la asociación en Amplify, obtenga inmediatamente `certificate_verification_dns_record` y `sub_domain[*].dns_record`, y cree los registros CNAME en Netlify sin deadlocks.
- **ACM en `us-east-1` para CloudFront:** CloudFront exige certificados ACM emitidos en `us-east-1`, resuelto mediante un provider AWS con `alias = "us_east_1"`.
- **Seguridad de Origen CloudFront → EC2:** Restringir el ingreso del puerto `8080` en `aws_security_group.backend_sg` a la Prefix List administrada `com.amazonaws.global.cloudfront.origin-facing` (`pl-82a045eb`) bloquea el acceso HTTP directo por IP pública (`44.224.188.254:8080`), forzando el tráfico público por HTTPS vía CloudFront sin afectar `http://localhost:8080/health` dentro de EC2.

## Verified Evidence
- **Provider & Terraform:**
  - `netlify/netlify` `v0.4.4` fijado en `.terraform.lock.hcl`.
  - `terraform fmt -check` y `terraform validate` exitosos.
  - `terraform apply`: `7 added, 4 changed, 0 destroyed` (EC2, RDS y CloudFront existentes preservados).
- **DNS en Netlify (`trycatch.tv`, Zone ID `62b0f7f0d39be72defd491a8`):**
  - `tryckers.trycatch.tv` → `CNAME drvkw5s5116n1.cloudfront.net`
  - `api.tryckers.trycatch.tv` → `CNAME dge26rgrdunzk.cloudfront.net`
  - Registros CNAME de validación ACM y Amplify administrados en Terraform (`6afd2384e3e16d4e45bfa24a`, `6afd238dd1a2ca055a9eeff5`).
- **Certificados & CloudFront:**
  - ACM `arn:aws:acm:us-east-1:570743710918:certificate/950aadcc-c200-4b4c-9e94-bb62c59f95f7`: `Status: ISSUED`.
  - CloudFront `EFOD1XGZJFR40`: `Status: Deployed`, `aliases: ["api.tryckers.trycatch.tv"]`.
  - Amplify Domain Association `d38qa3mp68ehik/trycatch.tv`: subdominio `tryckers` con `"verified": true`.
- **Validación HTTPS, Bundle & CORS en Vivo:**
  - `GET https://tryckers.trycatch.tv` → `HTTP 200`, bundle `chunk-BVXUVWND.js` contiene `https://api.tryckers.trycatch.tv/api/v1` y `0` referencias a `localhost:8080`.
  - `GET https://api.tryckers.trycatch.tv/health` → `HTTP 200` (`ACAO: https://tryckers.trycatch.tv`).
  - `GET https://api.tryckers.trycatch.tv/api/v1/health` → `HTTP 200` (`ACAO: https://tryckers.trycatch.tv`).
  - `POST https://api.tryckers.trycatch.tv/api/v1/register` (`Origin: https://tryckers.trycatch.tv`) → `HTTP 400` con validación de modelo y cabecera `Access-Control-Allow-Origin: https://tryckers.trycatch.tv`.
  - `GET https://api.tryckers.trycatch.tv/health` (`Origin: https://evil.com`) → `HTTP 403`.
  - Acceso directo `http://44.224.188.254:8080/health` → bloqueado por Security Group (`timed out`).

## Acceptance Criteria

### Terraform and Netlify
- [x] AC-01: El provider oficial Netlify está integrado en Terraform.
- [x] AC-02: La zona `trycatch.tv` existente se reutiliza sin recrearla.
- [x] AC-03: Los registros DNS de Tryckers se administran mediante Terraform.
- [x] AC-04: Ningún registro DNS ajeno a Tryckers resulta modificado.
- [x] AC-05: Las credenciales de Netlify no se encuentran versionadas.

### Frontend
- [x] AC-06: Amplify tiene asociado el dominio `tryckers.trycatch.tv`.
- [x] AC-07: Los registros de validación del certificado están administrados mediante Terraform.
- [x] AC-08: El certificado de Amplify está validado.
- [x] AC-09: El frontend es accesible mediante HTTPS en el dominio personalizado.

### Backend
- [x] AC-10: CloudFront reconoce `api.tryckers.trycatch.tv` como dominio alternativo.
- [x] AC-11: El certificado TLS del backend está correctamente configurado en ACM.
- [x] AC-12: Los registros de validación ACM están administrados mediante Terraform.
- [x] AC-13: El DNS del backend apunta a la distribución CloudFront existente.
- [x] AC-14: El backend responde mediante HTTPS desde el dominio personalizado.
- [x] AC-15: No se introduce una distribución CloudFront adicional.

### Application Configuration
- [x] AC-16: Amplify utiliza `https://api.tryckers.trycatch.tv/api/v1` como `API_BASE_URL`.
- [x] AC-17: El backend utiliza `https://tryckers.trycatch.tv` como origen permitido.
- [x] AC-18: Las configuraciones se propagan sin modificaciones manuales de URLs en el código.
- [x] AC-19: El desarrollo local conserva su funcionamiento.
- [x] AC-20: Login, registro, refresh token, perfiles y posts funcionan correctamente.
- [x] AC-21: No existen errores de CORS ni mixed content.

### Delivery
- [x] AC-22: `terraform fmt` y `terraform validate` finalizan correctamente.
- [x] AC-23: `terraform plan` no incluye destrucción o recreación innecesaria de recursos.
- [x] AC-24: La configuración se aplica correctamente mediante Terraform.
- [x] AC-25: Los certificados, registros DNS y URLs están validados con evidencia.
- [x] AC-26: La documentación y el registro de Kaddo están actualizados.
