path "kv/data/*" {
  capabilities = ["create", "update", "read", "delete", "sudo"]
}

path "kv/metadata/*" {
  capabilities = ["list", "read", "create", "update", "delete"]
}
