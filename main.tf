provider "ucloud" {
  private_key = "${var.private_key}"
  public_key = "${var.public_key}"
  project_id = "${var.project_id}"
  region = "${var.region}"
}

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
  operation = "${var.operation}"
  desired_blue_count = "${var.desired_blue_count}"
  desired_green_count = "${var.desired_green_count}"
  current_blue_count = "${local.current_blue_count}"
  current_green_count = "${local.current_green_count}"
}

locals {
  lb_direction = "${module.bg_module.lb_direction}"
}

locals {
  backend_count = "${local.lb_direction == "blue" ? module.bg_module.desired_blue_count : module.bg_module.desired_green_count}"
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
  backend_port = "${var.backend_port}"
  listener_port = "${var.listener_port}"
  lb_cidr = "${var.blue_lb_cidr}"
  lb_name = "${var.blue_tag}"
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
  lb_cidr = "${var.green_lb_cidr}"
  listener_port = "${var.listener_port}"
  backend_port = "${var.backend_port}"
  lb_name = "${var.green_tag}"
}

resource "ucloud_eip" "public_ip" {
  internet_type = "bgp"
  bandwidth = "200"
  charge_mode = "traffic"
  charge_type = "dynamic"
  name = "test-eip"
}

locals {
  target_lb_id = "${local.lb_direction == "blue" ? module.blue.lb_id : module.green.lb_id}"
}

resource "ucloud_eip_association" "eip_association" {
  eip_id = "${ucloud_eip.public_ip.id}"
  resource_id = "${local.target_lb_id}"
}

output "eip" {
  value = "${ucloud_eip.public_ip.public_ip}"
}

//locals {
//  server_ids = "${split(",", local.lb_direction == "blue" ? join(",", module.blue.ids) : join(",", module.green.ids))}"
//}
//
//module "load_balancer" {
//  source = "./loadbalance"
//  backend_port = "${var.backend_port}"
//  listener_port = "${var.listener_port}"
//  lb_cidr = "${var.lb_cidr}"
//  vpc_id = "${module.network.vpc_id}"
//  backend_count = "${local.backend_count}"
//  server_ids = "${local.server_ids}"
//  lb_direction = "${local.lb_direction}"
//}
//
//output "elb_ip" {
//  value = "${module.load_balancer.elb_eip}"
//}
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