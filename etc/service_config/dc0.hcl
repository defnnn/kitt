service {
  name = "dc0"
  id = "dc0"
  port = 8888

  connect { 
    sidecar_service {
      port = 20000
      
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
