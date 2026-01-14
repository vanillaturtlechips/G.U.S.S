# ==========================================
# 1. 운영 환경 (PRD) - Bastion Host
# ==========================================
resource "aws_instance" "prd_bastion" {
  ami                    = var.prd_ami_nginx # 베스천 전용 AMI가 없을 경우 대체
  instance_type          = "t3.micro"
  subnet_id              = aws_subnet.prd_nginx_pub_2a.id
  vpc_security_group_ids = [aws_security_group.prd_nginx_sg.id]

  tags = {
    Name = "GUSS-PRD-BASTION-2A"
  }
}

# ==========================================
# 2. 운영 환경 (PRD) - 로드밸런서 (ALB)
# ==========================================

# 외부 로드밸런서 (사용자 -> Nginx)
resource "aws_lb" "prd_ex_alb" {
  name               = "GUSS-PRD-EX-ALB"
  internal           = false
  load_balancer_type = "application"
  security_groups    = [aws_security_group.prd_nginx_sg.id]
  subnets            = [
    aws_subnet.prd_nginx_pub_2a.id,
    aws_subnet.prd_nginx_pub_2c.id
  ]
}

# 내부 로드밸런서 (Nginx -> Go Backend)
resource "aws_lb" "prd_in_alb" {
  name               = "GUSS-PRD-IN-ALB"
  internal           = true
  load_balancer_type = "application"
  security_groups    = [aws_security_group.prd_backend_sg.id]
  subnets            = [
    aws_subnet.prd_backend_pri_2a.id,
    aws_subnet.prd_backend_pri_2c.id
  ]
}

# ==========================================
# 3. 운영 환경 (PRD) - 타겟 그룹 및 리스너
# ==========================================

# Nginx 타겟 그룹 (관리대장 30번)
resource "aws_lb_target_group" "prd_tg_nginx" {
  name        = "GUSS-PRD-TG-NGINX"
  port        = 80
  protocol    = "HTTP"
  vpc_id      = aws_vpc.prd_vpc.id
  target_type = "instance"

  health_check {
    path = "/dashboard.tsx"
  }
}

# Backend 타겟 그룹 (관리대장 29번)
resource "aws_lb_target_group" "prd_tg_backend" {
  name        = "GUSS-PRD-TG-BACKEND"
  port        = 9000
  protocol    = "HTTP"
  vpc_id      = aws_vpc.prd_vpc.id
  target_type = "instance"

  health_check {
    path = "/health"
  }
}

# 외부 리스너 연결
resource "aws_lb_listener" "prd_ex_alb_http" {
  load_balancer_arn = aws_lb.prd_ex_alb.arn
  port              = "80"
  protocol          = "HTTP"

  default_action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.prd_tg_nginx.arn
  }
}

# 내부 리스너 연결
resource "aws_lb_listener" "prd_in_alb_backend" {
  load_balancer_arn = aws_lb.prd_in_alb.arn
  port              = "9000"
  protocol          = "HTTP"

  default_action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.prd_tg_backend.arn
  }
}

# ==========================================
# 4. 운영 환경 (PRD) - 오토스케일링 (ASG)
# ==========================================

# Nginx ASG (관리대장 40번)
resource "aws_autoscaling_group" "prd_asg_nginx" {
  name                = "GUSS-PRD-ASG-NGINX"
  desired_capacity    = 2
  max_size            = 2
  min_size            = 1
  vpc_zone_identifier = [
    aws_subnet.prd_nginx_pub_2a.id,
    aws_subnet.prd_nginx_pub_2c.id
  ]
  target_group_arns   = [aws_lb_target_group.prd_tg_nginx.arn]

  launch_template {
    id      = aws_launch_template.prd_lt_nginx.id
    version = "$Latest"
  }
}

# Backend ASG (관리대장 41번)
resource "aws_autoscaling_group" "prd_asg_backend" {
  name                = "GUSS-PRD-ASG-BACKEND"
  desired_capacity    = 2
  max_size            = 2
  min_size            = 1
  vpc_zone_identifier = [
    aws_subnet.prd_backend_pri_2a.id,
    aws_subnet.prd_backend_pri_2c.id
  ]
  target_group_arns   = [aws_lb_target_group.prd_tg_backend.arn]

  launch_template {
    id      = aws_launch_template.prd_lt_backend.id
    version = "$Latest"
  }
}

# 시작 템플릿 - Nginx
resource "aws_launch_template" "prd_lt_nginx" {
  name_prefix   = "GUSS-PRD-LT-NGINX"
  image_id      = var.prd_ami_nginx
  instance_type = "t3.micro"
  
  network_interfaces {
    security_groups = [aws_security_group.prd_nginx_sg.id]
  }
}

# 시작 템플릿 - Backend (IAM 권한 포함)
resource "aws_launch_template" "prd_lt_backend" {
  name_prefix   = "GUSS-PRD-LT-BACKEND"
  image_id      = var.prd_ami_backend
  instance_type = "t3.micro"
  
  iam_instance_profile {
    name = aws_iam_instance_profile.backend_profile.name
  }

  network_interfaces {
    security_groups = [aws_security_group.prd_backend_sg.id]
  }
}

# ==========================================
# 5. 개발 환경 (DEV) - EC2 인스턴스
# ==========================================

# 개발용 Nginx (관리대장 9번)
resource "aws_instance" "dev_nginx" {
  ami                    = var.dev_ami_nginx
  instance_type          = "t3.micro"
  subnet_id              = aws_subnet.dev_nginx_pub_2a.id
  vpc_security_group_ids = [aws_security_group.prd_nginx_sg.id] # 동일 보안그룹 재사용 가능

  tags = {
    Name = "GUSS-DEV-NGINX-2A"
  }
}

# 개발용 Backend (관리대장 10번)
resource "aws_instance" "dev_backend" {
  ami                    = var.prd_ami_backend
  instance_type          = "t3.micro"
  subnet_id              = aws_subnet.dev_backend_pri_2a.id
  vpc_security_group_ids = [aws_security_group.prd_backend_sg.id]
  
  # 개발 환경에서도 SQS/SecretsManager 접근이 필요할 경우 IAM 추가
  iam_instance_profile = aws_iam_instance_profile.backend_profile.name

  tags = {
    Name = "GUSS-DEV-BACKEND-2A"
  }
}

# 개발용 베스천 (관리대장 8번)
resource "aws_instance" "dev_bastion" {
  ami                    = var.prd_ami_nginx
  instance_type          = "t3.micro"
  subnet_id              = aws_subnet.dev_nginx_pub_2a.id
  vpc_security_group_ids = [aws_security_group.prd_nginx_sg.id]

  tags = {
    Name = "GUSS-DEV-BASTION-2A"
  }
}