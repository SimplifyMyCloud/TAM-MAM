# backend.pkr.hcl
source "amazon-ebs" "backend" {
  ami_name      = "mam-backend-smc-dev-{{timestamp}}"
  instance_type = "t3.small"
  region        = "us-west-2"

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
    Environment = "smc-dev"
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
    script          = "scripts/backend-init.sh"
    execute_command = "sudo sh -c '{{ .Vars }} {{ .Path }}'"
  }
}
