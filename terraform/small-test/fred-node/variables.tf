variable "instance_type" {
  type = string
  default = "t2.micro"
}

variable "name" {
  type = string
}

variable "key_name" {
  type = string
  default = "keypair"
}

variable "key_pub" {
  type = string
}

variable "key_prv" {
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

variable "instance_count" {
  type = number
  default = 3
}

variable "identifier" {
  type = string
}

variable "fred_flags" {
  type = list(string)
}