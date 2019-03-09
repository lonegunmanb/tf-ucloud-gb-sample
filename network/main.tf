resource "ucloud_vpc" "example" {
  name = "${var.vpc_name}"
  tag = "${var.vpc_tag}"
  cidr_blocks = ["${var.vpc_cidr}"]
  remark = "${var.vpc_remark}"
}

resource "ucloud_subnet" "green" {
  cidr_block = "${var.green_subnet_cidr}"
  vpc_id = "${ucloud_vpc.example.id}"
  name = "${var.green_subnet_name}"
  tag = "${var.green_subnet_tag}"
  remark = "${var.green_subnet_remark}"
}

resource "ucloud_subnet" "blue" {
  cidr_block = "${var.blue_subnet_cidr}"
  vpc_id = "${ucloud_vpc.example.id}"
  name = "${var.blue_subnet_name}"
  tag = "${var.blue_subnet_tag}"
  remark = "${var.blue_subnet_remark}"
}

resource "ucloud_security_group" "security_group" {
  name = "${var.security_group_name}"
  "rules" {
    port_range = "22"
    protocol = "tcp"
    cidr_block = "0.0.0.0/0"
    policy = "accept"
  }
  "rules" {
    port_range = "8080"
    protocol = "tcp"
    cidr_block = "0.0.0.0/0"
    policy = "accept"
  }
}