output "desired_blue_count" {
  value = "${data.external.arbiter.result.desiredBlueCount}"
}

output "desired_green_count" {
  value = "${data.external.arbiter.result.desiredGreenCount}"
}

output "lbDirection" {
  value = "${data.external.arbiter.result.loadBalanceDirection}"
}