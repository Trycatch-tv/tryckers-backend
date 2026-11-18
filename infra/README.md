# Tryckers Infrastructure

This directory contains the Terraform configuration to deploy and manage the AWS infrastructure and Netlify DNS records for the Tryckers platform.

## Architecture Overview

- **Authoritative DNS:** Netlify DNS (`trycatch.tv` zone reused via `data "netlify_dns_zone"`).
- **Frontend Custom Domain:** `https://tryckers.trycatch.tv` → AWS Amplify Hosting (`aws_amplify_domain_association`) with Amplify-managed TLS certificate.
- **Backend Custom Domain:** `https://api.tryckers.trycatch.tv` → Amazon CloudFront (`aws_cloudfront_distribution`) with ACM TLS 1.2+ certificate in `us-east-1`.
- **Backend Origin:** Amazon EC2 (`i-007cfaf0508858c13`, Elastic IP `44.224.188.254`) running Go/Gin in Docker Compose on port `8080`.
  - *Origin Security:* Public HTTPS terminates at CloudFront. CloudFront → EC2 origin traffic uses HTTP on port `8080`, restricted via `aws_security_group.backend_sg` to the AWS-managed CloudFront origin-facing prefix list (`com.amazonaws.global.cloudfront.origin-facing`).
- **Database & Storage:** RDS PostgreSQL (`tryckers-db`), S3 (`tryckers-media-storage`), ECR (`tryckers-backend`).
- **Configuration & Secrets:** SSM Parameter Store (`/tryckers/database/url`, `/tryckers/jwt/secret`, `/tryckers/frontend/origin`).

## File Structure

- `provider.tf`: AWS (`us-west-2` default + `us_east_1` alias for CloudFront ACM) and Netlify (`netlify/netlify ~> 0.4.4`) providers.
- `variables.tf`: Input variables for AWS, domains, and overrides.
- `main.tf`: Core compute, database, ECR, IAM, SSM, CloudFront distribution, and Amplify app.
- `amplify-domains.tf`: `aws_amplify_domain_association` for `tryckers.trycatch.tv`.
- `certificates.tf`: `aws_acm_certificate` and `aws_acm_certificate_validation` in `us-east-1` for `api.tryckers.trycatch.tv`.
- `netlify-dns.tf`: Reused Netlify DNS zone and CNAME records for Amplify, CloudFront, and certificate validation.
- `outputs.tf`: Exported infrastructure identifiers and public URLs.
- `terraform.tfvars.example`: Non-secret variable template.

## Prerequisites & Authentication

1. **AWS Credentials:** Configure AWS CLI or environment variables with permissions for EC2, RDS, ECR, IAM, SSM, Amplify, CloudFront, and ACM (`us-west-2` and `us-east-1`).
2. **Netlify Personal Access Token:** Export `NETLIFY_API_TOKEN` in your shell before running Terraform (never commit tokens to version control or `.tfvars` files):
   ```powershell
   # PowerShell (current session)
   $env:NETLIFY_API_TOKEN = "<your-netlify-personal-access-token>"
   ```
   ```bash
   # Bash / zsh
   export NETLIFY_API_TOKEN="<your-netlify-personal-access-token>"
   ```

## Input Variables & Outputs

### Key Variables
- `root_domain`: Root DNS zone in Netlify (default: `trycatch.tv`).
- `frontend_subdomain`: Subdomain prefix for Amplify (default: `tryckers`).
- `backend_subdomain`: Subdomain prefix for CloudFront API (default: `api.tryckers`).
- `api_public_origin`: Optional override; defaults to `https://<backend_subdomain>.<root_domain>`.
- `db_password`: Database admin password (sensitive).

### Key Outputs
- `frontend_origin`: `https://tryckers.trycatch.tv`
- `api_base_url`: `https://api.tryckers.trycatch.tv/api/v1`
- `cloudfront_domain_name`: CloudFront distribution hostname.
- `acm_certificate_arn`: ACM certificate ARN in `us-east-1`.
- `netlify_zone_id`: Reused Netlify DNS zone ID.
- `ec2_public_ip`, `ec2_instance_id`, `ecr_repository_url`, `deploy_role_arn`, `amplify_app_id`, `db_endpoint`.

## Usage

```bash
cd infra
terraform init
terraform fmt -check
terraform validate
terraform plan -var="db_password=YOUR_PASSWORD"
terraform apply -var="db_password=YOUR_PASSWORD"
```

## Validation & Troubleshooting

- **DNS Resolution:**
  ```bash
  dig tryckers.trycatch.tv CNAME
  dig api.tryckers.trycatch.tv CNAME
  ```
- **HTTPS & Health Checks:**
  ```bash
  curl -I https://tryckers.trycatch.tv
  curl -f https://api.tryckers.trycatch.tv/health
  ```
- **ACM Validation Timeout:** Ensure `NETLIFY_API_TOKEN` is valid and the `_*.api.tryckers.trycatch.tv` CNAME exists in Netlify DNS.
- **Amplify Domain Stuck in Pending:** Verify `wait_for_verification = false` is set on `aws_amplify_domain_association.frontend` and check both the routing CNAME (`tryckers.trycatch.tv`) and validation CNAME in Netlify DNS.
- **CORS Errors After Domain Change:** Trigger a backend redeploy (`scripts/deploy.sh`) so EC2 pulls the updated `/tryckers/frontend/origin` from SSM, and trigger an Amplify release build so Angular picks up the updated `API_BASE_URL`.
