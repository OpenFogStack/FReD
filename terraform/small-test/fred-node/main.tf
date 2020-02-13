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
      # Args:
      # $1: username to access gitlab registry
      # $2: password to access gitlab registry
      # $3: identifier to use the correct Docker container
      # $4: "host" argument for FReD, in this case the domain name returned from R53
      # $5: flags that should be passed to FReD
      "/tmp/script.sh ${var.gitlab_repo_username} ${var.gitlab_repo_password} ${var.identifier} ${aws_route53_record.dns_record[count.index].fqdn} ${var.fred_flags[count.index]}",
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