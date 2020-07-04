SHELL := /bin/bash

.PHONY: test backup

menu:
	@perl -ne 'printf("%10s: %s\n","$$1","$$2") if m{^([\w+-]+):[^#]+#\s(.+)$$}' Makefile

test: # Run tests
	cd test && $(MAKE)
	cd test && git diff

drone-test: # Run tests with drone specific setup
	mkdir /tmp/src
	rsync -ia . /tmp/src/.
	cd /tmp/src/test && $(MAKE)
	cd /tmp/src/test && git diff

setup once:
	if ! test -f etc/vault/vault/vault.hcl; then ln -nfs vault-consul etc/vault/vault; fi
	if ! test -d backup/.; then mkdir backup || true; fi
	exec/kitt-setup
	kitt recreate
	$(MAKE) unseal
	$(MAKE) wait

teardown:
	$(MAKE) seal
	kitt down

dc0:
	consul agent -config-file="$(PWD)/etc/consul_config/dc0.hcl" -data-dir="$(PWD)/etc/consul_dc0" -join-wan=169.254.32.1

dc0-gateway:
	set -a; . .env; set +a; exec/kitt-dc0 connect envoy -mesh-gateway -register -address 169.254.32.0:4444

dc0-proxy:
	set -a; . .env; set +a; exec/kitt-dc0 connect envoy -sidecar-for dc0 -admin-bind 127.0.0.1:19001

dc0-test:
	@curl http://localhost:9091

migrate-ddb migrate-s3 backup:
	$(MAKE) seal
	vault operator migrate -config etc/vault/vault/$@.hcl
	$(MAKE) restart-vault
	$(MAKE) wait

restart-vault:
	kitt restart vault
	$(MAKE) unseal

wait:
	@set -x; while true; do if [[ "$$(vault status -format json | jq -r '.sealed')" == "false" ]]; then break; fi; date; sleep 1; done

root-login:
	@vault login "$(shell pass moria/root-token)" >/dev/null

seal:
	$(MAKE) root-login
	vault operator seal

init:
	@pyinfra @local scripts/init.py
	$(MAKE) unseal

unseal:
	@pyinfra @local scripts/unseal.py

consul ddb s3 file-consul file-ddb file-s3 file:
	ln -nfs vault-$@ etc/vault/vault
	$(MAKE) restart-vault
	$(MAKE) wait
