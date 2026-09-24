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
