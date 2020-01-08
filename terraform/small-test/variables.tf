variable "aws_region" {
  type = string
  default = "eu-central-1"
}

variable "identifier" {
  type = string
}

variable "instance_count" {
  type = number
  default = 3
}