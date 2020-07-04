ui = true

storage "file" {
  path = "/vault/vault-ddb"
}

listener "tcp" {
  address = "0.0.0.0:8200"
  tls_disable = 1
}

seal "awskms" {
  region = "us-west-1"
  kms_key_id = "3e40c099-05ee-4462-b8dc-06ed3849d648"
}

disable_mlock = true
