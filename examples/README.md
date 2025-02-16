# Examples

## Ensure that the plugin is working well
Just to be sure that the provider is working fine, just run `go run main.go`. The message that you should get should be similar to this one:
```
This binary is a plugin. These are not meant to be executed directly.
Please execute the program that consumes these plugins, which will
load any plugins automatically
exit status 1
```

## Ensure that the plugin is installed correctly ([provider-install-verification](./provider-install-verification) dir)
Run `go install .` to have the plugin installed in your $GOBIN (or `go env GOBIN`). Running `ls -lah "$(go env GOBIN)" | grep "terraform-provider-"` should yield your newly installed plugin.

Add a new file `~/.terraformrc` with a content similar to this one:
```
provider_installation {

  dev_overrides {
      "registry.opentofu.org/edu/provs" = "<the path that $GOBIN is pointing to>"
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
│  - edu/provs in "<the path that $GOBIN is pointing to>"
│
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with published releases.
╵
data.provs_coffees.example: Reading...
data.provs_coffees.example: Read complete after 0s

No changes. Your infrastructure matches the configuration.

OpenTofu has compared your real infrastructure against your configuration and found no differences, so no changes are needed.
```
## Play around with data sources ([coffees](./coffees) dir
Run `tofu plan` and you should see something like this:
```

╷
│ Warning: Provider development overrides are in effect
│
│ The following provider development overrides are set in the CLI configuration:
│  - edu/provs in "<the path that $GOBIN is pointing to>"
│
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with published releases.
╵
data.provs_coffees.test_coffees: Reading...
data.provs_coffees.test_coffees: Read complete after 0s

Changes to Outputs:
  + coffees = {
      + coffees = [
          + {
              + description = "Description 9"
              + id          = 9
              + image       = ""
              + ingredients = [
                  + {
                      + id = 0
...                     
```

This is a data source, this is just getting some hardcoded information from the provider.
The data can be seen by running `ls -lah /var/tmp/custom_tf_provider/coffees/`.
## Create, update, destroy, import resources ([order](./order) dir)
### Create
First step is to apply the changes and see what happens.
* The `computed` values will be shown as "known after apply"
* `id` of the order is "known after apply"
* Checking the filesystem to see that everything is created correctly
  * `ls -lah /var/tmp/custom_tf_provider/order/`

### Update
Edit [order/main.tf](./order/main.tf) and update the id of the second coffe from the order:
```

resource "provs_order" "edu" {
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
  * Check the filesystem to ensure that the order was created `ls -lah /var/tmp/custom_tf_provider/order/`.

### Delete
Just run `tofu destroy` and the order should be deleted from the actual server.
Check this by checking again `ls -lah /var/tmp/custom_tf_provider/order/`. The directory should be empty.

### Import
Assuming that the guide in this file was followed step by step, you should have no resources created anymore.
Therefore, let's create one by running `tofu apply --auto-approve`.

By running `tofu show` can be seen that the resource is created and saved into our state.

Now, let's remove the resource from the state by running `tofu state rm provs_order.new_order`.
Having this removed from the state, by running `tofu show` can be seen that the resource is not tracked by OpenTofu anymore (only the output will be in the state).

To import a resource, we need its `id`. Do to so, run `ls -lah /var/tmp/custom_tf_provider/order/` and get the file name visible there.
With that id, run `tofu import provs_order.new_order <file name>`.

Running again `tofu show`, can be observed that in the state we now have referenced the order with the id provided above.

## Test how a provider can expose functions ([compute_tax](./compute_tax))
This is just a simple example on how a provider can expose some functions.
Just run `terraform apply -auto-approve` and you should see an output similar to this one:
```
%> tofu apply -auto-approve

╷
│ Warning: Provider development overrides are in effect
│
│ The following provider development overrides are set in the CLI configuration:
│  - edu/provs in <your GOBIN path>
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

