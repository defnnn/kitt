storage_source "file" {
  path = "backup/vault-ddb"
}

storage_destination "s3" {
  bucket = "defnnn"
  region = "us-west-1"
}
