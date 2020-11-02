SHELL := /bin/bash

menu:
	@perl -ne 'printf("%10s: %s\n","$$1","$$2") if m{^([\w+-]+):[^#]+#\s(.+)$$}' Makefile

clean:
	docker-compose down --remove-orphans

setup:
	$(MAKE) up

watch:
	./bin/dns-update

up:
	docker-compose up -d --remove-orphans

down:
	docker-compose rm -f -s

recreate:
	-$(MAKE) clean
	$(MAKE) up
