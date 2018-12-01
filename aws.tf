variable "aws_access_key" {}
variable "aws_secret_key" {}

variable "region" {
  default = "ap-northeast-1"
}

variable "tags" {
  default = "cheers-button"
}

variable "webhook_url" {}
variable "message" {}

provider "aws" {
  access_key = "${var.aws_access_key}"
  secret_key = "${var.aws_secret_key}"
  region     = "${var.region}"
}
