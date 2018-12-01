resource "aws_iam_role" "iot_button_metrics" {
  name = "iot_button_metrics"

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "lambda.amazonaws.com"
      },
      "Effect": "Allow",
      "Sid": ""
    }
  ]
}
EOF
}

resource "aws_iam_role_policy" "iot_button_metrics" {
  name = "iot_button_metrics"
  role = "${aws_iam_role.iot_button_metrics.id}"

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
       "Effect": "Allow",
       "Action": ["s3:ListBucket"],
       "Resource": ["${aws_s3_bucket.metrics.arn}"]
     },
     {
       "Effect": "Allow",
       "Action": [
         "s3:PutObject",
         "s3:GetObject",
         "s3:DeleteObject"
       ],
       "Resource": ["${aws_s3_bucket.metrics.arn}/*"]
     },
     {
       "Effect": "Allow",
       "Action": [
         "logs:*"        
       ],
       "Resource": "arn:aws:logs:*:*:*"
     }
  ]
}
EOF
}

resource "aws_lambda_function" "send_webhook" {
  filename         = "bin/send_webhook.zip"
  function_name    = "send_webhook"
  role             = "${aws_iam_role.iot_button_metrics.arn}"
  handler          = "bin/send_webhook"
  source_code_hash = "${base64sha256(file("bin/send_webhook.zip"))}"
  runtime          = "go1.x"

  environment {
    variables = {
      WEBHOOK_URL = "${var.webhook_url}"
      MESSAGE = "${var.message}"
    }
  }

  tags {
    Name        = "${var.tags}"
    Environment = "Production"
  }
}

