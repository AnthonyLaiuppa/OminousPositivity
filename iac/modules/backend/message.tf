resource "aws_lambda_function" "message" {
  function_name = "MessageFunction-${terraform.workspace}"
  handler       = "bootstrap"
  runtime       = "provided.al2"
  memory_size   = 128
  timeout       = 3
  role          = aws_iam_role.message.arn
  s3_bucket     = var.lambda_bucket_name
  s3_key        = var.message_lambda_package_name

  environment {
    variables = {
      TABLE_NAME   = aws_dynamodb_table.backend.name
      ALLOW_ORIGIN = "https://${var.root_dns_name}"
    }
  }
  architectures = ["x86_64"]

  tracing_config {
    mode = "Active"
  }
}

resource "aws_cloudwatch_log_group" "lambda_message" {
  name              = "/aws/lambda/${aws_lambda_function.message.function_name}"
  retention_in_days = 7
}

resource "aws_iam_role" "message" {
  name               = "message-function-${terraform.workspace}"
  assume_role_policy = data.aws_iam_policy_document.lambda_assume_role.json
}

resource "aws_iam_role_policy_attachment" "message" {
  policy_arn = aws_iam_policy.message.arn
  role       = aws_iam_role.message.name
}

resource "aws_iam_policy" "message" {
  name   = "message-function-${terraform.workspace}"
  policy = data.aws_iam_policy_document.message.json
}

data "aws_iam_policy_document" "message" {
  statement {
    sid = "AllowMessageDynamo"
    actions = [
      "dynamodb:GetItem"
    ]
    effect = "Allow"
    resources = [
      "${aws_dynamodb_table.backend.arn}/*",
      aws_dynamodb_table.backend.arn
    ]
  }
  statement {
    sid = "AllowMessageKMS"
    actions = [
      "kms:Decrypt"
    ]
    effect = "Allow"
    resources = [
      aws_kms_key.backend.arn
    ]
  }
  statement {
    sid = "EnableXrayTracing"
    actions = [
      "xray:PutTraceSegments",
      "xray:PutTelemetryRecords"
    ]
    effect    = "Allow"
    resources = ["*"]
  }
  statement {
    sid = "EnableCloudWatchLogging"
    actions = [
      "logs:CreateLogStream",
      "logs:PutLogEvents"
    ]
    effect = "Allow"
    resources = [
      aws_cloudwatch_log_group.lambda_message.arn,
      "${aws_cloudwatch_log_group.lambda_message.arn}:*"
    ]
  }
}