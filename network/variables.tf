variable "vpc_name" {}
variable "vpc_tag" {
  default = ""
}
variable "vpc_cidr" {}
variable "vpc_remark" {
  default = ""
}

variable "green_subnet_cidr" {}
variable "green_subnet_name" {
  default = "green_subnet"
}
variable "green_subnet_tag" {
  default = ""
}
variable "green_subnet_remark" {
  default = ""
}

variable "blue_subnet_cidr" {}
variable "blue_subnet_name" {
  default = "blue_subnet"
}
variable "blue_subnet_tag" {
  default = ""
}
variable "blue_subnet_remark" {
  default = ""
}
variable "security_group_name" {
  default = "example_sg"
}