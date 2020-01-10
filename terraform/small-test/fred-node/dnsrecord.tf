resource "aws_route53_record" "dns_record" {
  # Use the ID of the Hosted Zone we retrieved earlier
  zone_id = data.aws_route53_zone.hosted_zone.zone_id

  # Set the name of the record, e.g. pc.mydomain.com
  name = "${count.index}.${var.domain_name}"

  count = var.instance_count

  # We're pointing to an IP address so we need to use an A record
  type = "A"

  # We'll set the TTL of the record to 30 minutes (1800 seconds)
  ttl = "1800"

  records = [ aws_eip.test-eip[count.index].public_ip ]
}