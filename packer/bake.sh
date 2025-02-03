# build.sh
#!/bin/bash
set -e

# Build backend AMI
packer init backend-vm.pkr.hcl
packer build backend-vm.pkr.hcl

# Build frontend AMI
packer init frontend-vm.pkr.hcl
packer build frontend-vm.pkr.hcl