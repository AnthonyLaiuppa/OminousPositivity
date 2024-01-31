terraform {
  required_version = ">= 1.5.5"
  backend "s3" {
    key            = ""
    region         = ""
    bucket         = ""
    dynamodb_table = ""
    encrypt        = true
    kms_key_id     = ""
  }
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = ">= 5.0.0"
    }
  }
}
