output "server-ip" {
  value =  [ aws_eip.test-eip.*.public_ip ]
}

output "server-domain-name" {
  value =  [ aws_route53_record.dns_record.*.fqdn ]
}
