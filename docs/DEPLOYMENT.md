# Backend CI/CD Pipeline & Deployment Guide

This document describes the automated CI/CD pipeline for the Tryckers Backend, deployed to Amazon EC2 using Terraform, Amazon ECR, GitHub Actions, and AWS Systems Manager (SSM).

---

## Architecture Overview

```text
GitHub (Trycatch-tv/tryckers-backend)
        │
        ├── PR to main ───────> GitHub Actions (CI)
        │                         ├── Go Test (go test ./...)
        │                         ├── Go Build (go build ./src/cmd)
        │                         └── Docker Build verification
        │
        └── Push / Merge to main ──> GitHub Actions (CD)
                                      ├── Authenticate via AWS OIDC (no static keys)
                                      ├── Build & Push Docker image to Amazon ECR
                                      │     └── Tags: <commit-sha> and latest
                                      └── AWS SSM Run Command
                                            │
                                            ▼
                                      Amazon EC2 (i-007cfaf0508858c13)
                                            ├── IAM Instance Profile auth to ECR
                                            ├── Docker pull image (<commit-sha>)
                                            ├── Docker compose up -d --no-deps api
                                            └── Health check verification (GET /health)
```

---

## Key Infrastructure Resources (Terraform)

All AWS infrastructure is declaratively managed under `infra/`:

1. **Amazon ECR (`tryckers-backend`):**
   - Private repository with AWS-managed encryption (`AES256`).
   - Image tag immutability enabled (`IMMUTABLE`).
   - Lifecycle policy: keeps the last 10 images to control storage costs.
2. **GitHub Actions OIDC Deploy Role (`tryckers-github-deploy-role`):**
   - Reuses existing OIDC Provider (`arn:aws:iam::570743710918:oidc-provider/token.actions.githubusercontent.com`).
   - Strict trust policy restricted to `repo:Trycatch-tv/tryckers-backend:*` and audience `sts.amazonaws.com`.
   - Least privilege policy: push to ECR and `ssm:SendCommand` / `ssm:GetCommandInvocation` targeting the EC2 instance.
3. **EC2 IAM Role (`tryckers-ec2-role`):**
   - Managed policies: `AmazonSSMManagedInstanceCore` and `AmazonEC2ContainerRegistryReadOnly`.
   - Enables execution via SSM and ECR image pull without needing local AWS keys on EC2.

---

## GitHub Actions Workflows

### 1. CI Pipeline (`.github/workflows/backend-ci.yml`)
- **Trigger:** Pull Requests targeting `main`.
- **Steps:**
  - Installs Go 1.24.
  - Downloads module dependencies.
  - Runs unit tests (`go test -v -race ./src/...`).
  - Compiles the application binary.
  - Validates Dockerfile buildability (`docker/build-push-action` dry-run).

### 2. CD Pipeline (`.github/workflows/backend-cd.yml`)
- **Trigger:** Push / Merge to `main`.
- **Steps:**
  - Assumes IAM Role `arn:aws:iam::570743710918:role/tryckers-github-deploy-role` via GitHub OIDC token.
  - Logs into Amazon ECR.
  - Builds and tags the Docker image with `${{ github.sha }}` and `latest`.
  - Pushes the image to ECR.
  - Executes `scripts/deploy.sh` on the EC2 instance using AWS SSM Run Command.
  - Monitors SSM execution and fails the workflow if deployment or health check fails.

---

## Deployment & Verification

### Health Check Endpoint
The backend exposes a public health endpoint:
- **URL:** `http://<EC2_PUBLIC_IP>:8080/health`
- **Response:**
  ```json
  {
    "status": "healthy",
    "service": "tryckers-backend"
  }
  ```

---

## Manual Rollback Procedure

If a deployed version experiences issues in production, any previously published version can be restored in seconds without modifying infrastructure:

1. Identify the desired previous commit SHA from GitHub or Amazon ECR.
2. Execute the rollback command on EC2 via AWS CLI (or SSM Run Command):
   ```bash
   aws ssm send-command \
     --instance-ids "i-007cfaf0508858c13" \
     --document-name "AWS-RunShellScript" \
     --parameters 'commands=["/opt/tryckers/scripts/deploy.sh <PREVIOUS_COMMIT_SHA>"]'
   ```
3. Alternatively, connect via SSH and run:
   ```bash
   /opt/tryckers/scripts/deploy.sh <PREVIOUS_COMMIT_SHA>
   ```
4. Verify application status:
   ```bash
   curl http://localhost:8080/health
   ```
