SHELL := /bin/bash

menu:
	@perl -ne 'printf("%10s: %s\n","$$1","$$2") if m{^([\w+-]+):[^#]+#\s(.+)$$}' Makefile

setup once:
	kitt setup
	kitt recreate

env:
	@pyinfra @local scripts/env.py

init:
	@pyinfra @local scripts/init.py
	$(MAKE) unseal

unseal:
	@pyinfra @local scripts/unseal.py

api-tunnel:
	socat TCP4-LISTEN:8888,fork TCP4:traefik-kitt.$(KITT_DOMAIN):80
