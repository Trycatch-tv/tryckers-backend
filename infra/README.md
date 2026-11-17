# Tryckers Infrastructure

This directory contains the Terraform configuration to deploy the AWS infrastructure for the Tryckers platform.

## Resources Created
- **Amplify:** For the Tryckers frontend SPA with managed environment variables (`API_BASE_URL`).
- **EC2:** For the Tryckers backend API deployment.
- **RDS PostgreSQL:** Database for user profiles, posts, etc.
- **ECR:** Docker image repository for backend deployments.
- **S3:** Bucket for media storage.
- **SSM Parameter Store:** Parameters for database URL (`/tryckers/database/url`) and frontend origin (`/tryckers/frontend/origin`).
- **IAM:** Associated roles and profiles, including EC2 SSM parameter read policy and GitHub Actions OIDC deploy role.

## Environment Variables & Outputs

### Input Variables
- `api_public_origin`: Public HTTPS origin of the backend API (e.g. `https://api.tryckers.com`).
- `project_name`: Name prefix for project resources (default: `tryckers`).
- `aws_region`: AWS Region (default: `us-west-2`).
- `db_password`: Database admin password.
- `github_repo`: Target GitHub repository for OIDC trust.

### Outputs
- `api_base_url`: Public backend API URL (`${api_public_origin}/api/v1`).
- `frontend_origin`: Public origin of the frontend (`https://main.<amplify_app_domain>`).
- `ec2_instance_id`: EC2 Instance ID for SSM deployment.
- `ec2_public_ip`: Public IP of the EC2 instance.
- `ecr_repository_url`: ECR repository URL.
- `deploy_role_arn`: GitHub Actions OIDC deploy role ARN.
- `amplify_app_id`: AWS Amplify App ID.
- `db_endpoint`: RDS PostgreSQL endpoint.

## Usage

1. Initialize Terraform
   ```bash
   terraform init
   ```

2. Validate Configuration
   ```bash
   terraform validate
   ```

3. Plan Deployment
   ```bash
   terraform plan -var="db_password=YOUR_PASSWORD"
   ```

4. Apply
   ```bash
   terraform apply -var="db_password=YOUR_PASSWORD"
   ```
