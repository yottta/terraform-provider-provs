terraform {
  required_providers {
    hashicups = {
      source = "hashicorp.com/edu/hashicups"
    }
  }
}

provider "hashicups" {
  path = "/var/tmp/custom_tf_provider"
}

data "hashicups_coffees" "example" {}
