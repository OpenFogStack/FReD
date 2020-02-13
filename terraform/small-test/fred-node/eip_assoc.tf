resource "aws_eip_association" "test-eip-assoc" {
  count = var.instance_count
  allocation_id = aws_eip.test-eip[count.index].id
  instance_id   = aws_instance.fred_instance[count.index].id
}