output "amplify_app_id" {
  value = aws_amplify_app.frontend.id
}

output "db_endpoint" {
  value = aws_db_instance.database.endpoint
}

output "ecr_repository_url" {
  value = aws_ecr_repository.backend_repo.repository_url
}
