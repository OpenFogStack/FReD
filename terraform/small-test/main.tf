module "fred-node" {
  source = ".//fred-node"

  name            = "${var.identifier}-fred-node"
  key_name        = "${var.identifier}-terraform_key"
  key_pub         = file("terraform.key.pub")
  key_prv         = file("terraform.key")
  instance_count = var.instance_count
  security_groups = [
    aws_security_group.allow_ssh.name,
    aws_security_group.allow_outbound.name,
    aws_security_group.allow_fred_web.name,
    aws_security_group.allow_fred_zmq.name
  ]
  gitlab_repo_password = var.gitlab_repo_password
  gitlab_repo_username = var.gitlab_repo_username
  identifier = var.identifier
  instance_type = var.instance_type
  fred_flags = var.fred_flags
  domain_name = "nodes.${var.identifier}.${var.root_domain}"
  domain_hosted_zone = var.root_domain
}