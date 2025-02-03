# variables.pkr.hcl
variable "environment" {
  type    = string
  default = "smc-dev"
}

variable "aws_region" {
  type    = string
  default = "us-west-2"
}