variable "lb_name" {
  default = "example_bg_lb"
}
variable "lb_direction" {}
variable "vpc_id" {}
variable "lb_cidr" {}
variable "listener_port" {}
variable "backend_port" {}
variable "backend_count" {}
variable "server_ids" {
  type = "list"
}