# Tryckers Infrastructure

This directory contains the Terraform configuration to deploy the AWS infrastructure for the Tryckers platform.

## Resources Created
- **Amplify:** For the Tryckers frontend SPA.
- **EC2:** For the Tryckers backend API deployment.
- **RDS PostgreSQL:** Database for user profiles, posts, etc.
- **ECR:** Docker image repository for backend deployments.
- **S3:** Bucket for media storage.
- **SSM Parameter Store:** Secure string parameter for database connection URL.
- **IAM:** Associated roles and profiles.

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
