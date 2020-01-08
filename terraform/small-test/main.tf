resource "aws_key_pair" "terraform_key" {
  key_name   = "terraform_key"
  public_key = file("terraform.key.pub")
}

module "fred-node-0" {
  source = ".//fred-node"

  name            = "${var.identifier}-fred-node-0"
  key_pair        = aws_key_pair.terraform_key.key_name
  key_pair_key    = "terraform.key"
  security_groups = [
    aws_security_group.allow_ssh.name,
    aws_security_group.allow_outbound.name,
    aws_security_group.allow_fred_web.name,
    aws_security_group.allow_fred_zmq.name
  ]
  gitlab_repo_password = var.gitlab_repo_password
  gitlab_repo_username = var.gitlab_repo_username
}