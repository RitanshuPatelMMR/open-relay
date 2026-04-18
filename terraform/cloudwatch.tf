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

# SNS topic for alerts
resource "aws_sns_topic" "alerts" {
  name = "${var.app_name}-alerts"
}

resource "aws_sns_topic_subscription" "email" {
  topic_arn = aws_sns_topic.alerts.arn
  protocol  = "email"
  endpoint  = var.alert_email
}

# Metric filter — count failed deliveries from worker logs
resource "aws_cloudwatch_log_metric_filter" "failed_deliveries" {
  name           = "${var.app_name}-failed-deliveries"
  log_group_name = aws_cloudwatch_log_group.worker.name
  pattern        = "dead letter"

  metric_transformation {
    name      = "FailedDeliveries"
    namespace = "OpenRelay"
    value     = "1"
  }
}

# Alarm — too many failed deliveries
resource "aws_cloudwatch_metric_alarm" "high_failure_rate" {
  alarm_name          = "${var.app_name}-high-failure-rate"
  comparison_operator = "GreaterThanThreshold"
  evaluation_periods  = 1
  metric_name         = "FailedDeliveries"
  namespace           = "OpenRelay"
  period              = 300
  statistic           = "Sum"
  threshold           = 5
  alarm_description   = "More than 5 dead letter events in 5 minutes"
  alarm_actions       = [aws_sns_topic.alerts.arn]
  treat_missing_data  = "notBreaching"
}