variable "aws_region" {
  description = "AWS Region to deploy to"
  type        = string
  default     = "us-west-2"
}

variable "project_name" {
  description = "Name of the project"
  type        = string
  default     = "tryckers"
}

variable "db_password" {
  description = "Database admin password"
  type        = string
  sensitive   = true
}

variable "github_repo" {
  description = "GitHub repository (owner/repo) for OIDC trust"
  type        = string
  default     = "Trycatch-tv/tryckers-backend"
}

variable "api_public_origin" {
  description = "Public origin of the backend API. If empty, defaults to http://<ec2-elastic-ip>:8080."
  type        = string
  default     = ""

  validation {
    condition     = var.api_public_origin == "" || can(regex("^https?://", var.api_public_origin))
    error_message = "The api_public_origin must be empty or start with 'http://' or 'https://'."
  }
}

