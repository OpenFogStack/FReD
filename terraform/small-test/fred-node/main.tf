resource "aws_instance" "fred_instance" {
  ami             = data.aws_ami.amazonlinux2.id
  instance_type   = var.instance_type
  key_name        = aws_key_pair.keypair.key_name
  count           = var.instance_count

  security_groups = var.security_groups

  provisioner "file" {
    source      = "./fred-node/config.toml"
    destination = "/tmp/config.toml"
  }


  provisioner "file" {
    source      = "./fred-node/setup_node.sh"
    destination = "/tmp/script.sh"
  }

  provisioner "remote-exec" {
    inline = [
      "chmod +x /tmp/script.sh",
      "/tmp/script.sh ${var.gitlab_repo_username} ${var.gitlab_repo_password} ${var.identifier}",
    ]
  }

  connection {
    type          = "ssh"
    user          = "ec2-user"
    private_key   = var.key_prv
    host          = self.public_ip
  }

  tags = {
    Name = "${var.name}-${count.index}"
    type = "fred"
  }
}