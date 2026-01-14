# --- 1. 백엔드용 IAM Role ---
resource "aws_iam_role" "backend_role" {
  name = "GUSS-PRD-BACKEND-ROLE"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Action = "sts:AssumeRole"
      Effect = "Allow"
      Principal = { Service = "ec2.amazonaws.com" }
    }]
  })
}

# --- 2. 백엔드 정책 (SQS, SecretsManager 접근) ---
resource "aws_iam_role_policy" "backend_policy" {
  name = "GUSS-PRD-IAM-BACKEND-POLICY"
  role = aws_iam_role.backend_role.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = [
          "sqs:SendMessage",
          "sqs:ReceiveMessage",
          "sqs:DeleteMessage",
          "sqs:GetQueueAttributes",
          "secretsmanager:GetSecretValue",
          "kms:Decrypt"
        ]
        Effect   = "Allow"
        Resource = "*"
      },
      {
        Action = [
          "logs:CreateLogStream",
          "logs:PutLogEvents"
        ]
        Effect   = "Allow"
        Resource = "*"
      }
    ]
  })
}

# --- 3. EC2 인스턴스 프로파일 (이걸 인스턴스/ASG에 붙여야 함) ---
resource "aws_iam_instance_profile" "backend_profile" {
  name = "GUSS-PRD-BACKEND-PROFILE"
  role = aws_iam_role.backend_role.name
}