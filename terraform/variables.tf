variable "aws_region" {
  default = "ap-south-1"
}

variable "app_name" {
  default = "openrelay"
}

variable "environment" {
  default = "prod"
}

variable "db_password" {
  description = "RDS PostgreSQL password"
  type        = string
  sensitive   = true
}

variable "ssh_public_key" {
  description = "SSH public key for EC2 access"
  type        = string
}