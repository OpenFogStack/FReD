resource "aws_key_pair" "terraform_key" {
  key_name   = "terraform_key"
  public_key = file("terraform.key.pub")
}

module "fred-nodes" {
  source = ".//fred-node"

  name            = "${var.identifier}-fred-node"
  key_pair        = aws_key_pair.terraform_key.key_name
  key_pair_key    = "terraform.key"
  instance_count = var.instance_count
  security_groups = [
    aws_security_group.allow_ssh.name,
    aws_security_group.allow_outbound.name,
    aws_security_group.allow_fred_web.name,
    aws_security_group.allow_fred_zmq.name
  ]
  gitlab_repo_password = var.gitlab_repo_password
  gitlab_repo_username = var.gitlab_repo_username
}