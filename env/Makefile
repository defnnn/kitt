SHELL := /bin/bash

menu:
	@perl -ne 'printf("%10s: %s\n","$$1","$$2") if m{^([\w+-]+):[^#]+#\s(.+)$$}' Makefile

clean:
	docker-compose down --remove-orphans

up:
	docker-compose up -d --remove-orphans

down:
	docker-compose rm -f -s

recreate:
	$(MAKE) clean
	$(MAKE) auth
	$(MAKE) up

logs:
	docker-compose logs -f

auth:
	rm -f vault/token
	v read -field=role_id auth/approle/role/jenkins/role-id  > vault/role-id
	v write -wrap-ttl=180s -field=wrapping_token -f auth/approle/role/jenkins/secret-id > vault/secret-id
