output "desired_blue_count" {
  value = "${data.external.arbiter.result.desiredBlueCount}"
}

output "desired_green_count" {
  value = "${data.external.arbiter.result.desiredGreenCount}"
}

output "lb_direction" {
  value = "${data.external.arbiter.result.loadBalanceDirection}"
}