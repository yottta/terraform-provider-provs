terraform {
  required_providers {
    provs = {
      source = "registry.opentofu.org/edu/provs"
    }
  }
  required_version = ">= 1.8.0"
}

provider "provs" {
  path = "/var/tmp/custom_tf_provider"
}

output "total_price" {
  value = provider::provs::compute_tax(5.00, 0.085)
}

