service {
  name = "traefik"
  id = "traefik"
  address = "YYYY"
  port = 80

  checks {
    name = "Traefik port TCP connect"
    tcp = "YYYY:80"
    interval = "10s"
  }
  
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
