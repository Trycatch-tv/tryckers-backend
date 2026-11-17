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
  description = "Public HTTPS origin of the backend API"
  type        = string
  default     = "https://api.tryckers.com"

  validation {
    condition     = can(regex("^https://", var.api_public_origin))
    error_message = "The api_public_origin must start with 'https://'."
  }
}

