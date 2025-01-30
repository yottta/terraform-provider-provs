# Terraform Provider Scaffolding (Terraform Plugin Framework)
Started from [terraform-provider-scaffolding-framework](https://github.com/hashicorp/terraform-provider-scaffolding-framework).
Playing around with some examples and basic implementation.


To play around with this, you can follow the guide [here](https://developer.hashicorp.com/terraform/tutorials/providers-plugin-framework/providers-plugin-framework-provider).


# Quick start
* Run `cd docker_compose && docker compose up -d` - this is starting a service that the provider will talk with and where the actual "resources" will live.
* Run `curl -X POST localhost:19090/signup -d '{"username":"education", "password":"test123"}'` that will create a new user that is being used later in the tofu scripts.
  * Keep the response as you will need it later, especially for the `HASHICUPS_TOKEN` env var.
* Go over the examples in the [examples](./examples) dir and play around with different topics.
