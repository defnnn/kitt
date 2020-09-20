SHELL := /bin/bash

menu:
	@perl -ne 'printf("%10s: %s\n","$$1","$$2") if m{^([\w+-]+):[^#]+#\s(.+)$$}' Makefile

kitt:
	$(MAKE) setup
	$(MAKE) up
	$(MAKE) kumactl

clean:
	docker-compose down
	docker network rm kitt || true

setup:
	$(MAKE) network || true
	$(MAKE) build

build:
	docker-compose build

network:
	docker network create kitt

dummy:
	sudo ip link add dummy0 type dummy || true
	sudo ip addr add 169.254.32.1/32 dev dummy0 || true
	sudo ip link set dev dummy0 up

up:
	docker-compose rm -f -s
	docker-compose up -d --remove-orphans

restore:
	set -a; source .env; set +a; $(MAKE) restore-inner

restore-inner:
	mkdir -p etc/traefik/acme
	pass kitt/$(KITT_DOMAIN)/acme.json | base64 -d > etc/traefik/acme/acme.json
	chmod 0600 etc/traefik/acme/acme.json
	pass kitt/$(KITT_DOMAIN)/authtoken.secret | base64 -d | perl -pe 's{\s*$$}{}'  > etc/zerotier/zerotier-one/authtoken.secret
	pass kitt/$(KITT_DOMAIN)/identity.public | base64 -d | perl -pe 's{\s*$$}{}' > etc/zerotier/zerotier-one/identity.public
	pass kitt/$(KITT_DOMAIN)/identity.secret | base64 -d | perl -pe 's{\s*$$}{}' > etc/zerotier/zerotier-one/identity.secret
	pass kitt/$(KITT_DOMAIN)/hook-customize| base64 -d > etc/zerotier/hooks/hook-customize
	chmod 755 etc/zerotier/hooks/hook-customize
	pass kitt/$(KITT_DOMAIN)/cert.pem | base64 -d > etc/cloudflared/cert.pem
	pass kitt/$(KITT_DOMAIN)/env | base64 -d > .env

kumactl:
	kumactl config control-planes add --address http://169.254.32.1:5681 --name kitt --overwrite
	kumactl config control-planes switch --name kitt
