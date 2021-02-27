pi_file = "/vault/pid"

exit_after_auth = false

listener "tcp" {
  address = "0.0.0.0:8200"
  tls_disable = true
}

cache {
  use_auto_auth_token = true
}

auto_auth {
  method {
    type = "approle"
    config = {
      role_id_file_path = "/vault/role-id"
      secret_id_file_path = "/vault/secret-id"
      secret_id_response_wrapping_path = "auth/approle/role/jenkins/secret-id"
    }
  }
}

template {
  source      = "/vault-template/consul.tmpl"
  destination = "/vault-env/consul/.env"
}

template {
  source      = "/vault-template/cilium.tmpl"
  destination = "/vault-env/cilium/.env"
}

template {
  source      = "/vault-template/kitt.tmpl"
  destination = "/vault-env/kitt/.env"
}

template {
  source      = "/vault-template/zerotier.tmpl"
  destination = "/vault-env/zerotier/.env"
}

template {
  source      = "/vault-template/traefik.tmpl"
  destination = "/vault-env/traefik/.env"
}

template {
  source      = "/vault-template/traefik-auth.tmpl"
  destination = "/vault-env/traefik-auth/.env"
}
