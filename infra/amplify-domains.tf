# Custom domain association for AWS Amplify Frontend
resource "aws_amplify_domain_association" "frontend" {
  app_id                 = aws_amplify_app.frontend.id
  domain_name            = var.root_domain
  wait_for_verification  = false
  enable_auto_sub_domain = false

  sub_domain {
    branch_name = "main"
    prefix      = var.frontend_subdomain
  }
}
