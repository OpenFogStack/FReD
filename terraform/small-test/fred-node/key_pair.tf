resource "aws_key_pair" "keypair" {
  key_name   = var.key_name
  public_key = var.key_pub
}