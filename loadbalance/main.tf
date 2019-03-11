resource "ucloud_subnet" "lb_subnet" {
  cidr_block = "${var.lb_cidr}"
  vpc_id = "${var.vpc_id}"
}
resource "ucloud_lb" "load_balancer" {
  name = "${var.lb_name}"
  vpc_id = "${var.vpc_id}"
  subnet_id = "${ucloud_subnet.lb_subnet.id}"
  tag = "${var.lb_direction}"
  remark = "${var.lb_direction}"
}
resource "ucloud_lb_listener" "listener" {
  load_balancer_id = "${ucloud_lb.load_balancer.id}"
  protocol = "tcp"
  port = "${var.listener_port}"
}

resource "ucloud_lb_attachment" "backend" {
  count = "${var.backend_count}"
  listener_id = "${ucloud_lb_listener.listener.id}"
  load_balancer_id = "${ucloud_lb.load_balancer.id}"
  resource_id = "${element(var.server_ids, count.index)}"
  port = "${var.backend_port}"
}

resource "ucloud_eip" "public_ip" {
  internet_type = "bgp"
  bandwidth = "200"
  charge_mode = "traffic"
  charge_type = "dynamic"
  name = "test-eip"
}

resource "ucloud_eip_association" "eip_association" {
  eip_id = "${ucloud_eip.public_ip.id}"
  resource_id = "${ucloud_lb.load_balancer.id}"
}