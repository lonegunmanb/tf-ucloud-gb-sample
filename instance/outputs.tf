output "ids" {
  value = "${ucloud_instance.server.*.id}"
}