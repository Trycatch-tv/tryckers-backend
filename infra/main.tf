locals {
  frontend_fqdn        = "${var.frontend_subdomain}.${var.root_domain}"
  backend_fqdn         = "${var.backend_subdomain}.${var.root_domain}"
  effective_api_origin = var.api_public_origin != "" ? var.api_public_origin : "https://${local.backend_fqdn}"
  api_base_url         = "${trimsuffix(local.effective_api_origin, "/")}/api/v1"
  frontend_origin      = "https://${local.frontend_fqdn}"
}


# 1. IAM Role and Profile for EC2
resource "aws_iam_role" "ec2_role" {
  name = "${var.project_name}-ec2-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service = "ec2.amazonaws.com"
        }
      }
    ]
  })
}

resource "aws_iam_instance_profile" "ec2_profile" {
  name = "${var.project_name}-ec2-profile"
  role = aws_iam_role.ec2_role.name
}

# Attach SSM Managed Policy to EC2 Role
resource "aws_iam_role_policy_attachment" "ec2_ssm" {
  role       = aws_iam_role.ec2_role.name
  policy_arn = "arn:aws:iam::aws:policy/AmazonSSMManagedInstanceCore"
}

# Attach ECR Read-Only Policy to EC2 Role
resource "aws_iam_role_policy_attachment" "ec2_ecr_read" {
  role       = aws_iam_role.ec2_role.name
  policy_arn = "arn:aws:iam::aws:policy/AmazonEC2ContainerRegistryReadOnly"
}

# Attach SSM Parameter Store Read Policy for /tryckers/* to EC2 Role
resource "aws_iam_role_policy" "ec2_ssm_params" {
  name = "${var.project_name}-ec2-ssm-params"
  role = aws_iam_role.ec2_role.name

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "ssm:GetParameter",
          "ssm:GetParameters"
        ]
        Resource = [
          "arn:aws:ssm:${var.aws_region}:*:parameter/${var.project_name}/*"
        ]
      }
    ]
  })
}

# Data source for latest Amazon Linux 2 AMI
data "aws_ami" "amazon_linux_2" {
  most_recent = true
  owners      = ["amazon"]

  filter {
    name   = "name"
    values = ["amzn2-ami-hvm-*-x86_64-gp2"]
  }
}

# Data source for default VPC
data "aws_vpc" "default" {
  default = true
}

data "aws_security_group" "default" {
  vpc_id = data.aws_vpc.default.id
  name   = "default"
}

# Allow PostgreSQL access from backend EC2 to RDS in default SG
resource "aws_security_group_rule" "rds_ingress_from_backend" {
  type                     = "ingress"
  from_port                = 5432
  to_port                  = 5432
  protocol                 = "tcp"
  security_group_id        = data.aws_security_group.default.id
  source_security_group_id = aws_security_group.backend_sg.id
  description              = "Allow PostgreSQL access from backend EC2"
}

# AWS Managed Prefix List for CloudFront origin-facing servers
data "aws_ec2_managed_prefix_list" "cloudfront" {
  name = "com.amazonaws.global.cloudfront.origin-facing"
}

# Security Group for EC2 Backend
resource "aws_security_group" "backend_sg" {
  name        = "${var.project_name}-backend-sg"
  description = "Security group for backend EC2 instance"
  vpc_id      = data.aws_vpc.default.id

  ingress {
    description     = "HTTP API traffic from CloudFront only"
    from_port       = 8080
    to_port         = 8080
    protocol        = "tcp"
    prefix_list_ids = [data.aws_ec2_managed_prefix_list.cloudfront.id]
  }


  egress {
    description = "Allow all outbound traffic"
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name = "${var.project_name}-backend-sg"
  }
}

# 2. EC2 Instance for backend
resource "aws_instance" "backend" {
  ami                    = data.aws_ami.amazon_linux_2.id
  instance_type          = "t3.micro"
  iam_instance_profile   = aws_iam_instance_profile.ec2_profile.name
  key_name               = "tryckers-key"
  vpc_security_group_ids = [aws_security_group.backend_sg.id]

  tags = {
    Name = "${var.project_name}-backend"
  }

  lifecycle {
    ignore_changes = [
      ami
    ]
  }
}

# Elastic IP for EC2 backend (static public IP)
resource "aws_eip" "backend_eip" {
  instance = aws_instance.backend.id
  domain   = "vpc"

  tags = {
    Name = "${var.project_name}-backend-eip"
  }
}

