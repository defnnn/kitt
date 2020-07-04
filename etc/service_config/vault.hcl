service {
  name = "vault"
  id = "vault"
  address = "YYYY"
  port = 8200
  
  connect { 
    sidecar_service {
      port = 20000
      
      check {
        name = "Connect Envoy Sidecar"
        tcp = "YYYY:20000"
        interval ="10s"
      }
    }  
  }
}
