ui = true

storage "consul" {
  address = "consul:8500"
  scheme  = "http"
  path    = "vault"

  disable_registration = "true"
}

listener "tcp" {
  address = "0.0.0.0:8200"
  tls_disable = 1
}

disable_mlock = true
