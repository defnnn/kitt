service {
  name = "dc0"
  id = "dc0"
  port = 8888
  address = "169.254.32.0"

  connect { 
    sidecar_service {
      port = 20000
      address = "169.254.32.0"

      check {
        name = "Connect Envoy Sidecar"
        tcp = "169.254.32.0:20000"
        interval ="10s"
      }
      
      proxy {
        upstreams {
          destination_name = "goodbye"
          local_bind_address = "127.0.0.1"
          local_bind_port = 9091
        }
      }
    }
  }
}
