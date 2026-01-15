# 리전 설정
variable "region" {
  description = "AWS 리전"
  default     = "ap-northeast-2"
}

# 운영 환경 Nginx AMI (나중에 생성 후 실제 ID로 교체)
variable "prd_ami_nginx" {
  description = "GUSS-PRD-NGINX AMI ID"
  default     = "ami-xxxxxxxxxxxx" # 임시 값
}

# 운영 환경 Backend AMI
variable "prd_ami_backend" {
  description = "GUSS-PRD-BACKEND AMI ID"
  default     = "ami-yyyyyyyyyyyy" # 임시 값
}

# 개발 환경 Nginx AMI
variable "dev_ami_nginx" {
  description = "GUSS-DEV-NGINX AMI ID"
  default     = "ami-zzzzzzzzzzzz" # 임시 값
}