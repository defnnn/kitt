storage_source "consul" {
  address = "consul.kitt.run:443"
  scheme  = "https"
  path    = "vault"
}

storage_destination "file" {
  path = "backup/vault-consul"
}
