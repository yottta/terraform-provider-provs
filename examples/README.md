# Examples

## Ensure that the plugin is working well
Just to be sure that the provider is working fine, just run `go run main.go`. The message that you should get should be similar to this one:
```
This binary is a plugin. These are not meant to be executed directly.
Please execute the program that consumes these plugins, which will
load any plugins automatically
exit status 1
```

## Ensure that the plugin is installed correctly ([provider-install-verification](./examples/provider-install-verification) dir)
Run `go install .` to have the plugin installed in your $GOBIN (or `go env GOBIN`). Running `ls -lah "$(go env GOBIN)" | grep "terraform-provider-"` should yield your newly installed plugin.

Add a new file `~/.terraformrc` with a content similar to this one:
```
provider_installation {

  dev_overrides {
    "hashicorp.com/edu/hashicups" = "<the path that $GOBIN is pointing to>"
  }

  # For all other providers, install them directly from their origin provider
  # registries as normal. If you omit this, Terraform will _only_ use
  # the dev_overrides block, and so no other providers will be available.
  direct {}
}
```

Go to [provider-install-verification](./examples/provider-install-verification) and run `tofu plan`.

You should get a response like this:
```
╷
│ Warning: Provider development overrides are in effect
│
│ The following provider development overrides are set in the CLI configuration:
│  - hashicorp.com/edu/hashicups in <your GOBIN path>
│
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become
│ incompatible with published releases.
╵
data.hashicups_coffees.example: Reading...
data.hashicups_coffees.example: Read complete after 0s

No changes. Your infrastructure matches the configuration.
```
## Play around with data sources ([coffees](./examples/coffees) dir
Run `tofu plan` and you should see something like this:
```

╷
│ Warning: Provider development overrides are in effect
│
│ The following provider development overrides are set in the CLI configuration:
│  - hashicorp.com/edu/hashicups in <your GOBIN path>
│
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become
│ incompatible with published releases.
╵
data.hashicups_coffees.edu: Reading...
data.hashicups_coffees.edu: Read complete after 0s

Changes to Outputs:
  + edu_coffees = {
      + coffees = [
          + {
              + description = ""
              + id          = 1
              + image       = "/hashicorp.png"
              + ingredients = [
                  + {
                      + id = 6
                    },
                ]
              + name        = "HCP Aeropress"
              + price       = 200
              + teaser      = "Automation in a cup"
            },
          + {
              + description = ""
              + id          = 2
...
```

This is a data source, this is just getting some hardcoded information from the docker app ran earlier.
## Create, update, destroy, import resources ([order](./examples/order) dir)
### Create
First step is to apply the changes and see what happens.
* The `computed` values will be shown as "known after apply"
* `id` of the order is "known after apply"
* Checking the actual server, we can see the order in the response
  * `curl -X GET  -H "Authorization: ${HASHICUPS_TOKEN}" localhost:19090/orders`

### Update
Edit [order/main.tf](./order/main.tf) and update the id of the second coffe from the order:
```

resource "hashicups_order" "edu" {
  items = [{
    coffee = {
      id = 3
    }
    quantity = 2
    }, {
    coffee = {
      id = 2 <- here (from 1 to 2)
    }
    quantity = 2
    }
  ]
}
```

Run `tofu apply` and this is highlighting the following:
* The `id` of the order is known now and shown accordingly
* The values of the updated order will be "known after apply"
* Once applied, it succeeds:
  * Check the state with `tofu show`
  * Check the service to see that the order is updated `curl -X GET  -H "Authorization: ${HASHICUPS_TOKEN}" localhost:19090/orders`

### Delete
Just run `tofu destroy` and the order should be deleted from the actual server.
Check this by checking again `curl -X GET  -H "Authorization: ${HASHICUPS_TOKEN}" localhost:19090/orders`. The response should be an empty object.

### Import
Assuming that the guide in this file was followed step by step, you should have no resources created anymore.
Therefore, let's create one by running `tofu apply --auto-approve`.

By running `tofu show` can be seen that the resource is created and saved into our state (use the `curl` command above to be sure that the resource is created on the server too).

Now, let's remove the resource from the state by running `terraform state rm hashicups_order.edu`.
Having this removed from the state, by running `tofu show` can be seen that the resource is not tracked by OpenTofu anymore (only the output will be in the state).

To import a resource, we need its `id`. Do to so, run `curl -X GET -H "Authorization: ${HASHICUPS_TOKEN}" localhost:19090/orders` and get the id of the object in there (probably id=2).
With that id, run `tofu import hashicups_order.edu 2`.

Running again `tofu show`, can be observed that in the state we now have referenced the order with id 2.

## Test how a provider can expose functions ([compute_tax](./examples/compute_tax))
This is just a simple example on how a provider can expose some functions.
Just run `terraform apply -auto-approve` and you should see an output similar to this one:
```
%> tofu apply -auto-approve

╷
│ Warning: Provider development overrides are in effect
│
│ The following provider development overrides are set in the CLI configuration:
│  - hashicorp.com/edu/hashicups in <your GOBIN path>
│
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become
│ incompatible with published releases.
╵

No changes. Your infrastructure matches the configuration.

OpenTofu has compared your real infrastructure against your configuration and found no differences, so no changes are needed.

Apply complete! Resources: 0 added, 0 changed, 0 destroyed.

Outputs:

total_price = 5.43                                                                                                                    
```

