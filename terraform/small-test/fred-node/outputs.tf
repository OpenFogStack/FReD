output "server-ip" {
  value = aws_eip.test-eip.public_ip
}