data "aws_ami" "amazon_linux" {
  most_recent = true
  owners      = ["amazon"]

  filter {
    name   = "name"
    values = ["al2023-ami-*-x86_64"]
  }
}

resource "aws_key_pair" "deployer" {
  key_name   = "${var.app_name}-deployer-key"
  public_key = var.ssh_public_key
}

resource "aws_instance" "app" {
  ami                    = data.aws_ami.amazon_linux.id
  instance_type          = "t2.micro"
  key_name               = aws_key_pair.deployer.key_name
  subnet_id              = aws_subnet.public[0].id
  vpc_security_group_ids = [aws_security_group.app.id]
  iam_instance_profile   = aws_iam_instance_profile.ec2_profile.name

  user_data = templatefile("${path.module}/scripts/install.sh", {
    aws_region       = var.aws_region
    app_name         = var.app_name
    db_url           = "postgres://openrelay:${var.db_password}@${aws_db_instance.postgres.address}:5432/openrelay?sslmode=disable"
    redis_url        = "redis://${aws_elasticache_cluster.redis.cache_nodes[0].address}:6379"
    ecr_ingestion    = aws_ecr_repository.ingestion.repository_url
    ecr_worker       = aws_ecr_repository.worker.repository_url
    ecr_api          = aws_ecr_repository.api.repository_url
    ecr_dashboard    = aws_ecr_repository.dashboard.repository_url
  })

  tags = {
    Name = "${var.app_name}-app-server"
  }
}