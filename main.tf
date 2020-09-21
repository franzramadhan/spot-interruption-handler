terraform {
  required_version = ">= 0.12"
}

resource "aws_iam_role" "this" {
  name        = var.service_name
  path        = "/"
  description = "IAM Role of ${var.service_name}"

  assume_role_policy    = data.aws_iam_policy_document.lambda.json
  force_detach_policies = true
  max_session_duration  = 3600

  tags = {
    Name          = var.service_name
    Environment   = var.environment
    Description   = "IAM Role of ${var.service_name}"
    ManagedBy     = "terraform"
  }

}

resource "aws_iam_role_policy" "lambda_policy" {
  name   = var.service_name
  role   = aws_iam_role.this.name
  policy = data.aws_iam_policy_document.lambda_policy.json
}

resource "aws_cloudwatch_log_group" "this" {
  name              = "/aws/lambda/${var.service_name}"
  retention_in_days = var.log_retention_in_days

  tags = {
    Service       = var.service_name
    Description   = "Cloudwatch Log for ${var.service_name}"
    Environment   = var.environment
    ManagedBy     = "terraform"
  }
}

resource "aws_lambda_function" "this" {
  function_name    = var.service_name
  filename         = "${path.module}/build/main.zip"
  handler          = "main"
  source_code_hash = data.archive_file.lambda.output_base64sha256
  role             = aws_iam_role.this.arn
  runtime          = "go1.x"
  memory_size      = 128
  timeout          = var.lambda_timeout
  description      = var.description

  tags = {
    Service     = var.service_name
    Description = "Lambda Function for ${var.service_name}"
    Environment = var.environment
    ManagedBy   = "terraform"
  }
}

resource "aws_cloudwatch_event_rule" "spot" {
  name          = var.service_name
  event_pattern = <<PATTERN
  {
    "detail-type": [
      "EC2 Spot Instance Interruption Warning"
    ],
    "source": [
      "aws.ec2"
    ]
  }
  PATTERN
}

resource "aws_cloudwatch_event_target" "lambda" {
  rule = aws_cloudwatch_event_rule.spot.name
  arn  = aws_lambda_function.this.arn
}

resource "aws_lambda_permission" "this" {
  statement_id  = "AllowCloudwatchEventRuleInvoke"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.this.function_name
  principal     = "events.amazonaws.com"
  source_arn    = aws_cloudwatch_event_rule.spot.arn
}

resource "null_resource" "lambda" {
  triggers = {
    filebase64 = filebase64("${path.module}/lambda/main.go")
  }

  provisioner "local-exec" {
    command = "mkdir -p ${path.module}/build && GOARCH=amd64 GOOS=linux go build -ldflags='-w -s' -o ${path.module}/build/main ${path.module}/lambda/main.go"
  }
}
