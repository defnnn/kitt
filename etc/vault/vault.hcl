ui = true

storage "consul" {
  address = "consul.kitt.run:443"
  scheme  = "https"
  path    = "vault"

  disable_registration = "true"
}

listener "tcp" {
  address = "0.0.0.0:8200"
  tls_disable = 1
}

disable_mlock = true
