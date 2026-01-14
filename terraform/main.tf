# provider.tf
provider "aws" {
  region = "ap-northeast-2"
}

locals {
  common_tags = {
    Project = "GUSS"
    Team    = "2" # 2조 식별 태그
    ManagedBy = "Terraform"
  }
}

# variables.tf (AMI ID 등은 실제 생성된 ID로 변경 필요)
variable "prd_ami_nginx" { default = "ami-xxxxxxxxxxxx" } 
variable "prd_ami_backend" { default = "ami-yyyyyyyyyyyy" }
variable "dev_ami_nginx" { default = "ami-zzzzzzzzzzzz" }