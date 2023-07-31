<a href="https://pritunl.com">
    <img src="https://pritunl.com/img/logo.png" alt="Pritunl logo" title="Pritunl" align="right" height="100" />
</a>
<a href="https://terraform.io">
    <img src="https://dashboard.snapcraft.io/site_media/appmedia/2019/11/terraform.png" alt="Terraform logo" title="Terraform" align="right" height="100" />
</a>

# Terraform Provider for Pritunl VPN Server

[![Release](https://img.shields.io/github/v/release/matthisholleville/terraform-provider-pritunl)](https://github.com/matthisholleville/terraform-provider-pritunl/releases)
[![Registry](https://img.shields.io/badge/registry-doc%40latest-lightgrey?logo=terraform)](https://registry.terraform.io/providers/disc/pritunl/latest/docs)
[![License](https://img.shields.io/badge/License-MPL%202.0-brightgreen.svg)](https://github.com/matthisholleville/terraform-provider-pritunl/blob/master/LICENSE)  
[![Go Report Card](https://goreportcard.com/badge/github.com/matthisholleville/terraform-provider-pritunl)](https://goreportcard.com/report/github.com/disc/terraform-provider-pritunl)

- Website: https://www.terraform.io
- Pritunl VPN Server: https://pritunl.com/
- Provider: [matthisholleville/pritunl](https://registry.terraform.io/providers/matthisholleville/pritunl/latest)

## Requirements
-	[Terraform](https://www.terraform.io/downloads.html) >=0.13.x
-	[Go](https://golang.org/doc/install) 1.18.x (to build the provider plugin)

## Building The Provider

```sh
$ git clone git@github.com:matthisholleville/terraform-provider-pritunl
$ make build
```

## Example usage

Take a look at the examples in the [documentation](https://registry.terraform.io/providers/matthisholleville/pritunl/latest/docs) of the registry
or use the following example:


```hcl
# Set the required provider and versions
terraform {
  required_providers {
    pritunl = {
      source  = "matthisholleville/pritunl"
      version = "0.0.1"
    }
  }
}

# Configure the pritunl provider
provider "pritunl" {
  url    = "https://vpn.server.com"
  token  = "api-token"
  secret = "api-secret"
  insecure = false
}

# Create a pritunl organization resource
resource "pritunl_organization" "developers" {
  name = "Developers"
}

# Create a pritunl server resource
resource "pritunl_server" "example" {
  name      = "example"
  port      = 15500
  protocol  = "udp"
  network   = "192.168.1.0/24"
  groups    = [
    "admins",
    "developers",
  ]
  
  # Attach the organization to the server
  organization_ids = [
    pritunl_organization.developers.id,
  ]
}

# Create a pritunl server route
resource "pritunl_route" "example" {
  server_id = pritunl_server.example.id
  nat       = true
  comment   = "my custom route"
  network   = "34.56.43.32/32"
}
```

## Importing exist resources

Describe exist resource in the terraform file first and then import them:

Import an organization:
```hcl
# Describe a pritunl organization resource
resource "pritunl_organization" "developers" {
  name = "Developers"
}
```

Execute the shell command:
```sh
terraform import pritunl_organization.developers ${ORGANIZATION_ID}
terraform import pritunl_organization.developers 610e42d2a0ed366f41dfe6e8
```
The organization ID (as well as other resource IDs) can be found in the Pritunl API responses or in the HTML document response.

Import a server:

```hcl
# Describe a pritunl server resource
resource "pritunl_server" "example" {
  name      = "example"
  port      = 15500
  protocol  = "udp"
  network   = "192.168.1.0/24"
  groups    = [
    "developers",
  ]

  # Attach the organization to the server
  organization_ids = [
    pritunl_organization.developers.id,
  ]
}
```

Execute the shell command:
```sh
terraform import pritunl_server.example ${SERVER_ID}
terraform import pritunl_server.example 60cd0bfa7723cf3c911468a8
```

Import a route:

```hcl
# Describe a pritunl route resource
resource "pritunl_route" "example" {
  server_id = pritunl_server.example.id
  nat       = true
  comment   = "my custom route"
  network   = "34.56.43.32/32"
}
```

Execute the shell command:
```sh
terraform import pritunl_route.example ${ROUTE_ID}
terraform import pritunl_route.example 60cd0bfa7723cf3c911468a8
```

## License

The Terraform Pritunl Provider is available to everyone under the terms of the Mozilla Public License Version 2.0. [Take a look the LICENSE file](LICENSE).
