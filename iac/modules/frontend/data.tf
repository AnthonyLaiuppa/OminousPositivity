data "aws_partition" "current" {}
data "aws_region" "current" {}
data "aws_caller_identity" "current" {}

data "aws_iam_policy_document" "lambda_assume_role" {}

data "aws_route53_zone" "main" {
  name         = var.root_dns_name
  private_zone = false
}