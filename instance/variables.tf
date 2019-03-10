variable "count" {}
variable "az" {}
variable "image_id" {}
variable "instance_type" {}
variable "root_password" {}
variable "tag" {}
variable "name_prefix" {}
variable "charge_type" {
  default = "dynamic"
}
variable "security_group_id" {}
variable "subnet_id" {}
variable "vpc_id" {}
variable "duration" {
  default = 1
}
