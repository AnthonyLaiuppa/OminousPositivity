variable "root_dns_name" {
  description = "TLD of the Route53 DNS zone"
  type        = string
}

variable "lambda_bucket_name" {
  description = "Name of the bucket holding all the Lambda functions"
  type        = string
}

variable "message_lambda_package_name" {
  default = "Name of the S3 key corresponding to the message lambda function package"
  type    = string
}