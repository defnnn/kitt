SHELL := /bin/bash

.PHONY: test

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
	if ! test -d backup/.; then mkdir backup || true; fi
	exec/kitt-setup
	kitt recreate
	$(MAKE) wait-vault-ready
	$(MAKE) unseal
	$(MAKE) wait-vault-unseal

teardown:
	$(MAKE) seal
	kitt down

backup-consul:
	exec/kitt-dc1 snapshot save backup/consul/consul-$(shell date +%s)

restart-vault:
	kitt restart vault
	$(MAKE) unseal

wait-vault-unseal:
	@set -x; while true; do if [[ "$$(vault status -format json | jq -r '.sealed')" == "false" ]]; then break; fi; date; sleep 1; done

wait-vault-ready:
	@set -x; while true; do if curl -sS https://vault.$(KITT_DOMAIN) | grep /ui/; then break; fi; date; sleep 1; done

root-login:
	@vault login "$(shell pass moria/root-token)" >/dev/null

seal:
	$(MAKE) root-login
	vault operator seal

env:
	@pyinfra @local scripts/env.py

init:
	@pyinfra @local scripts/init.py
	$(MAKE) unseal

unseal:
	@pyinfra @local scripts/unseal.py

daemon.json: fixed-cidr-v6
	$(MAKE) daemon.json-inner

daemon.json-inner:
	@jq -n --arg cidr "$(shell cat fixed-cidr-v6)" '{debug: true, experimental: true, "default-address-pools": [{base:"172.31.0.0/16","size":24}], ipv6: true, "fixed-cidr-v6": $$cidr}' | jq . > daemon.json.1 && mv daemon.json.1 daemon.json
	@rm -f fixed-cidr-v6
	@cat daemon.json

fixed-cidr-v6:
	@echo $(shell docker-compose exec zerotier zerotier-cli listnetworks | tail -n +2 | head -1 | awk '{print $$9}' | cut -d, -f1 | cut -d/ -f1 | cut -b1-12)$(shell docker-compose exec zerotier cut -c 1-2 /var/lib/zerotier-one/identity.public):$(shell docker-compose exec zerotier cut -c 3-6 /var/lib/zerotier-one/identity.public):$(shell docker-compose exec zerotier cut -c 7-10 /var/lib/zerotier-one/identity.public)::/80 > fixed-cidr-v6.1 && mv fixed-cidr-v6.1 fixed-cidr-v6
