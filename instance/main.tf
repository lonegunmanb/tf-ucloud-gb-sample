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

resource "ucloud_subnet" "lb_subnet" {
  cidr_block = "${var.lb_cidr}"
  vpc_id = "${var.vpc_id}"
}
resource "ucloud_lb" "load_balancer" {
  name = "${var.lb_name}"
  vpc_id = "${var.vpc_id}"
  subnet_id = "${ucloud_subnet.lb_subnet.id}"
  tag = "${var.tag}"
}
resource "ucloud_lb_listener" "listener" {
  load_balancer_id = "${ucloud_lb.load_balancer.id}"
  protocol = "tcp"
  port = "${var.listener_port}"
}

resource "ucloud_lb_attachment" "backend" {
  count = "${var.count}"
  listener_id = "${ucloud_lb_listener.listener.id}"
  load_balancer_id = "${ucloud_lb.load_balancer.id}"
  //  resource_id = "${element(var.server_ids, count.index)}"
  resource_id = "${ucloud_instance.server.*.id[count.index]}"
  port = "${var.backend_port}"
}
