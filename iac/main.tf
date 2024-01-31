module "backend" {
  source                      = "./modules/backend"
  lambda_bucket_name          = var.lambda_bucket_name
  message_lambda_package_name = var.message_lambda_package_name
  root_dns_name               = var.root_dns_name
}

module "frontend" {
  source        = "./modules/frontend"
  root_dns_name = var.root_dns_name
}