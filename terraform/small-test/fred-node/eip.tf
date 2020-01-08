resource "aws_eip" "test-eip" {
  instance    = aws_instance.fred_instance.id
}