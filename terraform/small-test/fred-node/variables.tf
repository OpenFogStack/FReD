variable "instance_type" {
  type = string
  default = "t2.micro"
}

variable "name" {
  type = string
}

variable "key_pair" {
  type = string
  default = "my_test_key"
}

variable "key_pair_key" {
  type = string
}

variable "security_groups" {
  type = list(string)
  default = []
}

variable "gitlab_repo_username" {
  type = string
}

variable "gitlab_repo_password" {
  type = string
}