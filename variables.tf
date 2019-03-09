variable "private_key" {}
variable "public_key" {}
variable "project_id" {}
variable "region" {}

variable "blue_az" {}
variable "blue_subnet_cidr" {
  default = "10.0.1.0/24"
}
variable "blue_image_id" {}
variable "blue_tag" {}
variable "green_az" {}
variable "green_subnet_cidr" {
  default = "10.0.0.0/24"
}
variable "green_image_id" {}
variable "green_tag" {}
variable "vpc_cidr" {
  default = "10.0.0.0/16"
}
variable "vpc_name" {
  default = "example_vpc"
}
variable "security_group_name" {
  default = "example_sg"
}

//variable "server_count" {}
variable "blue_instance_type" {
  default = "n-basic-1"
}
variable "green_instance_type" {
  default = "n-basic-1"
}
variable "root_password" {}
variable "server_name_prefix" {
  default = "example"
}
