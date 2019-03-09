output "vpc_id" {
  value = "${ucloud_vpc.example.id}"
}

output "green_subnet_id" {
  value = "${ucloud_subnet.green.id}"
}

output "blue_subnet_id" {
  value = "${ucloud_subnet.blue.id}"
}

output "sg_id" {
  value = "${ucloud_security_group.security_group.id}"
}