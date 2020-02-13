output "server-ip" {
  value =  [ aws_eip_association.test-eip-assoc.*.public_ip ]
}

output "server-domain-name" {
  value =  [ aws_route53_record.dns_record.*.fqdn ]
}
