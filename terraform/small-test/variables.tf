variable "aws_region" {
  description = "AWS region to deploy infrastructure in."
  type = string
  default = "eu-central-1"
}

variable "identifier" {
  description = "Identifier for the deployment."
  type = string
}

variable "instance_count" {
  description = "Number of nodes to deploy."
  type = number
  default = 3
}

variable "instance_type" {
  description = "Type of EC2 instance to use."
  type = string
  default = "t2.micro"
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

variable "fred_flags" {
  description = "Flags to pass to FReD container."
  default = ["","",""]
  type = list(string)
}
