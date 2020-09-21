data "aws_caller_identity" "current" {}
data "aws_region" "current" {}

data "archive_file" "lambda" {
  type        = "zip"
  source_file = "${path.module}/build/main"
  output_path = "${path.module}/build/main.zip"
  depends_on = [
    null_resource.lambda
  ]
}

data "aws_iam_policy_document" "lambda" {
  statement {
    actions = ["sts:AssumeRole"]
    effect  = "Allow"

    principals {
      type        = "Service"
      identifiers = ["lambda.amazonaws.com"]
    }
  }
}

data "aws_iam_policy_document" "lambda_policy" {

  statement {
    effect = "Allow"

    actions = [
      "logs:CreateLogGroup",
      "logs:CreateLogStream",
      "logs:PutLogEvents",
    ]

    resources = [
      "${aws_cloudwatch_log_group.this.arn}",
    ]
  }

  statement {
    effect = "Allow"

    actions = [
      "autoscaling:DescribeAutoScalingInstances",
      "autoscaling:DetachInstances",
      "autoscaling:SetDesiredCapacity",
    ]

    resources = [
      "*",
    ]
  }
}
