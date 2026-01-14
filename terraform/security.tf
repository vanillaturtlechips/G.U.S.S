# 1. Nginx 보안 그룹 (Public) - 관리대장 60, 64번
resource "aws_security_group" "prd_nginx_sg" {
  name        = "GUSS-PRD-NGINX-PUB-2A-SG"
  description = "Security group for Nginx servers"
  vpc_id      = aws_vpc.prd_vpc.id

  ingress {
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name = "GUSS-PRD-NGINX-PUB-2A-SG"
  }
}

# 2. 백엔드/WAS 보안 그룹 (Private) - 관리대장 62, 65, 71번 (참조 오류 해결부)
resource "aws_security_group" "prd_backend_sg" {
  name        = "GUSS-PRD-BACKEND-PRI-2A-SG"
  description = "Security group for Go Backend and Internal ALB"
  vpc_id      = aws_vpc.prd_vpc.id

  # Go 앱 포트 (Internal ALB 및 Nginx로부터의 통신)
  ingress {
    from_port   = 9000
    to_port     = 9000
    protocol    = "tcp"
    cidr_blocks = ["10.1.0.0/16"] # VPC 내부 통신 허용
  }

  # SSH 접속 (Bastion으로부터의 접속 대비)
  ingress {
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = ["10.1.0.0/16"]
  }

  # RDS 포트 (DB 통신을 위해 이 SG를 RDS에도 적용)
  ingress {
    from_port   = 3306
    to_port     = 3306
    protocol    = "tcp"
    cidr_blocks = ["10.1.0.0/16"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name = "GUSS-PRD-BACKEND-PRI-2A-SG"
  }
}