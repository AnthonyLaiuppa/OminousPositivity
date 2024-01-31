variable "lambda_bucket_name" {
  description = "S3 Bucket Containing Backend Lambda function packages"
  type        = string
}

variable "message_lambda_package_name" {
  description = "S3 Key for the Backend Function Message Lambda package"
  type        = string
}

variable "root_dns_name" {
  description = "Root DNS Name of the Zone to create records in"
  type        = string
}