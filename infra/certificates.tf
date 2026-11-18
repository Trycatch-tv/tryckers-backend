# ACM Certificate for CloudFront Backend API (must be in us-east-1)
resource "aws_acm_certificate" "api" {
  provider          = aws.us_east_1
  domain_name       = local.backend_fqdn
  validation_method = "DNS"

  lifecycle {
    create_before_destroy = true
  }

  tags = {
    Name = "${var.project_name}-api-cert"
  }
}

resource "aws_acm_certificate_validation" "api" {
  provider                = aws.us_east_1
  certificate_arn         = aws_acm_certificate.api.arn
  validation_record_fqdns = [for record in netlify_dns_record.api_cert_validation : record.hostname]
}
