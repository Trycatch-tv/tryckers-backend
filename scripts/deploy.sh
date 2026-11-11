#!/usr/bin/env bash
set -euo pipefail

# Deploy script for Tryckers Backend on EC2 via SSM or local execution
# Usage: ./scripts/deploy.sh <IMAGE_TAG> [ECR_REPOSITORY_URL] [AWS_REGION]

IMAGE_TAG="${1:-latest}"
AWS_REGION="${3:-${AWS_REGION:-us-west-2}}"
ECR_REPOSITORY_URL="${2:-${ECR_REPOSITORY_URL:-570743710918.dkr.ecr.${AWS_REGION}.amazonaws.com/tryckers-backend}}"
DEPLOY_DIR="/opt/tryckers"

echo "=== Starting Tryckers Backend Deployment ==="
echo "Target Image: ${ECR_REPOSITORY_URL}:${IMAGE_TAG}"
echo "Region: ${AWS_REGION}"

# Ensure Docker is installed and running
if ! command -v docker &> /dev/null; then
  echo "Installing Docker..."
  sudo yum update -y
  sudo amazon-linux-extras install docker -y || sudo yum install docker -y
  sudo systemctl enable docker
  sudo systemctl start docker
  sudo usermod -aG docker ec2-user || true
fi

# Ensure Docker Compose plugin or binary is available
if ! docker compose version &> /dev/null; then
  echo "Installing Docker Compose..."
  sudo mkdir -p /usr/local/lib/docker/cli-plugins
  sudo curl -SL "https://github.com/docker/compose/releases/download/v2.24.5/docker-compose-linux-x86_64" -o /usr/local/lib/docker/cli-plugins/docker-compose
  sudo chmod +x /usr/local/lib/docker/cli-plugins/docker-compose
  sudo ln -sf /usr/local/lib/docker/cli-plugins/docker-compose /usr/bin/docker-compose || true
fi

# Ensure deployment directory exists
sudo mkdir -p "${DEPLOY_DIR}"
sudo chown -R $(whoami):$(whoami) "${DEPLOY_DIR}"

# Copy compose configuration if running from repo or create it
if [ -f "docker-compose.prod.yml" ]; then
  cp docker-compose.prod.yml "${DEPLOY_DIR}/docker-compose.prod.yml"
fi

cd "${DEPLOY_DIR}"

# Login to Amazon ECR using instance IAM credentials
echo "Logging into Amazon ECR..."
aws ecr get-login-password --region "${AWS_REGION}" | docker login --username AWS --password-stdin "${ECR_REPOSITORY_URL}"

# Fetch database URL from SSM Parameter Store if not present
if [ -z "${DATABASE_URL:-}" ]; then
  echo "Retrieving Database URL from SSM Parameter Store..."
  DATABASE_URL=$(aws ssm get-parameter --name "/tryckers/database/url" --with-decryption --region "${AWS_REGION}" --query "Parameter.Value" --output text 2>/dev/null || true)
  export DATABASE_URL
fi

# Pull the new image
echo "Pulling Docker image ${ECR_REPOSITORY_URL}:${IMAGE_TAG}..."
docker pull "${ECR_REPOSITORY_URL}:${IMAGE_TAG}"

# Deploy container
export ECR_REPOSITORY_URL
export IMAGE_TAG
echo "Updating container with Docker Compose..."
docker compose -f docker-compose.prod.yml up -d --no-deps api

# Health check verification
echo "Validating health check..."
MAX_ATTEMPTS=15
ATTEMPT=1
HEALTHY=false

while [ $ATTEMPT -le $MAX_ATTEMPTS ]; do
  if curl -sf http://localhost:8080/health > /dev/null; then
    echo "Health check passed on attempt ${ATTEMPT}!"
    HEALTHY=true
    break
  fi
  echo "Waiting for backend to be healthy (attempt ${ATTEMPT}/${MAX_ATTEMPTS})..."
  sleep 3
  ATTEMPT=$((ATTEMPT + 1))
done

if [ "$HEALTHY" = false ]; then
  echo "ERROR: Health check failed after ${MAX_ATTEMPTS} attempts!"
  docker compose -f docker-compose.prod.yml logs api --tail=50
  exit 1
fi

# Prune dangling images to save disk space
echo "Pruning unused images..."
docker image prune -f

echo "=== Deployment Successfully Completed ==="
echo "Deployed Image: ${ECR_REPOSITORY_URL}:${IMAGE_TAG}"
