data "external" "arbiter" {
  program = [
    "go",
    "run",
    "${path.module}/main.go",
    "-operation=${var.operation}",
    "-currentBlue=${var.current_blue_count}",
    "-currentGreen=${var.current_green_count}",
    "-desiredBlue=${var.desired_blue_count}",
    "-desiredGreen=${var.desired_green_count}"
  ]
}

output "desired_blue_count" {
  value = "${data.external.arbiter.result.desiredBlueCount}"
}

output "desired_green_count" {
  value = "${data.external.arbiter.result.desiredGreenCount}"
}

output "lbDirection" {
  value = "${data.external.arbiter.result.loadBalanceDirection}"
}