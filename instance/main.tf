resource "ucloud_instance" "server" {
  count = "${var.count}"
  availability_zone = "${var.az}"
  image_id = "${var.image_id}"
  instance_type = "${var.instance_type}"
  root_password = "${var.root_password}"
  name = "${var.name_prefix}-${count.index}"
  tag = "${var.tag}"
  charge_type = "${var.charge_type}"
  duration = "${var.duration}"
  vpc_id = "${var.vpc_id}"
  subnet_id = "${var.subnet_id}"
  security_group = "${var.security_group_id}"
}
