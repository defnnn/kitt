service {
  name = "traefik"
  id = "traefik"
  address = "YYYY"
  port = 80
  
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
