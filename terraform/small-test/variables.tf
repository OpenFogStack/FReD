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

variable "aws_access_key" {
  description = "The AWS access key."
}

variable "aws_secret_key" {
  description = "The AWS secret key."
}

variable "gitlab_repo_username" {
  description = "The username for the GitLab registry."
}

variable "gitlab_repo_password" {
  description = "The password for the GitLab registry."
}

