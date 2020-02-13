resource "aws_eip" "test-eip" {
  count = var.instance_count
}