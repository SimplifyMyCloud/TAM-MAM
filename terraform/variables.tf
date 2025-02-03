# variables.tf
variable "aws_region" {
  description = "AWS region"
  default     = "us-west-2"
}

variable "environment" {
  description = "Environment (production/staging)"
  default     = "smc-devs"
}

variable "db_username" {
  description = "Database username"
}

variable "db_password" {
  description = "Database password"
  sensitive   = true
}

variable "ami_id_backend" {
  description = "AMI ID for Backend EC2 instances"
}

variable "ami_id_frontend" {
  description = "AMI ID for Frontend EC2 instances"
}

# outputs.tf
output "alb_dns_name" {
  value = aws_lb.mam_alb.dns_name
}

output "database_endpoint" {
  value = aws_db_instance.mam_db.endpoint
}

output "redis_endpoint" {
  value = aws_elasticache_cluster.mam_redis.cache_nodes[0].address
}