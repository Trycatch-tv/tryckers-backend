# Production S3 Media Security, IAM, CloudFront OAC, and SSM Configuration (WI-008)

resource "aws_s3_bucket_public_access_block" "media_storage" {
  bucket = aws_s3_bucket.media_storage.id

  block_public_acls       = true
  block_public_policy     = true
  ignore_public_acls      = true
  restrict_public_buckets = true
}

resource "aws_s3_bucket_server_side_encryption_configuration" "media_storage" {
  bucket = aws_s3_bucket.media_storage.id

  rule {
    apply_server_side_encryption_by_default {
      sse_algorithm = "AES256"
    }
  }
}

resource "aws_s3_bucket_versioning" "media_storage" {
  bucket = aws_s3_bucket.media_storage.id

  versioning_configuration {
    status = "Enabled"
  }
}

resource "aws_s3_bucket_lifecycle_configuration" "media_storage" {
  bucket = aws_s3_bucket.media_storage.id

  rule {
    id     = "cleanup-tmp-objects"
    status = "Enabled"

    filter {
      prefix = "tmp/"
    }

    expiration {
      days = 1
    }
  }

  rule {
    id     = "cleanup-noncurrent-and-multipart"
    status = "Enabled"

    filter {}

    noncurrent_version_expiration {
      noncurrent_days = 30
    }

    abort_incomplete_multipart_upload {
      days_after_initiation = 1
    }
  }
}

# Least-privilege S3 policy for EC2 backend role
resource "aws_iam_role_policy" "ec2_s3_media" {
  name = "${var.project_name}-ec2-s3-media"
  role = aws_iam_role.ec2_role.name

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Sid    = "ListProdMediaBucket"
        Effect = "Allow"
        Action = [
          "s3:ListBucket"
        ]
        Resource = aws_s3_bucket.media_storage.arn
      },
      {
        Sid    = "ReadWriteDeleteProdMediaObjects"
        Effect = "Allow"
        Action = [
          "s3:GetObject",
          "s3:PutObject",
          "s3:DeleteObject"
        ]
        Resource = "${aws_s3_bucket.media_storage.arn}/*"
      }
    ]
  })
}

# CloudFront Origin Access Control (OAC) for private S3 media delivery
resource "aws_cloudfront_origin_access_control" "media_oac" {
  name                              = "${var.project_name}-media-oac"
  description                       = "OAC for private S3 media bucket delivery via CloudFront"
  origin_access_control_origin_type = "s3"
  signing_behavior                  = "always"
  signing_protocol                  = "sigv4"
}

data "aws_cloudfront_cache_policy" "caching_optimized" {
  name = "Managed-CachingOptimized"
}

resource "aws_s3_bucket_policy" "media_cloudfront_oac" {
  bucket = aws_s3_bucket.media_storage.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Sid    = "AllowCloudFrontServicePrincipalReadOnly"
        Effect = "Allow"
        Principal = {
          Service = "cloudfront.amazonaws.com"
        }
        Action   = "s3:GetObject"
        Resource = "${aws_s3_bucket.media_storage.arn}/*"
        Condition = {
          StringEquals = {
            "AWS:SourceArn" = aws_cloudfront_distribution.api_distribution.arn
          }
        }
      }
    ]
  })

  depends_on = [aws_s3_bucket_public_access_block.media_storage]
}

# SSM Parameters for Prod Backend Storage Configuration
resource "aws_ssm_parameter" "storage_driver" {
  name        = "/${var.project_name}/storage/driver"
  description = "Active media storage driver (s3)"
  type        = "String"
  value       = "s3"
}

resource "aws_ssm_parameter" "s3_bucket" {
  name        = "/${var.project_name}/storage/s3-bucket"
  description = "Production S3 bucket name for media storage"
  type        = "String"
  value       = aws_s3_bucket.media_storage.bucket
}

resource "aws_ssm_parameter" "media_public_base_url" {
  name        = "/${var.project_name}/storage/media-public-base-url"
  description = "Public HTTPS base URL (CloudFront) for resolving S3 media objects"
  type        = "String"
  value       = local.effective_api_origin
}
