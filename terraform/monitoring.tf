# 운영 환경 CPU 부하 알람
resource "aws_cloudwatch_metric_alarm" "prd_cpu_alarm" {
  alarm_name          = "GUSS-PRD-CW-ALM"
  comparison_operator = "GreaterThanOrEqualToThreshold"
  evaluation_periods  = "2"
  metric_name         = "CPUUtilization"
  namespace           = "AWS/EC2"
  period              = "60"
  statistic           = "Average"
  threshold           = "80"
  alarm_description   = "운영 서버 CPU 부하 감지 (80% 초과)"
  alarm_actions       = [aws_sns_topic.prd_sns_alm.arn]

  dimensions = {
    AutoScalingGroupName = aws_autoscaling_group.prd_asg_backend.name
  }
}

# 개발 환경 상태 확인 알람
resource "aws_cloudwatch_metric_alarm" "dev_status_alarm" {
  alarm_name          = "GUSS-DEV-CW-ALM"
  comparison_operator = "GreaterThanThreshold"
  evaluation_periods  = "1"
  metric_name         = "StatusCheckFailed"
  namespace           = "AWS/EC2"
  period              = "60"
  statistic           = "Maximum"
  threshold           = "0"
  alarm_description   = "개발 서버 상태 이상 감지"
  alarm_actions       = [aws_sns_topic.dev_sns_alm.arn] # 아래 data_messaging에서 생성

  dimensions = {
    InstanceId = aws_instance.dev_backend.id
  }
}