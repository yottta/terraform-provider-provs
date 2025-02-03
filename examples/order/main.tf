terraform {
  required_providers {
    provs = {
      source  = "registry.opentofu.org/edu/provs"
    }
  }
  required_version = ">= 1.1.0"
}

provider "provs" {
  path = "/var/tmp/custom_tf_provider"
}

resource "provs_order" "new_order" {
  items = [{
    coffee = {
      id = 3
    }
    quantity = 2
    }, {
    coffee = {
      id = 1
    }
    quantity = 2
    }
  ]
}

output "order" {
  value = provs_order.new_order
}

