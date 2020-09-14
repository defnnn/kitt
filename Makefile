SHELL := /bin/bash

menu:
	@perl -ne 'printf("%10s: %s\n","$$1","$$2") if m{^([\w+-]+):[^#]+#\s(.+)$$}' Makefile

setup:
	docker-compose rm -f -s
	docker-compose up -d --remove-orphans

api-tunnel:
	set -a; . .env; set +a; socat TCP4-LISTEN:8888,fork TCP4:traefik.$${KITT_DOMAIN}:80
