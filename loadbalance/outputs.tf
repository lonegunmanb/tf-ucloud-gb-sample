output "elb_eip" {
  value = "${ucloud_eip.public_ip.public_ip}"
}

output "server_ids" {
  value = "${var.server_ids}"
}