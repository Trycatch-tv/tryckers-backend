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

# 2. EC2 Instance for backend
resource "aws_instance" "backend" {
  ami                  = "ami-0c55b159cbfafe1f0" # Placeholder Amazon Linux 2 AMI
  instance_type        = "t3.micro"
  iam_instance_profile = aws_iam_instance_profile.ec2_profile.name

  tags = {
    Name = "${var.project_name}-backend"
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
  name       = "${var.project_name}-frontend"
  repository = "https://github.com/Trycatch-tv/tryckers-frontend"
  
  # Ensure OAuth token is provided in real deployments for private repos
}

# 5. S3 for media storage
resource "aws_s3_bucket" "media_storage" {
  bucket = "${var.project_name}-media-storage"
}

# 6. ECR for Docker images
resource "aws_ecr_repository" "backend_repo" {
  name                 = "${var.project_name}-backend"
  image_tag_mutability = "MUTABLE"
}

# 7. SSM Parameter Store
resource "aws_ssm_parameter" "db_url" {
  name  = "/${var.project_name}/database/url"
  type  = "SecureString"
  value = "postgres://postgres:${var.db_password}@${aws_db_instance.database.endpoint}/${var.project_name}"
}
