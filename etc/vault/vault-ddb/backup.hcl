storage_source "dynamodb" {
  table = "defnnn"
  region = "us-west-1"
}

storage_destination "file" {
  path = "backup/vault-ddb"
}
