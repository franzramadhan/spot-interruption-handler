variable "environment" {
  type        = string
  description = "The environment of the application. Valid value may varies from development, testing, staging, production, or management."
}

variable "description" {
  type        = string
  description = "Brief descriptive name of Lambda."
}

variable "service_name" {
  type        = string
  description = "Name of the service."
}

variable "role_max_session_duration" {
  type        = number
  description = "The maximum session duration (in seconds) that you want to set for the specified role. If you do not specify a value for this setting, the default maximum of one hour is applied. This setting can have a value from 1 hour to 12 hours."
  default     = 3600
}

variable "log_retention_in_days" {
  type        = number
  description = "How long in days lambda function log will be kept in cloudwatch."
  default     = 7
}

variable "lambda_timeout" {
  type        = number
  description = "The amount of time your Lambda Function has to run in seconds."
  default     = 60
}
