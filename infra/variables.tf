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
  description = "Optional override for the public origin of the backend API. If empty, defaults to https://<backend_subdomain>.<root_domain>."
  type        = string
  default     = ""

  validation {
    condition     = var.api_public_origin == "" || can(regex("^https?://", var.api_public_origin))
    error_message = "The api_public_origin must be empty or start with 'http://' or 'https://'."
  }
}

variable "root_domain" {
  description = "Root DNS domain managed in Netlify"
  type        = string
  default     = "trycatch.tv"
}

variable "frontend_subdomain" {
  description = "Subdomain prefix for the frontend application"
  type        = string
  default     = "tryckers"
}

variable "backend_subdomain" {
  description = "Subdomain prefix for the backend API"
  type        = string
  default     = "api.tryckers"
}


