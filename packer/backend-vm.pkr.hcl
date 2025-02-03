# backend.pkr.hcl
source "amazon-ebs" "backend" {
  ami_name      = "mam-backend-${var.environment}"
  instance_type = "t3.small"
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
    Component   = "MAM Backend"
  }
}

build {
  sources = ["source.amazon-ebs.backend"]

  provisioner "file" {
    source      = "mam-backend"
    destination = "/tmp/mam-backend"
  }

  provisioner "shell" {
    script = "scripts/backend-init.sh"
  }
}
