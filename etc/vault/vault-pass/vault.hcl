ui = true

storage "raft" {
  path = "/vault"
  node_id = "raft_node_1"
}

cluster_addr = "http://127.0.0.1:8201"
api_addr = "http://127.0.0.1:8200"

listener "tcp" {
  address = "0.0.0.0:8200"
  tls_disable = 1
}

disable_mlock = true