# CloudFront Data Sources for Managed Policies
data "aws_cloudfront_cache_policy" "caching_disabled" {
  name = "Managed-CachingDisabled"
}

data "aws_cloudfront_origin_request_policy" "all_viewer_except_host" {
  name = "Managed-AllViewerExceptHostHeader"
}

# CloudFront HTTPS Distribution for Backend API
resource "aws_cloudfront_distribution" "api_distribution" {
  enabled             = true
  is_ipv6_enabled     = true
  comment             = "CloudFront HTTPS distribution for Tryckers Backend API"
  price_class         = "PriceClass_100"
  wait_for_deployment = false
  aliases             = [local.backend_fqdn]


  origin {
    domain_name = aws_eip.backend_eip.public_dns
    origin_id   = "EC2-${aws_instance.backend.id}"

    custom_origin_config {
      http_port                = 8080
      https_port               = 443
      origin_protocol_policy   = "http-only"
      origin_ssl_protocols     = ["TLSv1.2"]
      origin_keepalive_timeout = 5
      origin_read_timeout      = 30
    }
  }

  origin {
    domain_name              = aws_s3_bucket.media_storage.bucket_regional_domain_name
    origin_id                = "S3-${aws_s3_bucket.media_storage.id}"
    origin_access_control_id = aws_cloudfront_origin_access_control.media_oac.id
  }

  default_cache_behavior {
    allowed_methods  = ["DELETE", "GET", "HEAD", "OPTIONS", "PATCH", "POST", "PUT"]
    cached_methods   = ["GET", "HEAD", "OPTIONS"]
    target_origin_id = "EC2-${aws_instance.backend.id}"

    viewer_protocol_policy = "redirect-to-https"

    cache_policy_id          = data.aws_cloudfront_cache_policy.caching_disabled.id
    origin_request_policy_id = data.aws_cloudfront_origin_request_policy.all_viewer_except_host.id
  }

  ordered_cache_behavior {
    path_pattern           = "avatars/*"
    allowed_methods        = ["GET", "HEAD", "OPTIONS"]
    cached_methods         = ["GET", "HEAD"]
    target_origin_id       = "S3-${aws_s3_bucket.media_storage.id}"
    viewer_protocol_policy = "redirect-to-https"
    cache_policy_id        = data.aws_cloudfront_cache_policy.caching_optimized.id
  }

  ordered_cache_behavior {
    path_pattern           = "banners/*"
    allowed_methods        = ["GET", "HEAD", "OPTIONS"]
    cached_methods         = ["GET", "HEAD"]
    target_origin_id       = "S3-${aws_s3_bucket.media_storage.id}"
    viewer_protocol_policy = "redirect-to-https"
    cache_policy_id        = data.aws_cloudfront_cache_policy.caching_optimized.id
  }

  ordered_cache_behavior {
    path_pattern           = "posts/*"
    allowed_methods        = ["GET", "HEAD", "OPTIONS"]
    cached_methods         = ["GET", "HEAD"]
    target_origin_id       = "S3-${aws_s3_bucket.media_storage.id}"
    viewer_protocol_policy = "redirect-to-https"
    cache_policy_id        = data.aws_cloudfront_cache_policy.caching_optimized.id
  }

  restrictions {
    geo_restriction {
      restriction_type = "none"
    }
  }

  viewer_certificate {
    acm_certificate_arn      = aws_acm_certificate_validation.api.certificate_arn
    ssl_support_method       = "sni-only"
    minimum_protocol_version = "TLSv1.2_2021"
  }


  tags = {
    Name = "${var.project_name}-api-cf"
  }
}

# 3. RDS PostgreSQL Single AZ
resource "aws_db_instance" "database" {
  identifier           = "${var.project_name}-db"
  allocated_storage    = 20
  engine               = "postgres"
  engine_version       = "14"
  instance_class       = "db.t3.micro"
  username             = "postgres"
  password             = var.db_password
  parameter_group_name = "default.postgres14"
  skip_final_snapshot  = true
  multi_az             = false
}

