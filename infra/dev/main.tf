terraform {
  required_version = ">= 1.5.0"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

provider "aws" {
  region = var.aws_region
}

variable "aws_region" {
  description = "AWS Region for Dev S3 media bucket"
  type        = string
  default     = "us-west-2"
}

variable "project_name" {
  description = "Project name prefix"
  type        = string
  default     = "tryckers"
}

# Dev S3 bucket ONLY — no EC2, RDS, Amplify, ECR, or CloudFront in Dev
resource "aws_s3_bucket" "media_dev" {
  bucket = "${var.project_name}-media-storage-dev"

  tags = {
    Name        = "${var.project_name}-media-storage-dev"
    Project     = var.project_name
    Environment = "dev"
    ManagedBy   = "terraform"
  }
}

resource "aws_s3_bucket_public_access_block" "media_dev" {
  bucket = aws_s3_bucket.media_dev.id

  block_public_acls       = true
  block_public_policy     = true
  ignore_public_acls      = true
  restrict_public_buckets = true
}

resource "aws_s3_bucket_server_side_encryption_configuration" "media_dev" {
  bucket = aws_s3_bucket.media_dev.id

  rule {
    apply_server_side_encryption_by_default {
      sse_algorithm = "AES256"
    }
  }
}

resource "aws_s3_bucket_lifecycle_configuration" "media_dev" {
  bucket = aws_s3_bucket.media_dev.id

  rule {
    id     = "cleanup-tmp-and-multipart"
    status = "Enabled"

    filter {
      prefix = "tmp/"
    }

    expiration {
      days = 1
    }

    abort_incomplete_multipart_upload {
      days_after_initiation = 1
    }
  }

  rule {
    id     = "abort-all-incomplete-multipart"
    status = "Enabled"

    filter {}

    abort_incomplete_multipart_upload {
      days_after_initiation = 1
    }
  }
}

resource "aws_iam_policy" "dev_media_policy" {
  name        = "${var.project_name}-dev-s3-media-policy"
  description = "Least-privilege S3 media access policy for local Dev backend"

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Sid    = "ListDevMediaBucket"
        Effect = "Allow"
        Action = [
          "s3:ListBucket"
        ]
        Resource = aws_s3_bucket.media_dev.arn
      },
      {
        Sid    = "ReadWriteDeleteDevMediaObjects"
        Effect = "Allow"
        Action = [
          "s3:GetObject",
          "s3:PutObject",
          "s3:DeleteObject"
        ]
        Resource = "${aws_s3_bucket.media_dev.arn}/*"
      }
    ]
  })

  tags = {
    Project     = var.project_name
    Environment = "dev"
    ManagedBy   = "terraform"
  }
}

output "dev_media_bucket_name" {
  value       = aws_s3_bucket.media_dev.bucket
  description = "Dev S3 media bucket name"
}

output "dev_media_bucket_arn" {
  value       = aws_s3_bucket.media_dev.arn
  description = "Dev S3 media bucket ARN"
}

output "dev_media_policy_arn" {
  value       = aws_iam_policy.dev_media_policy.arn
  description = "Least-privilege IAM policy ARN for Dev S3 access"
}
