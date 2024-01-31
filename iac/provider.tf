provider "aws" {
  region = "us-east-1"

  default_tags {
    tags = local.default_tags
  }
}

locals {
  default_tags = {
    "Environment" = terraform.workspace
  }
}