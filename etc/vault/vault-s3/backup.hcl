storage_source "s3" {
  bucket = "defnnn"
  region = "us-west-1"
}

storage_destination "file" {
  path = "backup/vault-s3"
}
