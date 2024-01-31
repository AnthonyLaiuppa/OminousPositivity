resource "aws_api_gateway_rest_api" "backend" {
  name        = "OminousPositivityAPI-${terraform.workspace}"
  description = "API Gateway for OminousPositivity Backend"
}

resource "aws_route53_record" "backend" {
  zone_id = data.aws_route53_zone.main.zone_id
  name    = "api.${var.root_dns_name}"
  type    = "A"

  alias {
    name                   = aws_api_gateway_domain_name.api_domain_name.cloudfront_domain_name
    zone_id                = aws_api_gateway_domain_name.api_domain_name.cloudfront_zone_id
    evaluate_target_health = false
  }
}

resource "aws_acm_certificate" "backend" {
  domain_name       = "api.${var.root_dns_name}"
  validation_method = "DNS"
  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_route53_record" "backend_validation" {
  for_each = {
    for dvo in aws_acm_certificate.backend.domain_validation_options : dvo.domain_name => {
      name   = dvo.resource_record_name
      record = dvo.resource_record_value
      type   = dvo.resource_record_type
    }
  }

  allow_overwrite = true
  name            = each.value.name
  records         = [each.value.record]
  ttl             = 60
  type            = each.value.type
  zone_id         = data.aws_route53_zone.main.zone_id
}

resource "aws_acm_certificate_validation" "backend" {
  certificate_arn         = aws_acm_certificate.backend.arn
  validation_record_fqdns = [for record in aws_route53_record.backend_validation : record.fqdn]
}

# Message Lambda Function
resource "aws_api_gateway_resource" "message" {
  rest_api_id = aws_api_gateway_rest_api.backend.id
  parent_id   = aws_api_gateway_rest_api.backend.root_resource_id
  path_part   = "message"
}

resource "aws_api_gateway_method" "message" {
  rest_api_id   = aws_api_gateway_rest_api.backend.id
  resource_id   = aws_api_gateway_resource.message.id
  http_method   = "GET"
  authorization = "NONE"
}

resource "aws_api_gateway_method" "message_cors" {
  rest_api_id   = aws_api_gateway_rest_api.backend.id
  resource_id   = aws_api_gateway_resource.message.id
  http_method   = "OPTIONS"
  authorization = "NONE"
}

resource "aws_api_gateway_integration" "message" {
  rest_api_id             = aws_api_gateway_rest_api.backend.id
  resource_id             = aws_api_gateway_resource.message.id
  http_method             = aws_api_gateway_method.message.http_method
  integration_http_method = "POST"
  type                    = "AWS_PROXY"
  uri                     = aws_lambda_function.message.invoke_arn
}

resource "aws_api_gateway_integration" "message_cors" {
  rest_api_id = aws_api_gateway_rest_api.backend.id
  resource_id = aws_api_gateway_resource.message.id
  http_method = aws_api_gateway_method.message_cors.http_method
  type        = "MOCK"
  request_templates = {
    "application/json" = "{\"statusCode\": 200}"
  }
}

resource "aws_lambda_permission" "api_gateway_message" {
  statement_id  = "AllowExecutionFromAPIGateway"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.message.function_name
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${aws_api_gateway_rest_api.backend.execution_arn}/*/GET/message"
}

resource "aws_api_gateway_deployment" "message" {
  depends_on = [
    aws_api_gateway_integration.message,
    aws_api_gateway_method.message,
    aws_api_gateway_integration_response.message_cors,
    aws_api_gateway_method_response.message_cors,
    aws_api_gateway_method_response.message_200
  ]

  rest_api_id = aws_api_gateway_rest_api.backend.id

  triggers = {
    redeployment = sha1(jsonencode({
      method_http_method      = aws_api_gateway_method.message.http_method,
      integration_http_method = aws_api_gateway_integration.message.integration_http_method,
      uri                     = aws_api_gateway_integration.message.uri
    }))
  }
  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_api_gateway_stage" "message" {
  stage_name    = terraform.workspace
  rest_api_id   = aws_api_gateway_rest_api.backend.id
  deployment_id = aws_api_gateway_deployment.message.id
}

resource "aws_api_gateway_domain_name" "api_domain_name" {
  depends_on = [
    aws_acm_certificate.backend,
    aws_acm_certificate_validation.backend
  ]
  domain_name     = "api.${var.root_dns_name}"
  certificate_arn = aws_acm_certificate.backend.arn
}

resource "aws_api_gateway_base_path_mapping" "base_path_mapping" {
  api_id      = aws_api_gateway_rest_api.backend.id
  stage_name  = aws_api_gateway_stage.message.stage_name
  domain_name = aws_api_gateway_domain_name.api_domain_name.domain_name
}

# Message CORS Configuration
resource "aws_api_gateway_method_response" "message_cors" {
  rest_api_id = aws_api_gateway_rest_api.backend.id
  resource_id = aws_api_gateway_resource.message.id
  http_method = aws_api_gateway_method.message_cors.http_method
  status_code = "200"

  response_models = {
    "application/json" = "Empty"
  }

  response_parameters = {
    "method.response.header.Access-Control-Allow-Headers" = true
    "method.response.header.Access-Control-Allow-Methods" = true
    "method.response.header.Access-Control-Allow-Origin"  = true
  }
}

resource "aws_api_gateway_integration_response" "message_cors" {
  rest_api_id = aws_api_gateway_rest_api.backend.id
  resource_id = aws_api_gateway_resource.message.id
  http_method = aws_api_gateway_integration.message_cors.http_method
  status_code = "200"

  response_parameters = {
    "method.response.header.Access-Control-Allow-Headers" = "'Content-Type,X-Amz-Date,Authorization,X-Api-Key,X-Amz-Security-Token'"
    "method.response.header.Access-Control-Allow-Methods" = "'GET,HEAD,OPTIONS'"
    "method.response.header.Access-Control-Allow-Origin"  = "'https://${var.root_dns_name}'"
  }

  response_templates = {
    "application/json" = ""
  }

  depends_on = [aws_api_gateway_integration.message_cors] // Ensure dependency is correctly defined.
}



resource "aws_api_gateway_method_response" "message_200" {
  rest_api_id = aws_api_gateway_rest_api.backend.id
  resource_id = aws_api_gateway_resource.message.id
  http_method = aws_api_gateway_method.message.http_method
  status_code = "200"

  response_models = {
    "application/json" = "Empty" //More for documentation purposes so Empty works fine as a placeholder
  }

  response_parameters = {
    "method.response.header.Access-Control-Allow-Origin" = true
  }
}

resource "aws_api_gateway_integration_response" "message_200" {
  rest_api_id = aws_api_gateway_rest_api.backend.id
  resource_id = aws_api_gateway_resource.message.id
  http_method = aws_api_gateway_method.message.http_method
  status_code = "200"

  response_parameters = {
    "method.response.header.Access-Control-Allow-Origin" = "'https://${var.root_dns_name}'"
  }

  depends_on = [
    aws_api_gateway_integration.message
  ]
}
