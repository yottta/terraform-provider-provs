terraform {
  required_providers {
    hashicups = {
      source = "hashicorp.com/edu/hashicups"
    }
  }
  required_version = ">= 1.8.0"
}

provider "hashicups" {
  path = "/var/tmp/custom_tf_provider"
}

output "total_price" {
  value = provider::hashicups::compute_tax(5.00, 0.085)
}

