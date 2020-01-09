resource "aws_eip" "test-eip" {
  count = var.instance_count
  instance    = aws_instance.fred_instance[count.index].id
}