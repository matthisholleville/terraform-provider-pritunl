provider "pritunl" {
  url      = "https://vpn.server.com"
  token    = "api-token"
  secret   = "api-secret"
  insecure = false
}

data "pritunl_server" "server" {
  server_name = "Devops_test"
}

resource "pritunl_route" "test" {
  server_id = data.pritunl_server.server.id
  nat       = true
  comment   = "my custom route"
  network   = "34.56.43.32/32"
}