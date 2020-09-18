SHELL := /bin/bash

menu:
	@perl -ne 'printf("%10s: %s\n","$$1","$$2") if m{^([\w+-]+):[^#]+#\s(.+)$$}' Makefile

setup:
	docker-compose rm -f -s
	docker-compose up -d --remove-orphans

restore:
	set -a; source .env; set +a; $(MAKE) restore-inner

restore-inner:
	mkdir -p etc/traefik/acme
	pass kitt/$(KITT_DOMAIN)/acme.json | base64 -d > etc/traefik/acme/acme.json
	chmod 0600 etc/traefik/acme/acme.json
