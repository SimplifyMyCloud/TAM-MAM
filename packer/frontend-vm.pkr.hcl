
# frontend.pkr.hcl
source "amazon-ebs" "frontend" {
  ami_name      = "mam-frontend-${var.environment}"
  instance_type = "t3.micro"
  region        = var.aws_region

  source_ami_filter {
    filters = {
      name                = "ubuntu/images/*ubuntu-jammy-22.04-amd64-server-*"
      root-device-type    = "ebs"
      virtualization-type = "hvm"
    }
    most_recent = true
    owners      = ["099720109477"]
  }

  ssh_username = "ubuntu"

  tags = {
    Environment = var.environment
    Component   = "MAM Frontend"
  }
}

build {
  sources = ["source.amazon-ebs.frontend"]

  provisioner "file" {
    source      = "frontend/build"
    destination = "/tmp/build"
  }

  provisioner "shell" {
    script = "scripts/frontend-init.sh"
  }
}

# variables.pkr.hcl
variable "environment" {
  type    = string
  default = "dev"
}

variable "aws_region" {
  type    = string
  default = "us-west-2"
}