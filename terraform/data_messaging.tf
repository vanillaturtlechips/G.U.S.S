# ==========================================
# 1. 운영 환경 (PRD) - 데이터 및 메시징
# ==========================================

# DB 서브넷 그룹
resource "aws_db_subnet_group" "prd_db_group" {
  name       = "guss-prd-db-group"
  subnet_ids = [
    aws_subnet.prd_backend_pri_2a.id,
    aws_subnet.prd_backend_pri_2c.id
  ]

  tags = {
    Name = "GUSS-PRD-DB-GROUP"
  }
}

# RDS MySQL 인스턴스
resource "aws_db_instance" "prd_db" {
  identifier           = "guss-prd-db"
  allocated_storage    = 20
  engine               = "mysql"
  engine_version       = "8.0"
  instance_class       = "db.t3.micro"
  username             = "admin"
  password             = "password123!" 
  db_subnet_group_name = aws_db_subnet_group.prd_db_group.name
  vpc_security_group_ids = [aws_security_group.prd_backend_sg.id]
  skip_final_snapshot  = true

  tags = {
    Name = "GUSS-PRD-RDS-2A"
  }
}

# DynamoDB 테이블 (PRD) - 세미콜론 제거 완료
resource "aws_dynamodb_table" "prd_ddb" {
  name         = "GUSS-PRD-DDB"
  billing_mode = "PAY_PER_REQUEST"
  hash_key     = "id"

  attribute {
    name = "id"
    type = "S"
  }

  tags = {
    Name = "GUSS-PRD-DDB"
  }
}

# SQS & SNS & Secrets Manager (PRD)
resource "aws_sqs_queue" "prd_sqs_rsv" {
  name = "GUSS-PRD-SQS-RSV"
}

resource "aws_sns_topic" "prd_sns_alm" {
  name = "GUSS-PRD-SNS-ALM"
}

resource "aws_secretsmanager_secret" "prd_rds_secret" {
  name = "GUSS-PRD-ASM-RDS"
}


# ==========================================
# 2. 개발 환경 (DEV) - 데이터 및 메시징
# ==========================================

# DB 서브넷 그룹 (DEV)
resource "aws_db_subnet_group" "dev_db_group" {
  name       = "guss-dev-db-group"
  subnet_ids = [aws_subnet.dev_backend_pri_2a.id]

  tags = {
    Name = "GUSS-DEV-DB-GROUP"
  }
}

# RDS MySQL 인스턴스 (DEV)
resource "aws_db_instance" "dev_db" {
  identifier           = "guss-dev-db"
  allocated_storage    = 20
  engine               = "mysql"
  instance_class       = "db.t3.micro"
  username             = "devadmin"
  password             = "devpassword123!"
  db_subnet_group_name = aws_db_subnet_group.dev_db_group.name
  vpc_security_group_ids = [aws_security_group.prd_backend_sg.id]
  skip_final_snapshot  = true

  tags = {
    Name = "GUSS-DEV-RDS-2A"
  }
}

# DynamoDB 테이블 (DEV) - 세미콜론 제거 완료 (오류 지점)
resource "aws_dynamodb_table" "dev_ddb" {
  name         = "GUSS-DEV-DDB"
  billing_mode = "PAY_PER_REQUEST"
  hash_key     = "id"

  attribute {
    name = "id"
    type = "S"
  }

  tags = {
    Name = "GUSS-DEV-DDB"
  }
}

# SQS & SNS (DEV)
resource "aws_sqs_queue" "dev_sqs_rsv" {
  name = "GUSS-DEV-SQS-RSV"
}

resource "aws_sns_topic" "dev_sns_alm" {
  name = "GUSS-DEV-SNS-ALM"
}