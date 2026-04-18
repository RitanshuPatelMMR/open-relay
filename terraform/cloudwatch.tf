resource "aws_cloudwatch_log_group" "ingestion" {
  name              = "/openrelay/ingestion"
  retention_in_days = 30

  tags = {
    Name = "${var.app_name}-ingestion-logs"
  }
}

resource "aws_cloudwatch_log_group" "worker" {
  name              = "/openrelay/worker"
  retention_in_days = 30

  tags = {
    Name = "${var.app_name}-worker-logs"
  }
}

resource "aws_cloudwatch_log_group" "api" {
  name              = "/openrelay/api"
  retention_in_days = 30

  tags = {
    Name = "${var.app_name}-api-logs"
  }
}

resource "aws_cloudwatch_log_group" "dashboard" {
  name              = "/openrelay/dashboard"
  retention_in_days = 30

  tags = {
    Name = "${var.app_name}-dashboard-logs"
  }
}