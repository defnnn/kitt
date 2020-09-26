SHELL := /bin/bash

menu:
	@perl -ne 'printf("%10s: %s\n","$$1","$$2") if m{^([\w+-]+):[^#]+#\s(.+)$$}' Makefile

thing:
	$(MAKE) setup
	$(MAKE) up
	sudo rsync -ia work/kuma/bin/. /usr/local/bin/.
	sleep 10
	$(MAKE) kuma

clean:
	docker-compose down

setup:
	$(MAKE) network || true
	$(MAKE) dummy
	$(MAKE) build

build:
	docker-compose build

network:
	docker network create kitt || true
	docker network create --subnet 172.25.0.0/16 --ip-range 172.25.1.0/24 kind || true

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
	pass kitt/$(KITT_DOMAIN)/authtoken.secret | base64 -d | perl -pe 's{\s*$$}{}'  > etc/zerotier/zerotier-one/authtoken.secret
	pass kitt/$(KITT_DOMAIN)/identity.public | base64 -d | perl -pe 's{\s*$$}{}' > etc/zerotier/zerotier-one/identity.public
	pass kitt/$(KITT_DOMAIN)/identity.secret | base64 -d | perl -pe 's{\s*$$}{}' > etc/zerotier/zerotier-one/identity.secret
	pass kitt/$(KITT_DOMAIN)/acme.json | base64 -d > etc/traefik/acme/acme.json
	chmod 0600 etc/traefik/acme/acme.json
	pass kitt/$(KITT_DOMAIN)/hook-customize| base64 -d > etc/zerotier/hooks/hook-customize
	chmod 755 etc/zerotier/hooks/hook-customize
	pass kitt/$(KITT_DOMAIN)/cert.pem | base64 -d > etc/cloudflared/cert.pem
	pass kitt/$(KITT_DOMAIN)/env | base64 -d > .env

restore-diff:
	set -a; source .env; set +a; $(MAKE) restore-diff-inner

restore-diff-inner:
	pdif kitt/$(KITT_DOMAIN)/authtoken.secret etc/zerotier/zerotier-one/authtoken.secret
	pdif kitt/$(KITT_DOMAIN)/identity.public etc/zerotier/zerotier-one/identity.public
	pdif kitt/$(KITT_DOMAIN)/identity.secret etc/zerotier/zerotier-one/identity.secret
	pdiff kitt/$(KITT_DOMAIN)/acme.json etc/traefik/acme/acme.json
	pdiff kitt/$(KITT_DOMAIN)/hook-customize etc/zerotier/hooks/hook-customize
	pdiff kitt/$(KITT_DOMAIN)/cert.pem etc/cloudflared/cert.pem
	pdiff kitt/$(KITT_DOMAIN)/env .env

kuma:
	$(MAKE) kumactl
	kumactl apply -f k/traffic-permission-allow-all-traffic.yaml
	kumactl apply -f k/mesh-default.yaml

kumactl:
	kumactl config control-planes add --address http://$(shell docker inspect kitt_kuma_1 | jq -r '.[].NetworkSettings.Networks.kind.IPAddress'):5681 --name kitt --overwrite
	kumactl config control-planes switch --name kitt
