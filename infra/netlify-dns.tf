# Reuse existing authoritative Netlify DNS zone for trycatch.tv
data "netlify_dns_zone" "trycatch" {
  name = var.root_domain
}

locals {
  amplify_cert_parts = split(" ", trimspace(aws_amplify_domain_association.frontend.certificate_verification_dns_record))
}

# Frontend (Amplify) routing CNAME: tryckers.trycatch.tv -> <amplify-cloudfront-target>
resource "netlify_dns_record" "frontend_cname" {
  for_each = {
    for sd in aws_amplify_domain_association.frontend.sub_domain : sd.prefix => sd
  }

  zone_id  = data.netlify_dns_zone.trycatch.id
  type     = "CNAME"
  hostname = "${each.value.prefix}.${var.root_domain}"
  value    = trimsuffix(split(" ", trimspace(each.value.dns_record))[2], ".")
  ttl      = 3600
}

# Frontend (Amplify) ACM certificate validation CNAME
resource "netlify_dns_record" "frontend_cert_validation" {
  zone_id  = data.netlify_dns_zone.trycatch.id
  type     = "CNAME"
  hostname = trimsuffix(local.amplify_cert_parts[0], ".")
  value    = trimsuffix(local.amplify_cert_parts[2], ".")
  ttl      = 3600
}

# Backend (ACM us-east-1) certificate validation CNAME
resource "netlify_dns_record" "api_cert_validation" {
  for_each = {
    for dvo in aws_acm_certificate.api.domain_validation_options : dvo.domain_name => {
      name  = dvo.resource_record_name
      type  = dvo.resource_record_type
      value = dvo.resource_record_value
    }
  }

  zone_id  = data.netlify_dns_zone.trycatch.id
  type     = each.value.type
  hostname = trimsuffix(each.value.name, ".")
  value    = trimsuffix(each.value.value, ".")
  ttl      = 3600
}

# Backend (CloudFront) routing CNAME: api.tryckers.trycatch.tv -> <cloudfront-distribution-domain>
resource "netlify_dns_record" "backend_cname" {
  zone_id  = data.netlify_dns_zone.trycatch.id
  type     = "CNAME"
  hostname = local.backend_fqdn
  value    = aws_cloudfront_distribution.api_distribution.domain_name
  ttl      = 3600
}
