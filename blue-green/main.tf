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