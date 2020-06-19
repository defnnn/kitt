data_dir = "/tmp/"
log_level = "DEBUG"

datacenter = "dc1"
primary_datacenter = "dc1"

server = true

bootstrap_expect = 1
ui = true

bind_addr = "172.31.188.99"
client_addr = "172.31.188.99"

ports {
  grpc = 8502
}

connect {
  enabled = true
}

enable_central_service_config = true
