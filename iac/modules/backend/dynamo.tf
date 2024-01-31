resource "aws_dynamodb_table" "backend" {
  name         = "ominous-positivty-${terraform.workspace}"
  billing_mode = "PAY_PER_REQUEST"
  hash_key     = "id"
  server_side_encryption {
    enabled     = true
    kms_key_arn = aws_kms_key.backend.arn
  }
  point_in_time_recovery {
    enabled = true
  }
  attribute {
    name = "id"
    type = "N"
  }
}

resource "aws_kms_key" "backend" {
  description             = "KMS Key used for Ominous Positivity Backend Dynamo Encryption"
  deletion_window_in_days = 7
  key_usage               = "ENCRYPT_DECRYPT"
  is_enabled              = true
  enable_key_rotation     = true
  policy                  = data.aws_iam_policy_document.backend_dynamo_kms_key_policy.json
}

resource "aws_kms_alias" "backend" {
  name          = "alias/OPBackend"
  target_key_id = aws_kms_key.backend.key_id
}

data "aws_iam_policy_document" "backend_dynamo_kms_key_policy" {
  statement {
    sid    = "Allow DynamoDB KMS Key Usage"
    effect = "Allow"
    actions = [
      "kms:Describe*",
      "kms:Get*",
      "kms:List*"
    ]
    resources = ["*"]
    principals {
      identifiers = [
        "dynamodb.amazonaws.com"
      ]
      type = "Service"
    }
  }
  statement {
    sid    = "Enable Account IAM access"
    effect = "Allow"
    principals {
      identifiers = ["arn:aws:iam::${data.aws_caller_identity.current.account_id}:root"]
      type        = "AWS"
    }
    resources = ["*"]
    actions   = ["kms:*"]
  }
}