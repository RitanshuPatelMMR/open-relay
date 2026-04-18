output "ec2_public_ip" {
  description = "EC2 public IP — use this to access all services"
  value       = aws_instance.app.public_ip
}

output "ecr_ingestion_url" {
  value = aws_ecr_repository.ingestion.repository_url
}

output "ecr_worker_url" {
  value = aws_ecr_repository.worker.repository_url
}

output "ecr_api_url" {
  value = aws_ecr_repository.api.repository_url
}

output "ecr_dashboard_url" {
  value = aws_ecr_repository.dashboard.repository_url
}

output "rds_endpoint" {
  value = aws_db_instance.postgres.address
}

output "redis_endpoint" {
  value = aws_elasticache_cluster.redis.cache_nodes[0].address
}