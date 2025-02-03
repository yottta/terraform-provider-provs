terraform {
  required_providers {
    provs = {
      source = "registry.opentofu.org/edu/provs"
    }
  }
}

provider "provs" {
  path = "/var/tmp/custom_tf_provider"
}

data "provs_coffees" "example" {}
