storage_source "file" {
  path = "backup/vault-ddb"
}

storage_destination "dynamodb" {
  table = "defnnn"
  region = "us-west-1"
}
