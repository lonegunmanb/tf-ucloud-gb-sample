provider "ucloud" {
  private_key = "${var.private_key}"
  public_key = "${var.public_key}"
  project_id = "${var.project_id}"
  region = "${var.region}"
}
//
resource "null_resource" "server_tag_check" {
  count = "${var.blue_tag == var.green_tag ? 1 : 0}"
  "Error: blue server tag AND green server tag MUST BE UNIQUE" = true
}

module "network" {
  source = "./network"
  green_subnet_cidr = "${var.green_subnet_cidr}"
  blue_subnet_cidr = "${var.blue_subnet_cidr}"
  vpc_cidr = "${var.vpc_cidr}"
  vpc_name = "${var.vpc_name}"
}

data "ucloud_instances" "current_blue" {
  availability_zone = "${var.blue_az}"
  tag = "${var.blue_tag}"
}

data "ucloud_instances" "current_green" {
  availability_zone = "${var.green_az}"
  tag = "${var.green_tag}"
}

locals {
  current_blue_count = "${length(data.ucloud_instances.current_blue.instances)}"
  current_green_count = "${length(data.ucloud_instances.current_green.instances)}"
}

module "bg_module" {
  source = "./blue-green"
  operation = "b2s"
  desired_blue_count = 1
  desired_green_count = 1
  current_blue_count = "${local.current_blue_count}"
  current_green_count = "${local.current_green_count}"
}

module "blue" {
  source = "./instance"
  vpc_id = "${module.network.vpc_id}"
  instance_type = "${var.blue_instance_type}"
  name_prefix = "${var.server_name_prefix}"
  subnet_id = "${module.network.blue_subnet_id}"
  image_id = "${var.blue_image_id}"
  tag = "${var.blue_tag}"
  count = "${module.bg_module.desired_blue_count}"
  root_password = "${var.root_password}"
  security_group_id = "${module.network.sg_id}"
  az = "${var.blue_az}"
}

module "green" {
  source = "./instance"
  vpc_id = "${module.network.vpc_id}"
  instance_type = "${var.green_instance_type}"
  name_prefix = "${var.server_name_prefix}"
  subnet_id = "${module.network.green_subnet_id}"
  image_id = "${var.green_image_id}"
  tag = "${var.green_tag}"
  count = "${module.bg_module.desired_green_count}"
  root_password = "${var.root_password}"
  security_group_id = "${module.network.sg_id}"
  az = "${var.green_az}"
}
//
//output "lb_target" {
//  value = "${module.bg_module.lb_target}"
//}
//
//output "desired_blue_count" {
//  value = "${module.bg_module.desired_blue_count}"
//}
//
//output "desired_green_count" {
//  value = "${module.bg_module.desired_green_count}"
//}