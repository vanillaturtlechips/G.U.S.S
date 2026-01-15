# VPC
resource "aws_vpc" "prd_vpc" {
  cidr_block           = "10.1.0.0/16"
  enable_dns_hostnames = true
  tags                 = { Name = "GUSS-PRD-VPC" }
}

# Gateways
resource "aws_internet_gateway" "prd_igw" {
  vpc_id = aws_vpc.prd_vpc.id
  tags   = { Name = "PRD-IGW" }
}

resource "aws_eip" "nat_eip" {
  domain = "vpc"
  
  tags = {
    Name = "PRD-NAT-EIP"
  }
}

# NAT 게이트웨이 설정 (기존과 동일)
resource "aws_nat_gateway" "prd_nat" {
  allocation_id = aws_eip.nat_eip.id
  subnet_id     = aws_subnet.prd_nginx_pub_2a.id
  tags          = { Name = "PRD-NAT" }
}

# Subnets (관리대장 수치 100% 반영)
resource "aws_subnet" "prd_nginx_pub_2a" {
  vpc_id            = aws_vpc.prd_vpc.id
  cidr_block        = "10.1.1.0/24"
  availability_zone = "ap-northeast-2a"
  tags              = { Name = "GUSS-PRD-NGINX-PUB-2A" }
}

resource "aws_subnet" "prd_nginx_pub_2c" {
  vpc_id            = aws_vpc.prd_vpc.id
  cidr_block        = "10.1.5.0/24"
  availability_zone = "ap-northeast-2c"
  tags              = { Name = "GUSS-PRD-NGINX-PUB-2C" }
}

resource "aws_subnet" "prd_backend_pri_2a" {
  vpc_id            = aws_vpc.prd_vpc.id
  cidr_block        = "10.1.3.0/24"
  availability_zone = "ap-northeast-2a"
  tags              = { Name = "GUSS-PRD-BACKEND-PRI-2A" }
}

resource "aws_subnet" "prd_backend_pri_2c" {
  vpc_id            = aws_vpc.prd_vpc.id
  cidr_block        = "10.1.6.0/24"
  availability_zone = "ap-northeast-2c"
  tags              = { Name = "GUSS-PRD-BACKEND-PRI-2C" }
}

# Route Tables (관리대장 46, 47번)
resource "aws_route_table" "prd_pub_rt" {
  vpc_id = aws_vpc.prd_vpc.id
  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.prd_igw.id
  }
  tags = { Name = "GUSS-PRD-PUB" }
}

resource "aws_route_table" "prd_pri_rt" {
  vpc_id = aws_vpc.prd_vpc.id
  route {
    cidr_block     = "0.0.0.0/0"
    nat_gateway_id = aws_nat_gateway.prd_nat.id
  }
  tags = { Name = "GUSS-PRD-PRI" }
}

# Associations
resource "aws_route_table_association" "pub_2a" {
  subnet_id      = aws_subnet.prd_nginx_pub_2a.id
  route_table_id = aws_route_table.prd_pub_rt.id
}

resource "aws_route_table_association" "pri_2a" {
  subnet_id      = aws_subnet.prd_backend_pri_2a.id
  route_table_id = aws_route_table.prd_pri_rt.id
}

# --- 개발 VPC (관리대장 2번) ---
resource "aws_vpc" "dev_vpc" {
  cidr_block           = "10.5.0.0/16"
  enable_dns_hostnames = true
  tags                 = { Name = "GUSS-DEV-VPC" }
}

resource "aws_internet_gateway" "dev_igw" {
  vpc_id = aws_vpc.dev_vpc.id
  tags   = { Name = "DEV-IGW" }
}

# 개발 서브넷 (관리대장 36~39번)
resource "aws_subnet" "dev_nginx_pub_2a" {
  vpc_id            = aws_vpc.dev_vpc.id
  cidr_block        = "10.5.1.0/24"
  availability_zone = "ap-northeast-2a"
  tags              = { Name = "GUSS-DEV-NGINX-PUB-2A" }
}

resource "aws_subnet" "dev_backend_pri_2a" {
  vpc_id            = aws_vpc.dev_vpc.id
  cidr_block        = "10.5.3.0/24"
  availability_zone = "ap-northeast-2a"
  tags              = { Name = "GUSS-DEV-BACKEND-PRI-2A" }
}

# 개발 라우팅 테이블
resource "aws_route_table" "dev_pub_rt" {
  vpc_id = aws_vpc.dev_vpc.id
  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.dev_igw.id
  }
  tags = { Name = "GUSS-DEV-PUB" }
}

resource "aws_route_table_association" "dev_pub_assoc" {
  subnet_id      = aws_subnet.dev_nginx_pub_2a.id
  route_table_id = aws_route_table.dev_pub_rt.id
}