# 4. Amplify for frontend
resource "aws_amplify_app" "frontend" {
  name = "${var.project_name}-frontend"

  environment_variables = {
    API_BASE_URL  = local.api_base_url
    _LIVE_UPDATES = "[{\"name\":\"Amplify CLI\",\"pkg\":\"@aws-amplify/cli\",\"type\":\"npm\",\"version\":\"latest\"}]"
  }

  lifecycle {
    ignore_changes = [
      repository,
      oauth_token,
      access_token,
      build_spec,
      custom_rule,
      iam_service_role_arn,
      auto_branch_creation_config,
      auto_branch_creation_patterns
    ]
  }
}

# 5. S3 for media storage
resource "aws_s3_bucket" "media_storage" {
  bucket = "${var.project_name}-media-storage"

  tags = {
    Name        = "${var.project_name}-media-storage"
    Project     = var.project_name
    Environment = "prod"
    ManagedBy   = "terraform"
  }

  lifecycle {
    prevent_destroy = true
  }
}

# 6. ECR for Docker images
resource "aws_ecr_repository" "backend_repo" {
  name                 = "${var.project_name}-backend"
  image_tag_mutability = "IMMUTABLE"

  image_scanning_configuration {
    scan_on_push = true
  }
}

resource "aws_ecr_lifecycle_policy" "backend_repo_policy" {
  repository = aws_ecr_repository.backend_repo.name

  policy = jsonencode({
    rules = [
      {
        rulePriority = 1
        description  = "Keep last 10 images"
        selection = {
          tagStatus   = "any"
          countType   = "imageCountMoreThan"
          countNumber = 10
        }
        action = {
          type = "expire"
        }
      }
    ]
  })
}

# 7. SSM Parameter Store
resource "aws_ssm_parameter" "db_url" {
  name  = "/${var.project_name}/database/url"
  type  = "SecureString"
  value = "postgres://postgres:${var.db_password}@${aws_db_instance.database.endpoint}/${var.project_name}"
}

resource "aws_ssm_parameter" "frontend_origin" {
  name        = "/${var.project_name}/frontend/origin"
  description = "Public origin of the frontend for backend CORS"
  type        = "String"
  value       = local.frontend_origin
}

# 8. GitHub Actions OIDC Authentication & Deploy Role
data "aws_iam_openid_connect_provider" "github" {
  url = "https://token.actions.githubusercontent.com"
}

resource "aws_iam_role" "github_deploy_role" {
  name = "${var.project_name}-github-deploy-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Principal = {
          Federated = data.aws_iam_openid_connect_provider.github.arn
        }
        Action = "sts:AssumeRoleWithWebIdentity"
        Condition = {
          StringEquals = {
            "token.actions.githubusercontent.com:aud" = "sts.amazonaws.com"
            "token.actions.githubusercontent.com:sub" = "repo:${var.github_repo}:ref:refs/heads/main"
          }
        }
      }
    ]
  })
}

resource "aws_iam_policy" "github_deploy_policy" {
  name        = "${var.project_name}-github-deploy-policy"
  description = "Policy allowing GitHub Actions to push images to ECR and execute deployments via SSM"

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Sid    = "ECRAuth"
        Effect = "Allow"
        Action = [
          "ecr:GetAuthorizationToken"
        ]
        Resource = "*"
      },
      {
        Sid    = "ECRPush"
        Effect = "Allow"
        Action = [
          "ecr:BatchCheckLayerAvailability",
          "ecr:GetDownloadUrlForLayer",
          "ecr:BatchGetImage",
          "ecr:PutImage",
          "ecr:InitiateLayerUpload",
          "ecr:UploadLayerPart",
          "ecr:CompleteLayerUpload"
        ]
        Resource = aws_ecr_repository.backend_repo.arn
      },
      {
        Sid    = "SSMDeployCommand"
        Effect = "Allow"
        Action = [
          "ssm:SendCommand"
        ]
        Resource = [
          aws_instance.backend.arn,
          "arn:aws:ssm:${var.aws_region}:*:document/AWS-RunShellScript"
        ]
      },
      {
        Sid    = "SSMDeployInvocation"
        Effect = "Allow"
        Action = [
          "ssm:GetCommandInvocation",
          "ssm:ListCommandInvocations",
          "ssm:DescribeInstanceInformation"
        ]
        Resource = "*"
      }
    ]
  })
}

resource "aws_iam_role_policy_attachment" "github_deploy_attach" {
  role       = aws_iam_role.github_deploy_role.name
  policy_arn = aws_iam_policy.github_deploy_policy.arn
}
