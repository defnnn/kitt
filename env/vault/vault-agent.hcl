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
  source      = "/vault-template/kitt.tmpl"
  destination = "/vault-env/kitt/.env"
}

template {
  source      = "/vault-template/cilium.tmpl"
  destination = "/vault-env/cilium/.env"
}
