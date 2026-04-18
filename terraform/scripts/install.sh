#!/bin/bash
set -e

# Install Docker
yum update -y
yum install -y docker
systemctl start docker
systemctl enable docker
usermod -aG docker ec2-user

# Install Docker Compose
curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
chmod +x /usr/local/bin/docker-compose

# Install AWS CLI v2
yum install -y awscli

# Login to ECR
aws ecr get-login-password --region ${aws_region} | docker login --username AWS --password-stdin ${ecr_ingestion}

# Create app directory
mkdir -p /app
cd /app

# Write docker-compose.prod.yml
cat > /app/docker-compose.prod.yml <<EOF
version: '3.9'

services:
  ingestion:
    image: ${ecr_ingestion}:latest
    ports:
      - "8080:8080"
    environment:
      - DB_URL=${db_url}
      - REDIS_URL=${redis_url}
      - PORT=8080
    restart: always

  worker:
    image: ${ecr_worker}:latest
    environment:
      - DB_URL=${db_url}
      - REDIS_URL=${redis_url}
      - WORKER_COUNT=5
      - MAX_RETRY_ATTEMPTS=5
      - REQUEST_TIMEOUT_SECONDS=10
    restart: always

  api:
    image: ${ecr_api}:latest
    ports:
      - "8081:8081"
    environment:
      - DB_URL=${db_url}
      - REDIS_URL=${redis_url}
      - PORT=8081
    restart: always

  dashboard:
    image: ${ecr_dashboard}:latest
    ports:
      - "3000:3000"
    environment:
      - NEXT_PUBLIC_API_URL=http://$(curl -s http://169.254.169.254/latest/meta-data/public-ipv4):8081
      - NEXT_PUBLIC_WS_URL=ws://$(curl -s http://169.254.169.254/latest/meta-data/public-ipv4):8081
    restart: always
EOF

# Pull and start all services
docker-compose -f /app/docker-compose.prod.yml pull
docker-compose -f /app/docker-compose.prod.yml up -d