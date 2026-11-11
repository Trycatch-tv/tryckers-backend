output "amplify_app_id" {
  value = aws_amplify_app.frontend.id
}

output "db_endpoint" {
  value = aws_db_instance.database.endpoint
}

output "ecr_repository_url" {
  value = aws_ecr_repository.backend_repo.repository_url
}

output "ec2_public_ip" {
  value = aws_instance.backend.public_ip
}

output "ec2_instance_id" {
  value       = aws_instance.backend.id
  description = "EC2 Instance ID for SSM deployment"
}

output "deploy_role_arn" {
  value       = aws_iam_role.github_deploy_role.arn
  description = "IAM Role ARN for GitHub Actions OIDC deploy"
}


