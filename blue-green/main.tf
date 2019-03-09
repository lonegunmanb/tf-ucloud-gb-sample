
locals {
  i2g = "i2g"
  i2b = "i2b"
  g2s = "g2s"
  b2s = "b2s"
  s2g = "s2g"
  s2b = "s2b"

  valid_opeartions = [
    "${local.i2g}",
    "${local.i2b}",
    "${local.s2b}",
    "${local.b2s}",
    "${local.g2s}",
    "${local.s2g}",
  ]
  blue = "blue"
  green = "green"
  init = "init"
  staging = "staging"
}

resource "null_resource" "state_migration_check" {
  count = "${contains(local.valid_opeartions, var.operation) ? 0 : 1}"
  "ERROR: cannot execute operation ${var.operation}, operation must in ${join(",", local.valid_opeartions)}" = true
}

locals {
  from_state = "${var.operation == local.i2b || var.operation == local.i2g ? local.init
                :
                 (var.operation == local.s2b || var.operation == local.s2g ? local.staging
                :(var.operation == local.b2s ? local.blue : local.green))}"
  to_state = "${var.operation == local.i2b || var.operation == local.s2b ? local.blue
                :
                 (var.operation == local.i2g || var.operation == local.s2g ? local.green
                : local.staging)}"
}

locals {
  operation_thumb = "${format("%s,%s,%s,%s,%s", var.operation, var.current_blue_count, var.current_green_count, var.desired_blue_count, var.desired_green_count)}"
}

resource "null_resource" "YOU_MUST_FROM_INIT_STATE" {
  provisioner "local-exec" {
    command = "${(local.from_state != local.init && var.current_blue_count == 0 && var.current_green_count == 0) ? "false" : "echo"}"
  }
  triggers {
    operation = "${local.operation_thumb}"
  }
  depends_on = ["null_resource.state_migration_check"]
}

resource "null_resource" "CHECK_FROM_INIT_STATE" {
  provisioner "local-exec" {
    command = "${(local.from_state == local.init && (var.current_green_count > 0 || var.current_blue_count > 0)) == true ? "false" : "echo"}"
  }
  triggers {
    operation = "${local.operation_thumb}"
  }
  depends_on = ["null_resource.YOU_MUST_FROM_INIT_STATE"]
}

resource "null_resource" "CHECK_FROM_BLUE_STATE" {
  provisioner "local-exec" {
    command = "${(local.from_state == local.blue && (var.current_green_count > 0 || var.current_blue_count == 0)) ? "false" : "echo"}"
  }
  triggers {
    operation = "${local.operation_thumb}"
  }
  depends_on = ["null_resource.CHECK_FROM_INIT_STATE"]
}

resource "null_resource" "CHECK_FROM_GREEN_STATE" {
  provisioner "local-exec" {
    command = "${(local.from_state == local.green && (var.current_blue_count > 0 || var.current_green_count == 0)) ? "false" : "echo"}"
  }
  triggers {
    operation = "${local.operation_thumb}"
  }
  depends_on = ["null_resource.CHECK_FROM_BLUE_STATE"]
}

resource "null_resource" "CHECK_FROM_STAGING_STAGE" {
  provisioner "local-exec" {
    command = "${(local.from_state == local.staging &&(var.current_blue_count == 0 || var.current_green_count == 0)) == true ? "false" : "echo"}"
  }
  triggers {
    operation = "${local.operation_thumb}"
  }
  depends_on = ["null_resource.CHECK_FROM_GREEN_STATE"]
}

resource "null_resource" "CHECK_TO_BLUE" {
  provisioner "local-exec" {
    command = "${(local.to_state == local.blue && var.desired_blue_count == 0) == true ? "false" : "echo"}"
  }
  triggers {
    operation = "${local.operation_thumb}"
  }
  depends_on = ["null_resource.CHECK_FROM_STAGING_STAGE"]
}

resource "null_resource" "CHECK_TO_GREEN" {
  provisioner "local-exec" {
    command = "${(local.to_state == local.green && var.desired_green_count == 0) == true ? "false" : "echo"}"
  }
  triggers {
    operation = "${local.operation_thumb}"
  }
  depends_on = ["null_resource.CHECK_TO_BLUE"]
}

resource "null_resource" "CHECK_TO_STAGING" {
  provisioner "local-exec" {
    command = "${local.to_state == local.staging && (var.desired_blue_count == 0 || var.desired_green_count == 0) == true ? "false" : "echo"}"
  }
  triggers {
    operation = "${local.operation_thumb}"
  }
  depends_on = ["null_resource.CHECK_TO_GREEN"]
}

locals  {
  lb_target = "${ local.to_state == local.blue ? "blue"
            :
            (local.to_state == local.green ? "green"
            :
            (var.operation == local.b2s ? "blue" : "green"))}"
}

output "lb_target" {
  value = "${local.lb_target}"
}

output "desired_blue_count" {
  value = "${ local.to_state == local.blue ? var.desired_blue_count
            :
            (local.to_state == local.green ? 0
            :
            (var.operation == local.b2s ? var.current_blue_count : var.desired_blue_count))}"
}

output "desired_green_count" {
  value = "${ local.to_state == local.blue ? 0
            :
            (local.to_state == local.green ? var.desired_green_count
            :
            (var.operation == local.g2s ? var.current_green_count : var.desired_green_count))}"
}

output "signal" {
  value = "${null_resource.CHECK_TO_STAGING.id}"
}
//
//output "current_blue_count" {
//  value = "${var.current_blue_count}"
//}
//
//output "current_green_count" {
//  value = "${var.current_green_count}"
//}